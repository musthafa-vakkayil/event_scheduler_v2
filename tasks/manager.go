package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

// TaskManager defines the interface for enqueuing Redis tasks
type TaskManager interface {
	EnqueueScheduleTaskIn(ctx context.Context, eventID int64, delay time.Duration) error
	EnqueueScheduleTaskAt(ctx context.Context, eventID int64, time time.Time) error
	EnqueueTestTaskAt(ctx context.Context, username string, time time.Time) error
	EnqueueTask(ctx context.Context, logID int, taskType string, queueType string, delay time.Duration) error
}

// RedisTaskManager is the production implementation of TaskManager
type RedisTaskManager struct {
	Client *asynq.Client
}

// NewRedisTaskManager creates a new RedisTaskManager instance
func NewRedisTaskManager(client *asynq.Client) *RedisTaskManager {
	return &RedisTaskManager{Client: client}
}

// EnqueueScheduleTaskIn enqueues a schedule task with a delay
func (tm *RedisTaskManager) EnqueueScheduleTaskIn(ctx context.Context, eventID int64, delay time.Duration) error {
	task, err := NewScheduleEventTask(eventID)
	if err != nil {
		return fmt.Errorf("failed to create schedule task: %w", err)
	}

	_, err = tm.Client.EnqueueContext(ctx, task, asynq.Queue("critical"), asynq.ProcessIn(delay))
	if err != nil {
		return fmt.Errorf("failed to enqueue schedule task: %w", err)
	}

	return nil
}

// EnqueueScheduleTaskAt enqueues a schedule task with a specific time
func (tm *RedisTaskManager) EnqueueScheduleTaskAt(ctx context.Context, eventID int64, time time.Time) error {
	task, err := NewScheduleEventTask(eventID)
	if err != nil {
		return fmt.Errorf("failed to create schedule task: %w", err)
	}

	_, err = tm.Client.EnqueueContext(ctx, task, asynq.Queue("critical"), asynq.ProcessAt(time))
	if err != nil {
		return fmt.Errorf("failed to enqueue schedule task: %w", err)
	}

	return nil
}

// EnqueueScheduleTaskAt enqueues a schedule task with a specific time
func (tm *RedisTaskManager) EnqueueTestTaskAt(ctx context.Context, username string, time time.Time) error {
	task, err := NewTestEventTask(username)
	if err != nil {
		return fmt.Errorf("failed to create test task: %w", err)
	}

	_, err = tm.Client.EnqueueContext(ctx, task, asynq.Queue("defualt"), asynq.ProcessAt(time))
	if err != nil {
		return fmt.Errorf("failed to enqueue test task: %w", err)
	}

	return nil
}

// EnqueueArchiveTask enqueues an archive task with a delay
func (tm *RedisTaskManager) EnqueueTask(ctx context.Context, logID int, taskType string, queueType string, delay time.Duration) error {
	task, err := NewLogTask(logID, taskType, queueType)
	if err != nil {
		return fmt.Errorf("failed to create %s task for %d: %w", taskType, logID, err)
	}

	_, err = tm.Client.EnqueueContext(ctx, task, asynq.Queue(queueType), asynq.ProcessIn(delay))
	if err != nil {
		return fmt.Errorf("failed to enqueue %s task for %d: %w", taskType, logID, err)
	}

	return nil
}
