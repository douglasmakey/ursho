package config

import (
	"encoding/json"
	"os"
	"io"
	"bytes"
)

type Config struct {
	Server     Server     `json:"server"`
	Redis      Redis      `json:"redis"`
	Postgres   Postgres   `json:"postgres"`
	Options    Options    `json:"options"`
}

type Server struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Redis struct {
	Host     string `json:"host"`
	Password string `json:"password"`
	DB       string `json:"db"`
}

type Postgres struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	DB       string `json:"db"`
}

type Options struct {
	Prefix string `json:"prefix"`
}

func ReadConfig() (*Config, error) {
	var objectConfig *Config
	var buf bytes.Buffer

	// open input file
	file, err := os.Open("./config/config.json")
	if err != nil {
		return nil, err
	}

	// close file
	defer file.Close()

	// copy to buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		return nil, err
	}

	// Unmarshal data
	if err := json.Unmarshal(buf.Bytes(), &objectConfig); err != nil {
		return nil, err
	}

	return objectConfig, nil
}