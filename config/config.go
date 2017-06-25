package config

import (
	"encoding/json"
	"os"
	"bufio"
	"io"
	"fmt"
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

func ReadConfig() (*Config, error) {
	var objectConfig *Config

	// open input file
	file, err := os.Open("./config/config.json")
	if err != nil {
		panic(err)
	}

	// close file on exit and check for its returned error
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	// make a read buffer
	r := bufio.NewReader(file)

	// make a buffer to keep chunks that are read
	buf := make([]byte, 1024)

	for {
		// read a chunk
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		// Unmarshal data
		if err := json.Unmarshal(buf[:n], &objectConfig); err != nil {
			return nil, err
		}

	}

	return objectConfig, nil
}