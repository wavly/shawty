package validate

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/tursodatabase/go-libsql"
	"github.com/wavly/shawty/internal/database"
)

// TestEvictOldLinks validates the EvictOldLinks function.
func TestEvictOldLinks(t *testing.T) {
	// Create an in-memory SQLite database
	db, err := sql.Open("libsql", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create in-memory database: %v", err)
	}
	defer db.Close()
	queries := database.New(db)

	err = queries.CreateUrlTable(context.Background())
	if err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	// Seed the database with test data
	oneMonthAgo := time.Now().AddDate(0, -1, 0)
	twoMonthsAgo := time.Now().AddDate(0, -2, 0)
	recent := time.Now()

	_, err = db.Exec(`
		INSERT INTO urls (original_url, code, created_at, accessed_count, last_accessed) VALUES
		('http://example1.com', 'abc123', ?, 5, ?),  -- Should get removed
		('http://example2.com', 'def456', ?, 3, ?),  -- Should get removed
		('http://example3.com', 'ghi789', ?, 10, ?); -- Should remain
	`, twoMonthsAgo, twoMonthsAgo, oneMonthAgo, oneMonthAgo, recent, recent)
	if err != nil {
		t.Fatalf("Failed to seed data: %v", err)
	}

	EvictOldLinks(db)

	// Validate that old links are evicted
	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM urls`).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query the database: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 remaining link, got %d", count)
	}

	// Validate remaining link
	var url string
	err = db.QueryRow(`SELECT original_url FROM urls`).Scan(&url)
	if err != nil {
		t.Fatalf("Failed to query remaining link: %v", err)
	}

	if url != "http://example3.com" {
		t.Errorf("Expected remaining link to be 'http://example3.com', got %s", url)
	}
}
