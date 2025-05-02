package configs

import (
	"os"
	"path/filepath"
)

type Config struct {
	WatchDir     string            `yaml:"watch_dir"`
	Categories   map[string]string `yaml:"categories"`
	FilePatterns map[string]string `yaml:"file_patterns"`
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
