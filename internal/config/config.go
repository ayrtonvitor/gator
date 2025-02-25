package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
	ConnString      string `json:"connection_string"`
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

func (c Config) SetUser(username string, id uuid.UUID) error {
	c.CurrentUserName = username

	err := write(c)
	if err != nil {
		return fmt.Errorf("Could not set user %s", username)
	}
	fmt.Printf("Successfully set %s as user\n", username)
	return nil
}

func write(conf Config) error {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	mConf, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(fullPath, mConf, 0666)
	if err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	fullPath := homeDir + configPath
	if err != nil {
		return "", err
	}
	return fullPath, nil
}
