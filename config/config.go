package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/wavly/shawty/asserts"
	. "github.com/wavly/shawty/env"
	prettylogger "github.com/wavly/shawty/pretty-logger"
)

var logger = prettylogger.GetLogger(nil)

func Init() {
	err := godotenv.Load(".env.local")
	asserts.NoErr(err, "Failed reading .env.local")

	port := os.Getenv("PORT")
	asserts.AssertEq(port == "", "Missing PORT number in .env.local")
	PORT = port

	environment := os.Getenv("ENVIRONMENT")
	asserts.AssertEq(environment == "", "Missing ENVIRONMENT in .env.local")
	MODE = environment

	// Only get [TURSO_TOKEN] and [TURSO_URL] in Prodution
	if MODE == "prod" {
		tursoToken := os.Getenv("TURSO_AUTH_TOKEN")
		asserts.AssertEq(tursoToken == "", "Missing TURSO_AUTH_TOKEN in .evn.local")
		TURSO_TOKEN = tursoToken

		tursoURL := os.Getenv("TURSO_DATABASE_URL")
		asserts.AssertEq(tursoURL == "", "Missing TURSO_URL in .evn.local")
		TURSO_URL = tursoURL
		return
	}
}
