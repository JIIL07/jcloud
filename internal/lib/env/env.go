package jenv

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	path := os.Getenv("ENV_PATH")
	if path == "" {
		log.Fatalf("ENV_PATH variable is not set")
	}
	err := godotenv.Load(path)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
