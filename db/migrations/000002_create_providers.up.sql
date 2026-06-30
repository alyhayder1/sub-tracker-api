CREATE TABLE providers (
    key               TEXT PRIMARY KEY,
    display_name      TEXT NOT NULL,
    category          TEXT NOT NULL,
    logo_url          TEXT,
    merchant_patterns TEXT[] NOT NULL DEFAULT '{}'
);