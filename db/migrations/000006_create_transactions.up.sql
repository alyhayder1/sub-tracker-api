CREATE TABLE transactions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    account_id          UUID REFERENCES bank_accounts(id) ON DELETE SET NULL,
    subscription_id     UUID REFERENCES subscriptions(id) ON DELETE SET NULL,
    external_id         TEXT NOT NULL,
    amount              NUMERIC(12,2) NOT NULL,
    currency            CHAR(3) NOT NULL DEFAULT 'EUR',
    description         TEXT,
    normalized_merchant TEXT,
    booked_at           DATE NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (account_id, external_id)
);