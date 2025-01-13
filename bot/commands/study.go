package commands

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/reddy-santhu/study-bot/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var ctx = context.TODO()

func HandleStudySet(s *discordgo.Session, m *discordgo.MessageCreate) {
	parts := strings.SplitN(m.Content, " ", 3)
	if len(parts) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Invalid command. Usage: /study set <goal>")
		return
	}

	goal := parts[2]

	userID := m.Author.ID

	err := SetStudyGoal(userID, goal, m.Author.Username)
	if err != nil {
		log.Printf("Error setting study goal for user %s: %v", userID, err)
		s.ChannelMessageSend(m.ChannelID, "Failed to set study goal. Please try again.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Study goal set to: %s", goal))
}

func SetStudyGoal(userID string, goal string, username string) error {
	usersCollection := db.GetCollection("users")

	user, err := GetUser(userID)
	if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}

	if user == nil {
		newUser := db.User{
			ID:         userID,
			Username:   username,
			StudyGoals: []string{goal},
			Streak:     0,
			LastStudy:  time.Now(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		_, err = usersCollection.InsertOne(ctx, newUser)
		if err != nil {
			return fmt.Errorf("error inserting new user: %w", err)
		}
		return nil
	}

	update := bson.M{"$push": bson.M{"study_goals": goal}}
	_, err = usersCollection.UpdateOne(ctx, bson.M{"_id": userID}, update)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}

func GetUser(userID string) (*db.User, error) {
	usersCollection := db.GetCollection("users")
	var user *db.User

	err := usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	return user, nil
}
