package config

import (
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/wavly/shawty/asserts"
	. "github.com/wavly/shawty/env"
	prettylogger "github.com/wavly/shawty/pretty-logger"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var logger = prettylogger.GetLogger(nil)

func Init(router *http.ServeMux) {
	err := godotenv.Load(".env.local")
	asserts.NoErr(err, "Failed reading .env.local")

	port := os.Getenv("PORT")
	asserts.AssertEq(port == "", "Missing PORT number in .env.local")
	PORT = port

	environment := os.Getenv("ENVIRONMENT")
	asserts.AssertEq(environment == "", "Missing ENVIRONMENT in .env.local")
	MODE = environment

	// Only get [TURSO_TOKEN] and [TURSO_URL] in Prodution
	if MODE == "prod" {
		tursoToken := os.Getenv("TURSO_AUTH_TOKEN")
		asserts.AssertEq(tursoToken == "", "Missing TURSO_AUTH_TOKEN in .evn.local")
		TURSO_TOKEN = tursoToken

		tursoURL := os.Getenv("TURSO_DATABASE_URL")
		asserts.AssertEq(tursoURL == "", "Missing TURSO_URL in .evn.local")
		TURSO_URL = tursoURL
		return
	}

	router.HandleFunc("GET /ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error("failed to start websocket connection", "error", err)
			return
		}
		defer conn.Close()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				logger.Error("Error reading message:", "error", err)
				break
			}
		}
	})
}
