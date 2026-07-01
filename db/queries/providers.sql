-- name: ListProviders :many
SELECT * FROM providers ORDER BY display_name;

-- name: GetProvider :one
SELECT * FROM providers WHERE key = $1;

-- name: UpsertProvider :one
INSERT INTO providers (key, display_name, category, logo_url, merchant_patterns)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (key) DO UPDATE
SET display_name = EXCLUDED.display_name,
    category = EXCLUDED.category,
    logo_url = EXCLUDED.logo_url,
    merchant_patterns = EXCLUDED.merchant_patterns
RETURNING *;
