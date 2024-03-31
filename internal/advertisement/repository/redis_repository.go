package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
)

//go:generate mockery --name=IAdCountRepository --structname=MockAdCountRepository --output=mocks --dir=./internal/advertisement/repository --inpackage --with-expecter --testonly
type IAdCountRepository interface {
	IncrByDate(ctx context.Context, key string) error
	GetByDate(ctx context.Context, key string) (int, error)
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
	countStr, err := r.rdb.Get(ctx, "ads:"+key).Result()
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
