package repository

import (
	"ad-service-api/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type IAdCountRepository interface {
	IncrByDate(ctx context.Context, key string) error
	GetByDate(ctx context.Context, key string) (int, error)
	GetAdsByKey(ctx context.Context, key string) ([]*models.Advertisement, error)
	SetAdsByKey(ctx context.Context, key string, ads []*models.Advertisement, expiration time.Duration) error
	DeleteAdsByPattern(ctx context.Context, pattern string) error
}

// AdCountRepositoryImpl implements the AdCountRepository interface.
type AdCountRepository struct {
	rdb *redis.Client
}

// NewAdCountRepository creates a new instance of AdCountRepositoryImpl.
func NewAdCountRepository(rdb *redis.Client) *AdCountRepository {
	return &AdCountRepository{
		rdb: rdb,
	}
}

// IncrByDate increments the count associated with the specified date key in Redis.
func (r *AdCountRepository) IncrByDate(ctx context.Context, key string) error {
	_, err := r.rdb.Incr(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to increment count for key %s: %w", key, err)
	}
	return nil
}

// GetByDate retrieves the count associated with the specified date key from Redis.
func (r *AdCountRepository) GetByDate(ctx context.Context, key string) (int, error) {
	countStr, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		// Key does not exist, return count as 0
		return 0, nil
	} else if err != nil {
		// Some other error occurred
		return 0, fmt.Errorf("failed to get count for key %s: %w", key, err)
	}

	// Convert the count from string to int
	count, err := strconv.Atoi(countStr)
	if err != nil {
		// Handle the case where the count is not a valid integer
		return 0, fmt.Errorf("failed to convert count to integer for key %s: %w", key, err)
	}

	return count, nil
}

// GetAdsByKey retrieves the advertisements associated with the specified key from Redis.
func (r *AdCountRepository) GetAdsByKey(ctx context.Context, key string) ([]*models.Advertisement, error) {
	adsData, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		// Key does not exist, return nil
		return nil, nil
	}
	if err != nil {
		// Some other error occurred
		return nil, fmt.Errorf("failed to get ads for key %s: %w", key, err)
	}

	var ads []*models.Advertisement
	err = json.Unmarshal([]byte(adsData), &ads)
	if err != nil {
		// Error occurred during unmarshalling
		return nil, fmt.Errorf("failed to unmarshal ads data: %w", err)
	}

	return ads, nil
}

// SetAdsByKey sets the advertisements associated with the specified key in Redis.
func (r *AdCountRepository) SetAdsByKey(ctx context.Context, key string, ads []*models.Advertisement, expiration time.Duration) error {
	adsData, err := json.Marshal(ads)
	if err != nil {
		// Error occurred during marshalling
		return fmt.Errorf("failed to marshal ads data: %w", err)
	}

	err = r.rdb.Set(ctx, key, adsData, expiration).Err()
	if err != nil {
		// Some other error occurred
		return fmt.Errorf("failed to set ads for key %s: %w", key, err)
	}

	return nil
}

// DeleteAdsByPattern deletes all keys that start with the specified pattern from Redis.
func (r *AdCountRepository) DeleteAdsByPattern(ctx context.Context, pattern string) error {
	keys, err := r.rdb.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys for pattern %s: %w", pattern, err)
	}

	if len(keys) > 0 {
		_, err = r.rdb.Del(ctx, keys...).Result()
		if err != nil {
			return fmt.Errorf("failed to delete keys for pattern %s: %w", pattern, err)
		}
	}

	return nil
}
