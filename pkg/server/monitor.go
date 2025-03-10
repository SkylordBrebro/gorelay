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
}

// NewMonitorServer creates a new monitoring server instance
func NewMonitorServer(port int) *MonitorServer {
	ms := &MonitorServer{
		port:      port,
		clients:   make(map[string]*ClientInfo),
		handlers:  make(map[string]http.HandlerFunc),
		startTime: time.Now(),
	}

	// Set up default handlers
	ms.handlers["/"] = ms.handleDashboard
	ms.handlers["/status"] = ms.handleStatus
	ms.handlers["/api/status"] = ms.handleAPIStatus
	ms.handlers["/static/"] = http.StripPrefix("/static/", http.FileServer(http.Dir("pkg/server/static"))).ServeHTTP

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
}

// RemoveClient removes a client from tracking
func (ms *MonitorServer) RemoveClient(alias string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	delete(ms.clients, alias)
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
		if mapName, ok := status["map"].(string); ok {
			client.CurrentMap = mapName
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

// handleAPIStatus returns status in JSON format
func (ms *MonitorServer) handleAPIStatus(w http.ResponseWriter, r *http.Request) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	// Create a copy of clients with updated status
	clientsCopy := make(map[string]*ClientInfo)
	for alias, client := range ms.clients {
		clientsCopy[alias] = &ClientInfo{
			Account:    client.Account,
			Connected:  client.Connected,
			LastSeen:   client.LastSeen,
			CurrentMap: client.CurrentMap,
			Stats:      client.Stats,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clientsCopy)
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
    <link rel="stylesheet" href="/static/style.css">
    <style>
        .client { 
            border: 1px solid var(--border-color);
            padding: 15px;
            margin: 10px 0;
            background: var(--card-bg);
            border-radius: 5px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.2);
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
            Promise.all([
                fetch('/api/status').then(response => response.json()),
                fetch('/status').then(response => response.json())
            ])
            .then(([clientData, systemData]) => {
                // Update system stats
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
                    const connectedTime = new Date(client.Connected).toLocaleString();
                    const lastSeenTime = lastSeen.toLocaleString();
                    
                    let statsHtml = '';
                    if (client.Stats) {
                        for (const [key, value] of Object.entries(client.Stats)) {
                            statsHtml += '<div class="stat-item">' + key + ': ' + value + '</div>';
                        }
                    }
                    
                    clientDiv.innerHTML = '<h3>' + alias + '</h3>' +
                        '<p>' +
                        'Connected: ' + connectedTime + '<br>' +
                        'Last Seen: ' + lastSeenTime +
                        '</p>' +
                        '<div class="stats">' + statsHtml + '</div>';
                        
                    clientsDiv.appendChild(clientDiv);
                }
            })
            .catch(error => {
                console.error('Error updating status:', error);
            });
            
            setTimeout(updateStatus, 1000);
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

        updateStatus();
    </script>
</head>
<body>
    <h1>GoRelay Monitor</h1>
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
