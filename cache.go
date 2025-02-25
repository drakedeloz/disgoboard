package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func (cfg *Config) commandCacheSound(soundName string) error {
	usrHome := getUserHome()
	cachePath := filepath.Join(usrHome, ".cache", "disgoboard")

	if err := os.MkdirAll(cachePath, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	soundPath := filepath.Join(cachePath, cfg.Sounds[soundName].ID+".ogg")
	if _, err := os.Stat(soundPath); err != nil {
		fmt.Printf("Caching %s...\n", soundName)
		err = cacheFile(cfg.Sounds[soundName].ID, cachePath)
		if err != nil {
			return fmt.Errorf("failed to cache file: %s", soundName)
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

func cacheFile(soundID, cachePath string) error {
	url := "https://cdn.discordapp.com/soundboard-sounds/" + soundID
	filePath := filepath.Join(cachePath, soundID+".ogg")
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		resBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("discord cdn error: status %d: %s", resp.StatusCode, string(resBody))
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}
