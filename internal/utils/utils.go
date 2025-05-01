package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// MoveFile moves a file to the destination directory
func MoveFile(sourcePath, destDir string) error {
	// Create directory if do not exist
	err := os.MkdirAll(destDir, os.ModeAppend)
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Get filename from source path
	fileName := filepath.Base(sourcePath)
	destPath := filepath.Join(destDir, fileName)

	// Check if file already exists in destination
	if _, err := os.Stat(destPath); err == nil {
		// File exists, rename it with a numerical suffix
		destPath = generateUniqueFilename(destDir, fileName)

	}
	// Move the file (rename in filesystem)
	err = os.Rename(sourcePath, destPath)
	if err != nil {
		return fmt.Errorf("failed to move file: %w", err)
	}

	fmt.Printf("Moved file from %s to %s\n", sourcePath, destPath)
	return nil
}

// Creates a unique filename by adding a numerical suffix
func generateUniqueFilename(dir, fileName string) string {
	fileExt := filepath.Ext(fileName)
	fileBase := fileName[:len(fileName)-len(fileExt)]

	counter := 1
	newPath := filepath.Join(dir, fileName)

	for {
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			// found a name that do not exists
			return newPath
		}
		// Try the next number
		newFileName := fmt.Sprintf("%s (%d)%s", fileBase, counter, fileExt)
		newPath = filepath.Join(dir, newFileName)
		counter++
	}
}
