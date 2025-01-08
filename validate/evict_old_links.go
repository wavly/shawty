package validate

import (
	"context"
	"database/sql"
	"time"

	"github.com/wavly/surf/internal/database"
	prettylogger "github.com/wavly/surf/pretty-logger"
)

// Evict Old links form the database
func EvictOldLinks(db *sql.DB) {
	logger := prettylogger.GetLogger(nil)
	queries := database.New(db)
	ctx := context.Background()

	rows, err := queries.GetLastAccessedTime(ctx)

	// Don't assert if failed to query last_accessed from the database because
	// it will then crash the whole server!
	if err != nil {
		logger.Error("Failed to query last_accessed", "error", err)
		return
	}

	for _, row := range rows {
		now := time.Now()
		oneMonthAgo := now.AddDate(0, -1, 0)

		if row.LastAccessed.Before(oneMonthAgo) {
			logger.Info("Link is older than a month", "link", row.OriginalUrl)
			if err = queries.DeleteLinkLastAccessed(ctx, row.LastAccessed); err != nil {
				logger.Error("Failed to delete the link from the database", "error", err)
				return
			}
		}
	}
}
