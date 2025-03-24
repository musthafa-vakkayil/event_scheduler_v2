package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TaskArchiveLog    = "log:archive"
	TaskDeleteLog     = "log:delete"
	TaskScheduleEvent = "event:schedule"
	TaskTestSchedule  = "event:test"
)

type LogPayload struct {
	LogID int `json:"log_id"`
}

type EventPayload struct {
	EventID int64 `json:"event_id"`
}

type TestEventPayload struct {
	Username string `json:"username"`
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

// Create Redis task for scheduling event
func NewScheduleEventTask(eventID int64) (*asynq.Task, error) {
	payload, err := json.Marshal(EventPayload{EventID: eventID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskScheduleEvent, payload, asynq.Queue("critical")), nil
}

// Create Redis task for test scheduling event
func NewTestEventTask(username string) (*asynq.Task, error) {
	payload, err := json.Marshal(TestEventPayload{Username: username})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskScheduleEvent, payload, asynq.Queue("default")), nil
}
