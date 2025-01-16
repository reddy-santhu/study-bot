package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         string    `bson:"_id" json:"id"` // Discord User ID
	Username   string    `bson:"username" json:"username"`
	StudyGoals []string  `bson:"study_goals" json:"study_goals"`
	Streak     int       `bson:"streak" json:"streak"`
	LastStudy  time.Time `bson:"last_study" json:"last_study"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at" json:"updated_at"`
}

type Reminder struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    string             `bson:"user_id" json:"user_id"`
	Task      string             `bson:"task" json:"task"`
	Time      time.Time          `bson:"time" json:"time"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type StudyLog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    string             `bson:"user_id" json:"user_id"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
	Activity  string             `bson:"activity" json:"activity"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
