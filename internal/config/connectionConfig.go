package config

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database
var ctx = context.Background()
var RDB *redis.Client

func ConnectMongoDB(uri string, dbName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	DB = client.Database(dbName)
}

func ConnectRedis() {
	// Initialize Redis client
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // replace with your Redis server address
		Password: "",               // replace with your password if needed
		DB:       0,                // use default DB
	})

	// Ping Redis to check connection
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
}
