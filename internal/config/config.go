package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const configFileName string = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigPath() string {
	// get config file path
	homePath, _ := os.UserHomeDir()
	return homePath + "/" + configFileName
}

func Read() (Config, error) {
	// initialize empty Config
	var config Config

	// get config file path
	configPath := getConfigPath()

	// open config file
	jsonFile, err := os.Open(configPath)
	if err != nil {
		return config, fmt.Errorf("error opening config file: %w", err)
	}
	defer jsonFile.Close()

	// read config file
	body, err := io.ReadAll(jsonFile)
	if err != nil {
		return config, fmt.Errorf("error reading config file: %w", err)
	}

	// unmarshal JSON to struct
	if err = json.Unmarshal(body, &config); err != nil {
		return config, fmt.Errorf("error unmarshalling data: %w", err)
	}

	return config, nil

}

func (c *Config) SetUser(user string) error {
	c.CurrentUserName = user

	// get config file path
	configPath := getConfigPath()

	// open file
	file, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	//marshal json
	data, err := json.MarshalIndent(c, "", "	")
	if err != nil {
		return fmt.Errorf("error marshalling data: %w", err)
	}

	return os.WriteFile(configPath, data, 0644)
}
