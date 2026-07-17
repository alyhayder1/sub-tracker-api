package common

import "testing"

func TestParseUUID(t *testing.T) {
	id, err := ParseUUID("11111111-1111-1111-1111-111111111111")
	if err != nil {
		t.Fatalf("ParseUUID: %v", err)
	}
	if !id.Valid {
		t.Error("expected parsed UUID to be valid")
	}
}

func TestParseUUID_Invalid(t *testing.T) {
	if _, err := ParseUUID("not-a-uuid"); err == nil {
		t.Error("expected an error for an invalid UUID string")
	}
}
