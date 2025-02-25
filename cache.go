package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func (cfg *Config) commandCacheSounds() error {
	usrHome := getUserHome()
	cachePath := filepath.Join(usrHome, ".cache", "disgoboard")
	if !directoryExists(cachePath) {
		err := os.Mkdir(cachePath, 0755)
		if err != nil {
			return fmt.Errorf("Error creating cache directory: %v\n", err)
		}
	}

	return nil

}

func getUserHome() string {
	usrHome, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("could not find users home directory")
		os.Exit(1)
	}
	return usrHome
}

func directoryExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
