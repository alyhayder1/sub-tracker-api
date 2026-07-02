-- name: CreateSubscription :one
INSERT INTO subscriptions (user_id, provider_key, name, category, amount, currency, billing_cycle, next_billing_date, status, source)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetSubscription :one
SELECT * FROM subscriptions WHERE id = $1 AND user_id = $2;

-- name: ListSubscriptionsByUser :many
SELECT * FROM subscriptions WHERE user_id = $1 ORDER BY next_billing_date NULLS LAST;

-- name: ListUpcomingRenewals :many
SELECT * FROM subscriptions
WHERE status = 'active'
  AND next_billing_date BETWEEN $1 AND $2
ORDER BY next_billing_date;

-- name: UpdateSubscription :one
UPDATE subscriptions
SET name = $3, category = $4, amount = $5, currency = $6, billing_cycle = $7, next_billing_date = $8, updated_at = now()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: UpdateSubscriptionStatus :exec
UPDATE subscriptions
SET status = $3, updated_at = now()
WHERE id = $1 AND user_id = $2;

-- name: DeleteSubscription :execrows
DELETE FROM subscriptions WHERE id = $1 AND user_id = $2;
