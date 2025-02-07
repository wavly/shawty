package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/mergestat/timediff"
	"github.com/wavly/surf/asserts"
	"github.com/wavly/surf/internal/database"
	"github.com/wavly/surf/static"
	"github.com/wavly/surf/utils"
	"github.com/wavly/surf/validate"
)

type AccessCount struct {
	ShortLink
	Count int64

	LastAccessed string
	OriginalUrl  string
}

func Stats(w http.ResponseWriter, r *http.Request) {
	inputCode := r.PathValue("code")

	err := validate.CustomCodeValidate(inputCode)
	if err != nil {
		Logger.Warn("failed to validate the input code", "code", inputCode, "user-agent", r.UserAgent(), "error", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	db := utils.ConnectDB()
	defer db.Close()
	queries := database.New(db)

	shortLinkInfo, err := queries.GetShortCodeInfo(r.Context(), inputCode)
	if err != nil {
		if err != sql.ErrNoRows {
			Logger.Error("failed to query to get the short url info", "code", inputCode, "user-agent", r.UserAgent(), "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			err := static.ServerError("An error occur when querying the database").Render(r.Context(), w)
			asserts.NoErr(err, "Failed to render server-internal-error page")
			return
		}

		Logger.Warn("Stats not found", "code", inputCode, "user-agent", r.UserAgent())
		w.WriteHeader(http.StatusNotFound)
		err = static.PageNotFound().Render(r.Context(), w)
		asserts.NoErr(err, "Failed to render 404-page template")
		return
	}

	count := strconv.Itoa(int(shortLinkInfo.AccessedCount))
	err = static.Layout(static.Stats(inputCode, shortLinkInfo.OriginalUrl, timediff.TimeDiff(shortLinkInfo.LastAccessed), count)).Render(r.Context(), w)
	asserts.NoErr(err, "Failed to render stats page")
}
