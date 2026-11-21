package gmail

import (
	"testing"
	"time"
)

func TestBuildPocketTransferDTO(t *testing.T) {
	date := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

	tests := []struct {
		name       string
		in         map[string]string
		wantSrc    string
		wantDest   string
		wantAmount string
	}{
		{
			name: "main to pocket",
			in: map[string]string{
				"Beneficiary Pocket": "Mam",
				"Source of Fund":     "0000****00",
				"Amount":             "IDR 2,488,000.00",
			},
			wantSrc: "Tahapan Gold", wantDest: "MyBCAPocket Mam", wantAmount: "2488000.00",
		},
		{
			name: "pocket to pocket",
			in: map[string]string{
				"Beneficiary Pocket": "Entertainment",
				"Source Pocket":      "Mam",
				"Amount":             "IDR 1,000.00",
			},
			wantSrc: "MyBCAPocket Mam", wantDest: "MyBCAPocket Entertainment", wantAmount: "1000.00",
		},
		{
			name: "pocket withdrawal to main account",
			in: map[string]string{
				"Beneficiary Account": "0000000000",
				"Source Pocket":       "Transport",
				"Amount":              "IDR 139,000.00",
			},
			wantSrc: "MyBCAPocket Transport", wantDest: "Tahapan Gold", wantAmount: "139000.00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto, err := buildPocketTransferDTO(tt.in, date)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dto.Type != "Transfer" || dto.Status != "Waiting" || dto.Source != "Email" {
				t.Errorf("unexpected type/status/source: %q %q %q", dto.Type, dto.Status, dto.Source)
			}
			if dto.Account != tt.wantSrc || dto.AccountDest != tt.wantDest {
				t.Errorf("accounts = %q -> %q, want %q -> %q", dto.Account, dto.AccountDest, tt.wantSrc, tt.wantDest)
			}
			if got := dto.Amount.Text('f', 2); got != tt.wantAmount {
				t.Errorf("amount = %s, want %s", got, tt.wantAmount)
			}
			if !dto.TradeDate.Equal(date) {
				t.Errorf("trade date not passed through")
			}
		})
	}
}

func TestBuildPocketTransferDTOErrors(t *testing.T) {
	date := time.Now()

	tests := []struct {
		name string
		in   map[string]string
	}{
		{"missing amount", map[string]string{
			"Beneficiary Pocket": "Mam", "Source of Fund": "x"}},
		{"short amount", map[string]string{
			"Beneficiary Pocket": "Mam", "Source of Fund": "x", "Amount": "IDR"}},
		{"invalid amount", map[string]string{
			"Beneficiary Pocket": "Mam", "Source of Fund": "x", "Amount": "IDR abcdef"}},
		{"empty pocket name", map[string]string{
			"Beneficiary Pocket": "  ", "Source of Fund": "x", "Amount": "IDR 1,000.00"}},
		{"same source and destination", map[string]string{
			"Beneficiary Pocket": "Mam", "Source Pocket": "Mam", "Amount": "IDR 1,000.00"}},
		{"unrecognized layout", map[string]string{
			"Amount": "IDR 1,000.00"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if dto, err := buildPocketTransferDTO(tt.in, date); err == nil {
				t.Fatalf("expected error, got dto %+v", dto)
			}
		})
	}
}
