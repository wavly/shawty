package handlers

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/patrickmn/go-cache"
	"github.com/wavly/shawty/asserts"
	. "github.com/wavly/shawty/cache"
	"github.com/wavly/shawty/internal/database"
	"github.com/wavly/shawty/utils"
	"github.com/wavly/shawty/validate"
)

// TODO: add support for links
// Current impl for this func only accepts the code the short link
func Unshort(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	errorTempl := template.Must(template.ParseFiles("./partial-html/short-link-error.html"))
	Logger.Info("Unshorten the URL using the code", "code", code, "user-agent", r.UserAgent())

	// Validate the Code
	err := validate.CustomCodeValidate(code)
	if err != nil {
		Logger.Warn("failed to validate the code", "url", code, "user-agent", r.UserAgent(), "error", err)
		asserts.NoErr(errorTempl.Execute(w, err), "Failed to execute template short-link-error.html")
		return
	}

	originalUrl, found := Cache.Get(code)
	if !found {
		Logger.Warn("Cache Miss: code for the URL isn't in the cache", "code", code, "user-agent", r.UserAgent())
		db := utils.ConnectDB()
		defer db.Close()
		queries := database.New(db)

		ret, err := queries.GetOriginalUrl(r.Context(), code)
		if err != nil {
			if err == sql.ErrNoRows {
				Logger.Warn("url doesn't exists in the database", "code", code, "user-agent", r.UserAgent(), "error", err)
				asserts.NoErr(errorTempl.Execute(w, "There is no destination URL for this short URL: "+code), "Failed to execute template short-link-error.html")
				return
			}

			Logger.Error("failed to query to get the original url", "code", code, "user-agent", r.UserAgent(), "error", err)
			asserts.NoErr(errorTempl.Execute(w, "Sorry, unexpected error occur when querying the database"), "Failed to execute short-link-error.html")
			return
		}

		Cache.Set(code, ret, cache.DefaultExpiration)
		templ := template.Must(template.ParseFiles("./partial-html/unshort-link.html"))
		asserts.NoErr(templ.Execute(w, ret), "Failed to execute template unshort-link.html")
		return
	}

	Logger.Info("Cache Hit: code for the URL found in the cache", "code", code, "url", originalUrl, "user-agent", r.UserAgent())
	templ := template.Must(template.ParseFiles("./partial-html/unshort-link.html"))
	asserts.NoErr(templ.Execute(w, originalUrl), "Failed to execute template unshort-link.html")
}
