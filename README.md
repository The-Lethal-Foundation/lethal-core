# ✨Lethal Core

A core module library for lethal mod manager.

## Core modules

1. Api:

   - Module for making requests to the external services.
   - `tsapi.go` makes requests to thunderstore api for checking mod versions / downloading mod packages.

1. Config:

   - Keeps track of the Lethal Mod Manager config.
   - `config.go` Initializes / Saves / Loads the config file.

1. Filesystem:

   - Makes sure the required file system is in place, all folders created.
   - `filesystem.go` Initializes the required directories for mod manager. A source of the defaul path for other modules.

1. Modmanager:

   - Takes care of installing / deleting / updating mods.
   - `bepinex.go` Makes sure that you have the latest version of BepInEx installed.
   - `modmanager.go` Installs / Updates / Deletes mods.
   - `unzipmod.go` Takes care of unzipping a mod zip into the plugins directory, and merging files.

1. Profile:

   - Takes care of creating, deleting, renaming profiles.
   - `profile.go` Profile interractions + installing the initial BepInEx into the profile.

1. Utils:
   - Random utilities
   - `game_launcher.go` Takes care of launching the actual game profile.
