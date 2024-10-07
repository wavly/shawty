package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/wavly/shawty/asserts"
)

var PORT string
var ENV string

var TURSO_TOKEN string
var TURSO_URL string

func Init() {
	err := godotenv.Load(".env.local")
	asserts.NoErr(err, "Failed reading .env.local")

	port := os.Getenv("PORT")
	asserts.AssertEq(port == "", "Please specify the PORT number in .env.local")
	PORT = port

	environment := os.Getenv("ENVIRONMENT")
	asserts.AssertEq(environment == "", "Please specify the ENVIRONMENT in .env.local")
	ENV = environment

	// Only get [TURSO_TOKEN] and [TURSO_URL] in Prodution
	if ENV == "prod" {
		tursoToken := os.Getenv("TURSO_AUTH_TOKEN")
		asserts.AssertEq(tursoToken == "", "Please provide the TURSO TOKEN in .evn.local")
		TURSO_TOKEN = tursoToken

		tursoURL := os.Getenv("TURSO_DATABASE_URL")
		asserts.AssertEq(tursoURL == "", "Please provide the TURSO URL in .evn.local")
		TURSO_URL = tursoURL
	}
}
