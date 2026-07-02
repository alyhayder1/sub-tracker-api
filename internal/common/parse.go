package common

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

// ParseUUID parses a URL path/query UUID string into a pgtype.UUID.
func ParseUUID(s string) (pgtype.UUID, error) {
	var id pgtype.UUID
	if err := id.Scan(s); err != nil {
		return pgtype.UUID{}, fmt.Errorf("invalid id %q: %w", s, err)
	}
	return id, nil
}
