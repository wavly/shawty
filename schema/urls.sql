CREATE TABLE IF NOT EXISTS urls (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  original_url TEXT NOT NULL,
  code TEXT UNIQUE NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  accessed_count INTEGER DEFAULT 0,
  last_accessed DATETIME
);

-- Creating an Index to retrieve the urls(code) from the database more quickly than otherwise
CREATE INDEX IF NOT EXISTS idx_short_code ON urls(code);
