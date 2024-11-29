package validate

import (
	"context"
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
	sleepDuration := time.Minute * mins
	queries := database.New(utils.ConnectDB())

	for {
		rows, err := queries.GetLastAccessedTime(context.Background())

		// Don't assert if failed to query last_accessed from the database because
		// it will then crash the whole server!
		if err != nil {
			logger.Error("Failed to query last_accessed", "error", err)
			time.Sleep(sleepDuration)
			continue
		}

		for _, row := range rows {
			now := time.Now()
			oneMonthAgo := now.AddDate(0, -1, 0)

			if row.LastAccessed.Before(oneMonthAgo) {
				logger.Info("Link is older than a month", "link", row.OriginalUrl)
				err = queries.DeleteLinkLastAccessed(context.Background(), row.LastAccessed)

				if err != nil {
					logger.Error("Failed to delete the link from the database", "error", err)
					continue
				}
			}
		}

		time.Sleep(sleepDuration)
	}
}
