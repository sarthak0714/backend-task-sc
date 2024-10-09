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

// Loads variables form Environment and returns Config struct
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
