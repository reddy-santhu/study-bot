package db

import (
	"context"
	"log"

	"time"

	"github.com/reddy-santhu/study-bot/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
		return
	}

	// Create client options
	log.Printf("MongoDB URI from config: %s", cfg.MongoDB.URI)       // Add this line
	log.Printf("MongoDB Dbname from config: %s", cfg.MongoDB.Dbname) // Add this line
	clientOptions := options.Client()

	// Apply URI first to establish basic connection parameters
	clientOptions = clientOptions.ApplyURI(cfg.MongoDB.URI)

	// The connect method does not accept options
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	DB = client.Database(cfg.MongoDB.Dbname)
}

func GetCollection(collectionName string) *mongo.Collection {
	return DB.Collection(collectionName)
}
