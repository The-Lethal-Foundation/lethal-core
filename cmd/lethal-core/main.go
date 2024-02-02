package main

import (
	"log"

	"github.com/The-Lethal-Foundation/lethal-core/pkg/filesystem"
	"github.com/The-Lethal-Foundation/lethal-core/pkg/modmanager"
	"github.com/The-Lethal-Foundation/lethal-core/pkg/profile"
	"github.com/The-Lethal-Foundation/lethal-core/pkg/utils"
)

func main() {
	err := filesystem.InitializeStructure()
	if err != nil {
		log.Fatal(err)
	}

	filesystem.InitializeStructure()

	// clone other profiles
	utils.CloneOtherProfiles(utils.KnownModManagersList)

	latestVer, err := modmanager.FetchLatestBepInExVersion()
	if err != nil {
		log.Fatal(err)
	}

	_, err = modmanager.DownloadAndCacheBepInEx(filesystem.GetDefaultPath(), latestVer)
	if err != nil {
		log.Fatal(err)
	}

	profile.CreateProfile("Default")

	log.Println("Done!")
}
