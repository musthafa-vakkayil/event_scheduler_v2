package cache

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/musthafa-vakkayil/event_scheduler_v2/models"
	"github.com/stretchr/testify/assert"
)

func setupRedis(t *testing.T) *redis.Client {
	// Connect to Redis running in Docker
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1,
	})

	// Verify connection
	_, err := client.Ping(context.Background()).Result()
	assert.NoError(t, err)

	return client
}

func TestRedisCacheIntegration(t *testing.T) {
	ctx := context.Background()
	redisClient := setupRedis(t)
	defer redisClient.Close()

	// Create cache instance
	cache := NewRedisCache(redisClient)

	// Test Set and Get operations
	logs := []models.LogsResponse{
		{ID: 1, Name: "Test log"},
	}

	// Set cache
	err := cache.Set(ctx, logs, false, true, 1, 10, 2*time.Minute)
	assert.NoError(t, err)

	// Get cache
	cachedLogs, err := cache.Get(ctx, false, true, 1, 10)
	assert.NoError(t, err)
	assert.NotNil(t, cachedLogs)
	assert.Equal(t, 1, len(cachedLogs))
	assert.Equal(t, "Test log", cachedLogs[0].Name)

	// Delete cache
	err = cache.DeletePattern(ctx, "event-scheduler*")
	assert.NoError(t, err)

	// Verify cache is deleted
	cachedLogs, err = cache.Get(ctx, false, true, 1, 10)
	assert.NoError(t, err)
	assert.Nil(t, cachedLogs)
}
