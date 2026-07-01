-- name: CreateBankAccount :one
INSERT INTO bank_accounts (connection_id, external_account_id, name, iban_last4, currency)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetBankAccount :one
SELECT * FROM bank_accounts WHERE id = $1;

-- name: ListBankAccountsByConnection :many
SELECT * FROM bank_accounts WHERE connection_id = $1 ORDER BY created_at;

-- name: DeleteBankAccount :exec
DELETE FROM bank_accounts WHERE id = $1;
