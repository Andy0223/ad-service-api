package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func ConnectRedis(addr string, password string, db int) (*redis.Client, error) {
	var RDB *redis.Client
	ctx := context.Background()
	// Initialize Redis client
	RDB = redis.NewClient(&redis.Options{
		Addr:     addr,     // replace with your Redis server address
		Password: password, // replace with your password if needed
		DB:       db,       // use default DB
	})

	// Ping Redis to check connection
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	return RDB, nil
}
