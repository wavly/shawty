package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/wavly/shawty/database"
)

func Main(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB()
	defer db.Close()

	inputUrl := r.FormValue("url")

	// Check if inputUrl contains valid schema "://" and if not then added it manually
	if !strings.Contains(inputUrl, "://") {
		inputUrl = "https://" + inputUrl
	} else if !strings.Contains(inputUrl, "http") { // Check if inputUrl schema is http(s)
		w.Write([]byte("Only http or https schema is allowed"))
		return
	}

	// Check if URL is valid by checking if it contains TLD (Top-Level Domain)
	if !strings.Contains(inputUrl, ".") {
		w.Write([]byte("Enter a valid URL"))
		return
	}

	// Hashing the URL
	hasher := sha256.New()
	hasher.Write([]byte(inputUrl))
	checksum := hasher.Sum(nil)

	// Truncate to 8 characters long hash
	hashUrl := hex.EncodeToString(checksum[:4])

	// Check if the url exists in the database
	row := db.QueryRow("select code from urls where code = ?", hashUrl)
	if err := row.Err(); err != nil {
		// Check if err doesn't equal to `sql.ErrNoRows`
		// And if true then log the error and return
		if err != sql.ErrNoRows {
			http.Error(w, "Sorry, an unexpected error occur when querying from the database", http.StatusInternalServerError)
      log.Printf("Database error when selecting original_url where code = %s, Error: %s\n", hashUrl, err)
			return
		}

		// Insert the URL in the database if doesn't exists
		row = db.QueryRow("insert into urls (original_url, code) values (?, ?)", inputUrl, hashUrl)
		if err := row.Err(); err != nil {
			http.Error(w, "Sorry, an unexpected error occur when saving the URL", http.StatusInternalServerError)
			log.Println("Failed to store URL in the database", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("Location: %s\n", inputUrl)))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Location: %s\n", inputUrl)))
}
