-- name: CreateEntry :one
INSERT INTO entries (
  username,
  video_name,
  amount
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
WHERE username = $1
ORDER BY id
LIMIT $2
OFFSET $3;

