package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// MonitorServer provides an HTTP server for debugging and monitoring
type MonitorServer struct {
	server   *http.Server
	port     int
	clients  map[string]interface{} // Using interface{} to avoid import cycle
	mu       sync.RWMutex
	handlers map[string]http.HandlerFunc
}

// NewMonitorServer creates a new monitoring server instance
func NewMonitorServer(port int) *MonitorServer {
	ms := &MonitorServer{
		port:     port,
		clients:  make(map[string]interface{}),
		handlers: make(map[string]http.HandlerFunc),
	}

	// Set up default handlers
	ms.handlers["/"] = ms.handleRoot
	ms.handlers["/clients"] = ms.handleClients
	ms.handlers["/status"] = ms.handleStatus

	return ms
}

// Start starts the monitoring server
func (ms *MonitorServer) Start() error {
	mux := http.NewServeMux()

	// Register handlers
	for path, handler := range ms.handlers {
		mux.HandleFunc(path, handler)
	}

	ms.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", ms.port),
		Handler: mux,
	}

	return ms.server.ListenAndServe()
}

// Stop stops the monitoring server
func (ms *MonitorServer) Stop() error {
	if ms.server != nil {
		return ms.server.Close()
	}
	return nil
}

// AddClient adds a client to the server's tracking
func (ms *MonitorServer) AddClient(alias string, client interface{}) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.clients[alias] = client
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

// Default handlers

func (ms *MonitorServer) handleRoot(w http.ResponseWriter, r *http.Request) {
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

func (ms *MonitorServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	status := map[string]interface{}{
		"uptime":        0, // TODO: Track uptime
		"clientCount":   len(ms.clients),
		"memoryUsage":   0, // TODO: Track memory usage
		"cpuUsage":      0, // TODO: Track CPU usage
		"lastHeartbeat": 0, // TODO: Track last heartbeat
	}

	json.NewEncoder(w).Encode(status)
}
