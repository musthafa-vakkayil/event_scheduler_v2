package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
)

// Create Redis task for archiving
func NewLogTask(logID int, taskType string, queueType string) (*asynq.Task, error) {
	payload, err := json.Marshal(models.LogPayload{LogID: logID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(taskType, payload, asynq.Queue(queueType)), nil
}

// Create Redis task for scheduling event
func NewEventTask(eventID int64, taskType string, queueType string) (*asynq.Task, error) {
	payload, err := json.Marshal(models.EventPayload{EventID: eventID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(taskType, payload, asynq.Queue(queueType)), nil
}
