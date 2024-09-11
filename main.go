package main

import (
  "fmt"
  "log"
  "net/http"

  "github.com/joho/godotenv"
  "github.com/wavly/shawty/asserts"
)

const PORT string = "1234"

func main() {
  // Creating the ServerMux router
  router := http.NewServeMux()

  // Loading the environment variables
  err := godotenv.Load()
  asserts.NoErr(err, "Failed to load environment variables")

  // Ping/Pong route
  router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("pong\n"))
  })

  fmt.Println("Listening on:", PORT)
  if err := http.ListenAndServe("0.0.0.0:"+PORT, router); err != nil {
    log.Fatalln("Failed to start the server:", err)
  }
}
