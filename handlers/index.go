package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/wavly/shawty/internal/database"
	prettylogger "github.com/wavly/shawty/pretty-logger"
	"github.com/wavly/shawty/utils"
	"github.com/wavly/shawty/validate"
)

var Logger = prettylogger.GetLogger(nil)

// TODO: write the name of the site in the response
func Main(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		Logger.Warn("POST request to unknown route", "route", r.URL.Path, "user-agent", r.UserAgent())
		w.Write([]byte("Not found"))
		return
	}
	inputUrl := r.FormValue("url")
	inputCode := r.FormValue("code")

	Logger.Info("POST request /", "input-url", inputUrl, "input-code", inputCode, "user-agent", r.UserAgent())

	// Check if inputUrl contains "://" and add "https://" if missing
	if !strings.Contains(inputUrl, "://") {
		inputUrl = "https://" + inputUrl
	}

	// Validate the URL
	err := validate.ValidateUrl(inputUrl)
	if err != nil {
		Logger.Warn("URL validation failed", "url", inputUrl, "user-agent", r.UserAgent(), "error", err, "input-url", inputUrl)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	db := utils.ConnectDB()
	defer db.Close()
	queries := database.New(db)

	if inputCode != "" {
		// Validate customCode
		err := validate.CustomCodeValidate(inputCode)
		if err != nil {
			Logger.Warn("Code validation failed", "code", inputUrl, "user-agent", r.UserAgent(), "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// Check if the url exists in the database
		_, err = queries.GetCode(r.Context(), inputCode)
		if err != nil {
			// Check if err doesn't equal to `sql.ErrNoRows`
			// And if true then log the error and return
			if err != sql.ErrNoRows {
				w.Write([]byte("An unexpected error occur when querying from the database"))
				Logger.Error("failed to query database to get the code", "error", err, "input-code", inputCode)
				return
			}

			// Insert the URL in the database if doesn't exists
			_, err = queries.CreateShortLink(r.Context(), database.CreateShortLinkParams{
				OriginalUrl: inputUrl,
				Code:        inputCode,
			})
			if err != nil {
				w.Write([]byte("An unexpected error occur when saving the URL to the database"))
				Logger.Error("failed to query database to store original url", "error", err, "original-url", inputUrl, "input-code", inputCode)
				return
			}

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(fmt.Sprintf("Location: /s/%s", inputCode)))
			return
		}

		w.WriteHeader(http.StatusFound)
		w.Write([]byte(fmt.Sprintf("Location: /s/%s", inputCode)))
		return
	}

	// Hashing the URL
	hasher := sha256.New()
	hasher.Write([]byte(inputUrl))
	checksum := hasher.Sum(nil)

	// Truncate to 8 characters long hash
	hashUrl := hex.EncodeToString(checksum[:4])

	// Check if the url exists in the database
	_, err = queries.GetCode(r.Context(), hashUrl)
	if err != nil {
		// Check if err doesn't equal to `sql.ErrNoRows`
		// And if true then log the error and return
		if err != sql.ErrNoRows {
			w.Write([]byte("An unexpected error occur when querying from the database"))
			Logger.Error("failed to query database to get the url code", "error", err, "input-code", inputCode)
			return
		}

		// Insert the URL in the database if doesn't exists
		_, err = queries.CreateShortLink(r.Context(), database.CreateShortLinkParams{
			OriginalUrl: inputUrl,
			Code:        hashUrl,
		})
		if err != nil {
			w.Write([]byte("An unexpected error occur when saving the URL to the database"))
			Logger.Error("failed to query to create the short url", "error", err, "original-url", inputUrl, "code", hashUrl)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("Location: /s/%s", hashUrl)))
		return
	}

	w.WriteHeader(http.StatusFound)
	w.Write([]byte(fmt.Sprintf("Location: /s/%s", hashUrl)))
}
