package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/wavly/shawty/asserts"
	"github.com/wavly/shawty/database"
	"github.com/wavly/shawty/handlers"
)

const PORT string = "1234"

func main() {
	// Creating the ServerMux router
	router := http.NewServeMux()

	// Reading the URLS-SQL schema file
	fileBytes, err := os.ReadFile("./schema/urls.sql")
	asserts.NoErr(err, "Failed to read URLS-SQL schema file")

	db := database.ConnectDB()
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

	// Route to handle redirection
	router.HandleFunc("GET /u/{url}", handlers.Redirection)

	fmt.Println("Listening on:", PORT)
	if err := http.ListenAndServe("0.0.0.0:"+PORT, router); err != nil {
		log.Fatalln("Failed to start the server:", err)
	}
}
