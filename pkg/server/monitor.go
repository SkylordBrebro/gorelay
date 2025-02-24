package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"runtime"
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
	ms.handlers["/static/"] = http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP
	ms.handlers["/reconnect/"] = ms.handleReconnect

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
    </style>
    <script>
        function updateStatus() {
            Promise.all([
                fetch('/api/status').then(response => response.json()),
                fetch('/status').then(response => response.json())
            ])
            .then(([clientData, systemData]) => {
                // Update server status indicator
                const statusDiv = document.getElementById('server-status');
                if (systemData.isOnline) {
                    statusDiv.className = 'server-status online';
                    statusDiv.textContent = 'Server Online';
                } else {
                    statusDiv.className = 'server-status offline';
                    statusDiv.textContent = 'Server Offline';
                }

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
            .catch(error => console.error('Error updating status:', error));
            
            setTimeout(updateStatus, 1000);
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
    </style>
    <script>
        function updateLogs() {
            fetch('/api/logs')
                .then(response => response.json())
                .then(data => {
                    const logsDiv = document.getElementById('logs');
                    logsDiv.innerHTML = '';
                    
                    for (const [alias, logs] of Object.entries(data)) {
                        const accountDiv = document.createElement('div');
                        accountDiv.className = 'account-logs';
                        
                        let terminalHtml = '<h3>' + alias + '</h3>' +
                            '<button class="reconnect-btn" onclick="reconnectClient(\'' + alias + '\')">Reconnect</button>' +
                            '<div class="terminal">';
                        logs.forEach(log => {
                            const timestamp = new Date(log.Timestamp).toLocaleString();
                            const color = {
                                'info': '#2196f3',
                                'warning': '#ff9800',
                                'error': '#f44336'
                            }[log.Level] || '#e0e0e0';
                            
                            terminalHtml += '<div class="terminal-line" style="color: ' + color + '">[' + timestamp + '] ' + log.Message + '</div>';
                        });
                        terminalHtml += '</div>';
                        
                        accountDiv.innerHTML = terminalHtml;
                        logsDiv.appendChild(accountDiv);

                        // Auto-scroll terminals to bottom
                        const terminals = accountDiv.getElementsByClassName('terminal');
                        for (const terminal of terminals) {
                            terminal.scrollTop = terminal.scrollHeight;
                        }
                    }
                })
                .catch(error => console.error('Error updating logs:', error));
            setTimeout(updateLogs, 1000);
        }
        updateLogs();
    </script>
</head>
<body>
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
