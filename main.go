package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gorelay/pkg/account"
	"gorelay/pkg/client"
	"gorelay/pkg/config"
	"gorelay/pkg/logger"
	"gorelay/pkg/plugin"
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

	// Load accounts
	accManager, err := account.LoadAccounts(*accountsPath)
	if err != nil {
		logger.Error("Main", "Failed to load accounts: %v", err)
		os.Exit(1)
	}

	// Create clients for each account
	clients := make([]*client.Client, 0)
	for _, acc := range accManager.Accounts {
		client := client.NewClient(acc, cfg, logger)
		clients = append(clients, client)

		// Create plugin manager for each client
		pluginManager := plugin.NewManager(client)

		// Load plugins if enabled
		if cfg.Plugins.Enabled {
			for _, pluginPath := range cfg.Plugins.List {
				if err := pluginManager.LoadPlugin(pluginPath); err != nil {
					logger.Error("Main", "Failed to load plugin %s: %v", pluginPath, err)
					continue
				}
			}
		}

		// Connect client
		if err := client.Connect(); err != nil {
			logger.Error("Main", "Failed to connect client %s: %v", acc.Alias, err)
			continue
		}
	}

	// Handle shutdown gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	logger.Info("Main", "Shutting down...")

	// Disconnect all clients
	for _, client := range clients {
		if client.IsConnected() {
			client.Disconnect()
		}
	}
}
