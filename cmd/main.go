package main

import (
	"fmt"
	"log"
	"neat-download/configs"
	"neat-download/internal/categorizer"
	"neat-download/internal/cloud"
	"neat-download/internal/watcher"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
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

	// Dropbox -------------------------------------------
	// Initialize Dropbox client if enabled
	var dropboxClient *cloud.DropboxClient
	var dropboxWatcher *cloud.DropboxWatcher

	if config.EnableDropbox {
		dropboxClient = cloud.NewDropboxClient(config.DropboxAppKey, config.DropboxAppSecret)

		// For initial setup, generate auth URL and handle token exchange
		// In a real app, you'd store the token and reuse it
		if dropboxClient.AccessToken == "" {
			fmt.Println("Please authorize the app. Visit this URL:")
			fmt.Println(dropboxClient.GetAuthURL())
			fmt.Print("Enter the authorization code: ")

			var code string
			fmt.Scanln(&code)

			err := dropboxClient.ExchangeCodeForToken(code)
			if err != nil {
				log.Fatalf("Failed to exchange token: %v", err)
			}

			// TODO: Save token for future use
		}

		// Start Dropbox watcher
		if config.SyncFromDropbox {
			dropboxWatcher = cloud.NewDropboxWatcher(
				dropboxClient,
				categorizer,
				config.DropboxFolder,
				config.SyncInterval,
			)
			dropboxWatcher.Start()
			fmt.Printf("Watching Dropbox folder: %s\n", config.DropboxFolder)
		}

		// Start sync worker for local -> Dropbox sync
		if config.SyncToDropbox {
			go syncWorker(dropboxClient, categorizer, config)
		}
	}

	// --------------------------------------------------
	fmt.Println("Press Ctrl+C to exit")

	// Wait for interupted signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Shutdown
	fileWatcher.Stop()
	if dropboxWatcher != nil {
		dropboxWatcher.Stop()
	}
	fmt.Println("File categorizer stopped")
}

// Worker to sync local files to Dropbox
func syncWorker(client *cloud.DropboxClient, cat *categorizer.Categorizer, config *configs.Config) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C

		// Get items to sync
		items := cat.GetSyncQueue()
		for _, item := range items {
			fileName := filepath.Base(item.LocalPath)
			dropboxPath := filepath.Join(config.DropboxFolder, item.Category, fileName)

			err := client.UploadFile(item.LocalPath, dropboxPath)
			if err != nil {
				log.Printf("Failed to upload %s to Dropbox: %v", item.LocalPath, err)
				continue
			}

			log.Printf("Uploaded %s to Dropbox: %s", item.LocalPath, dropboxPath)
		}
	}
}
