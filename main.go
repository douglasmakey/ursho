package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/douglasmakey/ursho/config"
	"github.com/douglasmakey/ursho/handler"
	"github.com/douglasmakey/ursho/storage/postgres"
)

func main() {
	configPath := flag.String("config", "./config/config.json", "path of the config file")

	flag.Parse()

	// Read config
	config, err := config.FromFile(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	// Set use storage, select [Postgres, Filesystem, Redis ...]
	svc, err := postgres.New(config.Postgres.Host, config.Postgres.Port, config.Postgres.User, config.Postgres.Password, config.Postgres.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer svc.Close()

	// Create a server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port),
		Handler: handler.New(config.Options.Prefix, svc),
	}

	go func() {
		// Start server
		log.Printf("Starting HTTP Server. Listening at %q", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("%v", err)
		} else {
			log.Println("Server closed!")
		}
	}()

	// Check for a closing signal
	// Graceful shutdown
	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, os.Interrupt, syscall.SIGTERM)

	sig := <-sigquit
	log.Printf("caught sig: %+v", sig)
	log.Printf("Gracefully shutting down server...")

	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("Unable to shut down server: %v", err)
	} else {
		log.Println("Server stopped")
	}
}
