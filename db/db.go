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
		log.Fatalf("Error loading the configuration : %v", err)
		return
	}

	clientOptions := options.Client()
	clientOptions = clientOptions.ApplyURI(cfg.MongoDB.URI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	DB = client.Database(cfg.MongoDB.DBname)

}
func GetCollection(collectionName string) *mongo.Collection {
	return DB.Collection(collectionName)
}
