package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"launchd_docker/pkg/config"
	"launchd_docker/pkg/service"
	"launchd_docker/pkg/vm"
)

func main() {
	// Parse command line arguments
	configPath := flag.String("config", "", "Path to configuration file")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	flag.Parse()

	if *configPath == "" {
		log.Fatal("Configuration file path is required. Use -config flag")
	}

	// Initialize logger
	log.SetOutput(os.Stderr)

	// Load and validate configuration
	if *verbose {
		log.Printf("Loading configuration from %s", *configPath)
	}
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	if *verbose {
		log.Printf("Configuration loaded successfully with %d services", len(cfg.Services))
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Initialize VM manager
	if *verbose {
		log.Printf("Initializing VM manager for Lima instance: %s", cfg.Hypervisor.LimaInstance)
	}
	vmManager := vm.NewManager(cfg.Hypervisor.LimaInstance)

	// Check and start VM if needed
	if *verbose {
		log.Println("Checking VM status...")
	}
	if err := vmManager.EnsureRunning(); err != nil {
		log.Fatalf("Failed to ensure VM is running: %v", err)
	}
	if *verbose {
		log.Println("VM is running")
	}

	// Wait for VM to be healthy
	if *verbose {
		log.Println("Waiting for VM to be healthy...")
	}
	if err := vmManager.WaitForHealthy(); err != nil {
		log.Fatalf("Failed waiting for VM to be healthy: %v", err)
	}
	if *verbose {
		log.Println("VM is healthy")
	}

	// Initialize service manager
	if *verbose {
		log.Println("Initializing service manager...")
	}
	serviceManager := service.NewManager(cfg.Services)
	serviceManager.SetVerbose(*verbose)

	// Start services
	if *verbose {
		log.Printf("Starting %d services...", len(cfg.Services))
	}
	if err := serviceManager.StartAll(); err != nil {
		log.Printf("Error starting services: %v", err)
		// Don't exit here, allow graceful shutdown
	}
	if *verbose {
		log.Println("All services started")
	}

	// Wait for shutdown signal
	if *verbose {
		log.Println("Waiting for shutdown signal...")
	}
	<-sigChan
	log.Println("Received shutdown signal, initiating graceful shutdown...")

	// Stop all services
	if *verbose {
		log.Println("Stopping all services...")
	}
	if err := serviceManager.StopAll(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	if *verbose {
		log.Println("Shutdown complete")
	}
} 