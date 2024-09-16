-- name: GetOriginalUrl :one
SELECT original_url FROM urls where code = ?;

-- name: UpdateAccessedAndLastCount :exec
UPDATE urls SET accessed_count = accessed_count + 1, last_accessed = ? WHERE code = ?;

-- name: CreateShortLink :one
INSERT INTO urls (
  original_url,
  code
) VALUES ( ?, ? )
RETURNING *;

-- name: GetShortCodeInfo :one
SELECT accessed_count, original_url, last_accessed FROM urls WHERE code = ?;
