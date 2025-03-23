package worker

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/musthafa-vakkayil/event_scheduler_v2/config"
	"github.com/musthafa-vakkayil/event_scheduler_v2/repo"
)

const (
	TaskArchiveLog = "log:archive"
	TaskDeleteLog  = "log:delete"
)

// Worker struct holds Redis server, Redis client, and GORM DB
type Worker struct {
	Config      config.Config
	RedisServer *asynq.Server
	RedisClient *asynq.Client // Shared client to enqueue new tasks
	Repo        repo.Repository
}

// StartWorker initializes and starts the worker
func StartWorker(redisAddr string, redisClient *asynq.Client, cfg config.Config, repo repo.Repository) *Worker {
	// Redis Server for task execution
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	worker := &Worker{
		RedisServer: srv,
		RedisClient: redisClient, // Inject shared Redis client
		Config:      cfg,
		Repo:        repo,
	}

	// Register task handlers
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskArchiveLog, worker.ArchiveLogHandler)
	mux.HandleFunc(TaskDeleteLog, worker.DeleteLogHandler)

	log.Println("🚀 Worker started with Redis server, Redis client, and DB connection...")

	// Run the worker in a goroutine
	go func() {
		if err := srv.Run(mux); err != nil {
			log.Fatalf("Failed to run worker: %v", err)
		}
	}()

	return worker
}
