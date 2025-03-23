package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TaskArchiveLog = "log:archive"
	TaskDeleteLog  = "log:delete"
)

type LogPayload struct {
	LogID int `json:"log_id"`
}

// Create Redis task for archiving
func NewArchiveLogTask(logID int) (*asynq.Task, error) {
	payload, err := json.Marshal(LogPayload{LogID: logID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskArchiveLog, payload, asynq.Queue("low")), nil
}

// Create Redis task for deletion
func NewDeleteLogTask(logID int) (*asynq.Task, error) {
	payload, err := json.Marshal(LogPayload{LogID: logID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskDeleteLog, payload, asynq.Queue("low")), nil
}
