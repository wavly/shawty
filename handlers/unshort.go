package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"github.com/patrickmn/go-cache"
	"github.com/wavly/shawty/asserts"
	. "github.com/wavly/shawty/cache"
	"github.com/wavly/shawty/internal/database"
	"github.com/wavly/shawty/utils"
	"github.com/wavly/shawty/validate"
)

func Unshort(w http.ResponseWriter, r *http.Request) {
	inputUrl := r.FormValue("url")
	errorTempl := template.Must(template.ParseFiles("./partial-html/short-link-error.html"))
	Logger.Info("Unshorten the URL using the code", "code", inputUrl, "user-agent", r.UserAgent())

	if len(inputUrl) > 1000 {
		asserts.NoErr(errorTempl.Execute(w, "URL is too long, Only 1000 characters are allowed"), "Failed to execute template short-link-error.html")
		return
	}

	if !strings.Contains(inputUrl, "://") {
		inputUrl = "https://"+inputUrl
	}

	parsedUrl, err := url.Parse(inputUrl)
	if err != nil {
		Logger.Warn("failed to parse the url", "url", inputUrl, "user-agent", r.UserAgent(), "error", err)
		asserts.NoErr(errorTempl.Execute(w, "Invalid URL format: "+inputUrl), "Failed to execute template short-link-error.html")
		return
	}

	// Check if domain matches
	if parsedUrl.Host != "wavly.shawty.com" || parsedUrl.Scheme != "https" {
		asserts.NoErr(errorTempl.Execute(w, "URL must use 'https://wavly.shawty.com'"), "Failed to execute template short-link-error.html")
		return
	}

	// Check if the path follows the expected structure
	pathParts := strings.Split(parsedUrl.Path, "/")
	if len(pathParts) != 3 || pathParts[1] != "s" || pathParts[2] == "" {
		asserts.NoErr(errorTempl.Execute(w, "URL must follow the pattern '/s/{code}'"), "Failed to execute template short-link-error.html")
		return
	}

	code := pathParts[2]
	err = validate.CustomCodeValidate(code)
	if err != nil {
		Logger.Warn("failed to validate custom code for the URL", "code", code, "original-url", inputUrl, "error", err, "user-agent", r.UserAgent())
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
				asserts.NoErr(errorTempl.Execute(w, "There is no destination URL for this short URL: "+"wavly.shawty.com/s/"+code), "Failed to execute template short-link-error.html")
				return
			}

			Logger.Error("failed to query to get the original url", "code", code, "user-agent", r.UserAgent(), "error", err)
			asserts.NoErr(errorTempl.Execute(w, "Sorry, unexpected error occur when querying the database"), "Failed to execute short-link-error.html")
			return
		}

		Cache.Set(inputUrl, ret, cache.DefaultExpiration)
		templ := template.Must(template.ParseFiles("./partial-html/unshort-link.html"))
		asserts.NoErr(templ.Execute(w, ret), "Failed to execute template unshort-link.html")
		return
	}

	Logger.Info("Cache Hit: code for the URL found in the cache", "code", code, "url", originalUrl, "user-agent", r.UserAgent())
	templ := template.Must(template.ParseFiles("./partial-html/unshort-link.html"))
	asserts.NoErr(templ.Execute(w, originalUrl), "Failed to execute template unshort-link.html")
}
