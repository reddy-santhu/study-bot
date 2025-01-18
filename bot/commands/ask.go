package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/reddy-santhu/study-bot/ai"
	"github.com/reddy-santhu/study-bot/config"
)

func HandleAsk(s *discordgo.Session, m *discordgo.MessageCreate) {
	parts := strings.SplitN(m.Content, "question: ", 2)
	if len(parts) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Invalid command. Usage: /ask question: <question>")
		return
	}

	question := parts[1]

	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		s.ChannelMessageSend(m.ChannelID, "Sorry, there was an error loading the configuration.")
		return
	}

	apiKey := cfg.Gemini.APIKey
	response, err := ai.AskGemini(apiKey, question)
	if err != nil {
		fmt.Println(err.Error())
		s.ChannelMessageSend(m.ChannelID, "Sorry, I couldn't answer that question right now. Please try again later.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, response)
}
