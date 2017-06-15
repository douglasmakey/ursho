package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/douglasmakey/ursho/config"
	"github.com/douglasmakey/ursho/handlers"
	"github.com/douglasmakey/ursho/storages"
)

func main() {
	// Set use storage, select [Postgres, Filesystem, Redis ...]
	storage := &storages.Postgres{}

	// Read config
	config, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}
	// Init storage
	if err = storage.Init(config); err != nil {
		log.Fatal(err)
	}

	// Defers
	defer storage.Close()

	// Handlers
	http.Handle("/encode/", handlers.EncodeHandler(storage))
	http.Handle("/", handlers.RedirectHandler(storage))
	http.Handle("/info/", handlers.DecodeHandler(storage))

	// Graceful shutdown
	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, os.Interrupt, os.Kill)

	// Wait signal
	close := make(chan bool, 1)

	// Create a server
	server := &http.Server{Addr: fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)}

	// Start server
	go func() {
		log.Printf("Starting HTTP Server. Listening at %q", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			if err := server.ListenAndServe(); err != nil {
				if err != http.ErrServerClosed {
					log.Println(err.Error())
				} else {
					log.Println("Server closed!")
				}
				close <- true
			}
		}

	}()

	// Check for a closing signal
	go func() {
		sig := <-sigquit
		log.Printf("caught sig: %+v", sig)
		log.Printf("Gracefully shutting down server...")

		if err := server.Shutdown(nil); err != nil {
			log.Println("Unable to shut down server: " + err.Error())
			close <- true
		} else {
			close <- true
			log.Println("Server stopped")
		}
	}()

	<-close
}
