package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading environment file")
	}

	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        ":" + os.Getenv("PORT"),
	}
}
