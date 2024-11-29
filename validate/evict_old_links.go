package validate

import (
	"context"
	"database/sql"
	"time"

	"github.com/wavly/shawty/asserts"
	"github.com/wavly/shawty/internal/database"
	prettylogger "github.com/wavly/shawty/pretty-logger"
	"github.com/wavly/shawty/utils"
)

// Evict Old links form the database
//
// The function runs periodically and blocks the thread forever, use a
// goroutine to run the function in a separate goroutine.
func EvictOldLinks(mins time.Duration) {
	asserts.AssertEq(mins <= 0, "The time to evict the old links from the database has to be greater than 0")
	logger := prettylogger.GetLogger(nil)

	db := utils.ConnectDB()
	defer db.Close()
	queries := database.New(db)

	for {
		rows, err := queries.GetLastAccessedTime(context.Background())

		// Don't assert if failed to query last_accessed from the database because
		// it will then crash the whole server!
		if err != nil {
			logger.Error("Failed to query last_accessed", "error", err)
			time.Sleep(time.Minute * mins)
			continue
		}

		// TODO: change the last_accessed to defualt to the current time
		for _, row := range rows {
			if row.LastAccessed.Valid {
				now := time.Now()
				oneMonthAgo := now.AddDate(0, -1, 0)

				if row.LastAccessed.Time.Before(oneMonthAgo) {
					logger.Info("Link is older than a month", "link", row.OriginalUrl, "time", row.LastAccessed.Time)
					err = queries.DeleteLinkTime(context.Background(), sql.NullTime{
						Valid: true,
						Time:  row.LastAccessed.Time,
					})
					if err != nil {
						logger.Error("Failed to delete the link from the database", "error", err)
						continue
					}
				}
			}
		}

		time.Sleep(time.Minute * mins)
	}
}
