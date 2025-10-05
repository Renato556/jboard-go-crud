package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func InitConnection() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI environment variable not set")
	}

	opts := options.Client().ApplyURI(uri).
		SetConnectTimeout(30 * time.Second).
		SetServerSelectionTimeout(30 * time.Second).
		SetMaxPoolSize(10).
		SetMinPoolSize(1).
		SetRetryWrites(true).
		SetRetryReads(true)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result bson.M
	if err := client.Database("admin").RunCommand(ctx, bson.D{{"ping", 1}}).Decode(&result); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}

	fmt.Println("MongoDB connected successfully")
	mongoClient = client
}

func GetClient() *mongo.Client {
	if mongoClient == nil {
		log.Printf("MongoDB client not initialized")
		return nil
	}
	return mongoClient
}

func CloseConnection() {
	if mongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Printf("MongoDB disconnect error: %v", err)
		}
		mongoClient = nil
	}
}

func GetCollection(dbName, collectionName string) *mongo.Collection {
	if mongoClient == nil {
		log.Printf("MongoDB client not initialized")
		return nil
	}
	return mongoClient.Database(dbName).Collection(collectionName)
}

func GetJobsCollection(dbName string) *mongo.Collection {
	jobsCollectionName := os.Getenv("MONGODB_JOB_COLLECTION")
	if jobsCollectionName == "" {
		jobsCollectionName = "jobs"
	}
	return GetCollection(dbName, jobsCollectionName)
}

func GetUsersCollection(dbName string) *mongo.Collection {
	usersCollectionName := os.Getenv("MONGODB_USER_COLLECTION")
	if usersCollectionName == "" {
		usersCollectionName = "users"
	}
	return GetCollection(dbName, usersCollectionName)
}
