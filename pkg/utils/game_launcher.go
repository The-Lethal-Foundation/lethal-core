package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/The-Lethal-Foundation/lethal-core/pkg/filesystem"
)

const GameId = "1966720"

// LaunchGameProfile launches the game with the specified profile.
func LaunchGameProfile(profile string) error {
	// Assuming `util.GetProfilePath` resolves the correct profile path.
	profilePath := filepath.Join(filesystem.GetDefaultPath(), "LethalCompany", "Profiles", profile)

	steamPath := `C:\Program Files (x86)\Steam\steam.exe`
	args := []string{
		"-applaunch",
		GameId,
		"--doorstop-enable", "true",
		"--doorstop-target", fmt.Sprintf("%s\\BepInEx\\core\\BepInEx.Preloader.dll", profilePath),
	}

	// Run the command
	cmd := exec.Command(steamPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to launch game: %w", err)
	}

	fmt.Println("Game launched successfully")
	return nil
}
