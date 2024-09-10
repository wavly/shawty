package main

import (
  "fmt"
  "log"
  "net/http"
)

const PORT string = "1234"

func main() {
  router := http.NewServeMux()

  router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("pong\n"))
  })

  fmt.Println("Listening on:", PORT)
  if err := http.ListenAndServe("0.0.0.0:"+PORT, router); err != nil {
    log.Fatalln("Failed to start the server:", err)
  }
}
