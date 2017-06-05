package handlers

import (
	"encoding/json"
	"github.com/douglasmakey/ursho/config"
	"log"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"response"`
}

type Error struct {
	Data    interface{} `json:"error"`
	Success bool        `json:"success"`
}

var prefix string

func init() {
	c, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if c.Options.Prefix == "" {
		prefix = ""
	} else {
		prefix = c.Options.Prefix
	}
}

func createError(e Error) []byte {
	d, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return d
}

func createResponse(r Response) []byte {
	d, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return d
}
