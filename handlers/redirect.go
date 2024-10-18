package handlers

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	. "github.com/wavly/shawty/cache"
	"github.com/wavly/shawty/internal/database"
	prettylogger "github.com/wavly/shawty/pretty-logger"
	"github.com/wavly/shawty/utils"
	"github.com/wavly/shawty/validate"
)

var logger = slog.New(prettylogger.NewHandler(nil))

func Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	// Validate the code
	err := validate.CustomCodeValidate(code)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	db := utils.ConnectDB()
	defer db.Close()
	queries := database.New(db)

	// Try to get the original URL from cache
	originalUrl, found := Cache.Get(code)

	// Update the cache if the doesn't exist
	if !found {
		originalUrl, err := queries.GetOriginalUrl(r.Context(), code)
		if err != nil {
			utils.ServerErrTempl(w, "Sorry, an unexpected error occur when querying the database for the URL")
			log.Println("Error: when querying the database for the URL", err)
			return
		}

		http.Redirect(w, r, originalUrl, http.StatusSeeOther)
		Cache.Set(code, originalUrl, cache.DefaultExpiration)
		return
	}

	http.Redirect(w, r, originalUrl.(string), http.StatusSeeOther)

	// Todo: update last-time/access count in cache
	// Update the accessed_count and last_accessed in one query
	err = queries.UpdateAccessedAndLastCount(r.Context(), database.UpdateAccessedAndLastCountParams{
		Code:         code,
		LastAccessed: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})

	if err != nil {
		logger.Error("Failed to update accessed_count and last_accessed", "error", err)
		return
	}
}
