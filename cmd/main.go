package main

import (
	"fmt"
	"log"
	"neat-download/configs"
	"neat-download/internal/categorizer"
	"neat-download/internal/watcher"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Load config
	config, err := configs.LoadConfig("./configs/config.yaml")
	if err != nil {
		log.Fatal("Failed to load configuration: %v", err)
	}

	// Initialize category
	categorizer := categorizer.NewCategorizer(config)

	// Create and start watcher
	fileWatcher, err := watcher.NewWatcher(config.WatchDir, categorizer)
	if err != nil {
		log.Fatalf("Failed to create watcher: %v", err)
	}

	// Satrt watching
	err = fileWatcher.Start()
	if err != nil {
		log.Fatalf("Failed to start watcher: %v", err)
	}

	fmt.Println("AI File Categorizer started")
	fmt.Printf("Watching directory: %s\n", config.WatchDir)
	fmt.Println("Press Ctrl+C to exit")

	// Wait for interupted signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Shutdown
	fileWatcher.Stop()
	fmt.Println("File categorizer stopped")

}
