package main

import (
  "net/http"
)

func redirectionHandler(w http.ResponseWriter, r *http.Request) {
  slug := r.URL.Path[8:]
  longURL, exist := getLongUrl(slug)

  if exist {
    http.Redirect(w, r, longURL, http.StatusFound)
  } else {
    w.Write([]byte("Couldn't find anything in your shawty, recheck..\n"))
  }
}
