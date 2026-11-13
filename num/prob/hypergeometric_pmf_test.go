package prob

import (
	"math"
	"math/big"
	"testing"
)

// exactHypergeometricPMF returns C(K,k)·C(N−K,n−k)/C(N,n) in exact rationals, rounded once.
func exactHypergeometricPMF(k, population, successes, draws int) float64 {
	num := new(big.Int).Binomial(int64(successes), int64(k))
	num.Mul(num, new(big.Int).Binomial(int64(population-successes), int64(draws-k)))
	f, _ := new(big.Rat).SetFrac(num, new(big.Int).Binomial(int64(population), int64(draws))).Float64()
	return f
}

// Ross §4.8 / Casella–Berger §3.2 (hypergeometric): one ace in a 5-card hand, N = 52, K = 4, n = 5,
// k = 1 → 4·C(48,4)/C(52,5) = 3243/10829 ≈ 0.2994736356080894 (exact fraction, hand-checked).
// Tolerance 1e-13 (three saddle-point terms).
func TestHypergeometricPMF_sourceExample(t *testing.T) {
	if got, err := HypergeometricPMF(1, 52, 4, 5); err != nil || !within(got, 3243.0/10829, 1e-13) {
		t.Errorf("HypergeometricPMF(1, 52, 4, 5) = (%v, %v), want (3243/10829, nil)", got, err)
	}
}

func TestHypergeometricPMF_cases(t *testing.T) {
	tests := []struct {
		name                            string
		k, population, successes, draws int
		want                            float64
	}{
		{"empty population", 0, 0, 0, 0, 1},
		{"no draws", 0, 10, 3, 0, 1},
		{"draw everything", 3, 10, 3, 10, 1},
		{"no successes", 0, 10, 0, 4, 1},
		{"all successes", 4, 10, 10, 4, 1},
		{"below support", 0, 10, 8, 5, 0},
		{"above support", 4, 10, 3, 5, 0},
		{"four aces", 4, 52, 4, 5, 48.0 / 2598960},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HypergeometricPMF(tt.k, tt.population, tt.successes, tt.draws)
			if err != nil || !within(got, tt.want, 1e-13) {
				t.Errorf("HypergeometricPMF(%d, %d, %d, %d) = (%v, %v), want %v", tt.k, tt.population, tt.successes, tt.draws, got, err, tt.want)
			}
		})
	}
	// Huge population: finite and in [0, 1].
	if got, err := HypergeometricPMF(5, math.MaxInt, math.MaxInt/2, 10); err != nil || got <= 0 || got > 1 {
		t.Errorf("HypergeometricPMF(huge) = (%v, %v), want in (0, 1]", got, err)
	}
}

func TestHypergeometricPMF_errors(t *testing.T) {
	for _, c := range [][3]int{{-1, 0, 0}, {10, 11, 3}, {10, -1, 3}, {10, 3, 11}, {10, 3, -1}} {
		_, err := HypergeometricPMF(0, c[0], c[1], c[2])
		requireErr(t, "HypergeometricPMF", err, ErrOutOfRange)
	}
}

// Oracle: exact big-integer ratio for random N ≤ 300. Tolerance 1e-12 (three saddle-point terms, each as
// BinomialPMF's oracle).
func TestHypergeometricPMF_oracle(t *testing.T) {
	rng := newTestRand()
	for trial := range 400 {
		population := 1 + rng.IntN(300)
		successes := rng.IntN(population + 1)
		draws := rng.IntN(population + 1)
		lo, hi := hypergeometricSupport(population, successes, draws)
		k := lo + rng.IntN(hi-lo+1)
		want := exactHypergeometricPMF(k, population, successes, draws)
		if want < 1e-290 {
			continue
		}
		got, err := HypergeometricPMF(k, population, successes, draws)
		if err != nil || !within(got, want, 1e-12) {
			t.Fatalf("trial %d: HypergeometricPMF(%d, %d, %d, %d) = (%v, %v), want %v", trial, k, population, successes, draws, got, err, want)
		}
	}
}

func TestHypergeometricPMF_noAlloc(t *testing.T) {
	requireNoAlloc(t, "HypergeometricPMF", func() { sinkFloat, _ = HypergeometricPMF(1, 52, 4, 5) })
}

func BenchmarkHypergeometricPMF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = HypergeometricPMF(1, 52, 4, 5)
	}
}
