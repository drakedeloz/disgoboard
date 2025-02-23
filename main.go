package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Config struct {
	Auth   Auth             `json:"auth"`
	Guild  Guild            `json:"guild"`
	Sounds map[string]Sound `json:"sounds"`
}

type Auth struct {
	Token string `json:"token"`
}

type Guild struct {
	ID       string   `json:"id"`
	Channels Channels `json:"channels"`
}

type Channels struct {
	PrimaryChannelID   string `json:"primaryChannelID"`
	SecondaryChannelID string `json:"secondaryChannelID"`
}

type Sound struct {
	ID      string `json:"id"`
	GuildID string `json:"sourceGuildID"`
}

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Println("Usage: disgoboard play <sound name>")
		os.Exit(1)
	}

	cfg := loadConfig()
	sound, found := cfg.Sounds[args[1]]
	if !found {
		fmt.Printf("Undefined sound: %s\n", args[1])
		os.Exit(1)
	}
	err := cfg.playSound(sound)
	if err != nil {
		fmt.Println(err)
	}
}

func (cfg *Config) playSound(sbItem Sound) error {
	url := fmt.Sprintf("https://discord.com/api/v10/channels/%s/send-soundboard-sound", cfg.Guild.Channels.SecondaryChannelID)
	jsonStr := fmt.Sprintf(`{"sound_id": "%s", "source_guild_id": "%s"}`, sbItem.ID, sbItem.GuildID)
	body := bytes.NewBufferString(jsonStr)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", cfg.Auth.Token)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode > 299 {
		resBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("discord api error: status %d: %s", resp.StatusCode, string(resBody))
	}
	defer resp.Body.Close()
	return nil
}

func loadConfig() Config {
	usrHome, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("could not find users home directory")
		os.Exit(1)
	}
	configPath := filepath.Join(usrHome, ".config", "disgoboard", "config.json")
	jsonConfig, err := os.Open(configPath)
	if err != nil {
		fmt.Println("could not load config file")
		os.Exit(1)
	}
	defer jsonConfig.Close()

	byteValue, err := io.ReadAll(jsonConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var cfg Config
	err = json.Unmarshal(byteValue, &cfg)
	if err != nil {
		fmt.Println(err)
	}
	return cfg
}
