package accounts

import (
	"errors"
	"testing"
)

func TestNormalizeAmount(t *testing.T) {
	ok := map[string]string{
		"0": "0.00", "5": "5.00", "12.5": "12.50", "-2500000.75": "-2500000.75",
		" 10 ": "10.00", "-0": "0.00",
	}
	for in, want := range ok {
		got, err := NormalizeAmount(in)
		if err != nil || got != want {
			t.Errorf("NormalizeAmount(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
	for _, in := range []string{"", "abc", "1.234", "1,5", "1e3", "--1", ".5", "1.", "12345678901234567890123"} {
		if _, err := NormalizeAmount(in); !errors.Is(err, ErrInvalidInput) {
			t.Errorf("NormalizeAmount(%q) should fail with ErrInvalidInput, got %v", in, err)
		}
	}
}

func TestValidateCreateAccount(t *testing.T) {
	name, bal, limit, err := validateCreateAccount(CreateAccountRequest{Name: "  CC Mandiri ", AccountGroupID: 3})
	if err != nil || name != "CC Mandiri" || bal != "0.00" || limit != nil {
		t.Fatalf("unexpected result: %q %q %v %v", name, bal, limit, err)
	}
	bad := []CreateAccountRequest{
		{Name: "", AccountGroupID: 1},
		{Name: "x", AccountGroupID: 0},
		{Name: "x", AccountGroupID: 1, InitialBalance: "1.234"},
		{Name: "x", AccountGroupID: 1, CreditLimit: "-5"},
	}
	for _, r := range bad {
		if _, _, _, err := validateCreateAccount(r); !errors.Is(err, ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput for %+v, got %v", r, err)
		}
	}
}