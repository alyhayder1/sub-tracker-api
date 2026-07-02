package common

import (
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Timestamp wraps pgtype.Timestamptz so it serializes as a plain RFC 3339
// string instead of leaking pgtype's internal struct fields via default
// JSON struct encoding.
type Timestamp struct {
	pgtype.Timestamptz
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(t.Time.Format(time.RFC3339))
}
