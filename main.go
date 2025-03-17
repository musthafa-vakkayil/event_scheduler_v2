package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/musthafa-vakkayil/event_scheduler_v2/config"
	db "github.com/musthafa-vakkayil/event_scheduler_v2/db/sqlc"
	"github.com/musthafa-vakkayil/event_scheduler_v2/handlers"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("unable to read config", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := handlers.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot start server", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
