package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"github.com/musthafa-vakkayil/event_scheduler_v2/tasks"
)

// LogPayload struct
type LogPayload struct {
	LogID int `json:"log_id"`
}

type EventPayload struct {
	EventID int64 `json:"event_id"`
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

// ScheduleEventHandler processes the schedule event task
func (w *Worker) ScheduleEventHandler(ctx context.Context, t *asynq.Task) error {
	var payload EventPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	log.Printf("✅ Executing event id: %d", payload.EventID)

	// Get the event from the DB
	event, err := w.Repo.GetEvent(int(payload.EventID))
	if err != nil {
		if err.Error() == "event not found" {
			// Skip the task if the event is deleted or not found
			log.Printf("⚠️ Event %d not found. Skipping task.", payload.EventID)
			return nil // Returning nil skips the retry
		}
		return fmt.Errorf("failed to get event %d: %v", payload.EventID, err)
	}

	// Execute the event
	logId, err := w.Repo.ExecuteEvent(event.CreatedBy, event.ID, "SCHEDULED_EVENT", "SUCCESS", event.ApiRequestBody)
	if err != nil {
		return fmt.Errorf("failed to execute event %d: %v", payload.EventID, err)
	}

	// Create an archive task
	archiveTask, err := tasks.NewArchiveLogTask(int(logId))
	if err != nil {
		return fmt.Errorf("failed to create archive task:%v", err)
	}

	// Enqueue the archive task with a 2-minute delay (for testing)
	_, err = w.RedisClient.Enqueue(archiveTask, asynq.Queue("low"), asynq.ProcessIn(w.Config.LogArchiveDuration))
	if err != nil {
		return fmt.Errorf("failed to enqueue archive task:%v", err)
	}

	if event.IsRecurring {
		// Create a new schedule event task
		scheduleEventTask, err := tasks.NewScheduleEventTask(event.ID)
		if err != nil {
			return fmt.Errorf("failed to create schedule event task:%v", err)
		}

		// Enqueue the schedule event task with a delay
		_, err = w.RedisClient.Enqueue(scheduleEventTask, asynq.Queue("critical"), asynq.ProcessAt(time.Now().Add(time.Duration(event.Interval)*time.Minute)))
		if err != nil {
			return fmt.Errorf("failed to enqueue schedule event task:%v", err)
		}

		log.Printf("✔️ Scheduled Event ID for Recurring: %d", payload.EventID)
	}

	return nil
}
