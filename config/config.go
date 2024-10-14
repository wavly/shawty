package config

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/wavly/shawty/asserts"
)

var PORT string
var ENV string

var TURSO_TOKEN string
var TURSO_URL string

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Init(router *http.ServeMux) {
	err := godotenv.Load(".env.local")
	asserts.NoErr(err, "Failed reading .env.local")

	port := os.Getenv("PORT")
	asserts.AssertEq(port == "", "Please specify the PORT number in .env.local")
	PORT = port

	environment := os.Getenv("ENVIRONMENT")
	asserts.AssertEq(environment == "", "Please specify the ENVIRONMENT in .env.local")
	ENV = environment

	// Only get [TURSO_TOKEN] and [TURSO_URL] in Prodution
	if ENV == "prod" {
		tursoToken := os.Getenv("TURSO_AUTH_TOKEN")
		asserts.AssertEq(tursoToken == "", "Please provide the TURSO TOKEN in .evn.local")
		TURSO_TOKEN = tursoToken

		tursoURL := os.Getenv("TURSO_DATABASE_URL")
		asserts.AssertEq(tursoURL == "", "Please provide the TURSO URL in .evn.local")
		TURSO_URL = tursoURL
		return
	}

	router.HandleFunc("GET /ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("failed to start websocket connection:", err)
			return
		}
		defer conn.Close()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				break
			}
		}
	})
}
