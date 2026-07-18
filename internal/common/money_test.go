package common

import "testing"

func TestMoneyMarshalJSON(t *testing.T) {
	amount, err := NewMoney("12.99")
	if err != nil {
		t.Fatalf("NewMoney: %v", err)
	}
	m := Money{Numeric: amount}

	b, err := m.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON: %v", err)
	}
	if got, want := string(b), `"12.99"`; got != want {
		t.Errorf("MarshalJSON() = %s, want %s", got, want)
	}
}

func TestMoneyUnmarshalJSON(t *testing.T) {
	var m Money
	if err := m.UnmarshalJSON([]byte(`"12.99"`)); err != nil {
		t.Fatalf("UnmarshalJSON: %v", err)
	}

	b, err := m.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON: %v", err)
	}
	if got, want := string(b), `"12.99"`; got != want {
		t.Errorf("round-trip = %s, want %s", got, want)
	}
}

func TestMoneyMarshalJSON_Null(t *testing.T) {
	var m Money
	b, err := m.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON: %v", err)
	}
	if got, want := string(b), "null"; got != want {
		t.Errorf("MarshalJSON() = %s, want %s", got, want)
	}
}
