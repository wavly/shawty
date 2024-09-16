package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/wavly/shawty/database"
)

func Main(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB()
	defer db.Close()

	inputUrl := r.FormValue("url")
	customCode := r.FormValue("code")

	// Check the lenght of the customCode
	if len(customCode) > 8 {
		w.Write([]byte("Max lenght of the custom code is 8"))
		return
	}

	// Check if inputUrl contains "://" and add "https://" if missing
	if !strings.Contains(inputUrl, "://") {
		inputUrl = "https://" + inputUrl
	}

	// Parse the URL to validate it and check its scheme
	parsedUrl, err := url.Parse(inputUrl)
	if err != nil || parsedUrl.Scheme != "https" {
		w.Write([]byte("Invalid URL. Only HTTPS schema is allowed"))
		return
	}

	// Check if URL contains a TLD (Top-Level Domain)
	if !strings.Contains(inputUrl, ".") {
		w.Write([]byte("The URL doesn't contain TLD (Top-Level Domain)"))
		return
	} else if split := strings.SplitN(inputUrl, ".", 2); split[1] == "" {
		w.Write([]byte("The URL doesn't contain TLD (Top-Level Domain)"))
		return
	}

	if customCode != "" {
		// Check if the url exists in the database
		err := db.QueryRow("select code from urls where code = ?", customCode).Err()
		if err != nil {
			// Check if err doesn't equal to `sql.ErrNoRows`
			// And if true then log the error and return
			if err != sql.ErrNoRows {
				w.Write([]byte("An unexpected error occur when querying from the database"))
				log.Printf("Database error when selecting original_url where code = %s, Error: %s\n", customCode, err)
				return
			}

			// Insert the URL in the database if doesn't exists
			err = db.QueryRow("insert into urls (original_url, code) values (?, ?)", inputUrl, customCode).Err()
			if err != nil {
				w.Write([]byte("An unexpected error occur when saving the URL to the database"))
				log.Println("Failed to store URL in the database", err)
				return
			}

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(fmt.Sprintf("Location: /s/%s", customCode)))
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("Location: /s/%s", customCode)))
		return
	}

	// Hashing the URL
	hasher := sha256.New()
	hasher.Write([]byte(inputUrl))
	checksum := hasher.Sum(nil)

	// Truncate to 8 characters long hash
	hashUrl := hex.EncodeToString(checksum[:4])

	// Check if the url exists in the database
	err = db.QueryRow("select code from urls where code = ?", hashUrl).Err()
	if err != nil {
		// Check if err doesn't equal to `sql.ErrNoRows`
		// And if true then log the error and return
		if err != sql.ErrNoRows {
			w.Write([]byte("An unexpected error occur when querying from the database"))
			log.Printf("Database error when selecting original_url where code = %s, Error: %s\n", hashUrl, err)
			return
		}

		// Insert the URL in the database if doesn't exists
		err = db.QueryRow("insert into urls (original_url, code) values (?, ?)", inputUrl, hashUrl).Err()
		if err != nil {
			w.Write([]byte("An unexpected error occur when saving the URL to the database"))
			log.Println("Failed to store URL in the database", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("Location: /s/%s", hashUrl)))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Location: /s/%s", hashUrl)))
}
