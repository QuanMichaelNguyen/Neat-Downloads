package configs

import (
	"os"

	"gopkg.in/yaml.v3"
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
	SyncInterval     int    `yaml:"sync_interval"`     // In minutes
}

func LoadConfig(configPath string) (*Config, error) {
	// Read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Parse the YAML
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// Expand environment variables in paths
	config.WatchDir = os.ExpandEnv(config.WatchDir)
	for category, path := range config.Categories {
		config.Categories[category] = os.ExpandEnv(path)
	}

	// Create category directories if they don't exist
	for _, dir := range config.Categories {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, err
		}
	}

	return &config, nil
}
