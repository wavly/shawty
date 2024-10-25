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
		logger.Debug("Request for static content", "resource", r.URL.Path, "user-agent", r.UserAgent())
	})))

	// Route for index page
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			logger.Warn("Page not found", "route", r.URL.Path, "user-agent", r.UserAgent())

			// NOTE: HTMX doesn't swap the elements if the returned status code isn't successful, e.g 4xx, 5xx
			// This is fine "here" because it isn't using the hx-trigger attribute to swap the elements
			w.WriteHeader(http.StatusNotFound)
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

	// Route for unshortening the URL
	router.HandleFunc("GET /unshort", func(w http.ResponseWriter, r *http.Request) {
		w.Write(utils.StaticFile("./static/unshort.html"))
	})

	// API route for shortening the URL
	router.HandleFunc("POST /shawty", handlers.Shawty)

	// API route for unshortening the URL
	router.HandleFunc("POST /unshort", handlers.Unshort)

	fmt.Printf("Listening on: %s\n\n", env.PORT)
	server := &http.Server{
		Addr:         ":" + env.PORT,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	asserts.NoErr(server.ListenAndServe(), "Failed to start the server")
}
