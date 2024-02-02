package utils

import (
	"os"
	"path/filepath"

	"github.com/The-Lethal-Foundation/lethal-core/pkg/filesystem"
	"github.com/otiai10/copy"
)

// CloneOtherProfiles accepts a map of mod manager names to their profiles directory paths.
// It copies all profiles for each specified mod manager to a new location.
func CloneOtherProfiles(modManagers map[string]string) error {
	for managerName, relativeProfilesPath := range modManagers {
		globalProfilesPath := filepath.Join(os.Getenv("APPDATA"), relativeProfilesPath)

		// List all profiles in the directory
		profiles, err := os.ReadDir(globalProfilesPath)
		if err != nil {
			return err
		}

		for _, profile := range profiles {
			if !profile.IsDir() {
				continue
			}

			// Create the new profile path in the manager's directory
			profileName := managerName + "-" + profile.Name()
			newProfilePath := filepath.Join(filesystem.GetDefaultPath(), "Profiles", profileName)
			err := os.MkdirAll(newProfilePath, os.ModePerm)
			if err != nil {
				return err
			}

			// Copy the mods in the profile to the new location.
			oldProfilePath := filepath.Join(globalProfilesPath, profile.Name())
			err = copy.Copy(oldProfilePath, newProfilePath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
