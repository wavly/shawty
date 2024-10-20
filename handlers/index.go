package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/wavly/shawty/internal/database"
	prettylogger "github.com/wavly/shawty/pretty-logger"
	"github.com/wavly/shawty/utils"
	"github.com/wavly/shawty/validate"
)

var Logger = prettylogger.GetLogger(nil)

func Main(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		Logger.Warn("POST request to unknown route", "route", r.URL.Path, "user-agent", r.UserAgent())
		w.Write([]byte("Not found"))
		return
	}
	inputUrl := r.FormValue("url")
	customCode := r.FormValue("code")

	Logger.Info("POST / request", "input-url", inputUrl, "input-code", customCode, "user-agent", r.UserAgent())

	// Validate customCode
	err := validate.CustomCodeValidate(customCode)
	if err != nil {
		Logger.Warn("Code validation failed", "code", inputUrl, "from-ip", r.RemoteAddr, "user-agent", r.UserAgent())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// Check if inputUrl contains "://" and add "https://" if missing
	if !strings.Contains(inputUrl, "://") {
		inputUrl = "https://" + inputUrl
	}

	// Validate the URL
	err = validate.ValidateUrl(inputUrl)
	if err != nil {
		Logger.Warn("URL validation failed", "url", inputUrl, "from-ip", r.RemoteAddr, "user-agent", r.UserAgent())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	db := utils.ConnectDB()
	defer db.Close()
	queries := database.New(db)

	if customCode != "" {
		// Check if the url exists in the database
		_, err := queries.GetCode(r.Context(), customCode)
		if err != nil {
			// Check if err doesn't equal to `sql.ErrNoRows`
			// And if true then log the error and return
			if err != sql.ErrNoRows {
				w.Write([]byte("An unexpected error occur when querying from the database"))
				log.Printf("Database error when selecting code where code = %s, Error: %s\n", customCode, err)
				return
			}

			// Insert the URL in the database if doesn't exists
			_, err = queries.CreateShortLink(r.Context(), database.CreateShortLinkParams{
				OriginalUrl: inputUrl,
				Code:        customCode,
			})
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
	_, err = queries.GetCode(r.Context(), customCode)
	if err != nil {
		// Check if err doesn't equal to `sql.ErrNoRows`
		// And if true then log the error and return
		if err != sql.ErrNoRows {
			w.Write([]byte("An unexpected error occur when querying from the database"))
			log.Printf("Database error when selecting original_url where code = %s, Error: %s\n", hashUrl, err)
			return
		}

		// Insert the URL in the database if doesn't exists
		_, err = queries.CreateShortLink(r.Context(), database.CreateShortLinkParams{
			OriginalUrl: inputUrl,
			Code:        hashUrl,
		})

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
