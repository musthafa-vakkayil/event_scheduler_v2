package main

import (
	"log"

	"github.com/musthafa-vakkayil/event_scheduler_v2/config"
	"github.com/musthafa-vakkayil/event_scheduler_v2/handlers"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("unable to read config", err)
	}

	// Initialize GORM DB
	gormDB, err := handlers.ConnectGORM(config)
	if err != nil {
		log.Fatal("unable to connect to GORM DB: %w", err)
	}
	server, err := handlers.NewServer(config, gormDB)
	if err != nil {
		log.Fatal("cannot start server", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
