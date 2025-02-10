package config

import (
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
		panic("Failed to load local env file: " + err.Error())
	}

	port := os.Getenv("PORT")
	if port == "" {
		panic("Missing PORT number in .env.local")
	}
	PORT = port

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		panic("Missing ENVIRONMENT number in .env")
	}
	MODE = environment

	// Only get [TURSO_TOKEN] and [TURSO_URL] in Prodution
	if MODE == "prod" {
		tursoToken := os.Getenv("TURSO_AUTH_TOKEN")
		if tursoToken == "" {
			panic("Missing TURSO_AUTH_TOKEN in .evn")
		}
		TURSO_TOKEN = tursoToken

		tursoURL := os.Getenv("TURSO_DATABASE_URL")
		if tursoURL == "" {
			panic("Missing TURSO_URL in .evn")
		}

		TURSO_URL = tursoURL
		return
	}
}
