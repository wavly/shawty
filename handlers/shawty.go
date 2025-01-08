package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"html/template"
	"net/http"

	"github.com/wavly/surf/asserts"
	"github.com/wavly/surf/internal/database"
	"github.com/wavly/surf/utils"
	"github.com/wavly/surf/validate"
)

type ShortLink struct {
	ShortUrl string
}

func Surf(w http.ResponseWriter, r *http.Request) {
	inputUrl := r.FormValue("url")
	Logger.Info("Shorten the URL", "url", inputUrl, "user-agent", r.UserAgent())

	// Validate the URL
	parsedUrl, err := validate.ValidateUrl(inputUrl)
	if err != nil {
		Logger.Warn("failed to validate URL", "url", parsedUrl, "user-agent", r.UserAgent(), "error", err)
		errorTempl := template.Must(template.ParseFiles("./partial-html/short-link-error.html"))
		asserts.NoErr(errorTempl.Execute(w, err), "Failed to execute template short-link-error.html")
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

	templ := template.Must(template.ParseFiles("./partial-html/short-link.html"))

	// Check if the url exists in the database
	code, err := queries.GetCode(r.Context(), hashUrl)
	if err != nil {
		Logger.Info("URL doesn't exists in the database", "url", parsedUrl, "user-agent", r.UserAgent())
		// Check if err doesn't equal to `sql.ErrNoRows`
		// And if true then log the error and return
		if err != sql.ErrNoRows {
			Logger.Error("failed to query the code for the URL", "error", err, "code", hashUrl, "input-url", parsedUrl, "user-agent", r.UserAgent())
			utils.ServerErrTempl(w, "An error occur when querying the database")
			return
		}

		// Insert the URL in the database if doesn't exists
		_, err = queries.CreateShortLink(r.Context(), database.CreateShortLinkParams{
			OriginalUrl: parsedUrl,
			Code:        hashUrl,
		})
		if err != nil {
			Logger.Error("failed to query to create short link", "original_url", parsedUrl, "code", hashUrl, "error", err)
			utils.ServerErrTempl(w, "An error occur when saving the URL to the database")
			return
		}

		w.WriteHeader(http.StatusCreated)
		data := ShortLink{
			ShortUrl: hashUrl,
		}
		asserts.NoErr(templ.Execute(w, data), "Failed to execute template short-link.html")
		return
	}

	Logger.Info("URL exists in the database", "url", parsedUrl, "code", hashUrl, "user-agent", r.UserAgent())
	w.WriteHeader(http.StatusCreated)
	data := ShortLink{
		ShortUrl: code,
	}
	asserts.NoErr(templ.Execute(w, data), "Failed to execute template short-link.html")
}
