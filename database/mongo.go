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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a connection string with the username and password
	connectionString := fmt.Sprintf("mongodb://%s:%s@%s/%s", username, password, host, dbname)
	fmt.Println(connectionString)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(dbname)
	col = db.Collection(collectionName)

	return col, nil
}
