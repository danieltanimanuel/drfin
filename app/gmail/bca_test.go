package gmail

import "testing"

func TestAccountFromSourceOfFund(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantAccount string
		wantOK      bool
	}{
		{"credit card", "BCA KKI 1234 - **** - **** - **00", "CC BCA", true},
		{"tahapan gold", "TAHAPAN - 0840 - 1234567890", "Tahapan Gold", true},
		{"surrounding whitespace", "  BCA KKI 1234  ", "CC BCA", true},
		{"missing falls back to tahapan gold", "", "Tahapan Gold", true},
		{"whitespace only falls back to tahapan gold", "   ", "Tahapan Gold", true},
		{"unrecognized", "SOMETHING ELSE 9999", "", false},
		{"other tahapan number is unrecognized", "TAHAPAN - 0999 - 1234567890", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := accountFromSourceOfFund(tt.input)
			if got != tt.wantAccount || ok != tt.wantOK {
				t.Errorf("accountFromSourceOfFund(%q) = (%q, %v), want (%q, %v)",
					tt.input, got, ok, tt.wantAccount, tt.wantOK)
			}
		})
	}
}