package series

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
	"time"
)

var sinkHourMean float64

// Doc example, by hand with the NIST §1.3.5.1 mean: 140 at 06:30, 150 at 11:59, 120 at 18:00 New York time;
// [5, 12) → (140 + 150)/2 = 145, count 2; [17, 24) → 120, count 1.
func TestLocalHourWindowMean_example(t *testing.T) {
	ny := newYork(t)
	d := func(h, m int) time.Time { return time.Date(2026, 6, 10, h, m, 0, 0, ny) }
	pts := []Point{{d(6, 30), 140}, {d(11, 59), 150}, {d(18, 0), 120}}
	if m, c, err := LocalHourWindowMean(pts, ny, 5, 12); err != nil || m != 145 || c != 2 {
		t.Errorf("[5, 12): (%v, %v, %v), want 145, 2", m, c, err)
	}
	if m, c, err := LocalHourWindowMean(pts, ny, 17, 24); err != nil || m != 120 || c != 1 {
		t.Errorf("[17, 24): (%v, %v, %v), want 120, 1", m, c, err)
	}
}

// parity: ledger C5 (diurnal.go:34): morning 05:00–11:59 and evening 17:00–23:59 means and the morning −
// evening difference. Source case built from its algorithm (no test file): systolic 138, 142 at 07:00 and
// 08:00; 128, 124 at 19:00 and 21:00; 133 at 14:00 (in neither window) → morning 140, evening 126, difference
// 14 (the source's > 20 "surge" flag and "non-dipping" flag are app rules, not ported).
func TestLocalHourWindowMean_parity(t *testing.T) {
	utc := time.UTC
	d := func(h int) time.Time { return time.Date(2026, 2, 3, h, 0, 0, 0, utc) }
	pts := []Point{{d(7), 138}, {d(8), 142}, {d(14), 133}, {d(19), 128}, {d(21), 124}}
	am, nAM, err1 := LocalHourWindowMean(pts, utc, 5, 12)
	pm, nPM, err2 := LocalHourWindowMean(pts, utc, 17, 24)
	if err1 != nil || err2 != nil || am != 140 || pm != 126 || nAM != 2 || nPM != 2 || am-pm != 14 {
		t.Errorf("morning (%v, %v, %v), evening (%v, %v, %v), want 140/2, 126/2, difference 14", am, nAM, err1, pm, nPM, err2)
	}
}

func TestLocalHourWindowMean_cases(t *testing.T) {
	utc := time.UTC
	d := func(day, h, m int) time.Time { return time.Date(2026, 1, day, h, m, 0, 0, utc) }
	pts := []Point{{d(2, 23, 0), 1}, {d(1, 0, 0), 2}, {d(1, 5, 59), 3}, {d(1, 6, 0), 4}, {d(1, 21, 59), 5}, {d(1, 22, 0), 6}}
	tests := []struct {
		name     string
		from, to int
		mean     float64
		count    int
	}{
		{"wraps midnight [22, 6), unsorted points", 22, 6, (1 + 2 + 3 + 6) / 4.0, 4},
		{"[6, 22)", 6, 22, 4.5, 2},
		{"whole day but one hour [0, 23)", 0, 23, (2 + 3 + 4 + 5 + 6) / 5.0, 5},
		{"[23, 24)", 23, 24, 1, 1},
		{"[0, 1)", 0, 1, 2, 1},
	}
	for _, tc := range tests {
		m, c, err := LocalHourWindowMean(pts, utc, tc.from, tc.to)
		if err != nil || !closeRel(m, tc.mean, 1e-15) || c != tc.count {
			t.Errorf("%s: (%v, %v, %v), want %v, %v", tc.name, m, c, err, tc.mean, tc.count)
		}
	}
}

// DST: on 2026-03-08 New York springs forward at 02:00; 01:30 EST and 03:30 EDT are one hour apart. The
// window [3, 4) holds only the 03:30 point, and the same instant read in UTC (07:30) is in [7, 8).
func TestLocalHourWindowMean_dst(t *testing.T) {
	ny := newYork(t)
	a := time.Date(2026, 3, 8, 1, 30, 0, 0, ny)
	b := a.Add(time.Hour) // 03:30 EDT
	pts := []Point{{a, 10}, {b, 20}}
	if m, c, err := LocalHourWindowMean(pts, ny, 3, 4); err != nil || m != 20 || c != 1 {
		t.Errorf("NY [3, 4): (%v, %v, %v), want 20, 1", m, c, err)
	}
	if m, c, err := LocalHourWindowMean(pts, time.UTC, 7, 8); err != nil || m != 20 || c != 1 {
		t.Errorf("UTC [7, 8): (%v, %v, %v), want 20, 1", m, c, err)
	}
}

// Overflow path: values near MaxFloat64 in the window still give a finite mean.
func TestLocalHourWindowMean_large(t *testing.T) {
	t0 := time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC)
	pts := []Point{{t0, math.MaxFloat64}, {t0, math.MaxFloat64}, {t0, math.MaxFloat64 / 2}}
	m, c, err := LocalHourWindowMean(pts, time.UTC, 9, 10)
	if want := math.MaxFloat64 / 6 * 5; err != nil || c != 3 || !closeRel(m, want, 1e-15) {
		t.Errorf("(%v, %v, %v), want %v", m, c, err, want)
	}
}

// Oracle: a filter-then-average loop on 200 random series (fixed seed, n ≤ 400, random windows).
// Tolerance 1e-12 relative.
func TestLocalHourWindowMean_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	ny := newYork(t)
	t0 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for range 200 {
		n := 1 + rng.IntN(400)
		pts := make([]Point, n)
		for i := range pts {
			pts[i] = Point{t0.Add(time.Duration(rng.Int64N(int64(365 * 24 * time.Hour)))), 100 + 20*rng.NormFloat64()}
		}
		from, to := rng.IntN(24), rng.IntN(25)
		if from == to {
			continue
		}
		sum, cnt := 0.0, 0
		for _, p := range pts {
			h := p.At.In(ny).Hour()
			if (from < to && h >= from && h < to) || (from > to && (h >= from || h < to)) {
				sum += p.Value
				cnt++
			}
		}
		m, c, err := LocalHourWindowMean(pts, ny, from, to)
		if cnt == 0 {
			if !errors.Is(err, ErrTooFew) {
				t.Fatalf("empty window: %v", err)
			}
			continue
		}
		if err != nil || c != cnt || !closeRel(m, sum/float64(cnt), 1e-12) {
			t.Fatalf("[%d, %d): (%v, %v, %v), want %v, %v", from, to, m, c, err, sum/float64(cnt), cnt)
		}
	}
}

func TestLocalHourWindowMean_errors(t *testing.T) {
	t0 := time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC)
	pts := []Point{{t0, 1}}
	if _, _, err := LocalHourWindowMean(pts, nil, 5, 12); !errors.Is(err, ErrNilLocation) {
		t.Errorf("nil loc: %v", err)
	}
	for _, h := range [][2]int{{-1, 5}, {24, 5}, {5, 25}, {5, -1}, {7, 7}} {
		if _, _, err := LocalHourWindowMean(pts, time.UTC, h[0], h[1]); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("hours %v: %v", h, err)
		}
	}
	if _, _, err := LocalHourWindowMean(pts, time.UTC, 10, 12); !errors.Is(err, ErrTooFew) {
		t.Errorf("empty window: %v", err)
	}
	if _, _, err := LocalHourWindowMean(nil, time.UTC, 0, 24); !errors.Is(err, ErrTooFew) {
		t.Errorf("no points: %v", err)
	}
	bad := []Point{{t0, math.NaN()}}
	if _, _, err := LocalHourWindowMean(bad, time.UTC, 9, 10); !errors.Is(err, ErrNotFinite) {
		t.Errorf("NaN in window: %v", err)
	}
	if _, _, err := LocalHourWindowMean(bad, time.UTC, 10, 11); !errors.Is(err, ErrTooFew) {
		t.Errorf("NaN outside window is not read: %v", err)
	}
}

func TestLocalHourWindowMean_noAlloc(t *testing.T) {
	t0 := time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC)
	pts := []Point{{t0, 1}, {t0.Add(time.Hour), 2}}
	if a := testing.AllocsPerRun(100, func() { sinkHourMean, _, _ = LocalHourWindowMean(pts, time.UTC, 5, 12) }); a != 0 {
		t.Errorf("allocates %v times, want 0", a)
	}
}

func BenchmarkLocalHourWindowMean(b *testing.B) {
	ny := newYork(b)
	t0 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for _, n := range []int{1000, 10000, 100000} {
		pts := make([]Point, n)
		for i := range pts {
			pts[i] = Point{t0.Add(time.Duration(i) * 17 * time.Minute), float64(100 + i%40)}
		}
		b.Run(seriesSizeName(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkHourMean, _, _ = LocalHourWindowMean(pts, ny, 5, 12)
			}
		})
	}
}
