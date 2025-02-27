package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"gorelay/pkg/models"
)

// MonitorServer provides an HTTP server for debugging and monitoring
type MonitorServer struct {
	server     *http.Server
	port       int
	clients    map[string]*ClientInfo
	mu         sync.RWMutex
	handlers   map[string]http.HandlerFunc
	logs       map[string][]LogEntry
	startTime  time.Time
	isRunning  bool
	lastUpdate time.Time
	cpuUsage   float64
	memUsage   uint64
	cmdHandler func(alias string, command string) string // Handler for client commands
}

// ClientInfo contains information about a connected client
type ClientInfo struct {
	Account    *models.Account
	Connected  time.Time
	LastSeen   time.Time
	CurrentMap string
	Stats      map[string]interface{}
	IsLoggedIn bool
}

// LogEntry represents a log message
type LogEntry struct {
	Timestamp time.Time
	Level     string // "info", "warning", "error"
	Message   string
	Account   string
}

// NewMonitorServer creates a new monitoring server instance
func NewMonitorServer(port int) *MonitorServer {
	ms := &MonitorServer{
		port:      port,
		clients:   make(map[string]*ClientInfo),
		handlers:  make(map[string]http.HandlerFunc),
		logs:      make(map[string][]LogEntry),
		startTime: time.Now(),
	}

	// Set up default handlers
	ms.handlers["/"] = ms.handleDashboard
	ms.handlers["/clients"] = ms.handleClients
	ms.handlers["/status"] = ms.handleStatus
	ms.handlers["/logs"] = ms.handleLogs
	ms.handlers["/api/logs"] = ms.handleAPILogs
	ms.handlers["/api/status"] = ms.handleAPIStatus
	ms.handlers["/api/accounts"] = ms.handleAPIAccounts
	ms.handlers["/static/"] = http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP
	ms.handlers["/reconnect/"] = ms.handleReconnect
	ms.handlers["/command/"] = ms.handleCommand

	return ms
}

// Start starts the monitoring server
func (ms *MonitorServer) Start() error {
	if ms.isRunning {
		return fmt.Errorf("monitor server is already running")
	}

	ms.startTime = time.Now()
	ms.lastUpdate = time.Now()

	// Start metrics collection in background
	go ms.collectMetrics()

	mux := http.NewServeMux()

	// Register handlers
	for path, handler := range ms.handlers {
		mux.HandleFunc(path, handler)
	}

	ms.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", ms.port),
		Handler: mux,
	}

	ms.isRunning = true
	fmt.Printf("Monitor server started at http://localhost:%d\n", ms.port)

	// Start server in a goroutine
	go func() {
		if err := ms.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Monitor server error: %v\n", err)
			ms.isRunning = false
		}
	}()

	return nil
}

// Stop stops the monitoring server
func (ms *MonitorServer) Stop() error {
	if ms.server != nil {
		return ms.server.Close()
	}
	return nil
}

// AddClient adds a client to the server's tracking
func (ms *MonitorServer) AddClient(alias string, account *models.Account) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.clients[alias] = &ClientInfo{
		Account:   account,
		Connected: time.Now(),
		LastSeen:  time.Now(),
		Stats:     make(map[string]interface{}),
	}

	ms.Log(alias, "info", "Client connected")
}

// RemoveClient removes a client from tracking
func (ms *MonitorServer) RemoveClient(alias string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	delete(ms.clients, alias)
}

// RegisterHandler registers a custom HTTP handler
func (ms *MonitorServer) RegisterHandler(path string, handler http.HandlerFunc) {
	ms.handlers[path] = handler
}

// Log adds a log entry for an account
func (ms *MonitorServer) Log(alias, level, message string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Account:   alias,
	}

	if _, exists := ms.logs[alias]; !exists {
		ms.logs[alias] = make([]LogEntry, 0)
	}
	ms.logs[alias] = append(ms.logs[alias], entry)

	// Keep only last 1000 logs per account
	if len(ms.logs[alias]) > 1000 {
		ms.logs[alias] = ms.logs[alias][len(ms.logs[alias])-1000:]
	}
}

// UpdateClientStatus updates a client's status
func (ms *MonitorServer) UpdateClientStatus(alias string, status map[string]interface{}) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if client, exists := ms.clients[alias]; exists {
		client.LastSeen = time.Now()
		for k, v := range status {
			client.Stats[k] = v
		}
	}
}

// handleDashboard serves the main dashboard page
func (ms *MonitorServer) handleDashboard(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("dashboard").Parse(dashboardHTML))
	ms.mu.RLock()
	data := struct {
		Clients map[string]*ClientInfo
		Port    int
	}{
		Clients: ms.clients,
		Port:    ms.port,
	}
	ms.mu.RUnlock()
	tmpl.Execute(w, data)
}

// handleLogs serves the logs page
func (ms *MonitorServer) handleLogs(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("logs").Parse(logsHTML))
	ms.mu.RLock()
	data := struct {
		Logs map[string][]LogEntry
		Port int
	}{
		Logs: ms.logs,
		Port: ms.port,
	}
	ms.mu.RUnlock()
	tmpl.Execute(w, data)
}

// handleAPILogs returns logs in JSON format
func (ms *MonitorServer) handleAPILogs(w http.ResponseWriter, r *http.Request) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ms.logs)
}

// handleAPIStatus returns status in JSON format
func (ms *MonitorServer) handleAPIStatus(w http.ResponseWriter, r *http.Request) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ms.clients)
}

// handleAPIAccounts returns a list of all configured account aliases
func (ms *MonitorServer) handleAPIAccounts(w http.ResponseWriter, r *http.Request) {
	// Get all account aliases from the models package
	aliases := models.GetAllAccountAliases()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(aliases)
}

// setupStaticFiles creates necessary static files
func (ms *MonitorServer) setupStaticFiles() error {
	// TODO: Create static directory and write CSS/JS files
	return nil
}

// collectMetrics periodically updates system metrics
func (ms *MonitorServer) collectMetrics() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if !ms.isRunning {
			return
		}

		ms.mu.Lock()
		ms.lastUpdate = time.Now()

		// Get memory stats
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		ms.memUsage = memStats.Alloc

		// Get CPU usage (simplified version)
		startCPU := time.Now()
		runtime.GC()
		cpuTime := time.Since(startCPU)
		ms.cpuUsage = float64(cpuTime.Microseconds()) / 1000.0 // Convert to milliseconds

		ms.mu.Unlock()
	}
}

func (ms *MonitorServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	uptime := time.Since(ms.startTime)
	hours := int(uptime.Hours())
	minutes := int(uptime.Minutes()) % 60
	seconds := int(uptime.Seconds()) % 60

	lastHeartbeat := time.Since(ms.lastUpdate)
	isOnline := lastHeartbeat < 30*time.Second

	status := map[string]interface{}{
		"uptime": fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds),
		"uptimeDetails": map[string]int{
			"hours":   hours,
			"minutes": minutes,
			"seconds": seconds,
		},
		"clientCount":   len(ms.clients),
		"memoryUsage":   fmt.Sprintf("%.2f MB", float64(ms.memUsage)/1024/1024),
		"cpuUsage":      fmt.Sprintf("%.2f ms", ms.cpuUsage),
		"lastHeartbeat": lastHeartbeat.String(),
		"isOnline":      isOnline,
	}

	json.NewEncoder(w).Encode(status)
}

const dashboardHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>GoRelay Monitor</title>
    <style>
        :root {
            --bg-color: #1a1a1a;
            --card-bg: #2d2d2d;
            --text-color: #e0e0e0;
            --accent-color: #3a3a3a;
            --border-color: #404040;
            --success-color: #4caf50;
            --error-color: #f44336;
            --info-color: #2196f3;
            --warning-color: #ff9800;
        }
        
        body { 
            font-family: 'Consolas', monospace;
            margin: 20px;
            background: var(--bg-color);
            color: var(--text-color);
        }
        
        .client { 
            border: 1px solid var(--border-color);
            padding: 15px;
            margin: 10px 0;
            background: var(--card-bg);
            border-radius: 5px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.2);
        }
        
        .status { 
            color: var(--success-color);
            font-weight: bold;
        }
        
        .status.offline { 
            color: var(--error-color);
        }
        
        .stats { 
            margin-left: 20px;
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
            gap: 10px;
        }
        
        .stat-item {
            padding: 5px;
            background: var(--accent-color);
            border-radius: 3px;
            font-family: 'Consolas', monospace;
        }
        
        nav {
            background: var(--card-bg);
            padding: 10px;
            margin-bottom: 20px;
            border-radius: 5px;
            border: 1px solid var(--border-color);
        }
        
        nav a {
            color: var(--text-color);
            text-decoration: none;
            padding: 5px 10px;
        }
        
        nav a:hover {
            background: var(--accent-color);
            border-radius: 3px;
        }
        
        h1, h2, h3 { 
            color: var(--text-color);
        }

        .terminal {
            background: #000;
            border: 1px solid var(--border-color);
            border-radius: 5px;
            padding: 10px;
            margin-top: 10px;
            font-family: 'Consolas', monospace;
            height: 200px;
            overflow-y: auto;
            color: #00ff00;
        }

        .reconnect-btn {
            background: var(--accent-color);
            color: var(--text-color);
            border: 1px solid var(--border-color);
            border-radius: 3px;
            padding: 5px 10px;
            cursor: pointer;
            font-family: 'Consolas', monospace;
            margin-top: 10px;
            transition: background-color 0.2s;
        }

        .reconnect-btn:hover {
            background: var(--card-bg);
        }

        .reconnect-btn:active {
            transform: translateY(1px);
        }

        .reconnect-btn.reconnecting {
            opacity: 0.7;
            cursor: wait;
        }

        .terminal::-webkit-scrollbar {
            width: 8px;
        }

        .terminal::-webkit-scrollbar-track {
            background: #000;
        }

        .terminal::-webkit-scrollbar-thumb {
            background: var(--accent-color);
            border-radius: 4px;
        }

        .terminal-line {
            margin: 2px 0;
            white-space: pre-wrap;
            word-wrap: break-word;
        }

        .system-stats {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
            gap: 10px;
            margin-bottom: 20px;
            padding: 15px;
            background: var(--card-bg);
            border-radius: 5px;
            border: 1px solid var(--border-color);
        }

        .stat-box {
            padding: 10px;
            background: var(--accent-color);
            border-radius: 3px;
            text-align: center;
        }

        .stat-label {
            font-size: 0.9em;
            opacity: 0.8;
        }

        .stat-value {
            font-size: 1.2em;
            margin-top: 5px;
            font-weight: bold;
        }

        .server-status {
            position: fixed;
            top: 10px;
            right: 10px;
            padding: 10px;
            border-radius: 5px;
            font-weight: bold;
            z-index: 1000;
        }

        .server-status.online {
            background: var(--success-color);
            color: #fff;
        }

        .server-status.offline {
            background: var(--error-color);
            color: #fff;
        }

        .uptime-detail {
            display: inline-block;
            padding: 0 5px;
            font-weight: bold;
        }
        
        .histograph-container {
            background: var(--card-bg);
            border: 1px solid var(--border-color);
            border-radius: 5px;
            padding: 15px;
            margin-bottom: 20px;
        }
        
        .histograph-title {
            margin-top: 0;
            margin-bottom: 10px;
            font-size: 1.2em;
        }
        
        .histograph {
            height: 150px;
            display: flex;
            align-items: flex-end;
            gap: 2px;
            margin-top: 10px;
            position: relative;
        }
        
        .histograph-bar {
            flex: 1;
            background: var(--info-color);
            min-width: 3px;
            transition: height 0.3s ease;
            position: relative;
        }
        
        .histograph-bar:hover::after {
            content: attr(data-value);
            position: absolute;
            bottom: 100%;
            left: 50%;
            transform: translateX(-50%);
            background: var(--accent-color);
            padding: 3px 6px;
            border-radius: 3px;
            font-size: 0.8em;
            white-space: nowrap;
            z-index: 10;
        }
        
        .histograph-axis {
            position: absolute;
            left: 0;
            right: 0;
            border-top: 1px dashed var(--border-color);
        }
        
        .histograph-axis-label {
            position: absolute;
            left: -5px;
            transform: translateY(-50%);
            font-size: 0.8em;
            color: var(--text-color);
            opacity: 0.7;
        }
        
        .histograph-legend {
            display: flex;
            justify-content: space-between;
            margin-top: 5px;
            font-size: 0.8em;
            color: var(--text-color);
            opacity: 0.7;
        }
    </style>
    <script>
        // Store historical data for CPU and RAM
        let cpuHistory = Array(60).fill(0);
        let ramHistory = Array(60).fill(0);
        let connectionStatus = true;
        let lastConnectionCheck = Date.now();
        
        // Load history data from localStorage if available
        function loadHistoryData() {
            try {
                const savedCpuHistory = localStorage.getItem('cpuHistory');
                const savedRamHistory = localStorage.getItem('ramHistory');
                const savedTimestamp = localStorage.getItem('historyTimestamp');
                
                if (savedCpuHistory && savedRamHistory && savedTimestamp) {
                    const timestamp = parseInt(savedTimestamp);
                    // Only use saved data if it's less than 1 hour old
                    if (Date.now() - timestamp < 60 * 60 * 1000) {
                        cpuHistory = JSON.parse(savedCpuHistory);
                        ramHistory = JSON.parse(savedRamHistory);
                        console.log('Loaded history data from localStorage');
                    } else {
                        console.log('Saved history data is too old, using defaults');
                    }
                }
            } catch (error) {
                console.error('Error loading history data:', error);
            }
        }
        
        // Save history data to localStorage
        function saveHistoryData() {
            try {
                localStorage.setItem('cpuHistory', JSON.stringify(cpuHistory));
                localStorage.setItem('ramHistory', JSON.stringify(ramHistory));
                localStorage.setItem('historyTimestamp', Date.now().toString());
            } catch (error) {
                console.error('Error saving history data:', error);
            }
        }
        
        // Load history data on page load
        loadHistoryData();
        
        function updateStatus() {
            // Check if we're still connected to the server
            const now = Date.now();
            if (now - lastConnectionCheck > 5000) {
                connectionStatus = false;
                updateConnectionStatus();
            }
            
            Promise.all([
                fetch('/api/status').then(response => {
                    connectionStatus = true;
                    lastConnectionCheck = now;
                    return response.json();
                }),
                fetch('/status').then(response => {
                    connectionStatus = true;
                    lastConnectionCheck = now;
                    return response.json();
                })
            ])
            .then(([clientData, systemData]) => {
                updateConnectionStatus();
                
                // Update system stats with detailed uptime
                const systemStatsDiv = document.getElementById('system-stats');
                const uptimeDetails = systemData.uptimeDetails;
                const uptimeHtml = 
                    '<div class="uptime-detail">' + 
                    String(uptimeDetails.hours).padStart(2, '0') + 
                    '</div>:' +
                    '<div class="uptime-detail">' + 
                    String(uptimeDetails.minutes).padStart(2, '0') + 
                    '</div>:' +
                    '<div class="uptime-detail">' + 
                    String(uptimeDetails.seconds).padStart(2, '0') + 
                    '</div>';

                systemStatsDiv.innerHTML = 
                    '<div class="stat-box">' +
                        '<div class="stat-label">Uptime</div>' +
                        '<div class="stat-value">' + uptimeHtml + '</div>' +
                    '</div>' +
                    '<div class="stat-box">' +
                        '<div class="stat-label">Memory Usage</div>' +
                        '<div class="stat-value">' + systemData.memoryUsage + '</div>' +
                    '</div>' +
                    '<div class="stat-box">' +
                        '<div class="stat-label">CPU Usage</div>' +
                        '<div class="stat-value">' + systemData.cpuUsage + '</div>' +
                    '</div>' +
                    '<div class="stat-box">' +
                        '<div class="stat-label">Connected Clients</div>' +
                        '<div class="stat-value">' + systemData.clientCount + '</div>' +
                    '</div>';
                
                // Update histographs
                // Extract numeric values from the strings
                const cpuValue = parseFloat(systemData.cpuUsage.replace(' ms', ''));
                const ramValue = parseFloat(systemData.memoryUsage.replace(' MB', ''));
                
                // Add new values to history arrays (shift left)
                cpuHistory.push(cpuValue);
                cpuHistory.shift();
                ramHistory.push(ramValue);
                ramHistory.shift();
                
                // Save updated history data
                saveHistoryData();
                
                updateHistographs();

                // Update clients
                const clientsDiv = document.getElementById('clients');
                clientsDiv.innerHTML = '';
                
                for (const [alias, client] of Object.entries(clientData)) {
                    const clientDiv = document.createElement('div');
                    clientDiv.className = 'client';
                    
                    const lastSeen = new Date(client.LastSeen);
                    const isOffline = Date.now() - lastSeen > 30000;
                    
                    const statusClass = isOffline ? 'offline' : '';
                    const statusText = isOffline ? 'Offline' : 'Online';
                    const connectedTime = new Date(client.Connected).toLocaleString();
                    const lastSeenTime = lastSeen.toLocaleString();
                    
                    let statsHtml = '';
                    if (client.Stats) {
                        for (const [key, value] of Object.entries(client.Stats)) {
                            statsHtml += '<div class="stat-item">' + key + ': ' + value + '</div>';
                        }
                    }
                    
                    clientDiv.innerHTML = '<h3>' + alias + '</h3>' +
                        '<p class="status ' + statusClass + '">' +
                        'Status: ' + statusText + '<br>' +
                        'Connected: ' + connectedTime + '<br>' +
                        'Last Seen: ' + lastSeenTime +
                        '</p>' +
                        '<div class="stats">' + statsHtml + '</div>' +
                        '<button class="reconnect-btn" onclick="reconnectClient(\'' + alias + '\')">Reconnect</button>' +
                        '<div class="terminal" id="terminal-' + alias + '"></div>';
                        
                    clientsDiv.appendChild(clientDiv);
                }

                // Update terminals
                updateTerminals();
            })
            .catch(error => {
                console.error('Error updating status:', error);
                connectionStatus = false;
                updateConnectionStatus();
            });
            
            setTimeout(updateStatus, 1000);
        }
        
        function updateConnectionStatus() {
            const statusDiv = document.getElementById('server-status');
            if (connectionStatus) {
                statusDiv.className = 'server-status online';
                statusDiv.textContent = 'Server Online';
            } else {
                statusDiv.className = 'server-status offline';
                statusDiv.textContent = 'Server Disconnected';
            }
        }
        
        function updateHistographs() {
            // Update CPU histograph
            const cpuHistograph = document.getElementById('cpu-histograph');
            if (cpuHistograph) {
                // Find max value for scaling
                const maxCpu = Math.max(...cpuHistory, 1); // Ensure at least 1 to avoid division by zero
                
                let cpuHtml = '';
                cpuHistory.forEach(value => {
                    const height = (value / maxCpu) * 100;
                    cpuHtml += '<div class="histograph-bar" style="height: ' + height + '%;" data-value="' + value.toFixed(2) + ' ms"></div>';
                });
                
                // Add axis lines at 25%, 50%, 75% and 100%
                cpuHtml += '<div class="histograph-axis" style="bottom: 25%;">' +
                           '<span class="histograph-axis-label">' + (maxCpu * 0.75).toFixed(1) + '</span></div>';
                cpuHtml += '<div class="histograph-axis" style="bottom: 50%;">' +
                           '<span class="histograph-axis-label">' + (maxCpu * 0.5).toFixed(1) + '</span></div>';
                cpuHtml += '<div class="histograph-axis" style="bottom: 75%;">' +
                           '<span class="histograph-axis-label">' + (maxCpu * 0.25).toFixed(1) + '</span></div>';
                cpuHtml += '<div class="histograph-axis" style="bottom: 100%;">' +
                           '<span class="histograph-axis-label">0</span></div>';
                
                cpuHistograph.innerHTML = cpuHtml;
            }
            
            // Update RAM histograph
            const ramHistograph = document.getElementById('ram-histograph');
            if (ramHistograph) {
                // Find max value for scaling
                const maxRam = Math.max(...ramHistory, 1); // Ensure at least 1 to avoid division by zero
                
                let ramHtml = '';
                ramHistory.forEach(value => {
                    const height = (value / maxRam) * 100;
                    ramHtml += '<div class="histograph-bar" style="height: ' + height + '%;" data-value="' + value.toFixed(2) + ' MB"></div>';
                });
                
                // Add axis lines at 25%, 50%, 75% and 100%
                ramHtml += '<div class="histograph-axis" style="bottom: 25%;">' +
                           '<span class="histograph-axis-label">' + (maxRam * 0.75).toFixed(1) + '</span></div>';
                ramHtml += '<div class="histograph-axis" style="bottom: 50%;">' +
                           '<span class="histograph-axis-label">' + (maxRam * 0.5).toFixed(1) + '</span></div>';
                ramHtml += '<div class="histograph-axis" style="bottom: 75%;">' +
                           '<span class="histograph-axis-label">' + (maxRam * 0.25).toFixed(1) + '</span></div>';
                ramHtml += '<div class="histograph-axis" style="bottom: 100%;">' +
                           '<span class="histograph-axis-label">0</span></div>';
                
                ramHistograph.innerHTML = ramHtml;
            }
        }

        function updateTerminals() {
            fetch('/api/logs')
                .then(response => response.json())
                .then(data => {
                    for (const [alias, logs] of Object.entries(data)) {
                        const terminal = document.getElementById('terminal-' + alias);
                        if (terminal) {
                            let terminalHtml = '';
                            logs.forEach(log => {
                                const timestamp = new Date(log.Timestamp).toLocaleString();
                                const color = {
                                    'info': '#2196f3',
                                    'warning': '#ff9800',
                                    'error': '#f44336'
                                }[log.Level] || '#e0e0e0';
                                
                                terminalHtml += '<div class="terminal-line" style="color: ' + color + '">[' + timestamp + '] ' + log.Message + '</div>';
                            });
                            terminal.innerHTML = terminalHtml;
                            terminal.scrollTop = terminal.scrollHeight;
                        }
                    }
                })
                .catch(error => console.error('Error updating terminals:', error));
        }

        function reconnectClient(alias) {
            const btn = document.querySelector('button[onclick="reconnectClient(\'' + alias + '\')"]');
            if (btn) {
                btn.classList.add('reconnecting');
                btn.textContent = 'Reconnecting...';
                btn.disabled = true;
            }

            fetch('/reconnect/' + encodeURIComponent(alias), { method: 'POST' })
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Reconnect failed');
                    }
                    return response.text();
                })
                .then(() => {
                    if (btn) {
                        btn.classList.remove('reconnecting');
                        btn.textContent = 'Reconnect';
                        btn.disabled = false;
                    }
                })
                .catch(error => {
                    console.error('Error reconnecting client:', error);
                    if (btn) {
                        btn.classList.remove('reconnecting');
                        btn.textContent = 'Reconnect Failed';
                        btn.disabled = false;
                        setTimeout(() => {
                            btn.textContent = 'Reconnect';
                        }, 2000);
                    }
                });
        }

        // Initialize connection status
        document.addEventListener('DOMContentLoaded', function() {
            updateConnectionStatus();
        });

        updateStatus();
    </script>
</head>
<body>
    <div id="server-status" class="server-status">Server Status</div>
    <h1>GoRelay Monitor</h1>
    <nav>
        <a href="/">Dashboard</a> |
        <a href="/logs">Logs</a>
    </nav>
    <div id="system-stats" class="system-stats"></div>
    
    <div class="histograph-container">
        <h3 class="histograph-title">CPU Usage (Last Hour)</h3>
        <div id="cpu-histograph" class="histograph"></div>
        <div class="histograph-legend">
            <span>60 minutes ago</span>
            <span>Now</span>
        </div>
    </div>
    
    <div class="histograph-container">
        <h3 class="histograph-title">Memory Usage (Last Hour)</h3>
        <div id="ram-histograph" class="histograph"></div>
        <div class="histograph-legend">
            <span>60 minutes ago</span>
            <span>Now</span>
        </div>
    </div>
    
    <div id="clients"></div>
</body>
</html>
`

const logsHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>GoRelay Logs</title>
    <style>
        :root {
            --bg-color: #1a1a1a;
            --card-bg: #2d2d2d;
            --text-color: #e0e0e0;
            --accent-color: #3a3a3a;
            --border-color: #404040;
            --success-color: #4caf50;
            --error-color: #f44336;
            --info-color: #2196f3;
            --warning-color: #ff9800;
        }
        
        body { 
            font-family: 'Consolas', monospace;
            margin: 20px;
            background: var(--bg-color);
            color: var(--text-color);
        }
        
        .log { 
            margin: 5px 0;
            padding: 5px;
            border-radius: 3px;
            font-family: 'Consolas', monospace;
        }
        
        .log.info { 
            color: var(--info-color);
            background: rgba(33, 150, 243, 0.1);
        }
        
        .log.warning { 
            color: var(--warning-color);
            background: rgba(255, 152, 0, 0.1);
        }
        
        .log.error { 
            color: var(--error-color);
            background: rgba(244, 67, 54, 0.1);
        }
        
        nav {
            background: var(--card-bg);
            padding: 10px;
            margin-bottom: 20px;
            border-radius: 5px;
            border: 1px solid var(--border-color);
        }
        
        nav a {
            color: var(--text-color);
            text-decoration: none;
            padding: 5px 10px;
        }
        
        nav a:hover {
            background: var(--accent-color);
            border-radius: 3px;
        }
        
        .account-logs {
            background: var(--card-bg);
            padding: 15px;
            margin: 10px 0;
            border-radius: 5px;
            border: 1px solid var(--border-color);
            box-shadow: 0 2px 4px rgba(0,0,0,0.2);
        }
        
        h1, h3 {
            color: var(--text-color);
        }

        .terminal {
            background: #000;
            border: 1px solid var(--border-color);
            border-radius: 5px;
            padding: 10px;
            margin-top: 10px;
            font-family: 'Consolas', monospace;
            height: 300px;
            overflow-y: auto;
        }

        .terminal::-webkit-scrollbar {
            width: 8px;
        }

        .terminal::-webkit-scrollbar-track {
            background: #000;
        }

        .terminal::-webkit-scrollbar-thumb {
            background: var(--accent-color);
            border-radius: 4px;
        }

        .terminal-line {
            margin: 2px 0;
            white-space: pre-wrap;
            word-wrap: break-word;
        }
        
        .client-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 10px;
        }
        
        .client-status {
            display: inline-block;
            padding: 5px 10px;
            border-radius: 3px;
            font-size: 0.9em;
            font-weight: bold;
        }
        
        .client-status.online {
            background: var(--success-color);
            color: #fff;
        }
        
        .client-status.offline {
            background: var(--error-color);
            color: #fff;
        }
        
        .client-status.connecting {
            background: var(--warning-color);
            color: #fff;
        }
        
        .reconnect-btn {
            background: var(--accent-color);
            color: var(--text-color);
            border: 1px solid var(--border-color);
            border-radius: 3px;
            padding: 5px 10px;
            cursor: pointer;
            font-family: 'Consolas', monospace;
            transition: background-color 0.2s;
        }

        .reconnect-btn:hover {
            background: var(--card-bg);
        }

        .reconnect-btn:active {
            transform: translateY(1px);
        }

        .reconnect-btn.reconnecting {
            opacity: 0.7;
            cursor: wait;
        }
        
        .log-controls {
            display: flex;
            gap: 10px;
            margin-bottom: 10px;
        }
        
        .log-filter {
            background: var(--accent-color);
            color: var(--text-color);
            border: 1px solid var(--border-color);
            border-radius: 3px;
            padding: 5px 10px;
            cursor: pointer;
            font-family: 'Consolas', monospace;
        }
        
        .log-filter.active {
            background: var(--info-color);
            color: #fff;
        }
        
        .server-status {
            position: fixed;
            top: 10px;
            right: 10px;
            padding: 10px;
            border-radius: 5px;
            font-weight: bold;
            z-index: 1000;
        }

        .server-status.online {
            background: var(--success-color);
            color: #fff;
        }

        .server-status.offline {
            background: var(--error-color);
            color: #fff;
        }
        
        .client-info {
            margin-top: 5px;
            font-size: 0.9em;
            color: var(--text-color);
            opacity: 0.8;
        }
        
        .command-input-container {
            display: flex;
            margin-top: 10px;
            gap: 5px;
        }
        
        .command-input {
            flex: 1;
            background: var(--accent-color);
            color: var(--text-color);
            border: 1px solid var(--border-color);
            border-radius: 3px;
            padding: 8px 10px;
            font-family: 'Consolas', monospace;
        }
        
        .command-input:focus {
            outline: none;
            border-color: var(--info-color);
        }
        
        .command-send {
            background: var(--accent-color);
            color: var(--text-color);
            border: 1px solid var(--border-color);
            border-radius: 3px;
            padding: 8px 15px;
            cursor: pointer;
            font-family: 'Consolas', monospace;
        }
        
        .command-send:hover {
            background: var(--card-bg);
        }
        
        .command-send:active {
            transform: translateY(1px);
        }
        
        .command-response {
            margin-top: 5px;
            padding: 5px;
            font-size: 0.9em;
            color: var(--info-color);
            background: rgba(33, 150, 243, 0.1);
            border-radius: 3px;
            display: none;
        }
    </style>
    <script>
        let connectionStatus = true;
        let lastConnectionCheck = Date.now();
        let activeFilters = {
            info: true,
            warning: true,
            error: true
        };
        
        function updateLogs() {
            // Check if we're still connected to the server
            const now = Date.now();
            if (now - lastConnectionCheck > 5000) {
                connectionStatus = false;
                updateConnectionStatus();
            }
            
            Promise.all([
                fetch('/api/logs').then(response => {
                    connectionStatus = true;
                    lastConnectionCheck = now;
                    return response.json();
                }),
                fetch('/api/status').then(response => {
                    connectionStatus = true;
                    lastConnectionCheck = now;
                    return response.json();
                }),
                fetch('/api/accounts').then(response => {
                    connectionStatus = true;
                    lastConnectionCheck = now;
                    return response.json();
                })
            ])
            .then(([logsData, clientsData, accountAliases]) => {
                updateConnectionStatus();
                
                const logsDiv = document.getElementById('logs');
                
                // Get existing client divs
                const existingClients = {};
                document.querySelectorAll('.account-logs').forEach(div => {
                    const alias = div.getAttribute('data-alias');
                    existingClients[alias] = div;
                });
                
                // Create or update client divs
                const newLogsDiv = document.createElement('div');
                
                // Track which aliases we've processed
                const processedAliases = new Set();
                
                // Process all clients with status data first (even if they don't have logs)
                for (const [alias, client] of Object.entries(clientsData)) {
                    let clientDiv;
                    
                    if (existingClients[alias]) {
                        // Update existing div
                        clientDiv = existingClients[alias].cloneNode(true);
                        delete existingClients[alias]; // Remove from the list of existing clients
                    } else {
                        // Create new div
                        clientDiv = createClientDiv(alias);
                    }
                    
                    // Update client status
                    updateClientStatus(clientDiv, alias, client);
                    
                    // Update terminal with logs if available
                    const logs = logsData[alias] || [];
                    updateClientLogs(clientDiv, alias, logs);
                    
                    newLogsDiv.appendChild(clientDiv);
                    processedAliases.add(alias);
                }
                
                // Now process any clients that have logs but no status data
                for (const [alias, logs] of Object.entries(logsData)) {
                    // Skip if we already processed this client
                    if (processedAliases.has(alias)) continue;
                    
                    let clientDiv;
                    
                    if (existingClients[alias]) {
                        // Update existing div
                        clientDiv = existingClients[alias].cloneNode(true);
                        delete existingClients[alias]; // Remove from the list of existing clients
                    } else {
                        // Create new div
                        clientDiv = createClientDiv(alias);
                    }
                    
                    // Update client status (as disconnected)
                    const statusSpan = clientDiv.querySelector('#status-' + alias);
                    const infoDiv = clientDiv.querySelector('#info-' + alias);
                    
                    if (statusSpan) {
                        statusSpan.className = 'client-status offline';
                        statusSpan.textContent = 'Disconnected';
                    }
                    
                    if (infoDiv) {
                        infoDiv.innerHTML = 'No connection data available';
                    }
                    
                    // Update terminal with logs
                    updateClientLogs(clientDiv, alias, logs);
                    
                    newLogsDiv.appendChild(clientDiv);
                    processedAliases.add(alias);
                }
                
                // Finally, process any configured accounts that don't have logs or status
                for (const alias of accountAliases) {
                    // Skip if we already processed this client
                    if (processedAliases.has(alias)) continue;
                    
                    let clientDiv;
                    
                    if (existingClients[alias]) {
                        // Update existing div
                        clientDiv = existingClients[alias].cloneNode(true);
                        delete existingClients[alias]; // Remove from the list of existing clients
                    } else {
                        // Create new div
                        clientDiv = createClientDiv(alias);
                    }
                    
                    // Update client status (as connecting)
                    const statusSpan = clientDiv.querySelector('#status-' + alias);
                    const infoDiv = clientDiv.querySelector('#info-' + alias);
                    
                    if (statusSpan) {
                        statusSpan.className = 'client-status connecting';
                        statusSpan.textContent = 'Connecting';
                    }
                    
                    if (infoDiv) {
                        infoDiv.innerHTML = 'Client is attempting to connect...';
                    }
                    
                    // Update terminal with empty logs
                    updateClientLogs(clientDiv, alias, []);
                    
                    newLogsDiv.appendChild(clientDiv);
                    processedAliases.add(alias);
                }
                
                // Replace the logs div content
                if (logsDiv.innerHTML !== newLogsDiv.innerHTML) {
                    logsDiv.innerHTML = newLogsDiv.innerHTML;
                }
            })
            .catch(error => {
                console.error('Error updating logs:', error);
                connectionStatus = false;
                updateConnectionStatus();
            });
            
            setTimeout(updateLogs, 1000);
        }
        
        // Helper function to create a new client div
        function createClientDiv(alias) {
            const clientDiv = document.createElement('div');
            clientDiv.className = 'account-logs';
            clientDiv.setAttribute('data-alias', alias);
            
            // Create header with client name and controls
            const headerDiv = document.createElement('div');
            headerDiv.className = 'client-header';
            
            const heading = document.createElement('h3');
            heading.textContent = alias;
            
            const statusSpan = document.createElement('span');
            statusSpan.className = 'client-status';
            statusSpan.id = 'status-' + alias;
            
            headerDiv.appendChild(heading);
            headerDiv.appendChild(statusSpan);
            
            // Create log controls
            const controlsDiv = document.createElement('div');
            controlsDiv.className = 'log-controls';
            
            const infoFilter = document.createElement('button');
            infoFilter.className = 'log-filter' + (activeFilters.info ? ' active' : '');
            infoFilter.textContent = 'Info';
            infoFilter.onclick = function() { toggleFilter('info'); };
            
            const warningFilter = document.createElement('button');
            warningFilter.className = 'log-filter' + (activeFilters.warning ? ' active' : '');
            warningFilter.textContent = 'Warning';
            warningFilter.onclick = function() { toggleFilter('warning'); };
            
            const errorFilter = document.createElement('button');
            errorFilter.className = 'log-filter' + (activeFilters.error ? ' active' : '');
            errorFilter.textContent = 'Error';
            errorFilter.onclick = function() { toggleFilter('error'); };
            
            const reconnectBtn = document.createElement('button');
            reconnectBtn.className = 'reconnect-btn';
            reconnectBtn.textContent = 'Reconnect';
            reconnectBtn.onclick = function() { reconnectClient(alias); };
            
            controlsDiv.appendChild(infoFilter);
            controlsDiv.appendChild(warningFilter);
            controlsDiv.appendChild(errorFilter);
            controlsDiv.appendChild(reconnectBtn);
            
            // Create client info div
            const infoDiv = document.createElement('div');
            infoDiv.className = 'client-info';
            infoDiv.id = 'info-' + alias;
            
            // Create terminal
            const terminal = document.createElement('div');
            terminal.className = 'terminal';
            terminal.id = 'terminal-' + alias;
            
            // Create command input
            const commandContainer = document.createElement('div');
            commandContainer.className = 'command-input-container';
            
            const commandInput = document.createElement('input');
            commandInput.type = 'text';
            commandInput.className = 'command-input';
            commandInput.id = 'command-input-' + alias;
            commandInput.placeholder = 'Enter command...';
            commandInput.addEventListener('keydown', function(e) {
                if (e.key === 'Enter') {
                    sendCommand(alias);
                }
            });
            
            const commandButton = document.createElement('button');
            commandButton.className = 'command-send';
            commandButton.textContent = 'Send';
            commandButton.onclick = function() { sendCommand(alias); };
            
            commandContainer.appendChild(commandInput);
            commandContainer.appendChild(commandButton);
            
            // Create response area
            const responseDiv = document.createElement('div');
            responseDiv.className = 'command-response';
            responseDiv.id = 'command-response-' + alias;
            
            clientDiv.appendChild(headerDiv);
            clientDiv.appendChild(controlsDiv);
            clientDiv.appendChild(infoDiv);
            clientDiv.appendChild(terminal);
            clientDiv.appendChild(commandContainer);
            clientDiv.appendChild(responseDiv);
            
            return clientDiv;
        }
        
        // Helper function to update client status
        function updateClientStatus(clientDiv, alias, client) {
            const statusSpan = clientDiv.querySelector('#status-' + alias);
            const infoDiv = clientDiv.querySelector('#info-' + alias);
            
            if (client) {
                const lastSeen = new Date(client.LastSeen);
                const isOffline = Date.now() - lastSeen > 30000;
                
                if (statusSpan) {
                    statusSpan.className = 'client-status ' + (isOffline ? 'offline' : 'online');
                    statusSpan.textContent = isOffline ? 'Offline' : 'Online';
                }
                
                if (infoDiv) {
                    const connectedTime = new Date(client.Connected).toLocaleString();
                    const lastSeenTime = lastSeen.toLocaleString();
                    
                    infoDiv.innerHTML = 
                        'Connected: ' + connectedTime + '<br>' +
                        'Last Seen: ' + lastSeenTime;
                        
                    if (client.CurrentMap) {
                        infoDiv.innerHTML += '<br>Current Map: ' + client.CurrentMap;
                    }
                }
            } else {
                if (statusSpan) {
                    statusSpan.className = 'client-status offline';
                    statusSpan.textContent = 'Disconnected';
                }
                
                if (infoDiv) {
                    infoDiv.innerHTML = 'No connection data available';
                }
            }
        }
        
        // Helper function to update client logs
        function updateClientLogs(clientDiv, alias, logs) {
            const terminal = clientDiv.querySelector('#terminal-' + alias);
            if (!terminal) return;
            
            // Get current scroll position and check if scrolled to bottom
            const wasAtBottom = terminal.scrollHeight - terminal.clientHeight <= terminal.scrollTop + 5;
            
            let terminalHtml = '';
            logs.forEach(log => {
                // Skip if this log level is filtered out
                if (!activeFilters[log.Level]) return;
                
                const timestamp = new Date(log.Timestamp).toLocaleString();
                const color = {
                    'info': '#2196f3',
                    'warning': '#ff9800',
                    'error': '#f44336'
                }[log.Level] || '#e0e0e0';
                
                terminalHtml += '<div class="terminal-line" style="color: ' + color + '">[' + timestamp + '] ' + log.Message + '</div>';
            });
            
            // If no logs, show a message
            if (terminalHtml === '') {
                terminalHtml = '<div class="terminal-line" style="color: #888;">No logs available for this client.</div>';
            }
            
            terminal.innerHTML = terminalHtml;
            
            // Restore scroll position if it was at the bottom
            if (wasAtBottom) {
                terminal.scrollTop = terminal.scrollHeight;
            }
        }
        
        function updateConnectionStatus() {
            const statusDiv = document.getElementById('server-status');
            if (connectionStatus) {
                statusDiv.className = 'server-status online';
                statusDiv.textContent = 'Server Online';
            } else {
                statusDiv.className = 'server-status offline';
                statusDiv.textContent = 'Server Disconnected';
            }
        }
        
        function toggleFilter(level) {
            activeFilters[level] = !activeFilters[level];
            
            // Update filter button styles
            document.querySelectorAll('.log-filter').forEach(btn => {
                if (btn.textContent.toLowerCase() === level) {
                    if (activeFilters[level]) {
                        btn.classList.add('active');
                    } else {
                        btn.classList.remove('active');
                    }
                }
            });
            
            // Force update of all terminals
            document.querySelectorAll('.terminal').forEach(terminal => {
                const alias = terminal.id.replace('terminal-', '');
                updateTerminal(alias);
            });
        }
        
        function updateTerminal(alias) {
            fetch('/api/logs')
                .then(response => response.json())
                .then(data => {
                    const logs = data[alias];
                    if (!logs) return;
                    
                    const terminal = document.getElementById('terminal-' + alias);
                    if (!terminal) return;
                    
                    // Get current scroll position and check if scrolled to bottom
                    const wasAtBottom = terminal.scrollHeight - terminal.clientHeight <= terminal.scrollTop + 5;
                    
                    let terminalHtml = '';
                    logs.forEach(log => {
                        // Skip if this log level is filtered out
                        if (!activeFilters[log.Level]) return;
                        
                        const timestamp = new Date(log.Timestamp).toLocaleString();
                        const color = {
                            'info': '#2196f3',
                            'warning': '#ff9800',
                            'error': '#f44336'
                        }[log.Level] || '#e0e0e0';
                        
                        terminalHtml += '<div class="terminal-line" style="color: ' + color + '">[' + timestamp + '] ' + log.Message + '</div>';
                    });
                    
                    terminal.innerHTML = terminalHtml;
                    
                    // Restore scroll position if it was at the bottom
                    if (wasAtBottom) {
                        terminal.scrollTop = terminal.scrollHeight;
                    }
                })
                .catch(error => console.error('Error updating terminal:', error));
        }
        
        function sendCommand(alias) {
            const inputElement = document.getElementById('command-input-' + alias);
            const responseElement = document.getElementById('command-response-' + alias);
            
            if (!inputElement || !responseElement) return;
            
            const command = inputElement.value.trim();
            if (!command) return;
            
            // Disable input while sending
            inputElement.disabled = true;
            
            // Create form data
            const formData = new FormData();
            formData.append('command', command);
            
            // Send command to server
            fetch('/command/' + encodeURIComponent(alias), {
                method: 'POST',
                body: formData
            })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Failed to send command');
                }
                return response.json();
            })
            .then(data => {
                // Show response
                responseElement.textContent = data.response;
                responseElement.style.display = 'block';
                
                // Clear input
                inputElement.value = '';
                
                // Hide response after 5 seconds
                setTimeout(() => {
                    responseElement.style.display = 'none';
                }, 5000);
            })
            .catch(error => {
                console.error('Error sending command:', error);
                responseElement.textContent = 'Error: ' + error.message;
                responseElement.style.display = 'block';
                
                // Hide error after 5 seconds
                setTimeout(() => {
                    responseElement.style.display = 'none';
                }, 5000);
            })
            .finally(() => {
                // Re-enable input
                inputElement.disabled = false;
                inputElement.focus();
            });
        }
        
        function reconnectClient(alias) {
            const btn = document.querySelector('.account-logs[data-alias="' + alias + '"] .reconnect-btn');
            if (btn) {
                btn.classList.add('reconnecting');
                btn.textContent = 'Reconnecting...';
                btn.disabled = true;
            }

            fetch('/reconnect/' + encodeURIComponent(alias), { method: 'POST' })
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Reconnect failed');
                    }
                    return response.text();
                })
                .then(() => {
                    if (btn) {
                        btn.classList.remove('reconnecting');
                        btn.textContent = 'Reconnect';
                        btn.disabled = false;
                    }
                })
                .catch(error => {
                    console.error('Error reconnecting client:', error);
                    if (btn) {
                        btn.classList.remove('reconnecting');
                        btn.textContent = 'Reconnect Failed';
                        btn.disabled = false;
                        setTimeout(() => {
                            btn.textContent = 'Reconnect';
                        }, 2000);
                    }
                });
        }
        
        // Initialize connection status
        document.addEventListener('DOMContentLoaded', function() {
            updateConnectionStatus();
        });
        
        updateLogs();
    </script>
</head>
<body>
    <div id="server-status" class="server-status">Server Status</div>
    <h1>GoRelay Logs</h1>
    <nav>
        <a href="/">Dashboard</a> |
        <a href="/logs">Logs</a>
    </nav>
    <div id="logs"></div>
</body>
</html>
`

// Default handlers

func (ms *MonitorServer) handleRoot(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "GoRelay Monitoring Server")
}

func (ms *MonitorServer) handleClients(w http.ResponseWriter, r *http.Request) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	clients := make([]string, 0, len(ms.clients))
	for alias := range ms.clients {
		clients = append(clients, alias)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"clients": clients,
		"count":   len(clients),
	})
}

// handleReconnect handles client reconnection requests
func (ms *MonitorServer) handleReconnect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	alias := r.URL.Path[len("/reconnect/"):]
	if alias == "" {
		http.Error(w, "Client alias required", http.StatusBadRequest)
		return
	}

	ms.mu.RLock()
	client, exists := ms.clients[alias]
	ms.mu.RUnlock()

	if !exists {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	// Log the reconnection attempt
	ms.Log(alias, "info", "Manual reconnection requested")

	// Signal reconnection through the client's account
	if client.Account != nil {
		client.Account.Reconnect = true
	}

	w.WriteHeader(http.StatusOK)
}

// SetCommandHandler sets a handler function for client commands
func (ms *MonitorServer) SetCommandHandler(handler func(alias string, command string) string) {
	ms.cmdHandler = handler
}

// handleCommand handles command requests for clients
func (ms *MonitorServer) handleCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract client alias from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid command URL", http.StatusBadRequest)
		return
	}
	alias := parts[2]

	// Check if client exists
	ms.mu.RLock()
	_, exists := ms.clients[alias]
	ms.mu.RUnlock()

	if !exists {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	// Parse command from request body
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	command := r.FormValue("command")
	if command == "" {
		http.Error(w, "Command is required", http.StatusBadRequest)
		return
	}

	// Log the command
	ms.Log(alias, "info", "Command sent: "+command)

	// Execute command if handler is set
	var response string
	if ms.cmdHandler != nil {
		response = ms.cmdHandler(alias, command)
	} else {
		response = "Command received, but no handler is configured"
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":   "success",
		"response": response,
	})
}
