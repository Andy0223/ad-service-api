package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB(username string, password string, host string, dbname string, collectionName string) (*mongo.Collection, error) {
	var col *mongo.Collection

	uri := fmt.Sprintf("mongodb://%s:%s@%s/%s", username, password, host, dbname)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a connection string with the username and password
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(dbname)
	col = db.Collection(collectionName)

	return col, nil
}
