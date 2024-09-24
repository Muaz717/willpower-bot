package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Username string `yaml:"username"`
	DBPort   string `yaml:"port"`
	Host     string `yaml:"host"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load env file: %s", err)
	}

	cfgPath := os.Getenv("CONFIG_PATH")

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", cfgPath)
	}

	var cfg Config

	err = cleanenv.ReadConfig(cfgPath, &cfg)
	if err != nil {
		log.Fatalf("failed to read config: %s", err)
	}

	return &cfg
}
