package main

import (
	"context"
	"fmt"
	"log"
	"neat-download/configs"
	"neat-download/internal/categorizer"
	"neat-download/internal/cloud"
	"neat-download/internal/watcher"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx           context.Context
	config        *configs.Config
	categorizer   *categorizer.Categorizer
	fileWatcher   *watcher.Watcher
	dropboxClient *cloud.DropboxClient
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Load configuration
	config, err := configs.LoadConfig("../configs/config.yaml")
	if err != nil {
		log.Printf("Error loading config: %v", err)
		return
	}

	a.config = config
	a.categorizer = categorizer.NewCategorizer(config)

	// Create file watcher
	watcher, err := watcher.NewWatcher(config.WatchDir, a.categorizer)
	if err != nil {
		log.Printf("Error creating watcher: %v", err)
		return
	}
	a.fileWatcher = watcher

	// Set up event emission for file categorization
	// We'll use a goroutine to monitor the sync queue for new files
	go func() {
		for {
			items := a.categorizer.GetSyncQueue()
			for _, item := range items {
				fileName := filepath.Base(item.LocalPath)
				log.Printf("Emitting event for file: %s, category: %s", fileName, item.Category)
				runtime.EventsEmit(a.ctx, "fileCategorized", map[string]string{
					"fileName": fileName,
					"category": item.Category,
				})
				log.Printf("Successfully emitted event for file: %s", fileName)
			}
			time.Sleep(100 * time.Millisecond) // Check every 100ms
		}
	}()
}

// GetSettings returns the current config settings
func (a *App) GetSettings() map[string]interface{} {
	dropboxEnabled := a.config.EnableDropbox
	// Check for a valid Dropbox token
	if a.dropboxClient == nil {
		a.dropboxClient = cloud.NewDropboxClient(
			a.config.DropboxAppKey,
			a.config.DropboxAppSecret,
		)
	}
	if a.dropboxClient.AccessToken != "" {
		dropboxEnabled = true
	}
	return map[string]interface{}{
		"watchDir":       a.config.WatchDir,
		"categories":     a.config.Categories,
		"filePatterns":   a.config.FilePatterns,
		"dropboxEnabled": dropboxEnabled,
		"dropboxFolder":  a.config.DropboxFolder,
	}
}

// StartWatching starts the file watcher
func (a *App) StartWatching() string {
	err := a.fileWatcher.Start()
	if err != nil {
		log.Printf("Error starting watcher: %v", err)
		return "Error starting watcher: " + err.Error()
	}
	return "Monitoring started successfully"
}

// GetDropboxAuthURL returns the Dropbox OAuth URL
func (a *App) GetDropboxAuthURL() string {
	if a.dropboxClient == nil {
		a.dropboxClient = cloud.NewDropboxClient(
			a.config.DropboxAppKey,
			a.config.DropboxAppSecret,
		)
	}
	return a.dropboxClient.GetAuthURL()
}

// ExchangeDropboxCode exchanges the authorization code for tokens
func (a *App) ExchangeDropboxCode(code string) (map[string]string, error) {
	if a.dropboxClient == nil {
		return nil, fmt.Errorf("dropbox client not initialized")
	}

	err := a.dropboxClient.ExchangeCodeForToken(code)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"accessToken":  a.dropboxClient.AccessToken,
		"refreshToken": a.dropboxClient.RefreshToken,
	}, nil
}
