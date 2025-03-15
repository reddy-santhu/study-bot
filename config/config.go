package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Bot struct {
		Token  string `yaml:"token"`
		Prefix string `yaml:"prefix"`
	} `yaml:"bot"`
	MongoDB struct {
		URI    string `yaml:"uri"`
		Dbname string `yaml:"dbname"`
	} `yaml:"mongodb"`
	Gemini struct {
		APIKey string `yaml:"api_key"`
	} `yaml:"gemini"`
}

var (
	cfg  *Config
	once sync.Once
)

func LoadConfig(filename string) (*Config, error) {
	var err error
	once.Do(func() {
		cfg, err = loadConfigFile(filename)
		if err != nil {
			log.Println("Failed to load config file, using environment variables only")
			cfg = &Config{}
		}
		applyEnvOverrides(cfg)
	})
	return cfg, err
}

func loadConfigFile(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %v", err)
	}
	defer file.Close()

	yamlFile, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	config := &Config{}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML: %v", err)
	}

	return config, nil
}

func applyEnvOverrides(cfg *Config) {
	_ = godotenv.Load() // Load from .env file if present

	if token := os.Getenv("BOT_TOKEN"); token != "" {
		cfg.Bot.Token = token
	}
	if uri := os.Getenv("MONGODB_URI"); uri != "" {
		cfg.MongoDB.URI = uri
	}
	if apiKey := os.Getenv("GEMINI_API_KEY"); apiKey != "" {
		cfg.Gemini.APIKey = apiKey
	}
}
