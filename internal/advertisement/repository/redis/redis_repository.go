package redis

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

// AdCountRepository defines the interface for interacting with Redis to manage advertisement counts.
type AdCountRepository interface {
	IncrByDate(ctx context.Context, key string) error
	GetByDate(ctx context.Context, key string) (int, error)
}

// AdCountRepositoryImpl implements the AdCountRepository interface.
type AdCountRepositoryImpl struct {
	client *redis.Client
}

// NewAdCountRepository creates a new instance of AdCountRepositoryImpl.
func NewAdCountRepository(client *redis.Client) *AdCountRepositoryImpl {
	return &AdCountRepositoryImpl{
		client: client,
	}
}

// IncrByDate increments the count associated with the specified date key in Redis.
func (r *AdCountRepositoryImpl) IncrByDate(ctx context.Context, key string) error {
	_, err := r.client.Incr(ctx, key).Result()
	return err
}

// GetByDate retrieves the count associated with the specified date key from Redis.
func (r *AdCountRepositoryImpl) GetByDate(ctx context.Context, key string) (int, error) {
	countStr, err := r.client.Get(ctx, "ads:"+key).Result()
	if err == redis.Nil {
		// Key does not exist, return count as 0
		return 0, nil
	} else if err != nil {
		// Some other error occurred
		return 0, err
	}

	// Convert the count from string to int
	count, err := strconv.Atoi(countStr)
	if err != nil {
		// Handle the case where the count is not a valid integer
		return 0, err
	}

	return count, nil
}
