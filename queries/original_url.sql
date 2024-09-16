-- name: GetOriginalUrl :one
SELECT original_url FROM urls where code = ?;

-- name: UpdateAccessedAndLastCount :exec
UPDATE urls SET accessed_count = accessed_count + 1, last_accessed = ? WHERE code = ?
