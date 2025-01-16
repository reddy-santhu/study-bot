package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/reddy-santhu/study-bot/bot/commands"
	"github.com/reddy-santhu/study-bot/config"
)

var Discord *discordgo.Session

func StartBot(cfg *config.Config) {
	dg, err := discordgo.New("Bot " + cfg.Bot.Token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	Discord = dg

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages

	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "/study") {
		switch {
		case strings.HasPrefix(m.Content, "/study set"):
			commands.HandleStudySet(s, m)
		case strings.HasPrefix(m.Content, "/study list"):
			commands.HandleStudyList(s, m)
		case strings.HasPrefix(m.Content, "/study remove"):
			commands.HandleStudyRemove(s, m)
		}
		return
	}

	// Pomo command
	if strings.HasPrefix(m.Content, "/pomo") {
		commands.HandlePomo(s, m)
		return
	}

	// Remind command
	if strings.HasPrefix(m.Content, "/remind") {
		commands.HandleRemind(s, m)
		return
	}

	// streak command
	if strings.HasPrefix(m.Content, "/streak") {
		commands.HandleStreak(s, m)
		return
	}

	if m.Content == "ping!" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}
