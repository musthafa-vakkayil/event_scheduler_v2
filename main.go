package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"github.com/musthafa-vakkayil/event_scheduler_v2/cache"
	"github.com/musthafa-vakkayil/event_scheduler_v2/config"
	"github.com/musthafa-vakkayil/event_scheduler_v2/handlers"
	"github.com/musthafa-vakkayil/event_scheduler_v2/repo"
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

	gormDB, err := handlers.ConnectGORM(cfg)
	if err != nil {
		log.Fatal("unable to connect to GORM DB:", err)
	}

	repository := repo.NewRepository(gormDB)

	redisClient := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisUrl,
		DB:           1,
		PoolSize:     100,
		MinIdleConns: 10,
		IdleTimeout:  5 * time.Minute,
	})

	queueClient := asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.RedisUrl, DB: 0})

	cacheInstance := cache.NewRedisCache(redisClient)

	server, err := handlers.NewServer(cfg, repository, cacheInstance, queueClient)
	if err != nil {
		log.Fatal("cannot start server", err)
	}

	go func() {
		if err := server.Start(cfg.ServerAddress); err != nil {
			log.Fatal("cannot start server:", err)
		}
	}()

	// Graceful shutdown handler
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gracefully...")

	redisClient.Close()
	queueClient.Close()
}
