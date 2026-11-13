package astronomy

import (
	"errors"
	"math"
	"testing"
)

// Meeus, Astronomical Algorithms, 2nd ed., Example 7.d read backwards: 1910-04-20 + 27689 days = 1986-02-09.
func TestGregorianAddDays_cases(t *testing.T) {
	cases := []struct {
		name   string
		y      int64
		m, d   int
		days   int64
		wy     int64
		wm, wd int
	}{
		{"meeus 7.d", 1910, 4, 20, 27689, 1986, 2, 9},
		{"backwards", 1986, 2, 9, -27689, 1910, 4, 20},
		{"zero", 2024, 2, 29, 0, 2024, 2, 29},
		{"into leap day", 2024, 2, 28, 1, 2024, 2, 29},
		{"year end", 2023, 12, 31, 1, 2024, 1, 1},
	}
	for _, c := range cases {
		y, m, d, err := GregorianAddDays(c.y, c.m, c.d, c.days)
		if err != nil || y != c.wy || m != c.wm || d != c.wd {
			t.Errorf("%s: got %d-%d-%d %v", c.name, y, m, d, err)
		}
	}
	for _, days := range []int64{math.MaxInt64, math.MinInt64, 400_000_000, -2_500_000} {
		if _, _, _, err := GregorianAddDays(2000, 1, 1, days); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("days %d: err %v", days, err)
		}
	}
	if _, _, _, err := GregorianAddDays(2000, 13, 1, 1); !errors.Is(err, ErrInvalidDate) {
		t.Errorf("bad start: %v", err)
	}
}

func BenchmarkGregorianAddDays(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _, _, _ = GregorianAddDays(1910, 4, 20, 27689)
	}
}
