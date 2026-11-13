package prob

import (
	"errors"
	"math"
	"testing"
)

// NIST §1.3.6.6.19 (the CDF is the sum of the PMF): λ = 2, k = 2 → 5e⁻² = 0.6766764161830635; λ = 10,
// k = 15 → 0.9512595966960213 (exact rational sum times e⁻¹⁰ in 50-digit arithmetic). Tolerance 1e-14.
func TestPoissonCDF_sourceExample(t *testing.T) {
	for _, c := range []struct {
		k      int
		lambda float64
		want   float64
	}{{2, 2, 0.6766764161830635}, {15, 10, 0.9512595966960213}} {
		if got, err := PoissonCDF(c.k, c.lambda); err != nil || !within(got, c.want, 1e-14) {
			t.Errorf("PoissonCDF(%d, %v) = (%v, %v), want (%v, nil)", c.k, c.lambda, got, err, c.want)
		}
	}
}

func TestPoissonCDF_cases(t *testing.T) {
	tests := []struct {
		name   string
		k      int
		lambda float64
		want   float64
	}{
		{"k < 0", -1, 2, 0},
		{"lambda = 0", 0, 0, 1},
		{"k = 0", 0, 2, math.Exp(-2)},
		{"far upper", 1000, 2, 1},
		{"k = MaxInt", math.MaxInt, 2, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := PoissonCDF(tt.k, tt.lambda); err != nil || !within(got, tt.want, 1e-14) {
				t.Errorf("PoissonCDF(%d, %v) = (%v, %v), want %v", tt.k, tt.lambda, got, err, tt.want)
			}
		})
	}
}

func TestPoissonCDF_errors(t *testing.T) {
	_, err := PoissonCDF(1, -1)
	requireErr(t, "lambda < 0", err, ErrOutOfRange)
	_, err = PoissonCDF(1, nan)
	requireErr(t, "NaN", err, ErrNotFinite)
	// λ = 1e18 at k = λ: σ = 1e9, more than 2²⁸ terms.
	if _, err = PoissonCDF(1e18, 1e18); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("PoissonCDF(huge λ) err = %v, want ErrOutOfRange", err)
	}
}

// Oracle: the compensated sum of PoissonPMF over 0…k. Tolerance 1e-12 relative plus 1e-15 absolute.
func TestPoissonCDF_oracle(t *testing.T) {
	rng := newTestRand()
	for trial := range 300 {
		lambda := 0.01 + 100*rng.Float64()
		k := rng.IntN(250)
		var s compensatedSum
		for i := 0; i <= k; i++ {
			v, _ := PoissonPMF(i, lambda)
			s.add(v)
		}
		want := min(1, s.total())
		if got, err := PoissonCDF(k, lambda); err != nil || absDiff(got, want) > 1e-12*want+1e-15 {
			t.Fatalf("trial %d: PoissonCDF(%d, %v) = (%v, %v), want %v", trial, k, lambda, got, err, want)
		}
	}
}

func TestPoissonCDF_noAlloc(t *testing.T) {
	requireNoAlloc(t, "PoissonCDF", func() { sinkFloat, _ = PoissonCDF(15, 10) })
}

func BenchmarkPoissonCDF(b *testing.B) {
	for _, size := range []struct {
		name   string
		lambda float64
	}{{"λ=1e3", 1e3}, {"λ=1e5", 1e5}, {"λ=1e7", 1e7}} {
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, _ = PoissonCDF(int(size.lambda), size.lambda)
			}
		})
	}
}
