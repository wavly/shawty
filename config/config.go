package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Working environment
var MODE string

// Server PORT number
var PORT string

// Access Token for Turso Database
var TURSO_TOKEN string

// Remote URL for Turso Database
var TURSO_URL string

func Init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Failed to load local env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalln("Missing PORT number in .env.local")
	}
	PORT = port

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		log.Fatalln("Missing ENVIRONMENT number in .env.local")
	}
	MODE = environment

	// Only get [TURSO_TOKEN] and [TURSO_URL] in Prodution
	if MODE == "prod" {
		tursoToken := os.Getenv("TURSO_AUTH_TOKEN")
		if tursoToken == "" {
			log.Fatalln("Missing TURSO_AUTH_TOKEN in .evn.local")
		}
		TURSO_TOKEN = tursoToken

		tursoURL := os.Getenv("TURSO_DATABASE_URL")
		if tursoURL == "" {
			log.Fatalln("Missing TURSO_URL in .evn.local")
		}
		TURSO_URL = tursoURL
		return
	}
}
