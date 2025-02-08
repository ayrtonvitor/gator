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

func (c Config) SetUser(username string) error {
	c.CurrentUserName = username
	return write(c)
}

func write(conf Config) error {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	mConf, err := json.Marshal(conf)
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
