package config

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg *Config) SetUser() error {
	// cfg.CurrentUserName = username
	return write(*cfg)
}

func Read() (cfg Config, err error) {
	path, err := getConfigFilepath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return Config{}, err
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func getConfigFilepath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, configFileName), nil
}

func write(cfg Config) error {
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	path, err := getConfigFilepath()
	if err != nil {
		return err
	}

	os.WriteFile(path, jsonData, 0644)

	return nil
}
