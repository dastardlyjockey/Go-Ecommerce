package database

import (
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

func DBClient() *mongo.Client {
	_ = godotenv.Load()

	dbURL := os.Getenv("MONGODB_URL")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURL))
	if err != nil {
		log.Fatal("Error connecting to Mongo database server: ", err)
	}

	return client
}

var Client = DBClient()

func Collection(client *mongo.Client, name string) *mongo.Collection {
	collection := client.Database("E-commerce").Collection(name)
	return collection
}

func CloseClient(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Fatal("Error disconnecting from the database ", err)
	}

	log.Println("Disconnecting from the database")
}
