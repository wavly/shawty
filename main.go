package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/joho/godotenv"
	"github.com/wavly/shawty/asserts"
	"github.com/wavly/shawty/handlers"
	"github.com/wavly/shawty/utils"
)

func main() {
	// Creating the ServerMux router
	router := http.NewServeMux()

	err := godotenv.Load(".env.local")
	asserts.NoErr(err, "Failed reading .env.local")

	port := os.Getenv("PORT")
	asserts.AssertEq(port == "", "Please specify the PORT number in .env.local")

	environment := os.Getenv("ENVIRONMENT")
	asserts.AssertEq(environment == "", "Please specify the ENVIRONMENT in .env.local")

	var mcClient *memcache.Client
	switch environment {
	case "dev":
		mcClient = memcache.New("0.0.0.0:11211")
		if err := mcClient.Ping(); err != nil {
			log.Println("Memcached listener is not up!")
		}
	case "prod":
		mcClient = memcache.New("0.0.0.0:11211")
		asserts.NoErr(mcClient.Ping(), "Memcached listener is required in production environment")
	default:
		log.Fatalln("Unrecognize environment value:", environment)
	}

	// Serving static files
	router.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Route for index page
	router.Handle("GET /", http.FileServer(http.Dir("./static/")))

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

	// Route for stats page
	router.HandleFunc("GET /stat/{code}", handlers.Stats)

	// Route to handle redirection
	router.HandleFunc("GET /s/{code}", handlers.Redirect)

	// API route for shortening the URL
	router.HandleFunc("POST /shawty", handlers.Shawty)

	fmt.Println("Listening on:", port)
	asserts.NoErr(http.ListenAndServe("0.0.0.0:"+port, router), "Failed to start the server:")
}
