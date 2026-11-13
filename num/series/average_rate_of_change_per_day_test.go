package series

import (
	"errors"
	"math"
	"testing"
	"time"
)

var sinkRate float64

// OpenStax Calculus Vol. 1 §3.4, (f(a + h) − f(a))/h, by hand: 50 → 53 over 30 days = 0.1 per day.
func TestAverageRateOfChangePerDay_example(t *testing.T) {
	t0 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	got, err := AverageRateOfChangePerDay(Point{t0, 50}, Point{t0.AddDate(0, 0, 30), 53})
	if err != nil || !closeRel(got, 0.1, 1e-15) {
		t.Errorf("= (%v, %v), want 0.1", got, err)
	}
}

// parity: ledger C8 (velocity.go:32): RatePerDay = (last − first)/deltaDays, rounded by the source to 4
// decimals. Case: 100 → 92.5 over 15 days = −0.5 per day. The source's 7-day floor and 30.44-day month
// are not ported (a 3-day span is computed here; the source returned nil).
func TestAverageRateOfChangePerDay_parity(t *testing.T) {
	t0 := time.Date(2026, 3, 1, 8, 0, 0, 0, time.UTC)
	got, err := AverageRateOfChangePerDay(Point{t0, 100}, Point{t0.Add(15 * 24 * time.Hour), 92.5})
	if err != nil || !closeRel(got, -0.5, 1e-15) {
		t.Errorf("= (%v, %v), want −0.5", got, err)
	}
	got, err = AverageRateOfChangePerDay(Point{t0, 10}, Point{t0.Add(72 * time.Hour), 16})
	if err != nil || !closeRel(got, 2, 1e-15) {
		t.Errorf("3-day span: (%v, %v), want 2", got, err)
	}
}

func TestAverageRateOfChangePerDay_cases(t *testing.T) {
	t0 := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name     string
		from, to Point
		want     float64
	}{
		{"half a day", Point{t0, 0}, Point{t0.Add(12 * time.Hour), 1}, 2},
		{"no change", Point{t0, 7}, Point{t0.Add(time.Hour), 7}, 0},
		{"one nanosecond", Point{t0, 0}, Point{t0.Add(time.Nanosecond), 1}, 86400e9},
		{"500 years (beyond Duration)", Point{t0, 0}, Point{t0.AddDate(500, 0, 0), 182622}, 1}, // 500·365 + 122 leap days,
		{"zones do not matter", Point{t0, 0}, Point{t0.Add(48 * time.Hour).In(time.FixedZone("x", 5*3600)), 4}, 2},
	}
	for _, tc := range tests {
		got, err := AverageRateOfChangePerDay(tc.from, tc.to)
		if err != nil || !closeRel(got, tc.want, 1e-12) {
			t.Errorf("%s: = (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

func TestAverageRateOfChangePerDay_errors(t *testing.T) {
	t0 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	t1 := t0.Add(24 * time.Hour)
	for _, v := range []float64{math.NaN(), math.Inf(1)} {
		if _, err := AverageRateOfChangePerDay(Point{t0, v}, Point{t1, 1}); !errors.Is(err, ErrNotFinite) {
			t.Errorf("from %v: %v", v, err)
		}
		if _, err := AverageRateOfChangePerDay(Point{t0, 1}, Point{t1, v}); !errors.Is(err, ErrNotFinite) {
			t.Errorf("to %v: %v", v, err)
		}
	}
	if _, err := AverageRateOfChangePerDay(Point{t1, 1}, Point{t0, 2}); !errors.Is(err, ErrUnsorted) {
		t.Errorf("reversed: %v", err)
	}
	if _, err := AverageRateOfChangePerDay(Point{t0, 1}, Point{t0, 2}); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("same instant: %v", err)
	}
	if _, err := AverageRateOfChangePerDay(Point{t0, -math.MaxFloat64}, Point{t1, math.MaxFloat64}); !errors.Is(err, ErrOverflow) {
		t.Errorf("change overflow: %v", err)
	}
	if _, err := AverageRateOfChangePerDay(Point{t0, 0}, Point{t0.Add(time.Nanosecond), math.MaxFloat64}); !errors.Is(err, ErrOverflow) {
		t.Errorf("rate overflow: %v", err)
	}
}

func TestAverageRateOfChangePerDay_noAlloc(t *testing.T) {
	t0 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	a, b := Point{t0, 50}, Point{t0.Add(30 * 24 * time.Hour), 53}
	if n := testing.AllocsPerRun(100, func() { sinkRate, _ = AverageRateOfChangePerDay(a, b) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkAverageRateOfChangePerDay(b *testing.B) {
	t0 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	p, q := Point{t0, 50}, Point{t0.Add(30 * 24 * time.Hour), 53}
	b.ReportAllocs()
	for b.Loop() {
		sinkRate, _ = AverageRateOfChangePerDay(p, q)
	}
}
