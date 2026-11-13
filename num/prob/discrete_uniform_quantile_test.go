package prob

import (
	"math"
	"testing"
)

// Inverting the discrete uniform CDF of Casella–Berger §3.2: a fair die, u = 0.5 → 3 (F(3) = 0.5).
// Hand-checked.
func TestDiscreteUniformQuantile_sourceExample(t *testing.T) {
	if got, err := DiscreteUniformQuantile(0.5, 1, 6); err != nil || got != 3 {
		t.Errorf("DiscreteUniformQuantile(0.5, 1, 6) = (%v, %v), want (3, nil)", got, err)
	}
}

func TestDiscreteUniformQuantile_cases(t *testing.T) {
	tests := []struct {
		name string
		u    float64
		a, b int
		want int
	}{
		{"u = 0", 0, 1, 6, 1},
		{"u = 1", 1, 1, 6, 6},
		{"one point", 0.7, 4, 4, 4},
		{"just above a step", 0.5000001, 1, 6, 4},
		{"negative range", 0.25, -3, 0, -3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := DiscreteUniformQuantile(tt.u, tt.a, tt.b); err != nil || got != tt.want {
				t.Errorf("DiscreteUniformQuantile(%v, %d, %d) = (%v, %v), want %v", tt.u, tt.a, tt.b, got, err, tt.want)
			}
		})
	}
}

// The widest range: the result satisfies the defining property against the CDF as computed (near the
// top, float64 cannot tell some of the last CDF values from 1).
func TestDiscreteUniformQuantile_widestRange(t *testing.T) {
	for _, u := range []float64{0, 0.5, 1} {
		k, err := DiscreteUniformQuantile(u, 1, math.MaxInt)
		if err != nil {
			t.Fatal(err)
		}
		hi, _ := DiscreteUniformCDF(k, 1, math.MaxInt)
		lo, _ := DiscreteUniformCDF(k-1, 1, math.MaxInt)
		if hi < u || (k > 1 && lo >= u) {
			t.Errorf("DiscreteUniformQuantile(%v, 1, MaxInt) = %d: F(k−1) = %v, F(k) = %v", u, k, lo, hi)
		}
	}
}

// The full int range (m = 2⁶⁴): u = 0 → MinInt; u = ½ → −513, the least k whose float64 CDF
// (2⁶³ − 512)/2⁶⁴ rounds to ½ (k = −514 gives 2⁶³ − 1024 in the numerator, below ½); and the defining
// property F(k−1) < u ≤ F(k), for this package's CDF, over a spread of u. Before the fix: ErrOverflow.
func TestDiscreteUniformQuantile_fullIntRange(t *testing.T) {
	for _, c := range []struct {
		u    float64
		want int
	}{{0.5, -513}, {0, math.MinInt}} {
		if got, err := DiscreteUniformQuantile(c.u, math.MinInt, math.MaxInt); err != nil || got != c.want {
			t.Errorf("DiscreteUniformQuantile(%v, MinInt, MaxInt) = (%v, %v), want (%d, nil)", c.u, got, err, c.want)
		}
	}
	for _, u := range []float64{1e-19, 0x1p-64, 0.1, 0.3, 0.5, 0.7, 0.999, 1 - 0x1p-53, 1} {
		k, err := DiscreteUniformQuantile(u, math.MinInt, math.MaxInt)
		if err != nil {
			t.Fatalf("DiscreteUniformQuantile(%v, MinInt, MaxInt) error %v", u, err)
		}
		hi, _ := DiscreteUniformCDF(k, math.MinInt, math.MaxInt)
		lo, _ := DiscreteUniformCDF(k-1, math.MinInt, math.MaxInt)
		if hi < u || (k > math.MinInt && lo >= u) {
			t.Errorf("DiscreteUniformQuantile(%v, MinInt, MaxInt) = %d: F(k−1) = %v, F(k) = %v", u, k, lo, hi)
		}
	}
}

func TestDiscreteUniformQuantile_errors(t *testing.T) {
	_, err := DiscreteUniformQuantile(0.5, 2, 1)
	requireErr(t, "a > b", err, ErrOutOfRange)
	_, err = DiscreteUniformQuantile(nan, 1, 6)
	requireErr(t, "NaN", err, ErrNotFinite)
	_, err = DiscreteUniformQuantile(-0.1, 1, 6)
	requireErr(t, "u < 0", err, ErrProbability)
}

// Oracle: brute force over the support, the least k with DiscreteUniformCDF(k) ≥ u.
func TestDiscreteUniformQuantile_oracle(t *testing.T) {
	rng := newTestRand()
	for trial := range 500 {
		a := rng.IntN(200) - 100
		b := a + rng.IntN(50)
		u := rng.Float64()
		want := b
		for k := a; k <= b; k++ {
			if f, _ := DiscreteUniformCDF(k, a, b); f >= u {
				want = k
				break
			}
		}
		if got, err := DiscreteUniformQuantile(u, a, b); err != nil || got != want {
			t.Fatalf("trial %d: DiscreteUniformQuantile(%v, %d, %d) = (%v, %v), want %d", trial, u, a, b, got, err, want)
		}
	}
}

func TestDiscreteUniformQuantile_noAlloc(t *testing.T) {
	requireNoAlloc(t, "DiscreteUniformQuantile", func() { sinkInt, _ = DiscreteUniformQuantile(0.5, 1, 6) })
}

func BenchmarkDiscreteUniformQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = DiscreteUniformQuantile(0.5, 1, 6)
	}
}
