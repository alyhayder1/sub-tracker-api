# Sub-tracker API

Go REST API for Sub-tracker, a cross-platform subscription manager.
EU-first (Tink / PSD2). Part of a 5-repo product: api, web, ios, android, infra.

## Stack
- Go + Chi router
- sqlc (type-safe SQL, pgx/v5) — NEVER hand-write the generated layer
- golang-migrate for migrations
- PostgreSQL 16 + Redis 7 (via docker-compose)
- Auth: Clerk. Bank: Tink. Push: FCM + APNs.

## Conventions
- Money is always NUMERIC(12,2) — never floats
- Store iban_last4 only, never full IBANs
- ON DELETE CASCADE from users (GDPR erasure)
- Migrations are .up.sql/.down.sql pairs in db/migrations/
- Queries go in db/queries/, generated code in db/generated/ (don't edit)

## Workflow
- Branch off `dev`, PR into `dev`. `main` is the protected release gate.
- Branch naming: feature/api-*, fix/api-*
- Run: `make db-up`, `make migrate-up`, `make generate`, `make test`

## Status
- Foundation + schema migrations 000001–000008 done (PR'd to dev)
- Next: db/queries/*.sql + `make generate` for the typed Go layer