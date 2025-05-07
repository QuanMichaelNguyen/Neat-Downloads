package main

import (
	"context"
	"log"
	"neat-download/configs"
	"neat-download/internal/categorizer"
	"neat-download/internal/watcher"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"path/filepath"
	"time"
)

// App struct
type App struct {
	ctx         context.Context
	config      *configs.Config
	categorizer *categorizer.Categorizer
	fileWatcher *watcher.Watcher
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
	return map[string]interface{}{
		"watchDir":       a.config.WatchDir,
		"categories":     a.config.Categories,
		"filePatterns":   a.config.FilePatterns,
		"dropboxEnabled": a.config.EnableDropbox,
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
