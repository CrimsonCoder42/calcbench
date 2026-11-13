package stream

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// pairSeries returns n fixed-seed pairs: y = 20 + 2·N(0, 1), x = 0.5·y + N(0, 1) + 10.
func pairSeries(n int) (xs, ys []float64) {
	rng := rand.New(rand.NewPCG(3, 4))
	xs, ys = make([]float64, n), make([]float64, n)
	for i := range xs {
		ys[i] = 20 + 2*rng.NormFloat64()
		xs[i] = 0.5*ys[i] + rng.NormFloat64() + 10
	}
	return xs, ys
}

// twoPassCo is the oracle: the means, then the centred sums M2x, M2y and C.
func twoPassCo(xs, ys []float64) (m2x, m2y, c float64) {
	mx, _ := twoPass(xs)
	my, _ := twoPass(ys)
	for i := range xs {
		m2x += (xs[i] - mx) * (xs[i] - mx)
		m2y += (ys[i] - my) * (ys[i] - my)
		c += (xs[i] - mx) * (ys[i] - my)
	}
	return m2x, m2y, c
}

func newWC(t *testing.T, w int) *WindowedCovariance {
	t.Helper()
	v, err := NewWindowedCovariance(w)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

// Hand-derived, window 3. (1,2),(2,4),(3,6): means 2, 4; C = 2+0+2 = 4, covariance 2, M2y = 8, beta 0.5,
// correlation 1. Then (4,5) replaces (1,2): window (2,4),(3,6),(4,5), means 3, 5; deviations x −1 0 1,
// y −1 1 0; C = 1, covariance 0.5; M2x = 2, M2y = 2; beta 0.5; correlation 1/√4 = 0.5.
func TestWindowedCovariance_example(t *testing.T) {
	v := newWC(t, 3)
	for _, p := range [][2]float64{{1, 2}, {2, 4}, {3, 6}} {
		_ = v.Add(p[0], p[1])
	}
	check := func(step string, cov, beta, corr, mx, my float64) {
		t.Helper()
		gc, _ := v.Covariance()
		gb, _ := v.Beta()
		gr, _ := v.Correlation()
		gx, _ := v.MeanX()
		gy, _ := v.MeanY()
		// Tolerance 1e-14: a handful of operations on small integers.
		if !closeTo(gc, cov, 1e-14) || !closeTo(gb, beta, 1e-14) || !closeTo(gr, corr, 1e-14) || !closeTo(gx, mx, 1e-14) || !closeTo(gy, my, 1e-14) {
			t.Errorf("%s: cov %v beta %v corr %v means %v %v; want %v %v %v %v %v", step, gc, gb, gr, gx, gy, cov, beta, corr, mx, my)
		}
	}
	check("filled", 2, 0.5, 1, 2, 4)
	if err := v.Add(4, 5); err != nil {
		t.Fatal(err)
	}
	check("slid", 0.5, 0.5, 0.5, 3, 5)
	if v.Len() != 3 || v.WindowLen() != 3 {
		t.Errorf("Len, WindowLen = %d, %d", v.Len(), v.WindowLen())
	}
}

// Oracle: at every step of a 1,000-pair fixed-seed series, against two passes over the same window.
// Tolerance 1e-9 of the covariance's scale √(M2x·M2y)/(n − 1) (the covariance itself can be near 0,
// where a relative bound means nothing), and of the scales √(M2x/M2y) and 1 for beta and r. For w ≤ 3
// the floor is 1e-9 of the series' scale √(σx²·σy²) = √8 ≈ 2.8 and beta and r are not compared: the
// sliding rounding is of the size of the data that has passed through, so a window of two or three
// nearly equal values has a drift-dominated spread (see TestWindowedVariance_oracle and the drift note).
func TestWindowedCovariance_oracle(t *testing.T) {
	xs, ys := pairSeries(1000)
	const seriesScale = 2.83 // √(σx²·σy²): σy² = 4, σx² = 0.25·4 + 1 = 2
	for _, w := range []int{1, 2, 3, 7, 50, 1000, 1001} {
		v := newWC(t, w)
		for i := range xs {
			if err := v.Add(xs[i], ys[i]); err != nil {
				t.Fatal(err)
			}
			lo := max(0, i+1-w)
			m2x, m2y, c := twoPassCo(xs[lo:i+1], ys[lo:i+1])
			n := float64(i + 1 - lo)
			if n < 2 {
				if _, ok := v.Covariance(); ok {
					t.Fatalf("w=%d step %d: Covariance defined for one pair", w, i)
				}
				continue
			}
			gc, _ := v.Covariance()
			scale := math.Sqrt(m2x*m2y) / (n - 1)
			if w <= 3 {
				scale = max(scale, seriesScale)
			}
			if math.Abs(gc-c/(n-1)) > 1e-9*scale {
				t.Fatalf("w=%d step %d: covariance %v, want %v", w, i, gc, c/(n-1))
			}
			if w <= 3 {
				continue
			}
			gb, okb := v.Beta()
			gr, okr := v.Correlation()
			if !okb || !okr || math.Abs(gb-c/m2y) > 1e-9*math.Sqrt(m2x/m2y) || math.Abs(gr-c/math.Sqrt(m2x*m2y)) > 1e-9 {
				t.Fatalf("w=%d step %d: beta %v r %v, want %v %v", w, i, gb, gr, c/m2y, c/math.Sqrt(m2x*m2y))
			}
		}
	}
}

func TestWindowedCovariance_noSpreadAndReset(t *testing.T) {
	v := newWC(t, 3)
	for _, x := range []float64{1, 2, 3, 4} {
		_ = v.Add(x, 7) // y constant
	}
	if _, ok := v.Beta(); ok {
		t.Error("Beta defined with no spread in y")
	}
	if _, ok := v.Correlation(); ok {
		t.Error("Correlation defined with no spread in y")
	}
	if c, ok := v.Covariance(); !ok || c != 0 {
		t.Errorf("Covariance = (%v, %v), want (0, true)", c, ok)
	}
	v.Resync()
	v.Reset()
	if v.Len() != 0 || v.WindowLen() != 3 {
		t.Errorf("after Reset: Len %d WindowLen %d", v.Len(), v.WindowLen())
	}
	var zero WindowedCovariance
	for _, s := range []*WindowedCovariance{&zero, v} {
		for name, f := range map[string]func() (float64, bool){"MeanX": s.MeanX, "MeanY": s.MeanY,
			"Covariance": s.Covariance, "Correlation": s.Correlation, "Beta": s.Beta} {
			if g, ok := f(); ok || g != 0 {
				t.Errorf("empty %s = (%v, %v), want (0, false)", name, g, ok)
			}
		}
	}
	zero.Reset()
	zero.Resync()
	if zero.Len() != 0 || zero.WindowLen() != 0 {
		t.Error("zero value: Len or WindowLen not 0")
	}
}

// Resync after an outlier pair left. Hand-derived: window (1,1),(2,3),(3,2) has means 2, 2;
// deviations x −1 0 1, y −1 1 0; C = 1, covariance 0.5.
func TestWindowedCovariance_resync(t *testing.T) {
	v := newWC(t, 3)
	for _, p := range [][2]float64{{1e12, -1e12}, {1, 1}, {2, 3}, {3, 2}} {
		_ = v.Add(p[0], p[1])
	}
	v.Resync()
	if c, _ := v.Covariance(); !closeTo(c, 0.5, 1e-15) {
		t.Errorf("after Resync: Covariance = %v, want 0.5", c)
	}
}

func TestWindowedCovariance_errors(t *testing.T) {
	for _, w := range []int{0, maxWindowLen + 1} {
		if _, err := NewWindowedCovariance(w); !errors.Is(err, ErrWindowLen) {
			t.Errorf("NewWindowedCovariance(%d) error = %v, want ErrWindowLen", w, err)
		}
	}
	var zero WindowedCovariance
	if err := zero.Add(1, 1); !errors.Is(err, ErrUninitialized) {
		t.Errorf("zero.Add error = %v, want ErrUninitialized", err)
	}
	v := newWC(t, 2)
	_ = v.Add(-1e308, 0)
	for _, p := range [][2]float64{{math.NaN(), 0}, {0, math.Inf(1)}} {
		if err := v.Add(p[0], p[1]); !errors.Is(err, ErrNotFinite) {
			t.Errorf("Add(%v) error = %v, want ErrNotFinite", p, err)
		}
	}
	if err := v.Add(1e308, 0); !errors.Is(err, ErrOverflow) {
		t.Errorf("Add(1e308, 0) error = %v, want ErrOverflow", err)
	}
	if v.Len() != 1 {
		t.Errorf("state changed on error: Len %d", v.Len())
	}
}

func TestWindowedCovariance_noAllocation(t *testing.T) {
	v := newWC(t, 50)
	xs, ys := pairSeries(1000)
	allocs := testing.AllocsPerRun(20, func() {
		for i := range xs {
			_ = v.Add(xs[i], ys[i])
		}
	})
	if allocs != 0 {
		t.Errorf("Add allocates %v times per 1000 pairs, want 0", allocs)
	}
}

func BenchmarkWindowedCovariance(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(size.name, func(b *testing.B) {
			xs, ys := pairSeries(size.n)
			v, _ := NewWindowedCovariance(50)
			b.ReportAllocs()
			for b.Loop() {
				for i := range xs {
					_ = v.Add(xs[i], ys[i])
				}
			}
		})
	}
}
