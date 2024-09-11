package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/wavly/shawty/database"
)

func Redirection(w http.ResponseWriter, r *http.Request) {
	// Truncate the route `/shawty/` from path
	slug := r.URL.Path[8:]

	db := database.ConnectDB()
	defer db.Close()

	var code string
	var originalUrl string

	// Get the url for the slug if exist
	row := db.QueryRow("SELECT code, original_url FROM urls WHERE code = ?", slug)
	result := row.Scan(&code, &originalUrl)

	if result == sql.ErrNoRows {
		w.Write([]byte("Couldn't find anything in your shawty, recheck..\n"))
	} else if result != nil {
		http.Error(w, "Sorry, an unexpected error occur while fetching full URL", http.StatusInternalServerError)
		log.Println("Failed to fetch original_url", result)
	} else {
		http.Redirect(w, r, originalUrl, http.StatusFound)

		// If url exists, update accessed count and last accessed time. after recirecting
		_, err := db.Exec("UPDATE urls SET accessed_count = accessed_count + 1, last_accessed = CURRENT_TIMESTAMP WHERE code = ?", slug)
		if err != nil {
			log.Println("Failed to update URL's accessed_count, last_accessed", err)
		}
	}
}
