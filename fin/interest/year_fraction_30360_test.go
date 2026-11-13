package interest

import (
	"testing"
	"time"
)

// date returns midnight UTC on y-m-d.
func date(y int, m time.Month, d int) time.Time { return time.Date(y, m, d, 0, 0, 0, 0, time.UTC) }

func TestYearFraction30360_SourceExample(t *testing.T) {
	// MSRB Rule G-33(e): Number of Days = (Y2 − Y1)360 + (M2 − M1)30 + (D2 − D1); D2 = 31 becomes 30
	// when D1 is 30 or 31; D1 = 31 becomes 30. Hand-checked: 2026-01-31 to 2026-03-31 →
	// D1 = 30, D2 = 30 → 60 days → 60/360.
	got, err := YearFraction30360(date(2026, 1, 31), date(2026, 3, 31))
	wantFloat(t, "G-33 both 31", got, err, 60.0/360, 1e-15)
}

func TestYearFraction30360_Table(t *testing.T) {
	cases := []struct {
		name       string
		start, end time.Time
		days       int
	}{
		{"same day", date(2026, 5, 15), date(2026, 5, 15), 0},
		{"one full year", date(2025, 7, 1), date(2026, 7, 1), 360},
		{"D2 = 31 kept when D1 < 30", date(2026, 1, 15), date(2026, 3, 31), 76},
		{"D2 = 31 cut when D1 = 30", date(2026, 4, 30), date(2026, 5, 31), 30},
		{"D1 = 31 cut to 30", date(2026, 1, 31), date(2026, 2, 15), 15},
		{"February has no end-of-month rule", date(2026, 2, 28), date(2026, 3, 31), 33},
		{"leap day start", date(2024, 2, 29), date(2024, 3, 1), 2},
		{"across years", date(2025, 12, 31), date(2026, 1, 1), 1},
	}
	for _, c := range cases {
		got, err := YearFraction30360(c.start, c.end)
		wantFloat(t, c.name, got, err, float64(c.days)/360, 1e-15)
	}
}

func TestYearFraction30360_LocationAndTimeOfDay(t *testing.T) {
	// Each date is read in its own Location; times of day are ignored.
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("tz database not available")
	}
	start := time.Date(2026, 3, 1, 23, 30, 0, 0, ny) // 2026-03-02 04:30 UTC
	end := time.Date(2026, 4, 1, 0, 15, 0, 0, ny)
	got, err := YearFraction30360(start, end)
	wantFloat(t, "local dates", got, err, 30.0/360, 1e-15)
}

func TestYearFraction30360_Errors(t *testing.T) {
	_, err := YearFraction30360(date(2026, 3, 2), date(2026, 3, 1))
	wantErr(t, "end before start", err, ErrDomain)
	_, err = YearFraction30360(date(2026, 1, 31), date(2026, 1, 30))
	wantErr(t, "end before start, both cut to 30", err, ErrDomain)
}

func TestYearFraction30360_Oracle(t *testing.T) {
	// Additivity holds when no date is the 31st: f(a,c) = f(a,b) + f(b,c).
	rng := testRNG()
	for k := 0; k < 300; k++ {
		a := date(2000+rng.IntN(50), time.Month(1+rng.IntN(12)), 1+rng.IntN(28))
		b := a.AddDate(0, rng.IntN(40), rng.IntN(28))
		c := b.AddDate(0, rng.IntN(40), rng.IntN(28))
		if a.Day() == 31 || b.Day() == 31 || c.Day() == 31 {
			continue
		}
		ab, err1 := YearFraction30360(a, b)
		bc, err2 := YearFraction30360(b, c)
		ac, err3 := YearFraction30360(a, c)
		if err1 != nil || err2 != nil || err3 != nil {
			t.Fatal(err1, err2, err3)
		}
		wantFloat(t, "additive", ac, nil, ab+bc, 1e-12)
	}
}

func TestYearFraction30360_NoAllocs(t *testing.T) {
	s, e := date(2026, 1, 31), date(2026, 3, 31)
	wantNoAllocs(t, "YearFraction30360", func() { sinkFloat, _ = YearFraction30360(s, e) })
}

func BenchmarkYearFraction30360(b *testing.B) {
	s, e := date(2026, 1, 31), date(2026, 3, 31)
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = YearFraction30360(s, e)
	}
}
