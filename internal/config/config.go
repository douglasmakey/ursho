package config

import (
	"encoding/json"
	"os"
)

// Config contains the configuration of the url shortener.
type Config struct {
	Server struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"server"`
	Redis struct {
		Host     string `json:"host"`
		Password string `json:"password"`
		DB       string `json:"db"`
	} `json:"redis"`
	Postgres struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DB       string `json:"db"`
	} `json:"postgres"`
	Sqlite struct {
		FilePath string `json:"filepath"`
	} `json:"sqlite"`
	Options struct {
		Prefix string `json:"prefix"`
	} `json:"options"`
}

// FromFile returns a configuration parsed from the given file.
func FromFile(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
