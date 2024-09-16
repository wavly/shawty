-- name: UpdateAccessedAndLastCount :exec
UPDATE urls SET accessed_count = accessed_count + 1, last_accessed = ? WHERE code = ?
