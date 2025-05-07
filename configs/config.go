package configs

import (
	"os"
	"path/filepath"
)

type Config struct {
	WatchDir     string            `yaml:"watch_dir"`
	Categories   map[string]string `yaml:"categories"`
	FilePatterns map[string]string `yaml:"file_patterns"`

	// Dropbox settings
	EnableDropbox    bool   `yaml:"enable_dropbox"`
	DropboxAppKey    string `yaml:"dropbox_app_key"`
	DropboxAppSecret string `yaml:"dropbox_app_secret"`
	DropboxFolder    string `yaml:"dropbox_folder"`    // Folder to monitor within Dropbox
	SyncToDropbox    bool   `yaml:"sync_to_dropbox"`   // Whether to upload local files to Dropbox
	SyncFromDropbox  bool   `yaml:"sync_from_dropbox"` // Whether to download Dropbox files locally
	SyncInterval     int    `yaml:"sync_interval"`     // In minutes		   ``
}

func LoadConfig(configPath string) (*Config, error) {
	// TODO: Load YAML config
	// For now, return a default config
	userProfile := os.Getenv("USERPROFILE")
	downloadsDir := filepath.Join(userProfile, "Downloads")

	return &Config{
		WatchDir: downloadsDir,
		Categories: map[string]string{
			"documents": filepath.Join(downloadsDir, "Documents"),
			"images":    filepath.Join(downloadsDir, "Media", "Images"),
			"videos":    filepath.Join(downloadsDir, "Media", "Videos"),
			"audio":     filepath.Join(downloadsDir, "Media", "Audio"),
			"software":  filepath.Join(downloadsDir, "Software"),
			"archives":  filepath.Join(downloadsDir, "Archives"),
		},
		FilePatterns: map[string]string{
			".pdf":  "documents",
			".doc":  "documents",
			".docx": "documents",
			".txt":  "documents",
			".xlsx": "documents",
			".csv":  "documents",
			".ppt":  "documents",
			".pptx": "documents",

			".jpg":  "images",
			".jpeg": "images",
			".png":  "images",
			".gif":  "images",
			".webp": "images",
			".svg":  "images",

			".mp4": "videos",
			".mkv": "videos",
			".avi": "videos",
			".mov": "videos",

			".mp3":  "audio",
			".wav":  "audio",
			".flac": "audio",

			".exe":      "software",
			".msi":      "software",
			".dmg":      "software",
			".appimage": "software",

			".zip": "archives",
			".rar": "archives",
			".7z":  "archives",
			".tar": "archives",
			".gz":  "archives",
		},
	}, nil
}
