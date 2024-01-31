# Lethal Core

A core module for lethal mod manager.

## Core modules

1. Config:

   - Keeps track of the Lethal Mod Manager config.
   - Can save / load / initialize the config.

2. Filesystem:

   - Makes sure the required file system is in place, all folders created.
   - Initializes the file system required for the mod manager to run.

3. Modmanager:
   - `bepinex.go` Makes sure that you have the latest version of BepInEx installed.

## Project File Structure

```
LethalModManager/
├── cmd/
│ └── lethalmodmanager/
│   └── main.go # Main application entry point.
├── pkg/
│ ├── config/
│ │ ├── config.go # Functions for managing the JSON config.
│ │ └── config_test.go # Tests for config functions.
│ ├── filesystem/
│ │ ├── filesystem.go # Functions for handling file operations.
│ │ └── filesystem_test.go # Tests for file operations.
│ ├── modmanager/
│ │ ├── modmanager.go # Core mod management functions.
│ │ └── modmanager_test.go # Tests for mod management.
│ ├── profile/
│ │ ├── profile.go # Functions for managing user profiles.
│ │ └── profile_test.go # Tests for profile management.
│ ├── api/
│ │ ├── api.go # Functions to interact with mod APIs.
│ │ └── api_test.go # Tests for API interactions.
│ └── utils/
│   ├── utils.go # Utility functions.
│   └── utils_test.go # Tests for utility functions.
├── .gitignore
├── go.mod
└── README.md
```
