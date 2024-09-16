package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	sqlc "github.com/wavly/shawty/sqlc_db"
	"github.com/wavly/shawty/utils"
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

	// Open a connection to the database
	db := utils.ConnectDB()
	defer db.Close()
	queries := sqlc.New(db)

	var originalUrl string
	cache, err := mcClient.Get(code)
	if err != nil {
		if err != memcache.ErrCacheMiss {
			log.Println("Memcache error:", err)
		}

		originalUrl, err = queries.GetOriginalUrl(r.Context(), code)
		if err != nil {
			if err != sql.ErrNoRows {
				utils.ServerErrTempl(w, "An error occur when querying the database")
				log.Println("Failed to retrive original_url from the database:", err)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		err = mcClient.Set(&memcache.Item{Key: code, Value: []byte(originalUrl)})
		if err != nil {
			log.Println("Memcache error when setting the key:", err)
		}
	} else if cache != nil && cache.Value != nil { // Check if original URL is in the cache
		cacheOriginalUrl := string(cache.Value)
		// Redirect to original URL
		http.Redirect(w, r, cacheOriginalUrl, http.StatusFound)

		// Update the accessed_count and last_accessed in one query
		err = queries.UpdateAccessedAndLastCount(r.Context(), sqlc.UpdateAccessedAndLastCountParams{
			Code:         code,
			LastAccessed: sql.NullTime{Time: time.Now().UTC(), Valid: true},
		})
		if err != nil {
			log.Println("Failed to update accessed_count and last_accessed:", err)
			return
		}
		return
	}

	// Redirect to original URL
	http.Redirect(w, r, originalUrl, http.StatusFound)

	// Update the accessed_count and last_accessed in one query
	err = queries.UpdateAccessedAndLastCount(r.Context(), sqlc.UpdateAccessedAndLastCountParams{
		Code:         code,
		LastAccessed: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})
	if err != nil {
		log.Println("Failed to update accessed_count and last_accessed:", err)
		return
	}
}
