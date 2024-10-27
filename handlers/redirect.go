package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/wavly/shawty/asserts"
	. "github.com/wavly/shawty/cache"
	"github.com/wavly/shawty/internal/database"
	prettylogger "github.com/wavly/shawty/pretty-logger"
	"github.com/wavly/shawty/utils"
	"github.com/wavly/shawty/validate"
)

var Logger = prettylogger.GetLogger(nil)

func Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	// Validate the code
	err := validate.CustomCodeValidate(code)
	if err != nil {
		Logger.Warn("Code validation failed", "code", code, "user-agent", r.UserAgent(), "error", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
			if err != sql.ErrNoRows {
				w.WriteHeader(http.StatusInternalServerError)
				utils.ServerErrTempl(w, "Sorry, an unexpected error occur when querying the database for the URL")
				Logger.Error("failed to query the database for the original URL", "error", err)
				return
			}

			Logger.Warn("Redirect code not found", "code", code, "user-agent", r.UserAgent())
			w.WriteHeader(http.StatusNotFound)
			asserts.NoErr(utils.Templ("./templs/404.html").Execute(w, nil), "Failed to execute 404 template in redirect route")
			return
		}

		Logger.Info("Cache Miss, redirect code not found", "code", code, "user-agent", r.UserAgent())
		http.Redirect(w, r, originalUrl, http.StatusSeeOther)
		Cache.Set(code, originalUrl, cache.DefaultExpiration)
		return
	}

	Logger.Info("Cache hit, redirect code found", "code", code, "user-agent", r.UserAgent())
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
