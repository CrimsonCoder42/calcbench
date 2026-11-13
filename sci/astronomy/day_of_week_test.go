package astronomy

import (
	"errors"
	"testing"
	"time"
)

// Meeus, Astronomical Algorithms, 2nd ed., Example 7.e: 1954 June 30 was a Wednesday.
func TestGregorianDayOfWeek_cases(t *testing.T) {
	cases := []struct {
		y    int64
		m, d int
		want time.Weekday
	}{
		{1954, 6, 30, time.Wednesday},
		{2000, 1, 1, time.Saturday},
		{1970, 1, 1, time.Thursday},
		{-4712, 1, 1, time.Thursday}, // JDN 38: (38 + 1) mod 7 = 4 = Thursday; 2440588 − 38 is a multiple of 7, as 1970-01-01
	}
	for _, c := range cases {
		got, err := GregorianDayOfWeek(c.y, c.m, c.d)
		if err != nil || got != c.want {
			t.Errorf("%d-%d-%d: got %v %v, want %v", c.y, c.m, c.d, got, err, c.want)
		}
	}
	if _, err := GregorianDayOfWeek(2023, 2, 29); !errors.Is(err, ErrInvalidDate) {
		t.Errorf("invalid date: %v", err)
	}
}

func BenchmarkGregorianDayOfWeek(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = GregorianDayOfWeek(1954, 6, 30)
	}
}
