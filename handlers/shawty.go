package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/wavly/shawty/asserts"
	"github.com/wavly/shawty/internal/database"
	"github.com/wavly/shawty/utils"
	"github.com/wavly/shawty/validate"
)

type ShortLink struct {
	ShortUrl    string
}

func Shawty(w http.ResponseWriter, r *http.Request) {
	inputUrl := r.FormValue("url")

	// Check if longUrl contains "://" and add "https://" if missing
	if !strings.Contains(inputUrl, "://") {
		inputUrl = "https://" + inputUrl
	}

	// Validate the URL
	err := validate.ValidateUrl(inputUrl)
	if err != nil {
		errorTempl := template.Must(template.ParseFiles("./partial-html/short-link-error.html"))
		asserts.NoErr(errorTempl.Execute(w, err), "Failed to execute template short-link-error.html")
		return
	}

	hasher := sha256.New()
	hasher.Write([]byte(inputUrl))
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
		// Check if err doesn't equal to `sql.ErrNoRows`
		// And if true then log the error and return
		if err != sql.ErrNoRows {
			utils.ServerErrTempl(w, "An error occur when querying the database")
			log.Printf("Database error when selecting original_url where code = %s, Error: %s\n", hashUrl, err)
			return
		}

		// Insert the URL in the database if doesn't exists
		_, err = queries.CreateShortLink(r.Context(), database.CreateShortLinkParams{
			OriginalUrl: inputUrl,
			Code:        hashUrl,
		})
		if err != nil {
			utils.ServerErrTempl(w, "An error occur when saving the URL to the database")
			log.Println("Failed to store URL in the database", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		data := ShortLink{
			ShortUrl:    hashUrl,
		}
		asserts.NoErr(templ.Execute(w, data), "Failed to execute template short-link.html")
		return
	}

	w.WriteHeader(http.StatusCreated)
	data := ShortLink{
		ShortUrl:    code,
	}
	asserts.NoErr(templ.Execute(w, data), "Failed to execute template short-link.html")
}
