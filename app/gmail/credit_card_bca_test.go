package gmail

import "testing"

func TestIsSkippedCurrency(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want bool
	}{
		{"MYR", "MYR 99,47", true},
		{"MYR lowercase", "myr 10,00", true},
		{"MYR leading space", "  MYR 5,00", true},
		{"IDR is not skipped", "IDR 150.000,00", false},
		{"other foreign currency is not skipped", "USD 12,50", false},
		{"empty", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSkippedCurrency(tt.in); got != tt.want {
				t.Errorf("isSkippedCurrency(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}