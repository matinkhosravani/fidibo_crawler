package app

import (
	"github.com/joho/godotenv"
	"log"
)

// LoadEnv use godot package to load/read the .env file and
// return the value of the key
func LoadEnv() {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
