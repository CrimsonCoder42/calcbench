package interest

import (
	"testing"
	"time"
)

func TestYearFractionActual365_SourceExample(t *testing.T) {
	// Kirk, Contemporary Mathematics (OpenStax 2023) §6.3 Example 6.30 uses 45 actual days over a
	// 365-day year ("Actual/365"); with SimpleInterest it gives the book's $81.12 before its round-up.
	yf, err := YearFractionActual365(date(2026, 1, 1), date(2026, 2, 15))
	wantFloat(t, "45 days", yf, err, 45.0/365, 1e-15)
	got, err := SimpleInterest(560000, 0.1175, yf)
	wantCents(t, "Example 6.30 interest", got, err, 8112)
}

func TestYearFractionActual365_Table(t *testing.T) {
	cases := []struct {
		name       string
		start, end time.Time
		days       int
	}{
		{"same day", date(2026, 5, 15), date(2026, 5, 15), 0},
		{"common year", date(2025, 1, 1), date(2026, 1, 1), 365},
		{"leap year has 366 days over 365", date(2024, 1, 1), date(2025, 1, 1), 366},
		{"across February in a leap year", date(2024, 2, 28), date(2024, 3, 1), 2},
		{"before 1970", date(1969, 12, 31), date(1970, 1, 1), 1},
		{"long span", date(1900, 1, 1), date(2000, 1, 1), 36524},
	}
	for _, c := range cases {
		got, err := YearFractionActual365(c.start, c.end)
		wantFloat(t, c.name, got, err, float64(c.days)/365, 1e-15)
	}
}

func TestYearFractionActual365_DaylightSaving(t *testing.T) {
	// 2026-03-08 is a 23-hour day in New York; the count is still calendar days.
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("tz database not available")
	}
	got, err := YearFractionActual365(time.Date(2026, 3, 7, 12, 0, 0, 0, ny), time.Date(2026, 3, 9, 0, 30, 0, 0, ny))
	wantFloat(t, "across DST", got, err, 2.0/365, 1e-15)
}

func TestYearFractionActual365_Errors(t *testing.T) {
	_, err := YearFractionActual365(date(2026, 3, 2), date(2026, 3, 1))
	wantErr(t, "end before start", err, ErrDomain)
}

func TestYearFractionActual365_Oracle(t *testing.T) {
	// Against time.Sub between UTC midnights, which has no daylight-saving shifts.
	rng := testRNG()
	for k := 0; k < 300; k++ {
		a := date(1800+rng.IntN(400), time.Month(1+rng.IntN(12)), 1+rng.IntN(28))
		b := a.AddDate(0, 0, rng.IntN(20000))
		got, err := YearFractionActual365(a, b)
		days := b.Sub(a).Hours() / 24
		wantFloat(t, "oracle", got, err, days/365, 1e-12)
	}
}

func TestYearFractionActual365_NoAllocs(t *testing.T) {
	s, e := date(2026, 1, 1), date(2026, 2, 15)
	wantNoAllocs(t, "YearFractionActual365", func() { sinkFloat, _ = YearFractionActual365(s, e) })
}

func BenchmarkYearFractionActual365(b *testing.B) {
	s, e := date(2026, 1, 1), date(2026, 2, 15)
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = YearFractionActual365(s, e)
	}
}
