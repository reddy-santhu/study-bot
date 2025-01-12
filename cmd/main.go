package main

import (
	"fmt"

	"github.com/reddy-santhu/study-bot/config"
	"github.com/reddy-santhu/study-bot/utils"
)

func main() {
	utils.InitializeLogger()
	utils.InfoLogger.Println("Starting the Study Bot...")

	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		utils.ErrorLogger.Fatalf("Error loading configuration: %v", err)
	}

	fmt.Printf("Bot Prefix: %s\n", cfg.Bot.Prefix)
	fmt.Printf("MongoDB URI: %s\n", cfg.MongoDB.URI)
	fmt.Printf("Gemini API Key: %s\n", cfg.Gemini.APIKey)

	fmt.Println("Bot running")
	// Bot logic will go here later
}
