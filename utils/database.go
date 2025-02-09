package utils

import (
	"database/sql"
	"fmt"

	"github.com/wavly/surf/asserts"
	"github.com/wavly/surf/config"

	_ "github.com/tursodatabase/go-libsql"
)

func ConnectDB() *sql.DB {
	if config.MODE == "prod" {
		db, err := sql.Open("libsql", fmt.Sprintf("%s?authToken=%s", config.TURSO_URL, config.TURSO_TOKEN))
		asserts.NoErr(err, "Failed to connect to Turso remote db")
		return db
	}

	db, err := sql.Open("libsql", "file:./data.db")
	asserts.NoErr(err, "Failed to open sqlite file")
	return db
}
