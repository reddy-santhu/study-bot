package config

import (
	"fmt"
	"io"
	"log"
	"os"

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
	}
}

func LoadConfig(filename string) (*Config, error) {

	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return nil, err
	}
	defer file.Close()

	yamlFile, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML: %v", err)
	}
	return config, nil

}
