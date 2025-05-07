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
	fileInfo, err := os.Stat(filePath) // Stat returns file name (info)
	if err != nil {
		return err
	}
	if fileInfo.IsDir() {
		return nil
	}
	// Take file extensions
	ext := strings.ToLower(filepath.Ext(filePath)) // Ext returns the extension used by path

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

	// Move file to category folder
	err = utils.MoveFile(filePath, destDir)
	if err != nil {
		return err
	}

	// If Dropbox sync is enabled, upload to Dropbox too
	if c.config.EnableDropbox && c.config.SyncToDropbox {
		fileName := filepath.Base(filePath)
		newLocalPath := filepath.Join(destDir, fileName)

		// This will be handled in main.go with a file queue
		// that main.go will process
		c.AddToSyncQueue(newLocalPath, category)
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
}

func (c *Categorizer) GetSyncQueue() []SyncQueueItem {
	syncQueueMutex.Lock()
	defer syncQueueMutex.Unlock()

	items := make([]SyncQueueItem, len(syncQueue))
	copy(items, syncQueue)
	syncQueue = []SyncQueueItem{}

	return items
}
