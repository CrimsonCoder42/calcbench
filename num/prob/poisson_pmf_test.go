package prob

import (
	"math"
	"testing"
)

// NIST §1.3.6.6.19 / Ross §4.7 (e^(−λ)λᵏ/k!): λ = 2, k = 0 → e⁻² = 0.1353352832366127; λ = 1000,
// k = 1000 → 0.0126146113487215 (computed for this test with 50-digit decimal arithmetic:
// exp(−1000 + 1000·ln 1000 − ln 1000!) = 0.012614611348721499718…; the naive λᵏ/k! overflows).
// Tolerance 1e-14 (Loader's form is accurate to a few ulps).
func TestPoissonPMF_sourceExample(t *testing.T) {
	for _, c := range []struct {
		k      int
		lambda float64
		want   float64
	}{{0, 2, 0.1353352832366127}, {1000, 1000, 0.012614611348721500}} {
		if got, err := PoissonPMF(c.k, c.lambda); err != nil || !within(got, c.want, 1e-14) {
			t.Errorf("PoissonPMF(%d, %v) = (%v, %v), want (%v, nil)", c.k, c.lambda, got, err, c.want)
		}
	}
}

func TestPoissonPMF_cases(t *testing.T) {
	tests := []struct {
		name   string
		k      int
		lambda float64
		want   float64
	}{
		{"k < 0", -1, 2, 0},
		{"lambda = 0 at 0", 0, 0, 1},
		{"lambda = 0 elsewhere", 3, 0, 0},
		{"k = 1", 1, 2, 2 * math.Exp(-2)},
		{"k = 2", 2, 2, 2 * math.Exp(-2)},
		{"far tail underflows", 0, 1e4, 0},
		{"small lambda", 1, 1e-3, 1e-3 * math.Exp(-1e-3)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := PoissonPMF(tt.k, tt.lambda); err != nil || !within(got, tt.want, 1e-14) {
				t.Errorf("PoissonPMF(%d, %v) = (%v, %v), want %v", tt.k, tt.lambda, got, err, tt.want)
			}
		})
	}
}

func TestPoissonPMF_errors(t *testing.T) {
	_, err := PoissonPMF(1, -1)
	requireErr(t, "lambda < 0", err, ErrOutOfRange)
	_, err = PoissonPMF(1, nan)
	requireErr(t, "NaN", err, ErrNotFinite)
	_, err = PoissonPMF(1, inf)
	requireErr(t, "Inf", err, ErrNotFinite)
}

// Oracle: the product form e^(−λ)·Πᵢ(λ/i) for k ≤ 60 and λ ≤ 30, where each factor is rounded once and
// nothing overflows. Tolerance 5e-13 relative: Loader's absolute error in the exponent grows with the
// deviance term (up to about 500·2⁻⁵³ in this range), and the oracle rounds about 60 times.
func TestPoissonPMF_oracle(t *testing.T) {
	rng := newTestRand()
	for trial := range 400 {
		lambda := 0.01 + 30*rng.Float64()
		k := rng.IntN(61)
		want := math.Exp(-lambda)
		for i := 1; i <= k; i++ {
			want *= lambda / float64(i)
		}
		if want < 1e-290 {
			continue
		}
		if got, err := PoissonPMF(k, lambda); err != nil || !within(got, want, 5e-13) {
			t.Fatalf("trial %d: PoissonPMF(%d, %v) = (%v, %v), want %v", trial, k, lambda, got, err, want)
		}
	}
}

func TestPoissonPMF_noAlloc(t *testing.T) {
	requireNoAlloc(t, "PoissonPMF", func() { sinkFloat, _ = PoissonPMF(1000, 1000) })
}

func BenchmarkPoissonPMF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = PoissonPMF(1000, 1000)
	}
}
