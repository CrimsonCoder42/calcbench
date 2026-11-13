package series

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"testing"
	"time"
)

func newYork(t testing.TB) *time.Location {
	t.Helper()
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	return loc
}

func TestTheilSenSlopePerDay_example(t *testing.T) {
	// Doc example, hand-checked: slopes 2, 0.5, −1 → 0.5 per day.
	t0 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	pts := []Point{{t0, 10}, {t0.Add(24 * time.Hour), 12}, {t0.Add(48 * time.Hour), 11}}
	got, err := TheilSenSlopePerDay(pts, newRand())
	if err != nil || got != 0.5 {
		t.Fatalf("got %v, %v; want 0.5", got, err)
	}
}

// parity: ledger TheilSenPerDay, source case 1 (irregular readings at days 0 … 13, slope 0.4722/day, ±5e-5 as
// printed; unchanged when the last reading becomes 85).
func TestTheilSenSlopePerDay_parity(t *testing.T) {
	loc := newYork(t)
	t0 := time.Date(2026, 10, 1, 8, 0, 0, 0, loc)
	days := []float64{0, 0.25, 0.5, 1, 1.25, 1.5, 7, 13}
	vals := []float64{93, 94, 93, 94, 95, 95, 97, 97}
	pts := make([]Point, len(days))
	for i := range days {
		pts[i] = Point{At: t0.Add(time.Duration(days[i] * 24 * float64(time.Hour))), Value: vals[i]}
	}
	got, err := TheilSenSlopePerDay(pts, newRand())
	if err != nil || math.Abs(got-0.4722) > 5e-5 {
		t.Fatalf("got %v, %v; want 0.4722", got, err)
	}
	pts[7].Value = 85
	got, err = TheilSenSlopePerDay(pts, newRand())
	if err != nil || math.Abs(got-0.4722) > 5e-5 {
		t.Fatalf("with artefact: got %v, %v; want 0.4722", got, err)
	}
}

func TestTheilSenSlopePerDay_cases(t *testing.T) {
	t0 := time.Date(2026, 3, 7, 12, 0, 0, 0, time.UTC)
	// Same instant twice: no slope between them; the other two pairs give 1 and 3 per day → mean 2.
	pts := []Point{{t0, 0}, {t0, 2}, {t0.Add(24 * time.Hour), 3}}
	if got, err := TheilSenSlopePerDay(pts, newRand()); err != nil || got != 2 {
		t.Errorf("same instant: got %v, %v; want 2", got, err)
	}
	// Elapsed time, not calendar days: across the 2026-03-08 US DST start, 23 h is 23/24 day.
	loc := newYork(t)
	a := time.Date(2026, 3, 8, 0, 0, 0, 0, loc)
	b := time.Date(2026, 3, 9, 0, 0, 0, 0, loc)
	got, err := TheilSenSlopePerDay([]Point{{a, 0}, {b, 23}}, newRand())
	if err != nil || math.Abs(got-24) > 1e-12 {
		t.Errorf("DST day: got %v, %v; want 24 (23 units over 23/24 day)", got, err)
	}
	// Far apart instants (beyond time.Sub's ±292 years) still give the right elapsed days.
	far := []Point{{time.Date(1500, 1, 1, 0, 0, 0, 0, time.UTC), 0}, {time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC), 1}}
	wantDays := float64(far[1].At.Unix()-far[0].At.Unix()) / 86400
	if got, err := TheilSenSlopePerDay(far, newRand()); err != nil || math.Abs(got-1/wantDays) > 1e-18 {
		t.Errorf("far apart: got %v, %v; want %v", got, err, 1/wantDays)
	}
}

func TestTheilSenSlopePerDay_errors(t *testing.T) {
	t0 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		name string
		pts  []Point
		rng  *rand.Rand
		want error
	}{
		{"nil rng", []Point{{t0, 1}, {t0.Add(time.Hour), 2}}, nil, ErrNilRandom},
		{"empty", nil, newRand(), ErrTooFew},
		{"one point", []Point{{t0, 1}}, newRand(), ErrTooFew},
		{"one instant", []Point{{t0, 1}, {t0, 2}}, newRand(), ErrTooFew},
		{"NaN", []Point{{t0, math.NaN()}, {t0.Add(time.Hour), 2}}, newRand(), ErrNotFinite},
		{"Inf", []Point{{t0, 1}, {t0.Add(time.Hour), math.Inf(-1)}}, newRand(), ErrNotFinite},
	}
	for _, c := range cases {
		if _, err := TheilSenSlopePerDay(c.pts, c.rng); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func benchPoints(n int) []Point {
	r := rand.New(rand.NewPCG(1, 2))
	t0 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	pts := make([]Point, n)
	for i := range pts {
		pts[i] = Point{At: t0.Add(time.Duration(i) * 7 * time.Minute), Value: 100 + r.NormFloat64()*10}
	}
	return pts
}

func BenchmarkTheilSenSlopePerDay(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		pts := benchPoints(n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			r := newRand()
			for b.Loop() {
				if _, err := TheilSenSlopePerDay(pts, r); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
