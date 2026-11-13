package prob

import (
	"math/big"
	"math/rand/v2"
	"testing"
)

// Ross §3.3 / Casella–Berger §1.2: two urns chosen with probability 0.5 each, red with 0.3
// and 0.6 → 0.45. Hand-checked; tolerance 1 ulp (compensated sum of exact-error products).
func TestTotalProbability_sourceExample(t *testing.T) {
	if got, err := TotalProbability([]float64{0.3, 0.6}, []float64{0.5, 0.5}); err != nil || !within(got, 0.45, 2.5e-16) {
		t.Errorf("TotalProbability = (%v, %v), want (0.45, nil)", got, err)
	}
}

func TestTotalProbability_cases(t *testing.T) {
	tests := []struct {
		name string
		cond []float64
		pB   []float64
		want float64
	}{
		{"one block", []float64{0.3}, []float64{1}, 0.3},
		{"certain", []float64{1, 1, 1}, []float64{0.25, 0.25, 0.5}, 1},
		{"impossible", []float64{0, 0}, []float64{0.5, 0.5}, 0},
		{"rounded decimals sum to 1", []float64{1, 1, 1}, []float64{0.1, 0.2, 0.7}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := TotalProbability(tt.cond, tt.pB); err != nil || !within(got, tt.want, 2.5e-16) {
				t.Errorf("TotalProbability(%v, %v) = (%v, %v), want %v", tt.cond, tt.pB, got, err, tt.want)
			}
		})
	}
}

func TestTotalProbability_errors(t *testing.T) {
	_, err := TotalProbability(nil, nil)
	requireErr(t, "empty", err, ErrOutOfRange)
	_, err = TotalProbability([]float64{0.5}, []float64{0.5, 0.5})
	requireErr(t, "length mismatch", err, ErrOutOfRange)
	_, err = TotalProbability([]float64{nan, 0.5}, []float64{0.5, 0.5})
	requireErr(t, "NaN", err, ErrNotFinite)
	_, err = TotalProbability([]float64{0.5, 0.5}, []float64{inf, 0.5})
	requireErr(t, "Inf in pB", err, ErrNotFinite)
	_, err = TotalProbability([]float64{1.2, 0.5}, []float64{0.5, 0.5})
	requireErr(t, "conditional > 1", err, ErrProbability)
	_, err = TotalProbability([]float64{0.5, 0.5}, []float64{0.5, 0.49})
	requireErr(t, "pB sums to 0.99", err, ErrNotDistribution)
}

// Oracle: the exact rational Σ cᵢ·pᵢ of the float64 inputs (big.Rat), with pB normalised to sum to 1
// exactly in rationals and then rounded. Tolerance 2 ulps: the products are exact and the sum compensated.
func TestTotalProbability_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := 1 + rng.IntN(50)
		cond, pB := make([]float64, n), make([]float64, n)
		for i := range n {
			cond[i] = rng.Float64()
			pB[i] = 1 / float64(n)
		}
		exact := new(big.Rat)
		for i := range n {
			exact.Add(exact, new(big.Rat).Mul(new(big.Rat).SetFloat64(cond[i]), new(big.Rat).SetFloat64(pB[i])))
		}
		want, _ := exact.Float64()
		got, err := TotalProbability(cond, pB)
		if err != nil || !within(got, want, 5e-16) {
			t.Fatalf("trial %d: TotalProbability = (%v, %v), want %v", trial, got, err, want)
		}
	}
}

func TestTotalProbability_noAlloc(t *testing.T) {
	cond, pB := []float64{0.3, 0.6}, []float64{0.5, 0.5}
	requireNoAlloc(t, "TotalProbability", func() { sinkFloat, _ = TotalProbability(cond, pB) })
}

func BenchmarkTotalProbability(b *testing.B) {
	for _, size := range []struct {
		name string
		n    int
	}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}} {
		cond, pB := make([]float64, size.n), make([]float64, size.n)
		for i := range size.n {
			cond[i], pB[i] = 0.5, 1/float64(size.n)
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, _ = TotalProbability(cond, pB)
			}
		})
	}
}
