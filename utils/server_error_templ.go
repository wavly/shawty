package utils

import (
	"html/template"
	"net/http"

	"github.com/wavly/shawty/asserts"
)

func ServerErrTempl(w http.ResponseWriter, msg string) {
	templ := template.Must(template.ParseFiles("templs/server-error.html"))
	w.WriteHeader(http.StatusInternalServerError)
	asserts.NoErr(templ.Execute(w, msg), "Failed to execute template server-error.html")
}
