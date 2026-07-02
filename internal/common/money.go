package common

import (
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

// Money wraps pgtype.Numeric so amounts always serialize as a JSON string
// (e.g. "12.99"), never a bare JSON number — a bare number would let a
// weakly-typed client (e.g. JS) silently round-trip it through float64.
type Money struct {
	pgtype.Numeric
}

func (m Money) MarshalJSON() ([]byte, error) {
	if !m.Valid {
		return []byte("null"), nil
	}

	raw, err := m.Numeric.MarshalJSON()
	if err != nil {
		return nil, err
	}
	if len(raw) > 0 && (raw[0] == '"' || string(raw) == "null") {
		return raw, nil
	}

	quoted := make([]byte, 0, len(raw)+2)
	quoted = append(quoted, '"')
	quoted = append(quoted, raw...)
	quoted = append(quoted, '"')
	return quoted, nil
}

func (m *Money) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("money: expected a JSON string amount: %w", err)
	}
	n, err := NewMoney(s)
	if err != nil {
		return err
	}
	m.Numeric = n
	return nil
}

// NewMoney parses a decimal string (e.g. "12.99") into a pgtype.Numeric.
func NewMoney(amount string) (pgtype.Numeric, error) {
	var n pgtype.Numeric
	if err := n.Scan(amount); err != nil {
		return pgtype.Numeric{}, fmt.Errorf("invalid amount %q: %w", amount, err)
	}
	return n, nil
}
