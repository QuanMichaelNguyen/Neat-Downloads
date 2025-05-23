package categorizer

import (
	"log"
	"neat-download/configs"
	"neat-download/internal/utils"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Categorizer struct {
	config *configs.Config
}

func NewCategorizer(config *configs.Config) *Categorizer {
	log.Printf("Creating new categorizer with config: %+v", config)
	return &Categorizer{
		config: config,
	}
}

func (c *Categorizer) GetCategoryForExtension(ext string) string {
	category, found := c.config.FilePatterns[ext]
	if !found {
		return ""
	}
	return category
}

func (c *Categorizer) CategorizeFile(filePath string) error {
	// Skip directories
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Printf("Error getting file info for %s: %v", filePath, err)
		return err
	}
	if fileInfo.IsDir() {
		log.Printf("Skipping directory: %s", filePath)
		return nil
	}

	// Take file extensions
	ext := strings.ToLower(filepath.Ext(filePath))
	log.Printf("Processing file: %s with extension: %s", filePath, ext)

	// Categorize based on extension
	category, found := c.config.FilePatterns[ext]
	if !found {
		log.Printf("No category found for extension %s, skipping file: %s", ext, filePath)
		return nil
	}

	// Get destination directory
	destDir, found := c.config.Categories[category]
	if !found {
		log.Printf("Category %s not configured, skipping file: %s", category, filePath)
		return nil
	}

	log.Printf("Moving file %s to category %s (directory: %s)", filePath, category, destDir)

	// Move file to category folder
	err = utils.MoveFile(filePath, destDir)
	if err != nil {
		log.Printf("Error moving file %s to %s: %v", filePath, destDir, err)
		return err
	}

	// Always add to sync queue for UI updates, regardless of Dropbox settings
	fileName := filepath.Base(filePath)
	newLocalPath := filepath.Join(destDir, fileName)
	log.Printf("Adding file to sync queue: %s (category: %s)", newLocalPath, category)
	c.AddToSyncQueue(newLocalPath, category)

	// If Dropbox sync is enabled, handle Dropbox sync separately
	if c.config.EnableDropbox && c.config.SyncToDropbox {
		log.Printf("Dropbox sync enabled, will sync file: %s", newLocalPath)
	}

	return nil
}

// For syncing files to Dropbox
type SyncQueueItem struct {
	LocalPath string
	Category  string
}

var syncQueue []SyncQueueItem
var syncQueueMutex sync.Mutex

func (c *Categorizer) AddToSyncQueue(localPath, category string) {
	syncQueueMutex.Lock()
	defer syncQueueMutex.Unlock()

	syncQueue = append(syncQueue, SyncQueueItem{
		LocalPath: localPath,
		Category:  category,
	})
	log.Printf("Added to sync queue: %s (category: %s)", localPath, category)
}

func (c *Categorizer) GetSyncQueue() []SyncQueueItem {
	syncQueueMutex.Lock()
	defer syncQueueMutex.Unlock()

	items := make([]SyncQueueItem, len(syncQueue))
	copy(items, syncQueue)
	if len(items) > 0 {
		log.Printf("Retrieved %d items from sync queue: %+v", len(items), items)
	}
	syncQueue = []SyncQueueItem{}

	return items
}
