package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/wavly/shawty/asserts"
	"github.com/wavly/shawty/database"
)

func Redirect(w http.ResponseWriter, r *http.Request) {
	// Get the URL-Path slug "url"
	code := r.PathValue("code")

	if len(code) > 8 {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	// MemcacheD Client
	mcClient := memcache.New("0.0.0.0:11211")
	asserts.NoErr(mcClient.Ping(), "Failed to ping MemcacheD")

	db := database.ConnectDB()
	defer db.Close()

	var originalUrl string
	cache, err := mcClient.Get(code)
	if err != nil {
		if err != memcache.ErrCacheMiss {
			log.Println("Memcache error:", err)
		}

		row := db.QueryRow("select original_url from urls where code = ?", code)
		if err := row.Scan(&originalUrl); err != nil {
			if err != sql.ErrNoRows {
				http.Error(w, "Sorry, an unexpected error occur when querying the database", http.StatusInternalServerError)
				log.Println("Failed to retrive original_url from the database:", err)
				return
			}

			http.Redirect(w, r, "/", http.StatusBadRequest)
			return
		}

		log.Println("Cache Miss:", originalUrl)
		err = mcClient.Set(&memcache.Item{Key: code, Value: []byte(originalUrl)})
		if err != nil {
			log.Println("Memcache error when setting the key:", err)
		}
	} else if cache != nil && cache.Value != nil { // Check if original URL is in the cache
		cacheOriginalUrl := string(cache.Value)
		log.Println("Cashe Hit", cacheOriginalUrl)

		// Redirect to original URL
		http.Redirect(w, r, cacheOriginalUrl, http.StatusFound)

		// Update the accessed_count and last_accessed in one query
		_, err = db.Exec("update urls set accessed_count = accessed_count + 1, last_accessed = ? where code = ?", time.Now().UTC(), code)
		if err != nil {
			log.Println("Failed to update accessed_count and last_accessed:", err)
			return
		}
		return
	}

	// Redirect to original URL
	http.Redirect(w, r, originalUrl, http.StatusFound)

	// Update the accessed_count and last_accessed in one query
	_, err = db.Exec("update urls set accessed_count = accessed_count + 1, last_accessed = ? where code = ?", time.Now().UTC(), code)
	if err != nil {
		log.Println("Failed to update accessed_count and last_accessed:", err)
		return
	}
}
