package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
)

type LogStore struct {
	rdb *redis.Client
}

func (s *LogStore) Get(ctx context.Context, archived bool, active bool, page int, size int) ([]models.LogsResponse, error) {
	const cachePrefix = "event-scheduler"
	cacheKey := fmt.Sprintf("%s:log-%v-%v-%d-%d", cachePrefix, archived, active, page, size)

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss, return nil without error
	} else if err != nil {
		return nil, fmt.Errorf("redis GET error: %w", err)
	}

	// Skip unmarshalling if empty string (safe guard)
	if data == "" {
		return nil, nil
	}

	var logs []models.LogsResponse
	err = json.Unmarshal([]byte(data), &logs)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal logs: %w", err)
	}

	return logs, nil
}

func (s *LogStore) Set(ctx context.Context, logs []models.LogsResponse, archived bool, active bool, page int, size int, expiry time.Duration) error {
	const cachePrefix = "event-scheduler"
	cacheKey := fmt.Sprintf("%s:log-%v-%v-%d-%d", cachePrefix, archived, active, page, size)

	jsonData, err := json.Marshal(logs)
	if err != nil {
		return fmt.Errorf("failed to marshal logs: %w", err)
	}

	pipe := s.rdb.TxPipeline()
	pipe.SetEX(ctx, cacheKey, jsonData, expiry)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

func (s *LogStore) DeletePattern(ctx context.Context, pattern string) error {
	iter := s.rdb.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := s.rdb.Del(ctx, iter.Val()).Err(); err != nil {
			return fmt.Errorf("failed to delete cache: %w", err)
		}
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("redis scan error: %w", err)
	}
	return nil
}
