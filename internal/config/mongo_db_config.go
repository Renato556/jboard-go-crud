package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var mongoClient *mongo.Client

func InitConnection() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable. ")
	}

	opts := options.Client().ApplyURI(uri).
		SetConnectTimeout(30 * time.Second).
		SetServerSelectionTimeout(30 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(opts)
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v", err)
		panic(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Failed to ping MongoDB: %v", err)
		panic(err)
	}

	fmt.Println("Successful connection to MongoDB")
	mongoClient = client
}

func GetMongoClient() *mongo.Client {
	return mongoClient
}

func CloseConnection(ctx context.Context) error {
	if mongoClient == nil {
		return nil
	}
	return mongoClient.Disconnect(ctx)
}
