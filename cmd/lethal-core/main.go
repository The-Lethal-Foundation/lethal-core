package main

import (
	"log"

	"github.com/The-Lethal-Foundation/lethal-core/pkg/filesystem"
	"github.com/The-Lethal-Foundation/lethal-core/pkg/utils"
)

func main() {
	err := filesystem.InitializeStructure()
	if err != nil {
		log.Fatal(err)
	}

	// clone other profiles
	utils.CloneOtherProfiles(utils.KnownModManagersList)

	log.Println("Done!")
}
