package config

import (
	"encoding/json"
	"io"
	"os"
)

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