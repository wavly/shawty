package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/wavly/shawty/internal/database"
	"github.com/wavly/shawty/utils"
	"github.com/wavly/shawty/validate"
)

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

	originalUrl, err := queries.GetOriginalUrl(r.Context(), code)
	if err != nil {
		utils.ServerErrTempl(w, "Sorry, an unexpected error occur when querying the database for the URL")
		log.Println("Error: when querying the database for the URL", err)
		return
	}

	http.Redirect(w, r, originalUrl, http.StatusSeeOther)

	// Update the accessed_count and last_accessed in one query
	err = queries.UpdateAccessedAndLastCount(r.Context(), database.UpdateAccessedAndLastCountParams{
		Code:         code,
		LastAccessed: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})

	if err != nil {
		log.Println("Failed to update accessed_count and last_accessed:", err)
		return
	}
}
