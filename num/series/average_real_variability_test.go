package series

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

var sinkARV float64

// Mena et al. 2017 (JAHA), Introduction eq. 1, applied by hand: 120, 130, 125, 135 → |10| + |−5| + |10| = 25,
// N − 1 = 3 → 25/3. Tolerance 1e-15 relative.
func TestAverageRealVariability_sourceExample(t *testing.T) {
	got, err := AverageRealVariability([]float64{120, 130, 125, 135})
	if err != nil || !closeRel(got, 25.0/3, 1e-15) {
		t.Errorf("= (%v, %v), want 25/3", got, err)
	}
}

// parity: ledger ARV (moments.go:108), C7 (variability.go:42) and bp_trend.go (:435): the same formula
// Σ|xᵢ₊₁ − xᵢ|/(n − 1) in all three. Cases written from the formula (the sources have no numeric ARV test;
// series_test.go only checks the empty error): 118, 126, 121, 131, 124 → (8 + 5 + 10 + 7)/4 = 7.5; the
// sources round to 2 (variability.go) or 1 (bp_trend.go) decimals, which is not ported.
func TestAverageRealVariability_parity(t *testing.T) {
	got, err := AverageRealVariability([]float64{118, 126, 121, 131, 124})
	if err != nil || got != 7.5 {
		t.Errorf("= (%v, %v), want 7.5", got, err)
	}
	// source: ARV(Series{}) returns ErrTooFew{Need: 2}
	if _, err := AverageRealVariability(nil); !errors.Is(err, ErrTooFew) {
		t.Errorf("empty: %v, want ErrTooFew", err)
	}
}

func TestAverageRealVariability_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   float64
	}{
		{"two", []float64{3, 7}, 4},
		{"constant", []float64{5, 5, 5}, 0},
		{"monotone", []float64{1, 2, 3, 4, 5}, 1},
		{"alternating", []float64{0, 10, 0, 10}, 10},
		{"negative", []float64{-1, -4, -2}, 2.5},
		{"order matters (sorted copy differs)", []float64{1, 3, 2}, 1.5},
	}
	for _, tc := range tests {
		got, err := AverageRealVariability(tc.values)
		if err != nil || !closeRel(got, tc.want, 1e-15) {
			t.Errorf("%s: = (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
	// Differences that overflow while the mean does not: ±MaxFloat64/2 alternating, 3 values → mean |d| = MaxFloat64.
	h := math.MaxFloat64 / 2
	if got, err := AverageRealVariability([]float64{-h, h, -h}); err != nil || !closeRel(got, math.MaxFloat64, 1e-15) {
		t.Errorf("large alternating: (%v, %v), want MaxFloat64", got, err)
	}
}

// Oracle: the direct loop in float64 on 300 random sequences (fixed seed, n ≤ 500). Tolerance 1e-12 relative.
func TestAverageRealVariability_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		n := 2 + rng.IntN(499)
		v := make([]float64, n)
		for i := range v {
			v[i] = 120 + 15*rng.NormFloat64()
		}
		want := 0.0
		for i := 1; i < n; i++ {
			want += math.Abs(v[i] - v[i-1])
		}
		want /= float64(n - 1)
		got, err := AverageRealVariability(v)
		if err != nil || !closeRel(got, want, 1e-12) {
			t.Fatalf("n %d: (%v, %v), want %v", n, got, err, want)
		}
	}
}

func TestAverageRealVariability_errors(t *testing.T) {
	if _, err := AverageRealVariability([]float64{1}); !errors.Is(err, ErrTooFew) {
		t.Errorf("one value: %v", err)
	}
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := AverageRealVariability([]float64{1, v}); !errors.Is(err, ErrNotFinite) {
			t.Errorf("%v: %v", v, err)
		}
	}
	if _, err := AverageRealVariability([]float64{-math.MaxFloat64, math.MaxFloat64}); !errors.Is(err, ErrOverflow) {
		t.Errorf("±MaxFloat64: %v, want ErrOverflow", err)
	}
}

func TestAverageRealVariability_noAlloc(t *testing.T) {
	v := []float64{120, 130, 125, 135}
	if a := testing.AllocsPerRun(100, func() { sinkARV, _ = AverageRealVariability(v) }); a != 0 {
		t.Errorf("allocates %v times, want 0", a)
	}
}

func BenchmarkAverageRealVariability(b *testing.B) {
	for _, n := range []int{1000, 10000, 100000} {
		rng := rand.New(rand.NewPCG(1, 2))
		v := make([]float64, n)
		for i := range v {
			v[i] = 120 + 15*rng.NormFloat64()
		}
		b.Run(seriesSizeName(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkARV, _ = AverageRealVariability(v)
			}
		})
	}
}

// seriesSizeName names a benchmark size 1e3, 1e4 or 1e5.
func seriesSizeName(n int) string {
	switch n {
	case 1000:
		return "1e3"
	case 10000:
		return "1e4"
	}
	return "1e5"
}
