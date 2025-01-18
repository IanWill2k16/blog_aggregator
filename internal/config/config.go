package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/IanWill2k16/blog_aggregator/internal/database"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

type State struct {
	Db  *database.Queries
	Cfg *Config
}

func Read() (Config, error) {
	user_dir, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(user_dir)
	if err != nil {
		return Config{}, err
	}

	config := Config{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, fmt.Errorf("error unmarshaling json data: %w", err)
	}

	return config, nil
}

func (c *Config) SetUser(user_name string) {
	c.CurrentUserName = user_name
	write(*c)
}

func getConfigFilePath() (string, error) {
	user_dir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error finding home directory: %w", err)
	}
	return user_dir + "/" + configFileName, nil
}

func write(cfg Config) error {
	filepath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(cfg); err != nil {
		return err
	}
	return nil
}
