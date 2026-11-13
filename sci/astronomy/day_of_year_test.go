package astronomy

import (
	"errors"
	"testing"
	"time"
)

// Meeus, Astronomical Algorithms, 2nd ed., Example 7.f: 1978 November 14 is day 318; 1988 April 22 is day 113.
func TestGregorianDayOfYear_cases(t *testing.T) {
	cases := []struct {
		y    int64
		m, d int
		want int
	}{
		{1978, 11, 14, 318},
		{1988, 4, 22, 113},
		{2023, 1, 1, 1},
		{2023, 12, 31, 365},
		{2024, 12, 31, 366},
		{2024, 3, 1, 61},
		{1900, 3, 1, 60},
	}
	for _, c := range cases {
		got, err := GregorianDayOfYear(c.y, c.m, c.d)
		if err != nil || got != c.want {
			t.Errorf("%d-%d-%d: got %d %v, want %d", c.y, c.m, c.d, got, err, c.want)
		}
	}
	if _, err := GregorianDayOfYear(2024, 2, 30); !errors.Is(err, ErrInvalidDate) {
		t.Errorf("invalid: %v", err)
	}
}

// Oracle: time.Time.YearDay over every day of 1999..2001 and 2099..2101.
func TestGregorianDayOfYear_oracle(t *testing.T) {
	for _, start := range []int{1999, 2099} {
		tm := time.Date(start, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(start+3, 1, 1, 0, 0, 0, 0, time.UTC)
		for ; tm.Before(end); tm = tm.AddDate(0, 0, 1) {
			got, err := GregorianDayOfYear(int64(tm.Year()), int(tm.Month()), tm.Day())
			if err != nil || got != tm.YearDay() {
				t.Fatalf("%v: got %d %v, want %d", tm, got, err, tm.YearDay())
			}
		}
	}
}

func BenchmarkGregorianDayOfYear(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = GregorianDayOfYear(1978, 11, 14)
	}
}
