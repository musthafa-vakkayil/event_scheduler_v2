package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/musthafa-vakkayil/event_scheduler_v2/config"
	"github.com/musthafa-vakkayil/event_scheduler_v2/handlers"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("unable to read config", err)
	}
	// conn, err := sql.Open(config.DBDriver, config.DBSource)
	// if err != nil {
	// 	log.Fatal("cannot connect to db:", err)
	// }

	// store := db.NewStore(conn)

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
