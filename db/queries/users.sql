-- name: CreateUser :one
INSERT INTO users (clerk_id, email, full_name, default_currency)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByClerkID :one
SELECT * FROM users WHERE clerk_id = $1;

-- name: UpdateUser :one
UPDATE users
SET full_name = $2, default_currency = $3, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
