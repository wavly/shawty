package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"net/http"

	"github.com/wavly/surf/asserts"
	"github.com/wavly/surf/internal/database"
	partialhtml "github.com/wavly/surf/partial-html"
	"github.com/wavly/surf/static"
	"github.com/wavly/surf/utils"
	"github.com/wavly/surf/validate"
)

type ShortLink struct {
	ShortUrl string
}

func Short(w http.ResponseWriter, r *http.Request) {
	inputUrl := r.FormValue("url")
	Logger.Info("Shorten the URL", "url", inputUrl, "user-agent", r.UserAgent())

	// Validate the URL
	parsedUrl, err := validate.ValidateUrl(inputUrl)
	if err != nil {
		Logger.Warn("failed to validate URL", "url", parsedUrl, "user-agent", r.UserAgent(), "error", err)
		err = partialhtml.ShortLinkError(err.Error()).Render(r.Context(), w)
		asserts.NoErr(err, "Failed to render partial-html short-link-error")
		return
	}

	hasher := sha256.New()
	hasher.Write([]byte(parsedUrl))
	checksum := hasher.Sum(nil)

	db := utils.ConnectDB()
	defer db.Close()
	queries := database.New(db)

	// Only get 8 characters long hash
	hashUrl := hex.EncodeToString(checksum[:4])

	// Check if the url exists in the database
	code, err := queries.GetCode(r.Context(), hashUrl)
	if err != nil {
		Logger.Info("URL doesn't exists in the database", "url", parsedUrl, "user-agent", r.UserAgent())
		// Check if err doesn't equal to `sql.ErrNoRows`
		// And if true then log the error and return
		if err != sql.ErrNoRows {
			Logger.Error("failed to query the code for the URL", "error", err, "code", hashUrl, "input-url", parsedUrl, "user-agent", r.UserAgent())
			err := static.ServerError("An error occur when querying the database").Render(r.Context(), w)
			asserts.NoErr(err, "Failed to render server-internal-error page")
			return
		}

		// Insert the URL in the database if doesn't exists
		_, err = queries.CreateShortLink(r.Context(), database.CreateShortLinkParams{
			OriginalUrl: parsedUrl,
			Code:        hashUrl,
		})
		if err != nil {
			Logger.Error("failed to query to create short link", "original_url", parsedUrl, "code", hashUrl, "error", err)
			err := static.ServerError("An error occur when saving the URL to the database").Render(r.Context(), w)
			asserts.NoErr(err, "Failed to render server-internal-error page")
			return
		}

		w.WriteHeader(http.StatusCreated)
		asserts.NoErr(partialhtml.ShortLink(hashUrl).Render(r.Context(), w), "Failed to render partial-html short-link")
		return
	}

	Logger.Info("URL exists in the database", "url", parsedUrl, "code", hashUrl, "user-agent", r.UserAgent())
	w.WriteHeader(http.StatusCreated)
	asserts.NoErr(partialhtml.ShortLink(code).Render(r.Context(), w), "Failed to redner partial-html short-link")
}
