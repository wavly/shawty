package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/mergestat/timediff"
	"github.com/wavly/shawty/asserts"
	sqlc "github.com/wavly/shawty/sqlc_db"
	"github.com/wavly/shawty/utils"
)

type AccessCount struct {
	ShortLink
	Count int64

	LastAccessed string
}

func Stats(w http.ResponseWriter, r *http.Request) {
	inputCode := r.PathValue("code")

	if len(inputCode) > 8 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	templ := template.Must(template.ParseFiles("./templs/stat.html"))
	db := utils.ConnectDB()
	defer db.Close()
	queries := sqlc.New(db)

	shortLinkInfo, err := queries.GetShortCodeInfo(r.Context(), inputCode)
	if err != nil {
		if err != sql.ErrNoRows {
			utils.ServerErrTempl(w, "An error occur when querying the database")
			log.Printf("Database error when selecting accessed_count and original_url where code = %s, Error %s\n", inputCode, err)
			return
		}

		notFoundTempl := template.Must(template.ParseFiles("./templs/404.html"))
		w.WriteHeader(http.StatusNotFound)
		notFoundTempl.Execute(w, nil)
		return
	}

	data := AccessCount{
		Count:        shortLinkInfo.AccessedCount,
		LastAccessed: timediff.TimeDiff(shortLinkInfo.LastAccessed.Time),
		ShortLink: ShortLink{
			ShortUrl:    inputCode,
			OriginalUrl: shortLinkInfo.OriginalUrl,
		},
	}

	// Checking if the last accessed timestamp is not null
	// And if true: set the LastAccessed value to "None"
	if !shortLinkInfo.LastAccessed.Valid {
		data.LastAccessed = "None"
	}

	asserts.NoErr(templ.Execute(w, data), "Failed to execute template stat.html")
}
