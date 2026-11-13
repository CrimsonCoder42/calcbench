package stream

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// relClose compares within tol·|want|.
func relClose(got, want, tol float64) bool {
	return math.Abs(got-want) <= tol*math.Abs(want)
}

// normalSeries returns n fixed-seed values 50 + 10·N(0, 1).
func normalSeries(n int) []float64 {
	rng := rand.New(rand.NewPCG(1, 2))
	s := make([]float64, n)
	for i := range s {
		s[i] = 50 + 10*rng.NormFloat64()
	}
	return s
}

func newWV(t *testing.T, w int) *WindowedVariance {
	t.Helper()
	v, err := NewWindowedVariance(w)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

// Hand-derived, window 3 over 1, 2, 3, 4, 10:
// {1}: mean 1, population variance 0, sample undefined. {1,2}: mean 1.5, sample 0.5.
// {1,2,3}: mean 2, sample 1. {2,3,4}: mean 3, sample 1.
// {3,4,10}: mean 17/3, deviations −8/3, −5/3, 13/3, M2 = (64+25+169)/9 = 258/9, sample 129/9, population 86/9.
func TestWindowedVariance_example(t *testing.T) {
	v := newWV(t, 3)
	want := []struct{ mean, sample float64 }{{1, math.NaN()}, {1.5, 0.5}, {2, 1}, {3, 1}, {17.0 / 3, 129.0 / 9}}
	for i, x := range []float64{1, 2, 3, 4, 10} {
		if err := v.Add(x); err != nil {
			t.Fatal(err)
		}
		m, _ := v.Mean()
		s, ok := v.Variance()
		// Tolerance 1e-14: a handful of operations on small values.
		if !closeTo(m, want[i].mean, 1e-14) || ok == math.IsNaN(want[i].sample) || (ok && !closeTo(s, want[i].sample, 1e-14)) {
			t.Errorf("step %d: mean %v var (%v, %v), want %v %v", i, m, s, ok, want[i].mean, want[i].sample)
		}
	}
	if pv, _ := v.PopulationVariance(); !closeTo(pv, 86.0/9, 1e-14) {
		t.Errorf("PopulationVariance = %v, want 86/9", pv)
	}
	if sd, _ := v.Std(); !closeTo(sd, math.Sqrt(129.0/9), 1e-14) {
		t.Errorf("Std = %v, want √(129/9)", sd)
	}
	if v.Len() != 3 || v.WindowLen() != 3 {
		t.Errorf("Len, WindowLen = %d, %d, want 3, 3", v.Len(), v.WindowLen())
	}
}

// Oracle: at every step of a 1,000-value fixed-seed series, the sliding update against a direct two-pass
// sample variance of the same window, computed locally by twoPass (algo cannot import num/stats; this is
// the computation of NUM-STATS-021 SampleVariance: the mean, then the centred sum of squares).
//
// Tolerance: 1e-9 relative for windows of 7 and more (a simulation of this update over 20 such series
// measured at most 4.1e-13 at w = 7 and 1.5e-14 at w = 50). For w ≤ 3 the bound is 1e-9 of the series'
// variance (σ² = 100) instead: the sliding update's rounding is of the size of the data that has passed
// through, not of the window's own variance, and two or three values that happen to agree to 1e-3 have
// a variance near 1e-6 with that absolute error on it (measured up to 3e-4 relative at w = 2). This is
// the drift the WindowedVariance doc comment states.
func TestWindowedVariance_oracle(t *testing.T) {
	values := normalSeries(1000)
	const seriesVariance = 100 // 10², the variance normalSeries draws from
	for _, w := range []int{1, 2, 3, 7, 50, 999, 1000, 1001} {
		floor := 0.0
		if w <= 3 {
			floor = 1e-9 * seriesVariance
		}
		within := func(got, want float64) bool { return math.Abs(got-want) <= max(1e-9*math.Abs(want), floor) }
		v := newWV(t, w)
		for i, x := range values {
			if err := v.Add(x); err != nil {
				t.Fatal(err)
			}
			win := values[max(0, i+1-w) : i+1]
			mean, m2 := twoPass(win)
			n := float64(len(win))
			gm, _ := v.Mean()
			gp, _ := v.PopulationVariance()
			if !relClose(gm, mean, 1e-9) || !within(gp, m2/n) {
				t.Fatalf("w=%d step %d: mean %v pop var %v, want %v %v", w, i, gm, gp, mean, m2/n)
			}
			if len(win) >= 2 {
				gs, ok := v.Variance()
				if !ok || !within(gs, m2/(n-1)) {
					t.Fatalf("w=%d step %d: sample var (%v, %v), want %v", w, i, gs, ok, m2/(n-1))
				}
			}
		}
	}
}

// After a 1e12 outlier leaves a window of small values, Resync restores the exact variance
// (hand-derived: {1, 2, 3} has sample variance 1).
func TestWindowedVariance_resync(t *testing.T) {
	v := newWV(t, 3)
	for _, x := range []float64{1e12, 1, 2, 3} {
		if err := v.Add(x); err != nil {
			t.Fatal(err)
		}
	}
	v.Resync()
	if s, _ := v.Variance(); !closeTo(s, 1, 1e-15) {
		t.Errorf("after Resync: Variance = %v, want 1", s)
	}
	if m, _ := v.Mean(); m != 2 {
		t.Errorf("after Resync: Mean = %v, want 2", m)
	}
	var zero WindowedVariance
	zero.Resync() // no-op, no panic
}

func TestWindowedVariance_emptyAndReset(t *testing.T) {
	var zero WindowedVariance
	v := newWV(t, 4)
	for _, s := range []*WindowedVariance{&zero, v} {
		for name, f := range map[string]func() (float64, bool){"Mean": s.Mean, "Variance": s.Variance,
			"PopulationVariance": s.PopulationVariance, "Std": s.Std} {
			if g, ok := f(); ok || g != 0 {
				t.Errorf("empty %s = (%v, %v), want (0, false)", name, g, ok)
			}
		}
	}
	if zero.Len() != 0 || zero.WindowLen() != 0 {
		t.Error("zero value: Len or WindowLen not 0")
	}
	zero.Reset()
	for _, x := range []float64{5, 6, 7} {
		_ = v.Add(x)
	}
	v.Reset()
	if v.Len() != 0 || v.WindowLen() != 4 {
		t.Errorf("after Reset: Len %d WindowLen %d, want 0 4", v.Len(), v.WindowLen())
	}
	_ = v.Add(-3.5)
	if m, _ := v.Mean(); m != -3.5 {
		t.Errorf("after Reset and one Add: Mean = %v, want -3.5", m)
	}
	// Constant window: variance exactly 0.
	for range 6 {
		_ = v.Add(2.25)
	}
	if s, _ := v.Variance(); s != 0 {
		t.Errorf("constant window: Variance = %v, want 0", s)
	}
}

func TestWindowedVariance_errors(t *testing.T) {
	for _, w := range []int{0, -1, maxWindowLen + 1} {
		if _, err := NewWindowedVariance(w); !errors.Is(err, ErrWindowLen) {
			t.Errorf("NewWindowedVariance(%d) error = %v, want ErrWindowLen", w, err)
		}
	}
	var zero WindowedVariance
	if err := zero.Add(1); !errors.Is(err, ErrUninitialized) {
		t.Errorf("zero.Add error = %v, want ErrUninitialized", err)
	}
	v := newWV(t, 2)
	_ = v.Add(-1e308)
	for _, x := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if err := v.Add(x); !errors.Is(err, ErrNotFinite) {
			t.Errorf("Add(%v) error = %v, want ErrNotFinite", x, err)
		}
	}
	// 1e308 − (−1e308) overflows.
	if err := v.Add(1e308); !errors.Is(err, ErrOverflow) {
		t.Errorf("Add(1e308) error = %v, want ErrOverflow", err)
	}
	if m, _ := v.Mean(); v.Len() != 1 || m != -1e308 {
		t.Errorf("state changed on error: Len %d Mean %v", v.Len(), m)
	}
}

func TestWindowedVariance_noAllocation(t *testing.T) {
	v := newWV(t, 50)
	values := benchFloats(1000)
	allocs := testing.AllocsPerRun(20, func() {
		for _, x := range values {
			_ = v.Add(x)
		}
	})
	if allocs != 0 {
		t.Errorf("Add allocates %v times per 1000 values, want 0", allocs)
	}
}

func BenchmarkWindowedVariance(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(size.name, func(b *testing.B) {
			values := benchFloats(size.n)
			v, _ := NewWindowedVariance(50)
			b.ReportAllocs()
			for b.Loop() {
				for _, x := range values {
					_ = v.Add(x)
				}
			}
		})
	}
}
