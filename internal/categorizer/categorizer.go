package categorizer

import (
	"log"
	"neat-download/configs"
	"neat-download/internal/utils"
	"os"
	"path/filepath"
	"strings"
)

type Categorizer struct {
	config *configs.Config
}

func NewCategorizer(config *configs.Config) *Categorizer {
	return &Categorizer{
		config: config,
	}
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
	return utils.MoveFile(filePath, destDir)

}
