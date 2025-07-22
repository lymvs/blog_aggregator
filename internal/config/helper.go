package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"
const username = "lysandros"

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