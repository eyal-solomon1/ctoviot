-- name: CreateVideo :one
INSERT INTO videos (
  owner,
  video_identifier,
  video_name,
  video_length,
  video_remote_path,
  video_decs  
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;


-- name: GetVideo :one
SELECT * FROM videos
WHERE video_identifier = $1 AND owner = $2 LIMIT 1;

-- name: ListVideos :many
SELECT * FROM videos
WHERE owner = $1
ORDER BY created_at
LIMIT $2
OFFSET $3;


-- name: GetUsersVideosCount :one
SELECT COUNT(*) FROM videos
WHERE owner = $1;

-- name: DeleteVideo :one
DELETE FROM videos
WHERE video_identifier = $1
RETURNING *;
