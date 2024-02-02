package profile

import (
	"io"
	"os"
	"path/filepath"

	"github.com/The-Lethal-Foundation/lethal-core/pkg/filesystem"
	"github.com/The-Lethal-Foundation/lethal-core/pkg/modmanager"
)

const ProfilesDirName = "Profiles"

func CreateProfile(profileName string) error {
	profilePath := filepath.Join(filesystem.GetDefaultPath(), "LethalCompany", ProfilesDirName, profileName)
	if err := os.Mkdir(profilePath, 0755); err != nil {
		return err
	}

	// Unpack the BepInEx zip file into the profile directory.
	// Assuming you have a function in modmanager to handle this.
	err := modmanager.UnpackBepInEx(profilePath)
	if err != nil {
		return err
	}

	// create plugins dir at profilePath/BepInEx/plugins
	return copyDefaultConfig(profilePath)
}

// DeleteProfile deletes an existing profile.
func DeleteProfile(profileName string) error {
	profilePath := filepath.Join(filesystem.GetDefaultPath(), "LethalCompany", ProfilesDirName, profileName)
	return os.RemoveAll(profilePath)
}

// RenameProfile renames an existing profile.
func RenameProfile(oldName, newName string) error {
	oldPath := filepath.Join(filesystem.GetDefaultPath(), "LethalCompany", ProfilesDirName, oldName)
	newPath := filepath.Join(filesystem.GetDefaultPath(), "LethalCompany", ProfilesDirName, newName)
	return os.Rename(oldPath, newPath)
}

// copyDefaultConfig copies the default BepInEx configuration file into the profile.
func copyDefaultConfig(profilePath string) error {
	configDir := filepath.Join(profilePath, "BepInEx", "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	defaultConfigPath := filepath.Join("..", "..", "assets", "BepInEx.cfg")
	destConfigPath := filepath.Join(configDir, "BepInEx.cfg")

	input, err := os.Open(defaultConfigPath)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(destConfigPath)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	return err
}
