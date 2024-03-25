package redis

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
	}
}

func (r *RedisRepository) IncrAdCountsByDate(ctx context.Context, key string) error {
	_, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisRepository) GetAdCountsByDate(ctx context.Context, key string) (int, error) {
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
