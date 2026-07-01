-- name: CreateTransaction :one
INSERT INTO transactions (user_id, account_id, subscription_id, external_id, amount, currency, description, normalized_merchant, booked_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetTransaction :one
SELECT * FROM transactions WHERE id = $1;

-- name: ListTransactionsByUser :many
SELECT * FROM transactions
WHERE user_id = $1
ORDER BY booked_at DESC
LIMIT $2 OFFSET $3;

-- name: ListTransactionsBySubscription :many
SELECT * FROM transactions WHERE subscription_id = $1 ORDER BY booked_at DESC;

-- name: LinkTransactionToSubscription :exec
UPDATE transactions
SET subscription_id = $2
WHERE id = $1;
