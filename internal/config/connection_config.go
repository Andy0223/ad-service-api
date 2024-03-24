package config

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var RDB *redis.Client

func ConnectMongoDB(uri string, dbname string) (*mongo.Database, error) {
	var db *mongo.Database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	db = client.Database(dbname)

	return db, nil

}

func ConnectRedis(addr string, password string, db int) (*redis.Client, error) {
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
