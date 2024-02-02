package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Please don't the config values from inside the package.
// It's supposed to be used by the main package.
type Config struct {
	LastUsedProfile      string `json:"last_used_profile"`
	CachedBepInExVersion string `json:"cached_bepinex_version"`
	OtherProfilesCloned  bool   `json:"other_profiles_cloned"`
}

const ConfigFileName = "config.json"

// LoadConfig reads the configuration from the file.
func LoadConfig(configPath string) (*Config, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig writes the configuration to the file.
func SaveConfig(configPath string, config *Config) error {
	configFile, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, configFile, 0644)
}

// InitializeConfig creates a new configuration file with default settings.
func InitializeConfig(basePath string) error {
	configPath := filepath.Join(basePath, ConfigFileName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := Config{
			LastUsedProfile: "Default",
		}
		return SaveConfig(configPath, &defaultConfig)
	}
	return nil
}
