package handlers

import (
	"database/sql"
	"net/http"
	"net/url"
	"strings"

	"github.com/patrickmn/go-cache"
	"github.com/wavly/surf/asserts"
	. "github.com/wavly/surf/cache"
	"github.com/wavly/surf/internal/database"
	partialhtml "github.com/wavly/surf/partial-html"
	"github.com/wavly/surf/utils"
	"github.com/wavly/surf/validate"
)

func Unshort(w http.ResponseWriter, r *http.Request) {
	inputUrl := r.FormValue("url")
	Logger.Info("Unshorten the URL using the code", "code", inputUrl, "user-agent", r.UserAgent())

	if len(inputUrl) > 1000 {
		err := partialhtml.ShortLinkError("URL is too long, Only 1000 characters are allowed").Render(r.Context(), w)
		asserts.NoErr(err, "Failed to render partial-html short-link-error")
		return
	}

	if !strings.Contains(inputUrl, "://") {
		inputUrl = "https://" + inputUrl
	}

	parsedUrl, err := url.Parse(inputUrl)
	if err != nil {
		Logger.Warn("failed to parse the url", "url", inputUrl, "user-agent", r.UserAgent(), "error", err)
		err = partialhtml.ShortLinkError("Invalid URL format: "+inputUrl).Render(r.Context(), w)
		asserts.NoErr(err, "Failed to render partial-html short-link-error")
		return
	}

	// Check if domain matches
	if parsedUrl.Host != "surf.wavly.tech" || parsedUrl.Scheme != "https" {
		err = partialhtml.ShortLinkError("URL must use 'https://surf.wavly.tech'").Render(r.Context(), w)
		asserts.NoErr(err, "Failed to render partial-html short-link-error")
		return
	}

	// Check if the path follows the expected structure
	pathParts := strings.Split(parsedUrl.Path, "/")
	if len(pathParts) != 3 || pathParts[1] != "s" || pathParts[2] == "" {
		err = partialhtml.ShortLinkError("URL must follow the pattern '/s/{code}'").Render(r.Context(), w)
		asserts.NoErr(err, "Failed to redner partial-html short-link-error")
		return
	}

	code := pathParts[2]
	err = validate.CustomCodeValidate(code)
	if err != nil {
		Logger.Warn("failed to validate custom code for the URL", "code", code, "original-url", inputUrl, "error", err, "user-agent", r.UserAgent())
		err = partialhtml.ShortLinkError("Failed to render partial-html short-link-error").Render(r.Context(), w)
		asserts.NoErr(err, "")
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
				err = partialhtml.ShortLinkError("There is no destination URL for this short URL: "+"wavly.surf.com/s/"+code).Render(r.Context(), w)
				asserts.NoErr(err, "Failed to render partial-html short-link-error")
				return
			}

			Logger.Error("failed to query to get the original url", "code", code, "user-agent", r.UserAgent(), "error", err)
			err = partialhtml.ShortLinkError("Sorry, unexpected error occur when querying the database").Render(r.Context(), w)
			asserts.NoErr(err, "Failed to render partial-html short-link-error")
			return
		}

		Cache.Set(inputUrl, ret, cache.DefaultExpiration)
		err = partialhtml.UnShort(ret).Render(r.Context(), w)
		asserts.NoErr(err, "Failed to render partial-html unshort-link")
		return
	}

	Logger.Info("Cache Hit: code for the URL found in the cache", "code", code, "url", originalUrl, "user-agent", r.UserAgent())
	err = partialhtml.UnShort(originalUrl.(string)).Render(r.Context(), w)
	asserts.NoErr(err, "Failed to render partial-html unshort-link")
}
