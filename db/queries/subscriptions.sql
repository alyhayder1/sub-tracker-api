-- name: CreateSubscription :one
INSERT INTO subscriptions (user_id, provider_key, name, category, amount, currency, billing_cycle, next_billing_date, status, source)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetSubscription :one
SELECT * FROM subscriptions WHERE id = $1;

-- name: ListSubscriptionsByUser :many
SELECT * FROM subscriptions WHERE user_id = $1 ORDER BY next_billing_date NULLS LAST;

-- name: ListUpcomingRenewals :many
SELECT * FROM subscriptions
WHERE status = 'active'
  AND next_billing_date BETWEEN $1 AND $2
ORDER BY next_billing_date;

-- name: UpdateSubscription :one
UPDATE subscriptions
SET name = $2, category = $3, amount = $4, currency = $5, billing_cycle = $6, next_billing_date = $7, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateSubscriptionStatus :exec
UPDATE subscriptions
SET status = $2, updated_at = now()
WHERE id = $1;

-- name: DeleteSubscription :exec
DELETE FROM subscriptions WHERE id = $1;
