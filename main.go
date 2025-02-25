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
	UserID   string `json:"userID"`
	Token    string `json:"token"`
	BotToken string `json:"botToken"`
}

type Guild struct {
	ID string `json:"id"`
}

type Channel struct {
	ChannelID string `json:"channel_id"`
	DeafState bool   `json:"self_deaf"`
}

type Sound struct {
	ID      string `json:"id"`
	GuildID string `json:"sourceGuildID"`
}

func main() {
	cfg := loadConfig()
	args := os.Args[1:]

	if args[0] == "add" && len(args) == 3 {
		name, err := cfg.validateSound(args[1], args[2])
		if err != nil {
			fmt.Print(err)
			return
		}
		cfg.addToSounds(name, args[1], args[2])
		err = cfg.commandCacheSound(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("config sounds cached successfuly")
		return
	}

	if len(args) != 2 {
		fmt.Println("Usage: disgoboard play <sound name>")
		os.Exit(1)
	}

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
	usrChannel, err := cfg.getUserChannel()
	if err != nil {
		return err
	}

	if usrChannel.DeafState {
		return fmt.Errorf("cannot play sound while deafened")
	}

	url := fmt.Sprintf("https://discord.com/api/v10/channels/%s/send-soundboard-sound", usrChannel.ChannelID)
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

func (cfg *Config) getUserChannel() (Channel, error) {
	url := fmt.Sprintf("https://discord.com/api/v10/guilds/%s/voice-states/%s", cfg.Guild.ID, cfg.Auth.UserID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Channel{}, err
	}
	authString := fmt.Sprintf("Bot %s", cfg.Auth.BotToken)
	req.Header.Set("Authorization", authString)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Channel{}, err
	}

	if resp.StatusCode > 299 {
		return Channel{}, fmt.Errorf("could not get user channel: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Channel{}, err
	}

	var usrChannel Channel
	err = json.Unmarshal(body, &usrChannel)
	if err != nil {
		return Channel{}, err
	}
	return usrChannel, nil
}

func loadConfig() Config {
	usrHome := getUserHome()
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

func (cfg *Config) validateSound(soundID, sourceGuildID string) (string, error) {
	url := fmt.Sprintf("https://discord.com/api/v10/guilds/%s/soundboard-sounds/%s", sourceGuildID, soundID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", cfg.Auth.Token)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode > 299 {
		return "", fmt.Errorf("could not get sound info: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	validSound := struct {
		Name      string `json:"name"`
		Available bool   `json:"available"`
	}{}

	err = json.Unmarshal(body, &validSound)
	if err != nil {
		return "", err
	}

	if !validSound.Available {
		return "", fmt.Errorf("sound not available: %s", validSound.Name)
	}
	return validSound.Name, nil
}

func (cfg *Config) addToSounds(name, soundID, sourceGuildID string) error {
	addSound := Sound{
		ID:      soundID,
		GuildID: sourceGuildID,
	}

	cfg.Sounds[name] = addSound
	jsonData, err := json.MarshalIndent(cfg, "", "	")
	if err != nil {
		return err
	}

	configPath := filepath.Join(getUserHome(), ".config", "disgoboard", "config.json")
	err = os.WriteFile(configPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing json to file: %v", err)
	}

	fmt.Printf("%s added to Config\n", name)
	return nil
}
