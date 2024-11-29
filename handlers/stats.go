package handlers

import (
	"database/sql"
	"net/http"

	"github.com/mergestat/timediff"
	"github.com/wavly/shawty/asserts"
	"github.com/wavly/shawty/internal/database"
	"github.com/wavly/shawty/utils"
	"github.com/wavly/shawty/validate"
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

	templ := utils.Templ("./templs/stat.html")
	db := utils.ConnectDB()
	defer db.Close()
	queries := database.New(db)

	shortLinkInfo, err := queries.GetShortCodeInfo(r.Context(), inputCode)
	if err != nil {
		if err != sql.ErrNoRows {
			Logger.Error("failed to query to get the short url info", "code", inputCode, "user-agent", r.UserAgent(), "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			utils.ServerErrTempl(w, "An error occur when querying the database")
			return
		}

		Logger.Warn("Stats not found", "code", inputCode, "user-agent", r.UserAgent())
		w.WriteHeader(http.StatusNotFound)
		asserts.NoErr(utils.Templ("./templs/404.html").Execute(w, nil), "Failed to execute 404 template in stats page")
		return
	}

	data := AccessCount{
		Count:        shortLinkInfo.AccessedCount,
		LastAccessed: timediff.TimeDiff(shortLinkInfo.LastAccessed),
		OriginalUrl:  shortLinkInfo.OriginalUrl,

		ShortLink: ShortLink{
			ShortUrl: inputCode,
		},
	}

	asserts.NoErr(templ.Execute(w, data), "Failed to execute template stat.html")
}
