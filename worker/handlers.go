package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"github.com/musthafa-vakkayil/event_scheduler_v2/tasks"
)

// LogPayload struct
type LogPayload struct {
	LogID int `json:"log_id"`
}

// ArchiveLogHandler processes the archive task
func (w *Worker) ArchiveLogHandler(ctx context.Context, t *asynq.Task) error {
	var payload LogPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	log.Printf("✅ Archiving log ID: %d", payload.LogID)

	// Mark the log as archived in DB
	err := w.Repo.MarkLogAsArchived(payload.LogID)
	if err != nil {
		return fmt.Errorf("failed to archive log %d: %v", payload.LogID, err)
	}

	log.Printf("✔️ Archived log ID: %d", payload.LogID)

	// 👉 Enqueue the delete task using the shared Redis client
	deleteTask, err := tasks.NewDeleteLogTask(payload.LogID)
	if err != nil {
		return fmt.Errorf("failed to create delete task: %v", err)
	}

	// Schedule the delete task in 46 hours
	_, err = w.RedisClient.Enqueue(deleteTask, asynq.Queue("low"), asynq.ProcessIn(w.Config.LogDeleteDuration))
	if err != nil {
		return fmt.Errorf("failed to enqueue delete task: %v", err)
	}

	log.Printf("🗑️ Enqueued delete task for log ID: %d", payload.LogID)

	return nil
}

// DeleteLogHandler processes the delete task
func (w *Worker) DeleteLogHandler(ctx context.Context, t *asynq.Task) error {
	var payload LogPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	log.Printf("✅ Deleting log ID: %d", payload.LogID)

	// Delete the log from the DB
	err := w.Repo.DeleteLog(payload.LogID)
	if err != nil {
		return fmt.Errorf("failed to delete log %d: %v", payload.LogID, err)
	}

	log.Printf("✔️ Deleted log ID: %d", payload.LogID)
	return nil
}
