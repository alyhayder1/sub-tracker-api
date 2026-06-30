-include .env
export

MIGRATE = migrate -path db/migrations -database "$(DATABASE_URL)"

.PHONY: db-up db-down dev generate migrate-up migrate-down migrate-create test lint

db-up:        ; docker compose up -d
db-down:      ; docker compose down
dev:          ; air
generate:     ; sqlc generate
migrate-up:   ; $(MIGRATE) up
migrate-down: ; $(MIGRATE) down 1
migrate-create: ; migrate create -ext sql -dir db/migrations -seq $(name)
test:         ; go test ./...
lint:         ; golangci-lint run