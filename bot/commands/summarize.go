package commands

import (
	"fmt"
	"strings"

	"github.com/reddy-santhu/study-bot/ai"
	"github.com/reddy-santhu/study-bot/config"

	"github.com/bwmarrin/discordgo"
)

func HandleSummarize(s *discordgo.Session, m *discordgo.MessageCreate) {
	parts := strings.SplitN(m.Content, "text: ", 2)
	if len(parts) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Invalid command. Usage: /summarize text: <long text>")
		return
	}

	longText := parts[1]

	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		s.ChannelMessageSend(m.ChannelID, "Sorry, there was an error loading the configuration.")
		return
	}

	apiKey := cfg.Gemini.APIKey
	prompt := fmt.Sprintf(
		"Summarize the following text in a concise and clear manner while preserving the key points. "+
			"Ensure that the summary is easy to understand and structured logically.\n\n"+
			"Output the summary in bullet points if the content is factual or technical. "+
			"If the text is narrative, provide a paragraph summary.\n\n"+
			"Text:\n%s", longText)

	response, err := ai.AskGemini(apiKey, prompt)
	if err != nil {
		fmt.Println(err.Error())
		s.ChannelMessageSend(m.ChannelID, "Sorry, I couldn't summarize that text right now. Please try again later.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, response)
}
