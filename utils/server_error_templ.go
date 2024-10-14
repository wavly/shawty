package utils

import (
	"net/http"

	"github.com/wavly/shawty/asserts"
)

func ServerErrTempl(w http.ResponseWriter, msg string) {
	templ := Templ("templs/server-error.html")
	asserts.NoErr(templ.Execute(w, msg), "Failed to execute template server-error.html")
}
