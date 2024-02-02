# ✨Lethal Core

A core module library for lethal mod manager.

## Core modules

1. Api:

   - Module for making requests to the external services.
   - `tsapi.go` makes requests to thunderstore api for checking mod versions / downloading mod packages.

2. Config:

   - Keeps track of the Lethal Mod Manager config.
   - `config.go` Initializes / Saves / Loads the config file.

3. Filesystem:

   - Makes sure the required file system is in place, all folders created.
   - `filesystem.go` Initializes the required directories for mod manager. A source of the defaul path for other modules.

4. Modmanager:

   - Takes care of installing / deleting / updating mods.
   - `bepinex.go` Makes sure that you have the latest version of BepInEx installed.
   - `modmanager.go` Installs / Updates / Deletes mods.
   - `unzipmod.go` Takes care of unzipping a mod zip into the plugins directory, and merging files.

5. Profile:

   - Takes care of creating, deleting, renaming profiles.
   - `profile.go` Profile interractions + installing the initial BepInEx into the profile.

6. Utils:
   - Random utilities
   - `constants.go` Contains constants definitions like known mod managers.
   - `game_launcher.go` Takes care of launching the actual game profile.
   - `utils.go` For now primarily contains utilities for cloning profiles from other mod managers.
