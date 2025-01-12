package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/reddy-santhu/study-bot/config"
)

var Discord *discordgo.Session

func StartBot(cfg *config.Config) {

	dg, err := discordgo.New("Bot " + cfg.Bot.Token)
	if err != nil {
		log.Fatalf("Error creating Bot Session %v", err)
	}

	Discord = dg

	dg.AddHandler(messageCreate)

	// dg.Identify.Intents = discordgo.IntentsGuildMessages

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening a websocket session %v", err)
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

	if m.Content == "Hi" {
		s.ChannelMessageSend(m.ChannelID, "Hello ")
	}

}
