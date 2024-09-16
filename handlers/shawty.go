package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/wavly/shawty/asserts"
	"github.com/wavly/shawty/database"
	"github.com/wavly/shawty/utils"
)

type ShortLink struct {
	ShortUrl    string
	OriginalUrl string
}

func Shawty(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB()
	defer db.Close()

	longUrl := r.FormValue("url")

	// Check if longUrl contains "://" and add "https://" if missing
	if !strings.Contains(longUrl, "://") {
		longUrl = "https://" + longUrl
	}

	errorTempl := template.Must(template.ParseFiles("./partial-html/short-link-error.html"))

	// Parse the URL to validate it and check its scheme
	parsedUrl, err := url.Parse(longUrl)
	if err != nil || parsedUrl.Scheme != "https" {
		asserts.NoErr(errorTempl.Execute(w, "Invalid URL. Only HTTPS schema is allowed"), "Failed to execute template short-link-error.html")
		return
	}

	// Check if URL contains a TLD (Top-Level Domain)
	if !strings.Contains(longUrl, ".") {
		asserts.NoErr(errorTempl.Execute(w, "The URL doesn't contain TLD (Top-Level Domain)"), "Failed to execute template short-link-error.html")
		return
	} else if split := strings.SplitN(longUrl, ".", 2); split[1] == "" {
		asserts.NoErr(errorTempl.Execute(w, "The URL doesn't contain TLD (Top-Level Domain)"), "Failed to execute template short-link-error.html")
		return
	}

	// Check the lenght of the URL
	if len(longUrl) > 1000 {
		asserts.NoErr(errorTempl.Execute(w, "The URL is too long, Max URL lenght is 1000 characters"), "Failed to execute template short-link-error.html")
		return
	}

	hasher := sha256.New()
	hasher.Write([]byte(longUrl))
	checksum := hasher.Sum(nil)

	// Only get 8 characters long hash
	hashUrl := hex.EncodeToString(checksum[:4])

	templ := template.Must(template.ParseFiles("./partial-html/short-link.html"))

	// Check if the url exists in the database
	row := db.QueryRow("select code from urls where code = ?", hashUrl)
	var code string
	if err := row.Scan(&code); err != nil {
		// Check if err doesn't equal to `sql.ErrNoRows`
		// And if true then log the error and return
		if err != sql.ErrNoRows {
			utils.ServerErrTempl(w, "An error occur when querying the database")
			log.Printf("Database error when selecting original_url where code = %s, Error: %s\n", hashUrl, err)
			return
		}

		// Insert the URL in the database if doesn't exists
		row = db.QueryRow("insert into urls (original_url, code) values (?, ?)", longUrl, hashUrl)
		if err := row.Err(); err != nil {
			utils.ServerErrTempl(w, "An error occur when saving the URL to the database")
			log.Println("Failed to store URL in the database", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		data := ShortLink{
			ShortUrl:    hashUrl,
			OriginalUrl: longUrl,
		}
		asserts.NoErr(templ.Execute(w, data), "Failed to execute template short-link.html")
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := ShortLink{
		ShortUrl:    code,
		OriginalUrl: longUrl,
	}
	asserts.NoErr(templ.Execute(w, data), "Failed to execute template short-link.html")
}
