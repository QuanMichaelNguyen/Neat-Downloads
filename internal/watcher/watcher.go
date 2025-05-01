package watcher

import (
	"log"
	"neat-download/internal/categorizer"

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
				// Only process file creations and moves
				if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write {
					log.Printf("New file detected: %s", event.Name)
					// Don't categorize the file immediately, give it some time to finish downloading
					go func(filePath string) {
						err := w.categorizer.CategorizeFile(filePath)
						if err != nil {
							log.Printf("Error categorizing file %s: %v", filePath, err)
						}
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
