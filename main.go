package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/wavly/shawty/asserts"
	"github.com/wavly/shawty/config"
	"github.com/wavly/shawty/env"
	"github.com/wavly/shawty/handlers"
	"github.com/wavly/shawty/utils"
)

func main() {
	// Creating the ServerMux router
	router := http.NewServeMux()

	// Get the env variables and other config options
	config.Init(router)

	// Serving static files
	// TODO: make static content cache
	router.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Route for index page
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write(utils.StaticFile("./static/index.html"))
	})

	// Reading the URLS-SQL schema file
	fileBytes, err := os.ReadFile("./schema/urls.sql")
	asserts.NoErr(err, "Failed to read ./schema/urls.sql schema file")

	db := utils.ConnectDB()
	// Create the URLs table in the database
	_, err = db.Exec(string(fileBytes))
	db.Close()
	asserts.NoErr(err, "Failed to creating the URLs table in the database")

	// Route for shortening the URL
	router.HandleFunc("POST /", handlers.Main)

	// Route for stats page
	router.HandleFunc("GET /stat/{code}", handlers.Stats)

	// Route to handle redirection
	router.HandleFunc("GET /s/{code}", handlers.Redirect)

	// API route for shortening the URL
	router.HandleFunc("POST /shawty", handlers.Shawty)

	fmt.Println("Listening on:", env.PORT)
	asserts.NoErr(http.ListenAndServe("0.0.0.0:"+env.PORT, router), "Failed to start the server")
}
