package main

import (
	"log"

	"github.com/The-Lethal-Foundation/lethal-core/pkg/filesystem"
)

func main() {
	err := filesystem.InitializeStructure()
	if err != nil {
		log.Fatal(err)
	}

	filesystem.InitializeStructure()

	// clone other profiles
	// utils.CloneOtherProfiles(utils.KnownModManagersList)

	log.Println("Done!")
}
