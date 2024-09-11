package main

func getLongUrl(slug string) (string, bool) {
  // Simulate the db for now
  urlDB := map[string]string{
    "abc123": "http://testabc.com",
    "discord": "https://discord.com",
    "wavly": "https://github.com/wavly",
  }

  longURL, exist := urlDB[slug]
  return longURL, exist
}
