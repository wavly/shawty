package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/wavly/surf/asserts"
	. "github.com/wavly/surf/env"
	prettylogger "github.com/wavly/surf/pretty-logger"
)

var logger = prettylogger.GetLogger(nil)

func Init() {
	err := godotenv.Load(".env")
	asserts.NoErr(err, "Failed reading .env")

	port := os.Getenv("PORT")
	asserts.AssertEq(port == "", "Missing PORT number in .env")
	PORT = port

	environment := os.Getenv("ENVIRONMENT")
	asserts.AssertEq(environment == "", "Missing ENVIRONMENT in .env")
	MODE = environment

	// Only get [TURSO_TOKEN] and [TURSO_URL] in Prodution
	if MODE == "prod" {
		tursoToken := os.Getenv("TURSO_AUTH_TOKEN")
		asserts.AssertEq(tursoToken == "", "Missing TURSO_AUTH_TOKEN in .env")
		TURSO_TOKEN = tursoToken

		tursoURL := os.Getenv("TURSO_DATABASE_URL")
		asserts.AssertEq(tursoURL == "", "Missing TURSO_URL in .env")
		TURSO_URL = tursoURL
		return
	}
}
