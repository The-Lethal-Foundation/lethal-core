package profile

import (
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

	// Create a "plugins" dir in BepInEx
	pluginsPath := filepath.Join(profilePath, "BepInEx", "plugins")
	if err := os.Mkdir(pluginsPath, 0755); err != nil {
		return err
	}

	return nil
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

// ListProfiles returns a list of all profiles.
func ListProfiles() ([]string, error) {
	profilesPath := filepath.Join(filesystem.GetDefaultPath(), "LethalCompany", ProfilesDirName)
	dirEntries, err := os.ReadDir(profilesPath)
	if err != nil {
		return nil, err
	}

	var outProfiles []string
	for _, entry := range dirEntries {
		if !entry.IsDir() {
			continue
		}

		outProfiles = append(outProfiles, entry.Name())
	}

	return outProfiles, nil
}
