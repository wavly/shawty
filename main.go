package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/wavly/surf/asserts"
	"github.com/wavly/surf/config"
	"github.com/wavly/surf/env"
	"github.com/wavly/surf/handlers"
	"github.com/wavly/surf/internal/database"
	prettylogger "github.com/wavly/surf/pretty-logger"
	"github.com/wavly/surf/static"
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
	})))

	// Route for index page
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			logger.Warn("Page not found", "route", r.URL.Path, "user-agent", r.UserAgent())

			// NOTE: HTMX doesn't swap the elements if the returned status code isn't successful, e.g 4xx, 5xx
			// This is fine "here" because it isn't using the hx-trigger attribute to swap the elements
			w.WriteHeader(http.StatusNotFound)
			err := static.PageNotFound().Render(r.Context(), w)
			asserts.NoErr(err, "Failed to render 404-page template")
			return
		}

		err := static.Layout(static.Index()).Render(r.Context(), w)
		asserts.NoErr(err, "Failed to render index-template")
	})

	// Stats page
	router.HandleFunc("GET /stat/{code}", handlers.Stats)

	// URL-Info page
	router.Handle("GET /url-info", templ.Handler(static.Layout(static.UrlInfo())))

	// Handle redirection
	router.HandleFunc("GET /s/{code}", handlers.Redirect)

	// Unshort Page (get the destination URL for the short link)
	router.Handle("GET /unshort", templ.Handler(static.Layout(static.UnShort())))

	// API Route: shorten the URL
	router.HandleFunc("POST /short", handlers.Short)

	// API Route: unshorten the URL
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
