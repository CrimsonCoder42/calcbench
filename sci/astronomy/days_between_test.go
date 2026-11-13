package astronomy

import (
	"errors"
	"testing"
)

// Meeus, Astronomical Algorithms, 2nd ed., Example 7.d: Halley's Comet, 1910 April 20 to 1986 February 9,
// 27689 days.
func TestDaysBetweenDates_cases(t *testing.T) {
	cases := []struct {
		name   string
		fy     int64
		fm, fd int
		ty     int64
		tm, td int
		want   int64
	}{
		{"meeus 7.d", 1910, 4, 20, 1986, 2, 9, 27689},
		{"reversed", 1986, 2, 9, 1910, 4, 20, -27689},
		{"same day", 2024, 2, 29, 2024, 2, 29, 0},
		{"across leap day", 2024, 2, 28, 2024, 3, 1, 2},
		{"across non-leap century", 1900, 2, 28, 1900, 3, 1, 1},
		{"one leap year", 2024, 1, 1, 2025, 1, 1, 366},
	}
	for _, c := range cases {
		got, err := DaysBetweenDates(c.fy, c.fm, c.fd, c.ty, c.tm, c.td)
		if err != nil || got != c.want {
			t.Errorf("%s: got %d %v, want %d", c.name, got, err, c.want)
		}
	}
	if _, err := DaysBetweenDates(2023, 2, 29, 2024, 1, 1); !errors.Is(err, ErrInvalidDate) {
		t.Errorf("bad from: %v", err)
	}
	if _, err := DaysBetweenDates(2023, 1, 1, 2000000, 1, 1); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("bad to: %v", err)
	}
}

func BenchmarkDaysBetweenDates(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = DaysBetweenDates(1910, 4, 20, 1986, 2, 9)
	}
}
