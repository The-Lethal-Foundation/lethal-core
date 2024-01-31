package main

import (
	"log"

	"github.com/The-Lethal-Foundation/lethal-core/pkg/filesystem"
	"github.com/The-Lethal-Foundation/lethal-core/pkg/modmanager"
)

func main() {
	err := filesystem.InitializeStructure()
	if err != nil {
		log.Fatal(err)
	}

	// err = modmanager.InstallMod("Test", "Sligili", "More_Emotes")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	mods, err := modmanager.ListMods("Test")
	if err != nil {
		log.Fatal(err)
	}

	for _, mod := range mods {
		log.Println(mod)
	}

	log.Println("Done!")
}
