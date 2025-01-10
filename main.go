package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/wavly/surf/asserts"
	"github.com/wavly/surf/config"
	"github.com/wavly/surf/env"
	"github.com/wavly/surf/handlers"
	"github.com/wavly/surf/internal/database"
	prettylogger "github.com/wavly/surf/pretty-logger"
	"github.com/wavly/surf/utils"
	"github.com/wavly/surf/validate"
)

func main() {
	router := http.NewServeMux()
	logger := prettylogger.GetLogger(nil)

	// Get the env variables and other config options
	config.Init()

	// Create the URLs table in the database
	db := utils.ConnectDB()
	err := database.New(db).CreateUrlTable(context.Background())
	asserts.NoErr(err, "Failed to create the URLs table in the database")
	db.Close()

	// Find and remove links that are older than a month every hour
	go func() {
		for {
			validate.EvictOldLinks(db)
			time.Sleep(60 * time.Minute)
		}
	}()

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

	// Route for stats page
	router.HandleFunc("GET /stat/{code}", handlers.Stats)

	// Route for stats page
	router.HandleFunc("GET /url-info", func(w http.ResponseWriter, r *http.Request) {
		templ := utils.Templ("./templs/url-info.html")
		asserts.NoErr(templ.Execute(w, nil), "Failed to execute template url-info.html")
	})

	// Route to handle redirection
	router.HandleFunc("GET /s/{code}", handlers.Redirect)

	// Route for unshortening the URL
	router.HandleFunc("GET /unshort", func(w http.ResponseWriter, r *http.Request) {
		w.Write(utils.StaticFile("./static/unshort.html"))
	})

	// API route for shortening the URL
	router.HandleFunc("POST /short", handlers.Short)

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
