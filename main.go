package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/wavly/shawty/asserts"
	"github.com/wavly/shawty/config"
	"github.com/wavly/shawty/env"
	"github.com/wavly/shawty/handlers"
	"github.com/wavly/shawty/internal/database"
	prettylogger "github.com/wavly/shawty/pretty-logger"
	"github.com/wavly/shawty/utils"
)

func main() {
	router := http.NewServeMux()
	logger := prettylogger.GetLogger(nil)

	// Get the env variables and other config options
	config.Init(router)

	// Create the URLs table in the database
	db := utils.ConnectDB()
	err := database.New(db).CreateUrlTable(context.Background())
	asserts.NoErr(err, "Failed to creating the URLs table in the database")
	db.Close()

	// Serving static files with caching
	router.Handle("GET /static/", http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set Cache-Control headers for 10 days
		w.Header().Set("Cache-Control", "public, max-age=864000")

		// Set Last-Modified header
		w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))

		// Serve the static content
		http.FileServer(http.Dir("./static")).ServeHTTP(w, r)
		logger.Debug("Request for static content", "resource", r.URL.Path, "from-ip", r.RemoteAddr)
	})))

	// Route for index page
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			logger.Warn("Page not found", "route", r.URL.Path, "user-agent", r.UserAgent())
			utils.Templ("./templs/404.html").Execute(w, nil)
			return
		}
		w.Write(utils.StaticFile("./static/index.html"))
	})

	// Route for shortening the URL
	router.HandleFunc("POST /", handlers.Main)

	// Route for stats page
	router.HandleFunc("GET /stat/{code}", handlers.Stats)

	// Route to handle redirection
	router.HandleFunc("GET /s/{code}", handlers.Redirect)

	// API route for shortening the URL
	router.HandleFunc("POST /shawty", handlers.Shawty)

	fmt.Printf("Listening on: %s\n\n", env.PORT)
	asserts.NoErr(http.ListenAndServe("0.0.0.0:"+env.PORT, router), "Failed to start the server")
}
