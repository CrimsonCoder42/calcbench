package prob

import (
	"math/big"
	"testing"
)

// Ross §4.8 / Casella–Berger §3.2 (hypergeometric): at most one ace in a 5-card hand →
// (C(48,5) + 4·C(48,4))/C(52,5) ≈ 0.958315633945886 (exact fraction, hand-checked). Tolerance 1e-13.
func TestHypergeometricCDF_sourceExample(t *testing.T) {
	if got, err := HypergeometricCDF(1, 52, 4, 5); err != nil || !within(got, 0.958315633945886, 1e-13) {
		t.Errorf("HypergeometricCDF(1, 52, 4, 5) = (%v, %v), want (0.958315633945886, nil)", got, err)
	}
}

func TestHypergeometricCDF_cases(t *testing.T) {
	tests := []struct {
		name                            string
		k, population, successes, draws int
		want                            float64
	}{
		{"below support", -1, 52, 4, 5, 0},
		{"at the top", 4, 52, 4, 5, 1},
		{"above the top", 9, 52, 4, 5, 1},
		{"no aces", 0, 52, 4, 5, 1712304.0 / 2598960},
		{"single support point", 0, 10, 0, 4, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HypergeometricCDF(tt.k, tt.population, tt.successes, tt.draws)
			if err != nil || !within(got, tt.want, 1e-13) {
				t.Errorf("HypergeometricCDF(%d, %d, %d, %d) = (%v, %v), want %v", tt.k, tt.population, tt.successes, tt.draws, got, err, tt.want)
			}
		})
	}
}

func TestHypergeometricCDF_errors(t *testing.T) {
	_, err := HypergeometricCDF(0, -1, 0, 0)
	requireErr(t, "population < 0", err, ErrOutOfRange)
	_, err = HypergeometricCDF(0, 10, 11, 0)
	requireErr(t, "successes > population", err, ErrOutOfRange)
}

// Oracle: exact big-rational sum for random N ≤ 200. Tolerance 2e-12 relative plus 1e-15 absolute.
func TestHypergeometricCDF_oracle(t *testing.T) {
	rng := newTestRand()
	for trial := range 300 {
		population := 1 + rng.IntN(200)
		successes := rng.IntN(population + 1)
		draws := rng.IntN(population + 1)
		lo, hi := hypergeometricSupport(population, successes, draws)
		k := lo + rng.IntN(hi-lo+1)
		den := new(big.Int).Binomial(int64(population), int64(draws))
		num := new(big.Int)
		for i := lo; i <= k; i++ {
			t1 := new(big.Int).Binomial(int64(successes), int64(i))
			num.Add(num, t1.Mul(t1, new(big.Int).Binomial(int64(population-successes), int64(draws-i))))
		}
		want, _ := new(big.Rat).SetFrac(num, den).Float64()
		got, err := HypergeometricCDF(k, population, successes, draws)
		if err != nil || absDiff(got, want) > 2e-12*want+1e-15 {
			t.Fatalf("trial %d: HypergeometricCDF(%d, %d, %d, %d) = (%v, %v), want %v", trial, k, population, successes, draws, got, err, want)
		}
	}
}

func TestHypergeometricCDF_noAlloc(t *testing.T) {
	requireNoAlloc(t, "HypergeometricCDF", func() { sinkFloat, _ = HypergeometricCDF(1, 52, 4, 5) })
}

func BenchmarkHypergeometricCDF(b *testing.B) {
	for _, size := range []struct {
		name string
		n    int
	}{{"N=1e3", 1_000}, {"N=1e5", 100_000}, {"N=1e7", 10_000_000}} {
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, _ = HypergeometricCDF(size.n/8, size.n, size.n/2, size.n/4)
			}
		})
	}
}
