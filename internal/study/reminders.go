package study

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron"
)

var scheduler *gocron.Scheduler
var schedulerInit sync.Once

func init() {
	schedulerInit.Do(func() {
		location, err := time.LoadLocation("Asia/Kolkata")
		if err != nil {
			panic(err)
		}
		scheduler = gocron.NewScheduler(location)
		scheduler.StartAsync()
	})
}

func ScheduleReminder(s *discordgo.Session, m *discordgo.MessageCreate, task string, remindTime string) error {
	userID := m.Author.ID
	channelID := m.ChannelID

	parts := strings.Split(remindTime, ":")
	if len(parts) != 2 {
		s.ChannelMessageSend(channelID, "Invalid time format. Use HH:MM (e.g., 14:30)")
		return fmt.Errorf("invalid time format")
	}

	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		s.ChannelMessageSend(channelID, "Invalid hour. Use a number between 0 and 23.")
		return fmt.Errorf("invalid hour")
	}
	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		s.ChannelMessageSend(channelID, "Invalid minute. Use a number between 0 and 59.")
		return fmt.Errorf("invalid minute")
	}

	_, err = scheduler.Every(1).Day().At(remindTime).Do(func() {
		s.ChannelMessageSend(channelID, fmt.Sprintf("<@%s>, Reminder: %s", userID, task))
	})

	if err != nil {
		log.Printf("Error scheduling reminder for user %s: %v", userID, err)
		s.ChannelMessageSend(channelID, "Failed to schedule reminder. Please try again.")
		return err
	}

	log.Printf("Scheduled reminder for user %s at %02d:%02d to: %s", userID, hour, minute, task)
	s.ChannelMessageSend(channelID, fmt.Sprintf("Reminder scheduled for %02d:%02d to: %s", hour, minute, task))

	return nil
}

func ClearSchedule() {
	scheduler.Clear()
}
