package handlers

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/mergestat/timediff"
	"github.com/wavly/shawty/asserts"
	"github.com/wavly/shawty/database"
)

type AccessCount struct {
	ShortLink
	Count int

	LastAccessed string
}

func Stats(w http.ResponseWriter, r *http.Request) {
	inputCode := r.PathValue("code")

	// Don't query database if input URL lenght is less than 8
	if len(inputCode) < 8 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else if len(inputCode) > 8 { // Else check if input URL lenght is greater than 8
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	templ := template.Must(template.ParseFiles("./static/stat.html"))
	db := database.ConnectDB()
	defer db.Close()

	rows, err := db.Query("select accessed_count, original_url, last_accessed from urls where code = ?", inputCode)
	if err != nil {
		http.Error(w, "Sorry, an unexpected error occur", http.StatusInternalServerError)
		log.Printf("Database error when selecting accessed_count and original_url where code = %s, Error %s\n", inputCode, err)
		return
	}
	defer rows.Close()

	var accessedCount int
	var originalUrl string
	var lastAccessed time.Time

	// Redirect if no result is found
	if !rows.Next() {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Scan the result
	err = rows.Scan(&accessedCount, &originalUrl, &lastAccessed)
	if err != nil {
		http.Error(w, "Sorry, an unexpected error occur", http.StatusInternalServerError)
		log.Printf("Error scanning the result: %s", err)
		return
	}

	data := AccessCount{
		Count: accessedCount,
		LastAccessed: timediff.TimeDiff(lastAccessed),
		ShortLink: ShortLink{
			ShortUrl:    inputCode,
			OriginalUrl: originalUrl,
		},
	}

	asserts.NoErr(templ.Execute(w, data), "Failed to execute template stat.html")
}
