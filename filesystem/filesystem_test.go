package filesystem

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestFileSystem(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "testLethalModManager")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	log.Println(tempDir)
	defer os.RemoveAll(tempDir) // Clean up after the test.

	GetDefaultPath = func() string {
		return tempDir
	}

	if err := InitializeStructure(); err != nil {
		t.Fatalf("InitializeStructure() failed: %v", err)
	}

	requiredDirs := []string{"Profiles"}
	for _, dir := range requiredDirs {
		dirPath := filepath.Join(tempDir, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("Directory %s was not created", dir)
		}
	}

	setUp, err := IsFileSystemSetUp()
	if err != nil {
		t.Fatalf("IsFileSystemSetUp() failed: %v", err)
	}
	if !setUp {
		t.Errorf("IsFileSystemSetUp() returned false, want true")
	}
}
