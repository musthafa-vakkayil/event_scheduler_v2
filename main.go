package main

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/musthafa-vakkayil/event_scheduler_v2/config"
	"github.com/musthafa-vakkayil/event_scheduler_v2/handlers"
	"github.com/musthafa-vakkayil/event_scheduler_v2/repo"
	"github.com/musthafa-vakkayil/event_scheduler_v2/worker"
)

// @title Event Scheduler API
// @version 1.0
// @description This is an event scheduling service with JWT authentication
// @termsOfService http://swagger.io/terms/

// @contact.name Musthafa
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT token authorization format: "Bearer {token}"
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

	// Create Redis client (shared by server and worker)
	redisClient := asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.RedisUrl})
	defer redisClient.Close()

	server, err := handlers.NewServer(cfg, repository)
	if err != nil {
		log.Fatal("cannot start server", err)
	}

	server.RedisClient = redisClient

	// Start worker with shared Redis client and DB
	worker.StartWorker(cfg.RedisUrl, redisClient, cfg, repository)

	err = server.Start(cfg.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
