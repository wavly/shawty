package utils

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/wavly/shawty/asserts"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func ConnectDB() *sql.DB {
	// Loading the environment variables
	err := godotenv.Load(".env.local")
	asserts.NoErr(err, "Failed to load environment variables")

	tursoURL := os.Getenv("TURSO_DATABASE_URL")
	tursoToken := os.Getenv("TURSO_AUTH_TOKEN")

	db, err := sql.Open("libsql", fmt.Sprintf("%s?authToken=%s", tursoURL, tursoToken))
	asserts.NoErr(err, "Failed to connect to Turso remote db")
	return db
}
