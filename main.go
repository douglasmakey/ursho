package main

import (
	"log"
	"net/http"
	"runtime"
	"github.com/douglasmakey/ursho/handlers"
	"github.com/douglasmakey/ursho/storages"
	"github.com/douglasmakey/ursho/config"
)

func main() {
	// Sets the maximum number of CPUs
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Set use storage, select [Postgres, Filesystem, Redis ...]
	storage := &storages.Postgres{}

	// Read config
	config, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}
	// Init storage
	err = storage.Init(config)
	if err != nil {
		log.Fatal(err)
	}

	// Defers
	defer storage.Close()

	// Handlers
	http.Handle("/encode/", handlers.EncodeHandler(storage))
	http.Handle("/", handlers.RedirectHandler(storage))
	http.Handle("/info/", handlers.DecodeHandler(storage))

	// Init server
	err = http.ListenAndServe(config.Server.Host + ":" + config.Server.Port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
