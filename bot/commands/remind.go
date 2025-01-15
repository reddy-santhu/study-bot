package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/reddy-santhu/study-bot/internal/study"
)

func HandleRemind(s *discordgo.Session, m *discordgo.MessageCreate) {
	parts := strings.SplitN(m.Content, " ", 3)
	if len(parts) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Invalid command. Usage: /remind <time> <task>")
		return
	}

	remindTime := parts[1]
	task := parts[2]

	err := study.ScheduleReminder(s, m, task, remindTime)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
