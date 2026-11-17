package gmail

import (
	"testing"
	"time"
)

func TestParseDateBluIndonesianMonths(t *testing.T) {
	cases := map[string]time.Time{
		"16 Agu 2026 16:42:09 WIB": time.Date(2026, time.August, 16, 16, 42, 9, 0, time.UTC),
		"16 Agt 2026 16:42:09 WIB": time.Date(2026, time.August, 16, 16, 42, 9, 0, time.UTC),
		"05 Okt 2026 10:00:00 WIB": time.Date(2026, time.October, 5, 10, 0, 0, 0, time.UTC),
		"25 Des 2026 23:59:59 WIB": time.Date(2026, time.December, 25, 23, 59, 59, 0, time.UTC),
		"01 Mei 2026 09:00:00 WIB": time.Date(2026, time.May, 1, 9, 0, 0, 0, time.UTC),
		"10 Jan 2026 08:00:00 WIB": time.Date(2026, time.January, 10, 8, 0, 0, 0, time.UTC),
	}
	for in, want := range cases {
		got, err := ParseDate(in, EmailTitleTransaksiBlu)
		if err != nil {
			t.Errorf("ParseDate(%q) error: %v", in, err)
			continue
		}
		if !got.Equal(want) {
			t.Errorf("ParseDate(%q) = %v, want %v", in, got, want)
		}
	}
}