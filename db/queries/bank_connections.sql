-- name: CreateBankConnection :one
INSERT INTO bank_connections (user_id, provider, external_id, status, consent_expires_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetBankConnection :one
SELECT * FROM bank_connections WHERE id = $1;

-- name: ListBankConnectionsByUser :many
SELECT * FROM bank_connections WHERE user_id = $1 ORDER BY created_at DESC;

-- name: UpdateBankConnectionStatus :exec
UPDATE bank_connections
SET status = $2, updated_at = now()
WHERE id = $1;

-- name: DeleteBankConnection :exec
DELETE FROM bank_connections WHERE id = $1;
