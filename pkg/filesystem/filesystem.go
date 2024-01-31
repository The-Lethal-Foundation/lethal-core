package filesystem

import (
	"os"
	"path/filepath"

	"github.com/The-Lethal-Foundation/lethal-core/pkg/config"
)

const defaultBasePath = "Lethal Foundation/Lethal Mod Manager"
const DefaultCacheDir = "Caches"

// getDefaultPath gets the full default path under %AppData%.
var GetDefaultPath = func() string {
	return filepath.Join(os.Getenv("APPDATA"), defaultBasePath)
}

// InitializeStructure sets up the required directory structure in the default path.
func InitializeStructure() error {
	basePath := GetDefaultPath()

	requiredDirs := []string{
		"Profiles",
		"Caches",
	}

	for _, dir := range requiredDirs {
		if err := createDirIfNotExist(filepath.Join(basePath, dir)); err != nil {
			return err
		}
	}

	// Initialize the configuration file.
	configPath := filepath.Join(basePath, config.ConfigFileName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config.InitializeConfig(basePath)
	}

	return nil
}

// createDirIfNotExist creates a directory if it does not exist.
func createDirIfNotExist(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0755)
	}
	return nil
}

// IsFileSystemSetUp checks if the required directory structure is in place in the default path.
func IsFileSystemSetUp() (bool, error) {
	basePath := GetDefaultPath()
	requiredDirs := []string{
		"Profiles",
		"Caches",
	}

	for _, dir := range requiredDirs {
		if _, err := os.Stat(filepath.Join(basePath, dir)); os.IsNotExist(err) {
			return false, nil
		} else if err != nil {
			return false, err
		}
	}

	// Check if configuration file exists.
	configPath := filepath.Join(basePath, config.ConfigFileName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
