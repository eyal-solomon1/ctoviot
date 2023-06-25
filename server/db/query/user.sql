-- name: CreateUser :one
INSERT INTO users (
  username,
  hashed_password,
  full_name,
  email,
  balance
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;


-- name: UpdateUserBalance :one
UPDATE users
SET balance = balance - $1
WHERE username = $2
RETURNING *;
