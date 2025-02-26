package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"gorelay/pkg/account"
	"gorelay/pkg/client"
	"gorelay/pkg/config"
	"gorelay/pkg/logger"
	"gorelay/pkg/models"
	"gorelay/pkg/plugin"
	"gorelay/pkg/server"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config.json", "Path to config file")
	accountsPath := flag.String("accounts", "accounts.json", "Path to accounts file")
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger, err := logger.New("gorelay.log", *debug || cfg.Debug)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	// Initialize monitor server
	monitor := server.NewMonitorServer(8080)
	if err := monitor.Start(); err != nil {
		logger.Error("Main", "Failed to start monitor server: %v", err)
		os.Exit(1)
	}
	defer monitor.Stop()

	// Load accounts
	accManager, err := account.LoadAccounts(*accountsPath)
	if err != nil {
		logger.Error("Main", "Failed to load accounts: %v", err)
		os.Exit(1)
	}

	// Initialize account aliases for the monitor server
	aliases := make([]string, len(accManager.Accounts))
	for i, acc := range accManager.Accounts {
		aliases[i] = acc.Alias
	}
	models.SetAccountAliases(aliases)

	// Create wait group for managing client goroutines
	var wg sync.WaitGroup
	clients := make([]*client.Client, len(accManager.Accounts))
	clientMutex := sync.RWMutex{}

	// Create and connect clients concurrently
	for i, acc := range accManager.Accounts {
		wg.Add(1)
		go func(index int, acc *account.Account) {
			defer wg.Done()

			// Create client
			client := client.NewClient(acc, cfg, logger)

			// Create plugin manager
			pluginManager := plugin.NewManager(client)

			// Load plugins if enabled
			if cfg.Plugins.Enabled {
				for _, pluginPath := range cfg.Plugins.List {
					if err := pluginManager.LoadPlugin(pluginPath); err != nil {
						logger.Error("Main", "Failed to load plugin %s for account %s: %v", pluginPath, acc.Alias, err)
						continue
					}
				}
			}

			// Add client to slice with proper synchronization
			clientMutex.Lock()
			clients[index] = client
			clientMutex.Unlock()

			// Connect client with retries
			maxRetries := 3
			for retry := 0; retry < maxRetries; retry++ {
				if err := client.Connect(); err != nil {
					logger.Error("Main", "Failed to connect client %s (attempt %d/%d): %v",
						acc.Alias, retry+1, maxRetries, err)
					if retry < maxRetries-1 {
						time.Sleep(time.Second * 2) // Wait before retry
						continue
					}
				} else {
					logger.Info("Main", "Successfully connected client %s", acc.Alias)
					break
				}
			}

			// Start client update loop in a separate goroutine
			go func() {
				ticker := time.NewTicker(time.Millisecond * 50) // 20 updates per second
				defer ticker.Stop()

				for {
					select {
					case <-ticker.C:
						if client.IsConnected() {
							client.Update()
						}
					}
				}
			}()

		}(i, acc)
	}

	// Wait for all clients to be created and connected
	wg.Wait()
	logger.Info("Main", "All clients initialized")

	// Handle shutdown gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Info("Main", "Shutting down...")

	// Disconnect all clients concurrently
	var shutdownWg sync.WaitGroup
	for _, c := range clients {
		if c == nil {
			continue
		}
		shutdownWg.Add(1)
		go func(c *client.Client) {
			defer shutdownWg.Done()
			if c.IsConnected() {
				c.Disconnect()
			}
		}(c)
	}
	shutdownWg.Wait()
}
