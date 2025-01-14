package commands

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/reddy-santhu/study-bot/internal/study"
)

func HandlePomo(s *discordgo.Session, m *discordgo.MessageCreate) {
	parts := strings.SplitN(m.Content, " ", 2)
	if len(parts) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Invalid command. Usage: /pomo start <work time> <break time> <task> or /pomo stop or /pomo status")
		return
	}

	subcommand := parts[1]

	switch {
	case strings.HasPrefix(subcommand, "start"):
		trimmedContent := strings.TrimPrefix(m.Content, "/pomo start ")
		stringed := strings.SplitN(trimmedContent, " ", 3)
		if len(stringed) < 3 {
			s.ChannelMessageSend(m.ChannelID, "Invalid command. Usage: /pomo start <work time> <break time> <task>")
			return
		}

		workTime, err := strconv.Atoi(stringed[0])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "invalid Worktime command. Usage: /pomo start <work time> <break time> <task>")
			return
		}
		breakTime, err := strconv.Atoi(stringed[1])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "invalid breakTime command. Usage: /pomo start <work time> <break time> <task>")
			return
		}
		task := stringed[2]

		study.StartPomodoro(s, m, workTime, breakTime, task)

	case strings.HasPrefix(subcommand, "stop"):
		study.StopPomodoro(s, m)
	case strings.HasPrefix(subcommand, "status"):
		study.GetPomodoroStatus(s, m)
	default:
		s.ChannelMessageSend(m.ChannelID, "Invalid command. Usage: /pomo start <work time> <break time> <task> or /pomo stop or /pomo status")
	}
}
