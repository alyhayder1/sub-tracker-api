# sub-tracker-api
REST APIs for subtracker built with Go and Chi. Handles subscriptions, bank feed ingestion, renewal reminders, and spend analytics. PostgreSQL + sqlc + Redis.

## Getting started

```bash
cp .env.example .env
make db-up
make migrate-up
make generate
go run ./cmd/server
```

- `make test` — run tests
- `make lint` — run golangci-lint
- `make dev` — hot reload via air
