package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Config struct {
	Auth   Auth    `json:"auth"`
	Guild  Guild   `json:"guild"`
	Sounds []Sound `json:"sounds"`
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
	Keybind string `json:"keybind"`
}

func main() {
	cfg := loadConfig()
	err := cfg.playSound(cfg.Sounds[0])
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
	if resp.StatusCode != http.StatusOK {
		resBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("discord api error: status %d: %s", resp.StatusCode, string(resBody))
	}
	defer resp.Body.Close()
	return nil
}

func loadConfig() Config {
	jsonConfig, err := os.Open(".config.json")
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
