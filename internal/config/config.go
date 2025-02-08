package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configPath string = "/.config/experiments/.gatorconfig.json"

func Read() (Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}
	fullPath := homeDir + configPath
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return Config{}, err
	}
	var conf Config
	if err = json.Unmarshal(content, &conf); err != nil {
		return Config{}, err
	}
	return conf, nil
}
