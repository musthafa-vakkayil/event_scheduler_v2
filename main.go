package main

import (
	"log"

	"github.com/musthafa-vakkayil/event_scheduler_v2/config"
	"github.com/musthafa-vakkayil/event_scheduler_v2/handlers"
	"github.com/musthafa-vakkayil/event_scheduler_v2/repo"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("unable to read config", err)
	}

	// Initialize GORM DB
	gormDB, err := handlers.ConnectGORM(cfg)
	if err != nil {
		log.Fatal("unable to connect to GORM DB: %w", err)
	}

	// Initialize repository
	repository := repo.NewRepository(gormDB)

	server, err := handlers.NewServer(cfg, repository)
	if err != nil {
		log.Fatal("cannot start server", err)
	}

	err = server.Start(cfg.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
