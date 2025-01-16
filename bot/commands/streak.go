package commands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/reddy-santhu/study-bot/db"
)

func HandleStreak(s *discordgo.Session, m *discordgo.MessageCreate) {
	userID := m.Author.ID

	user, err := db.GetUser(userID)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		s.ChannelMessageSend(m.ChannelID, "Failed to retrieve streak information. Please try again.")
		return
	}

	if user == nil {
		s.ChannelMessageSend(m.ChannelID, "You haven't started studying yet! Set a study goal to begin.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You have studied for a total of %d days!", user.TotalDaysStudied))
}
