package utils

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/wavly/shawty/asserts"
	"github.com/wavly/shawty/env"
)

const script = `<script src="/static/live-reload.js"></script>`

// Injects a script into the [Template File] to live-reload the page if in `DEV` mode,
// else just returns the file as it is.
func Templ(path string) *template.Template {
	if env.MODE != "dev" {
		return template.Must(template.ParseFiles(path))
	}
	content, err := os.ReadFile(path)
	asserts.NoErr(err, "Failed to read file")

	content = append(content, []byte(script)...)
	ret, err := template.New(filepath.Base(path)).Parse(string(content))
	asserts.NoErr(err, "Faield to parse template")
	return ret
}

// Injects a script into the [Static File] to live-reload the page if in `DEV` mode,
// else just return the file as it is.
func StaticFile(path string) []byte {
	if env.MODE != "dev" {
		ret, err := os.ReadFile(path)
		asserts.NoErr(err, "Failed to read file")
		return ret
	}

	content, err := os.ReadFile(path)
	asserts.NoErr(err, "Failed to read file")

	content = append(content, []byte(script)...)
	return content
}
