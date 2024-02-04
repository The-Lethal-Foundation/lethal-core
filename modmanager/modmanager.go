package modmanager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/The-Lethal-Foundation/lethal-core/api"
	"github.com/The-Lethal-Foundation/lethal-core/filesystem"
)

type ModManifest struct {
	Name         string   `json:"name"`
	Version      string   `json:"version_number"`
	Description  string   `json:"description"`
	Dependencies []string `json:"dependencies"`
}

type ModDetails struct {
	Author   string      `json:"author"`
	Manifest ModManifest `json:"manifest"`
}

// InstallMod installs or updates the specified mod in the given profile.
func InstallMod(profileName, modAuthor, modName string) error {

	// Retrieve the latest mod version.
	modInfo, err := api.FetchModDetails(modAuthor, modName)
	if err != nil {
		return fmt.Errorf("error getting mod info: %w", err)
	}
	modVersion := modInfo.LatestVersion

	// Check if the mod already exists or if the version is outdated.
	exists, err := isLocalModExists(profileName, modAuthor, modName, modVersion)
	if err != nil {
		return fmt.Errorf("error checking if mod exists: %w", err)
	}

	outdated, err := isLocalModOutdated(profileName, modAuthor, modName, modVersion)
	if err != nil {
		return fmt.Errorf("error checking if mod is outdated: %w", err)
	}

	if exists && !outdated {
		return nil
	}

	// Download the mod to a temporary folder.
	zipName, err := api.DownloadModPackage(modAuthor, modName, modVersion)
	if err != nil {
		return fmt.Errorf("error downloading mod: %w", err)
	}

	// Unzip the mod to the profile folder.
	finalModPath := filepath.Join(filesystem.GetDefaultPath(), "LethalCompany", "profiles", profileName, "BepInEx", "plugins", fmt.Sprintf("%s-%s-%s", modAuthor, modName, modVersion))
	err = UnzipMod(zipName, finalModPath)
	if err != nil {
		return fmt.Errorf("error unzipping mod: %w", err)
	}

	// Read the mod manifest.
	var modDetails ModDetails
	modDetails.Author = modAuthor

	manifestPath := filepath.Join(finalModPath, "manifest.json")
	modDetails.Manifest, err = ReadModManifest(manifestPath)
	if err != nil {
		return fmt.Errorf("error reading mod manifest: %w", err)
	}

	// Install or update dependencies.
	return installDependencies(profileName, modDetails)
}

// installDependencies handles the installation or updating of mod dependencies.
func installDependencies(profileName string, mod ModDetails) error {
	for _, dep := range mod.Manifest.Dependencies {

		// Split dep into modAuthor and modName and modVersion.
		depSplit := strings.Split(dep, "-")

		// Skip the dependency if it's BepInEx.
		if depSplit[0] == "BepInEx" {
			continue
		}

		if len(depSplit) != 3 {
			return fmt.Errorf("invalid dependency format: %s", dep)
		}

		if err := InstallMod(profileName, depSplit[0], depSplit[1]); err != nil {
			return fmt.Errorf("error installing dependency: %w", err)
		}
	}
	return nil
}

// isLocalModExists checks if the local mod exists in the specified profile.
func isLocalModExists(profileName, modAuthor, modName, modVersion string) (bool, error) {
	modPath := filepath.Join(filesystem.GetDefaultPath(), "LethalCompany", "Profiles", profileName, "BepInEx", "plugins", modName)

	// Reading the directory content
	files, err := os.ReadDir(modPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // Mod directory doesn't exist
		}
		return false, fmt.Errorf("error reading mods directory: %w", err)
	}

	// Loop over all files/directories in the mods directory
	for _, file := range files {
		if file.IsDir() {
			dirName := file.Name()
			// Check if the directory name matches the modName, modAuthor, and modVersion
			if dirName == fmt.Sprintf("%s-%s-%s", modAuthor, modName, modVersion) {
				return true, nil
			}
		}
	}

	return false, nil
}

// isLocalModOutdated checks if the local mod version is older than the latest available version.
func isLocalModOutdated(profileName, modAuthor, modName, modVersion string) (bool, error) {
	modInfo, err := api.FetchModDetails(modAuthor, modName)
	if err != nil {
		return false, fmt.Errorf("error getting mod info: %w", err)
	}

	// Compare the local mod version with the latest version.
	isOutdated := modVersion != modInfo.LatestVersion
	return isOutdated, nil
}

// DeleteMod deletes a mod.
func DeleteMod(profileName, modAuthor, modName, modVersion string) error {

	modDirPath := filepath.Join(filesystem.GetDefaultPath(), "LethalCompany", "Profiles", profileName, "BepInEx", "plugins", fmt.Sprintf("%s-%s-%s", modAuthor, modName, modVersion))

	err := os.RemoveAll(modDirPath)
	if err != nil {
		return fmt.Errorf("error deleting mod: %w", err)
	}

	return nil
}

// EnableMod enables a mod.
func EnableMod(modName, profileName string) error {
	// Implementation for enabling a mod.

	return nil
}

// DisableMod disables a mod.
func DisableMod(modName, profileName string) error {
	// Implementation for disabling a mod.

	return nil
}

// ListMods returns a list of all mods.
func ListMods(profileName string) ([]ModDetails, error) {

	pluginsDir := filepath.Join(filesystem.GetDefaultPath(), "LethalCompany", "Profiles", profileName, "BepInEx", "plugins")

	// Reading the directory content
	files, err := os.ReadDir(pluginsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // Mod directory doesn't exist
		}
		return nil, fmt.Errorf("error reading mods directory: %w", err)
	}

	var modDetails []ModDetails

	// Loop over all files/directories in the mods directory
	for _, file := range files {
		if file.IsDir() {
			dirName := file.Name()
			// Check if the directory name matches the modName, modAuthor, and modVersion, and if so, check if the maifest.json file exists
			if strings.Contains(dirName, "-") {
				manifestPath := filepath.Join(pluginsDir, dirName, "manifest.json")
				if _, err := os.Stat(manifestPath); err == nil {
					var modDetail ModDetails
					modDetail.Author = strings.Split(dirName, "-")[0]
					modDetail.Manifest, err = ReadModManifest(manifestPath)
					if err != nil {
						return nil, fmt.Errorf("error reading mod manifest: %w", err)
					}
					modDetails = append(modDetails, modDetail)
				}
			}
		}
	}

	return modDetails, nil
}

func ReadModManifest(manifestPath string) (ModManifest, error) {
	var modManifest ModManifest

	// Open the manifest file.
	manifestFile, err := os.Open(manifestPath)
	if err != nil {
		return modManifest, fmt.Errorf("error opening manifest file: %w", err)
	}
	defer manifestFile.Close()

	// Decode the manifest file.
	if err := json.NewDecoder(manifestFile).Decode(&modManifest); err != nil {
		return modManifest, fmt.Errorf("error decoding manifest file: %w", err)
	}

	return modManifest, nil
}
