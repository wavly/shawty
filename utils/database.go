package utils

import (
	"database/sql"
	"fmt"

	"github.com/wavly/surf/asserts"
	"github.com/wavly/surf/env"

	_ "github.com/tursodatabase/go-libsql"
)

func ConnectDB() *sql.DB {
	if env.MODE == "prod" {
		db, err := sql.Open("libsql", fmt.Sprintf("%s?authToken=%s", env.TURSO_URL, env.TURSO_TOKEN))
		asserts.NoErr(err, "Failed to connect to Turso remote db")
		return db
	}

	db, err := sql.Open("libsql", "file:./data.db")
	asserts.NoErr(err, "Failed to open sqlite file")
	return db
}
