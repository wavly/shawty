package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/wavly/shawty/database"
)

func Redirection(w http.ResponseWriter, r *http.Request) {
	// Get the URL-Path slug "url"
	code := r.PathValue("code")

	db := database.ConnectDB()
	defer db.Close()

	// Get the url for the slug if exist
	row := db.QueryRow("select original_url from urls where code = ?", code)
	var originalUrl string
	if err := row.Scan(&originalUrl); err != nil {
		if err != sql.ErrNoRows {
			http.Error(w, "Sorry, an unexpected error occur when querying the database", http.StatusInternalServerError)
			log.Println("Failed to retrive original_url from the database:", err)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 not found, try again"))
		return
	}

	// Update the accessed_count
	_, err := db.Exec("update urls set accessed_count = accessed_count + 1 where code = ?", code)
	if err != nil {
		http.Error(w, "Sorry, an unexpected error occur when updating the access count from the database", http.StatusInternalServerError)
		log.Println("Failed to update the accessed_count from the database:", err)
		return
	}

	// Update the last_accessed
	_, err = db.Exec("update urls set last_accessed = ?", time.Now().UTC())
	if err != nil {
		http.Error(w, "Sorry, an unexpected error occur when updating the last_accessed count from the database", http.StatusInternalServerError)
		log.Println("Failed to update the last_accessed from the database:", err)
		return
	}

	// Redirect to original URL
	http.Redirect(w, r, originalUrl, http.StatusFound)
}
