package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

// TaskManager defines the interface for enqueuing Redis tasks
type TaskManager interface {
	EnqueueTask(ctx context.Context, logID int, taskType string, queueType string, delay time.Duration) error
	EnqueueTaskIn(ctx context.Context, eventID int64, taskType string, queueType string, delay time.Duration) error
	EnqueueTaskAt(ctx context.Context, eventID int64, taskType string, queueType string, time time.Time) error
}

// RedisTaskManager is the production implementation of TaskManager
type RedisTaskManager struct {
	Client *asynq.Client
}

// NewRedisTaskManager creates a new RedisTaskManager instance
func NewRedisTaskManager(client *asynq.Client) *RedisTaskManager {
	return &RedisTaskManager{Client: client}
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

// EnqueueTaskIn enqueues a task with a delay
func (tm *RedisTaskManager) EnqueueTaskIn(ctx context.Context, eventID int64, taskType string, queueType string, delay time.Duration) error {
	task, err := NewEventTask(eventID, taskType, queueType)
	if err != nil {
		return fmt.Errorf("failed to create %s task for %d: %w", queueType, eventID, err)
	}

	_, err = tm.Client.EnqueueContext(ctx, task, asynq.Queue(queueType), asynq.ProcessIn(delay))
	if err != nil {
		return fmt.Errorf("failed to enqueue %s task for %d: %w", queueType, eventID, err)
	}

	return nil
}

// EnqueueTaskAt enqueues a task with a specific time
func (tm *RedisTaskManager) EnqueueTaskAt(ctx context.Context, eventID int64, taskType string, queueType string, time time.Time) error {
	task, err := NewEventTask(eventID, taskType, queueType)
	if err != nil {
		return fmt.Errorf("failed to create %s task for %d: %w", taskType, eventID, err)
	}

	_, err = tm.Client.EnqueueContext(ctx, task, asynq.Queue(queueType), asynq.ProcessAt(time))
	if err != nil {
		return fmt.Errorf("failed to enqueue %s task for %d: %w", taskType, eventID, err)
	}

	return nil
}
