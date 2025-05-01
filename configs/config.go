package configs

type Config struct {
    WatchDir     string            `yaml:"watch_dir"`
    Categories   map[string]string `yaml:"categories"`
    FilePatterns map[string]string `yaml:"file_patterns"`
}

func LoadConfig(configPath string) (*Config, error) {
    // TODO: Load YAML config
    // For now, return a default config
    return &Config{
        WatchDir: "C:/Users/nguye/Downloads",
        Categories: map[string]string{
            "documents": "C:/Users/nguye/Downloads/Documents",
            "images":    "C:/Users/nguye/Downloads/Media/Images",
            "videos":    "C:/Users/nguye/Downloads/Media/Videos",
            "audio":     "C:/Users/nguye/Downloads/Media/Audio",
            "software":  "C:/Users/nguye/Downloads/Software",
            "archives":  "C:/Users/nguye/Downloads/Archives",
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
            
            ".mp4":  "videos",
            ".mkv":  "videos",
            ".avi":  "videos",
            ".mov":  "videos",
            
            ".mp3":  "audio",
            ".wav":  "audio",
            ".flac": "audio",
            
            ".exe":  "software",
            ".msi":  "software",
            ".dmg":  "software",
            ".appimage": "software",
            
            ".zip":  "archives",
            ".rar":  "archives",
            ".7z":   "archives",
            ".tar":  "archives",
            ".gz":   "archives",
        },
    }, nil
}