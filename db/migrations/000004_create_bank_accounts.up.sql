CREATE TABLE bank_accounts (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    connection_id       UUID NOT NULL REFERENCES bank_connections(id) ON DELETE CASCADE,
    external_account_id TEXT NOT NULL,
    name                TEXT,
    iban_last4          CHAR(4),
    currency            CHAR(3) NOT NULL DEFAULT 'EUR',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (connection_id, external_account_id)
);