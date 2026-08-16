package service

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	tests := []struct {
		name string
		dob  string
		now  string
		want int
	}{
		{"birthday today", "1990-05-10", "2025-05-10", 35},
		{"day before birthday", "1990-05-10", "2025-05-09", 34},
		{"day after birthday", "1990-05-10", "2025-05-11", 35},
		{"born today", "2025-05-10", "2025-05-10", 0},
		{"month before birthday", "1990-05-10", "2025-04-30", 34},
		{"month after birthday", "1990-05-10", "2025-06-01", 35},
		{"leap day birth, Feb 28 of a common year", "2000-02-29", "2025-02-28", 24},
		{"leap day birth, Mar 1 of a common year", "2000-02-29", "2025-03-01", 25},
		{"leap day birth, Feb 29 of a leap year", "2000-02-29", "2024-02-29", 24},
		{"year boundary, Dec 31 to Jan 1", "2000-12-31", "2025-01-01", 24},
		{"day before a July birthday, in a leap year", "1990-07-15", "2024-07-14", 33},
		{"born Mar 1 of a common year, checked Feb 29 of a leap year", "1999-03-01", "2024-02-29", 24},
		{"exact century", "1925-08-16", "2025-08-16", 100},
		{"future date of birth yields a negative age", "2030-01-01", "2025-01-01", -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateAge(parseDate(t, tt.dob), parseDate(t, tt.now))
			if got != tt.want {
				t.Errorf("CalculateAge(%s, %s) = %d, want %d", tt.dob, tt.now, got, tt.want)
			}
		})
	}
}

func parseDate(t *testing.T, s string) time.Time {
	t.Helper()
	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		t.Fatalf("invalid test date %q: %v", s, err)
	}
	return d
}
