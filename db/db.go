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

func LogStudyActivity(userID string, activity string) error {
	collection := DB.Collection("study_logs")

	studyLog := StudyLog{
		UserID:    userID,
		Timestamp: time.Now(),
		Activity:  activity,
		CreatedAt: time.Now(),
	}

	_, err := collection.InsertOne(context.TODO(), studyLog)
	if err != nil {
		return fmt.Errorf("failed to log study activity: %w", err)
	}
	usersCollection := DB.Collection("users")

	today := GetTodayDate()

	update := bson.M{"$set": bson.M{"last_study": today}, "$inc": bson.M{"total_days_studied": 1}}
	_, err = usersCollection.UpdateOne(context.TODO(), bson.M{"_id": userID, "last_study": bson.M{"$ne": today}}, update)
	if err != nil {
		return fmt.Errorf("error updating user last_study and total_days_studied: %w", err)
	}

	return nil
}

func GetTodayDate() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
}

func LogPDFData(userID string, filename string, text string) error {
	collection := DB.Collection("pdf_data")

	pdfData := PDFData{
		UserID:    userID,
		Filename:  filename,
		Text:      text,
		CreatedAt: time.Now(),
	}

	_, err := collection.InsertOne(context.TODO(), pdfData)
	if err != nil {
		return fmt.Errorf("failed to log PDF data: %w", err)
	}

	return nil
}

func GetPDFsByUser(userID string) ([]PDFData, error) {
	collection := DB.Collection("pdf_data")

	filter := bson.M{"user_id": userID}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("error finding PDFs: %w", err)
	}
	defer cursor.Close(context.TODO())
	var pdfs []PDFData
	for cursor.Next(context.TODO()) {
		var pdf PDFData
		err := cursor.Decode(&pdf)
		if err != nil {
			return nil, fmt.Errorf("error decoding PDF: %w", err)
		}
		pdfs = append(pdfs, pdf)
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}
	return pdfs, nil
}
