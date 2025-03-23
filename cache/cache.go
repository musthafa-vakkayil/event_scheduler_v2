package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
)

type Cache interface {
	Get(ctx context.Context, archived bool, active bool, page int, size int) ([]models.LogsResponse, error)
	Set(ctx context.Context, logs []models.LogsResponse, archived bool, active bool, page int, size int, expiry time.Duration) error
	DeletePattern(ctx context.Context, pattern string) error
}

func NewRedisCache(rdb *redis.Client) Cache {
	return &LogStore{
		rdb: rdb,
	}
}
