package commands

import (
	"context"
	"fmt"
	"log"
	"strconv"
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

		err = db.LogStudyActivity(userID, "set_goal")
		if err != nil {
			log.Printf("Error logging study activity: %v", err)
		}

		return nil
	}

	update := bson.M{"$push": bson.M{"study_goals": goal}}
	_, err = usersCollection.UpdateOne(ctx, bson.M{"_id": userID}, update)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	err = db.LogStudyActivity(userID, "set_goal")
	if err != nil {
		log.Printf("Error logging study activity: %v", err)
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

// list

func HandleStudyList(s *discordgo.Session, m *discordgo.MessageCreate) {
	userID := m.Author.ID

	user, err := GetUser(userID)
	if err != nil {
		log.Printf("Error getting study goals for user %s: %v", userID, err)
		s.ChannelMessageSend(m.ChannelID, "Failed to retrieve study goals. Please try again.")
		return
	}

	if user == nil || len(user.StudyGoals) == 0 {
		s.ChannelMessageSend(m.ChannelID, "You have no study goals set.")
		return
	}

	var goalsList string
	for i, goal := range user.StudyGoals {
		goalsList += fmt.Sprintf("%d. %s\n", i+1, goal)
	}

	s.ChannelMessageSend(m.ChannelID, "Your study goals:\n"+goalsList)
}

// remove
func HandleStudyRemove(s *discordgo.Session, m *discordgo.MessageCreate) {
	parts := strings.SplitN(m.Content, " ", 3)
	if len(parts) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Invalid command. Usage: /study remove <goal number>")
		return
	}

	userID := m.Author.ID

	goalNumberStr := parts[2]
	goalNumber, err := strconv.Atoi(goalNumberStr)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Invalid goal number. Please enter a valid number.")
		return
	}

	err = RemoveStudyGoal(userID, goalNumber)
	if err != nil {
		log.Printf("Error removing study goal for user %s: %v", userID, err)
		s.ChannelMessageSend(m.ChannelID, "Failed to remove study goal. Please try again.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Study goal removed successfully.")
}

func RemoveStudyGoal(userID string, goalNumber int) error {
	usersCollection := db.GetCollection("users")

	user, err := GetUser(userID)
	if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}

	if user == nil || len(user.StudyGoals) == 0 {
		return fmt.Errorf("no study goals to remove")
	}

	if goalNumber < 1 || goalNumber > len(user.StudyGoals) {
		return fmt.Errorf("invalid goal number")
	}

	goalIndex := goalNumber - 1
	goalToRemove := user.StudyGoals[goalIndex]

	user.StudyGoals = append(user.StudyGoals[:goalIndex], user.StudyGoals[goalIndex+1:]...)

	update := bson.M{"$set": bson.M{"study_goals": user.StudyGoals}}
	_, err = usersCollection.UpdateOne(ctx, bson.M{"_id": userID}, update)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	log.Printf("Removed goal '%s' for user %s", goalToRemove, userID)

	return nil
}
