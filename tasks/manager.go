package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

// TaskManager defines the interface for enqueuing Redis tasks
type TaskManager interface {
	EnqueueArchiveTask(ctx context.Context, logID int, delay time.Duration) error
	EnqueueDeleteTask(ctx context.Context, logID int, delay time.Duration) error
	EnqueueScheduleTaskIn(ctx context.Context, eventID int64, delay time.Duration) error
	EnqueueScheduleTaskAt(ctx context.Context, eventID int64, time time.Time) error
	EnqueueTestTaskAt(ctx context.Context, username string, time time.Time) error
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
func (tm *RedisTaskManager) EnqueueArchiveTask(ctx context.Context, logID int, delay time.Duration) error {
	task, err := NewArchiveLogTask(logID)
	if err != nil {
		return fmt.Errorf("failed to create archive task: %w", err)
	}

	_, err = tm.Client.EnqueueContext(ctx, task, asynq.Queue("low"), asynq.ProcessIn(delay))
	if err != nil {
		return fmt.Errorf("failed to enqueue archive task: %w", err)
	}

	return nil
}

// EnqueueDeleteTask enqueues a delete task with a delay
func (tm *RedisTaskManager) EnqueueDeleteTask(ctx context.Context, logID int, delay time.Duration) error {
	task, err := NewDeleteLogTask(logID)
	if err != nil {
		return fmt.Errorf("failed to create delete task: %w", err)
	}

	_, err = tm.Client.EnqueueContext(ctx, task, asynq.Queue("low"), asynq.ProcessIn(delay))
	if err != nil {
		return fmt.Errorf("failed to enqueue delete task: %w", err)
	}

	return nil
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
