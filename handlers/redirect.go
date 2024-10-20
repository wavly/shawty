package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	. "github.com/wavly/shawty/cache"
	"github.com/wavly/shawty/internal/database"
	"github.com/wavly/shawty/utils"
	"github.com/wavly/shawty/validate"
)

func Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	// Validate the code
	err := validate.CustomCodeValidate(code)
	if err != nil {
		Logger.Warn("Code validation failed", "code", code, "from-ip", r.RemoteAddr, "user-agent", r.UserAgent())
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
		Logger.Info("Cache Miss, redirect code not found", "code", code, "from-ip", r.RemoteAddr, "user-agent", r.UserAgent())
		originalUrl, err := queries.GetOriginalUrl(r.Context(), code)
		if err != nil {
			utils.ServerErrTempl(w, "Sorry, an unexpected error occur when querying the database for the URL")
			Logger.Error("failed to query the database for the original URL", "error", err)
			return
		}

		http.Redirect(w, r, originalUrl, http.StatusSeeOther)
		Cache.Set(code, originalUrl, cache.DefaultExpiration)
		return
	}

	Logger.Info("Cache hit, redirect code found", "code", code, "from-ip", r.RemoteAddr, "user-agent", r.UserAgent())
	http.Redirect(w, r, originalUrl.(string), http.StatusSeeOther)

	// Todo: update last-time/access count in cache
	// Update the accessed_count and last_accessed in one query
	err = queries.UpdateAccessedAndLastCount(r.Context(), database.UpdateAccessedAndLastCountParams{
		Code:         code,
		LastAccessed: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})

	if err != nil {
		Logger.Error("failed to update accessed_count and last_accessed", "error", err)
		return
	}
}
