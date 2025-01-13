package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"errors"

	"github.com/reddy-santhu/study-bot/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB            *mongo.Database
	ctx           context.Context
	clientOptions *options.ClientOptions
)

func ConnectDB() {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
		return
	}
	clientOptions = options.Client()

	clientOptions = clientOptions.ApplyURI(cfg.MongoDB.URI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

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

func GetUser(userID string) (*User, error) {
	collection := DB.Collection("users")
	var user User

	filter := bson.M{"_id": userID}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return &user, nil
}
