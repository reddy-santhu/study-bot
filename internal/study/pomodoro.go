package study

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/reddy-santhu/study-bot/db"
)

type PomodoroSession struct {
	UserID      string
	ChannelID   string
	WorkTime    int
	BreakTime   int
	CurrentTask string
	Timer       *time.Timer
	Stage       string
	StartTime   time.Time
	EndTime     time.Time
}

var activePomodoros = make(map[string]*PomodoroSession)
var mu sync.Mutex

func StartPomodoro(s *discordgo.Session, m *discordgo.MessageCreate, workTime int, breakTime int, task string) {
	userID := m.Author.ID
	channelID := m.ChannelID

	mu.Lock()
	defer mu.Unlock()

	if _, ok := activePomodoros[userID]; ok {
		s.ChannelMessageSend(channelID, "A Pomodoro timer is already running for you. Please stop it first.")
		return
	}

	session := &PomodoroSession{
		UserID:      userID,
		ChannelID:   channelID,
		WorkTime:    workTime,
		BreakTime:   breakTime,
		CurrentTask: task,
		Stage:       "work",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(time.Duration(workTime) * time.Minute),
	}

	session.Timer = time.NewTimer(time.Duration(workTime) * time.Minute)
	activePomodoros[userID] = session

	err := db.LogStudyActivity(userID, "start_pomodoro")
	if err != nil {
		log.Printf("Error logging study activity: %v", err)
	}
	s.ChannelMessageSend(channelID, fmt.Sprintf("Pomodoro timer started! Work for %d minutes, then break for %d minutes. Task: %s", workTime, breakTime, task))

	go runPomodoro(s, session)
}

func StopPomodoro(s *discordgo.Session, m *discordgo.MessageCreate) {
	userID := m.Author.ID
	channelID := m.ChannelID

	mu.Lock()
	defer mu.Unlock()

	session, ok := activePomodoros[userID]
	if !ok {
		s.ChannelMessageSend(channelID, "No Pomodoro timer is running for you.")
		return
	}

	session.Timer.Stop()
	delete(activePomodoros, userID)

	s.ChannelMessageSend(channelID, "Pomodoro timer stopped.")
}

func runPomodoro(s *discordgo.Session, session *PomodoroSession) {
	userID := session.UserID
	channelID := session.ChannelID

	defer func() {
		mu.Lock()
		defer mu.Unlock()
		delete(activePomodoros, userID)
	}()

	for {
		<-session.Timer.C

		mu.Lock()
		if session.Stage == "work" {

			session.Stage = "break"
			s.ChannelMessageSend(channelID, fmt.Sprintf("Time for a %d-minute break!", session.BreakTime))
			session.Timer = time.NewTimer(time.Duration(session.BreakTime) * time.Minute)
			session.EndTime = time.Now().Add(time.Duration(session.BreakTime) * time.Minute)

		} else {

			session.Stage = "work"
			s.ChannelMessageSend(channelID, fmt.Sprintf("Break over! Time to work for %d minutes on %s!", session.WorkTime, session.CurrentTask))
			session.Timer = time.NewTimer(time.Duration(session.WorkTime) * time.Minute)
			session.EndTime = time.Now().Add(time.Duration(session.WorkTime) * time.Minute)
		}
		mu.Unlock()

		if _, ok := activePomodoros[userID]; !ok {
			log.Printf("Pomodoro timer stopped for user %s", userID)
			return
		}

		log.Printf("Pomodoro timer running for user %s stage %s", userID, session.Stage)
	}
}

func ParsePomodoroArgs(s *discordgo.Session, m *discordgo.MessageCreate, content string) (int, int, string, error) {
	parts := strings.SplitN(content, " ", 4)
	if len(parts) < 4 {
		return 0, 0, "", fmt.Errorf("invalid command. Usage: /pomo start <work time> <break time> <task>")
	}

	workTimeStr := parts[1]
	breakTimeStr := parts[2]
	task := parts[3]

	workTime, err := strconv.Atoi(workTimeStr)
	if err != nil {
		return 0, 0, "", fmt.Errorf("invalid work time. Please enter a valid number")
	}

	breakTime, err := strconv.Atoi(breakTimeStr)
	if err != nil {
		return 0, 0, "", fmt.Errorf("invalid break time. Please enter a valid number")
	}

	return workTime, breakTime, task, nil
}

func GetPomodoroStatus(s *discordgo.Session, m *discordgo.MessageCreate) {
	userID := m.Author.ID
	channelID := m.ChannelID

	mu.Lock()
	defer mu.Unlock()

	session, ok := activePomodoros[userID]
	if !ok {
		s.ChannelMessageSend(channelID, "No Pomodoro timer is running for you.")
		return
	}
	now := time.Now()
	remaining := session.EndTime.Sub(now)

	minutes := int(remaining.Minutes())
	seconds := int(remaining.Seconds()) % 60

	s.ChannelMessageSend(channelID, fmt.Sprintf("Remaining time in the %s session: %dm %ds", session.Stage, minutes, seconds))
}
