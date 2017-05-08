package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Server     Server     `json:"server"`
	Redis      Redis      `json:"redis"`
	Postgres   Postgres   `json:"postgres"`
	Filesystem Filesystem `json:"files_system"`
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

type Filesystem struct {
	Path string `json:"path"`
}

type Options struct {
	Prefix string `json:"prefix"`
}

func ReadConfig() (Config, error) {
	file, e := ioutil.ReadFile("./config/config.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		return Config{}, e
	}

	var objectConfig Config
	json.Unmarshal(file, &objectConfig)
	return objectConfig, nil
}
