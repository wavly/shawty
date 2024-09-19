package handlers

import (
    "crypto/sha256"
    "database/sql"
    "encoding/hex"
    "fmt"
    "log"
    "net/http"
    "strings"

    "github.com/wavly/shawty/internal/database"
    "github.com/wavly/shawty/utils"
    "github.com/wavly/shawty/validate"
)

func Main(w http.ResponseWriter, r *http.Request) {
    inputUrl := r.FormValue("url")
    customCode := r.FormValue("code")

    // Validate customCode
    err := validate.CustomCodeValidate(customCode)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    // Check if inputUrl contains "://" and add "https://" if missing
    if !strings.Contains(inputUrl, "://") {
        inputUrl = "https://" + inputUrl
    }

    // Validate the URL
    err = validate.ValidateUrl(inputUrl)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    db := utils.ConnectDB()
    defer db.Close()
    queries := database.New(db)

    if customCode != "" {
        // Check if the custom code already exists in the database
        _, err := queries.GetCode(r.Context(), customCode)
        if err != nil {
            if err != sql.ErrNoRows {
                w.Write([]byte("An unexpected error occurred when querying from the database"))
                log.Printf("Database error when selecting code where code = %s, Error: %s\n", customCode, err)
                return
            }

            // Insert the URL in the database if it doesn't exist
            _, err = queries.CreateShortLink(r.Context(), database.CreateShortLinkParams{
                OriginalUrl: inputUrl,
                Code:        customCode,
            })
            if err != nil {
                w.Write([]byte("An unexpected error occurred when saving the URL to the database"))
                log.Println("Failed to store URL in the database", err)
                return
            }

            w.WriteHeader(http.StatusCreated)
            w.Write([]byte(fmt.Sprintf("Location: /s/%s", customCode)))
            return
        }

        w.WriteHeader(http.StatusConflict)
        w.Write([]byte("Custom code already exists"))
        return
    }

    // Hashing the URL
    hasher := sha256.New()
    hasher.Write([]byte(inputUrl))
    checksum := hasher.Sum(nil)

    // Truncate to 8 characters long hash
    hashUrl := hex.EncodeToString(checksum[:4])

    // Check if the hash URL already exists in the database
    _, err = queries.GetCode(r.Context(), hashUrl)
    if err != nil {
        if err != sql.ErrNoRows {
            w.Write([]byte("An unexpected error occurred when querying from the database"))
            log.Printf("Database error when selecting code where code = %s, Error: %s\n", hashUrl, err)
            return
        }

        // Insert the URL in the database if it doesn't exist
        _, err = queries.CreateShortLink(r.Context(), database.CreateShortLinkParams{
            OriginalUrl: inputUrl,
            Code:        hashUrl,
        })
        if err != nil {
            w.Write([]byte("An unexpected error occurred when saving the URL to the database"))
            log.Println("Failed to store URL in the database", err)
            return
        }

        w.WriteHeader(http.StatusCreated)
        w.Write([]byte(fmt.Sprintf("Location: /s/%s", hashUrl)))
        return
    }

    w.WriteHeader(http.StatusConflict)
    w.Write([]byte("Generated hash URL already exists"))
}