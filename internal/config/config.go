package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func Read() (Config, error) {
	var config Config
	file_dir, err := getConfigFile()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(file_dir)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func (c *Config) SetUser(name string) error {
	c.Current_user_name = name
	return write(*c)
}

func write(cfg Config) error {
	json_file, err := getConfigFile()
	if err != nil {
		return err
	}

	file, err := os.Create(json_file)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(&cfg)
	if err != nil {
		return err
	}

	return nil
}

func getConfigFile() (string, error) {
	user_home_dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	file_dir := user_home_dir + "/.gatorconfig.json"

	return file_dir, nil
}
