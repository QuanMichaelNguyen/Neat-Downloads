package cloud

import (
	"log"
	"neat-download/internal/categorizer"
	"path/filepath"
	"strings"
	"time"
)

type DropboxWatcher struct {
	client      *DropboxClient
	categorizer *categorizer.Categorizer
	watchFolder string
	interval    time.Duration
	stopChan    chan struct{}
	knownFiles  map[string]bool
}

func NewDropboxWatcher(client *DropboxClient, categorizer *categorizer.Categorizer, watchFolder string, intervalMinutes int) *DropboxWatcher {
	return &DropboxWatcher{
		client:      client,
		categorizer: categorizer,
		watchFolder: watchFolder,
		interval:    time.Duration(intervalMinutes) * time.Minute,
		stopChan:    make(chan struct{}),
		knownFiles:  make(map[string]bool),
	}
}

func (w *DropboxWatcher) Start() {
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		// Do an initial check
		w.checkDropboxFolder()

		for {
			select {
			case <-ticker.C:
				w.checkDropboxFolder()
			case <-w.stopChan:
				return
			}
		}
	}()

	log.Printf("Started watching Dropbox folder: %s", w.watchFolder)
}

func (w *DropboxWatcher) Stop() {
	close(w.stopChan)
	log.Println("Dropbox watcher stopped")
}

func (w *DropboxWatcher) checkDropboxFolder() {
	files, err := w.client.ListFiles(w.watchFolder)
	if err != nil {
		log.Printf("Error listing Dropbox files: %v", err)
		return
	}

	for _, filePath := range files {
		// Skip if we've seen this file before
		if w.knownFiles[filePath] {
			continue
		}

		// Mark it as known
		w.knownFiles[filePath] = true

		// Get the file extension
		ext := strings.ToLower(filepath.Ext(filePath))
		fileName := filepath.Base(filePath)

		// Get category
		category := w.categorizer.GetCategoryForExtension(ext)
		if category == "" {
			log.Printf("No category for file: %s", filePath)
			continue
		}

		// Move file to category folder in Dropbox
		destPath := filepath.Join(w.watchFolder, category, fileName)
		err := w.client.MoveFile(filePath, destPath)
		if err != nil {
			log.Printf("Error moving Dropbox file: %v", err)
			continue
		}

		log.Printf("Categorized Dropbox file: %s -> %s", filePath, destPath)
	}
}
