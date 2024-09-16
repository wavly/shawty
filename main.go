package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/wavly/shawty/asserts"
	"github.com/wavly/shawty/utils"
	"github.com/wavly/shawty/handlers"
)

const PORT string = "1234"

func main() {
	// Creating the ServerMux router
	router := http.NewServeMux()

	// Check if memcache is up
	mcClient := memcache.New("0.0.0.0:11211")
	asserts.NoErr(mcClient.Ping(), "Failed to ping MemcacheD")

	// Serving static files
	router.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Reading the URLS-SQL schema file
	fileBytes, err := os.ReadFile("./schema/urls.sql")
	asserts.NoErr(err, "Failed to read URLS-SQL schema file")

	db := utils.ConnectDB()
	defer db.Close()

	// Create the URLs table in the database
	_, err = db.Exec(string(fileBytes))
	asserts.NoErr(err, "Error creating the URLs table in the database")

	// Ping/Pong route
	router.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong\n"))
	})

	// Route for shortening the URL
	router.HandleFunc("POST /", handlers.Main)

	// Route for index page
	router.Handle("GET /", http.FileServer(http.Dir("./static/")))

	// Route for stats page
	router.HandleFunc("GET /stat/{code}", handlers.Stats)

	// Route to handle redirection
	router.HandleFunc("GET /s/{code}", handlers.Redirect)

	// API route for shortening the URL
	router.HandleFunc("POST /shawty", handlers.Shawty)

	fmt.Println("Listening on:", PORT)
	asserts.NoErr(http.ListenAndServe("0.0.0.0:"+PORT, router), "Failed to start the server:")
}
