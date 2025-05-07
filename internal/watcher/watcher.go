package watcher

import (
	"log"
	"neat-download/internal/categorizer"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	watchDir    string
	categorizer *categorizer.Categorizer
	watcher     *fsnotify.Watcher
}

func NewWatcher(watchDir string, categorizer *categorizer.Categorizer) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &Watcher{
		watchDir:    watchDir,
		categorizer: categorizer,
		watcher:     fsWatcher,
	}, nil
}

func (w *Watcher) Start() error {
	// Add the directory to watch
	err := w.watcher.Add(w.watchDir)
	if err != nil {
		return err
	}
	log.Printf("Started watching directory: %s", w.watchDir)

	// Process events
	go func() {
		for {
			select {
			case event, ok := <-w.watcher.Events:
				if !ok {
					return
				}
				log.Printf("File system event detected: %s (operation: %v)", event.Name, event.Op)

				// Skip temporary download files
				if strings.HasSuffix(event.Name, ".tmp") || strings.HasSuffix(event.Name, ".crdownload") {
					log.Printf("Skipping temporary download file: %s", event.Name)
					continue
				}

				// Only process file creations and writes
				if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write {
					log.Printf("New file detected: %s", event.Name)

					// Retry up to 5 times to ensure the file is available and is not a directory
					go func(filePath string) {
						for i := 0; i < 5; i++ {
							time.Sleep(500 * time.Millisecond)
							info, err := os.Stat(filePath)
							if err == nil && !info.IsDir() {
								log.Printf("Attempting to categorize file: %s", filePath)
								err := w.categorizer.CategorizeFile(filePath)
								if err != nil {
									log.Printf("Error categorizing file %s: %v", filePath, err)
								} else {
									log.Printf("Successfully categorized file: %s", filePath)
								}
								return
							}
							if err == nil && info.IsDir() {
								log.Printf("Skipping directory: %s", filePath)
								return
							}
							// If not found, retry
						}
						log.Printf("File not available after retries: %s", filePath)
					}(event.Name)
				}
			case err, ok := <-w.watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Watcher error: %v", err)
			}
		}
	}()
	return nil
}

func (w *Watcher) Stop() {
	w.watcher.Close()
}
