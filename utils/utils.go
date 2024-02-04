package utils

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/The-Lethal-Foundation/lethal-core/filesystem"
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
			newProfilePath := filepath.Join(filesystem.GetDefaultPath(), "LethalCompany", "Profiles", profileName)
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

// Retusn mod author and name from the Thunderstore mod URL.
func ParseThunderstoreModUrl(modUrl string) (string, string, error) {

	// Remove the trailing slash if exists
	modUrl = strings.TrimSuffix(modUrl, "/")

	// Parse the URL
	parsedUrl, err := url.Parse(modUrl)
	if err != nil {
		return "", "", fmt.Errorf("error parsing URL: %w", err)
	}

	// Split the path into segments
	segments := strings.Split(path.Clean(parsedUrl.Path), "/")[1:]

	// Assuming the URL format is like https://thunderstore.io/c/lethal-company/p/namespace/modname
	// and that there are at least 5 segments ("/c/lethal-company/p/namespace/modname")
	if len(segments) < 5 {
		return "", "", fmt.Errorf("invalid mod URL format")
	}

	return segments[len(segments)-2], segments[len(segments)-1], nil
}
