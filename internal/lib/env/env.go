package env

import (
	"log"

	"github.com/joho/godotenv"
)

const path = "../../secrets/.env"

func LoadEnv() {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
