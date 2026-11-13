package prob

import (
	"math/big"
	"testing"
)

// exactNegativeBinomialPMF returns C(x−1, r−1)·pʳ·(1−p)^(x−r) in exact rationals, rounded once.
func exactNegativeBinomialPMF(x, r int, p float64) float64 {
	pr := new(big.Rat).SetFloat64(p)
	qr := new(big.Rat).Sub(big.NewRat(1, 1), pr)
	v := new(big.Rat).SetInt(new(big.Int).Binomial(int64(x-1), int64(r-1)))
	v.Mul(v, ratPow(pr, r))
	v.Mul(v, ratPow(qr, x-r))
	f, _ := v.Float64()
	return f
}

// Casella–Berger §3.2 (negative binomial, X = trial of the r-th success, C(x−1,r−1)pʳ(1−p)^(x−r)):
// r = 2, p = 0.5, x = 4 → 3·0.25·0.25 = 0.1875. Hand-checked; tolerance 1e-14.
func TestNegativeBinomialTrialsPMF_sourceExample(t *testing.T) {
	if got, err := NegativeBinomialTrialsPMF(4, 2, 0.5); err != nil || !within(got, 0.1875, 1e-14) {
		t.Errorf("NegativeBinomialTrialsPMF(4, 2, 0.5) = (%v, %v), want (0.1875, nil)", got, err)
	}
}

func TestNegativeBinomialTrialsPMF_cases(t *testing.T) {
	tests := []struct {
		name string
		x, r int
		p    float64
		want float64
	}{
		{"x < r", 1, 2, 0.5, 0},
		{"x = r", 2, 2, 0.5, 0.25},
		{"r = 1 is geometric", 3, 1, 0.2, 0.128},
		{"p = 1 at x = r", 3, 3, 1, 1},
		{"p = 1 beyond r", 4, 3, 1, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := NegativeBinomialTrialsPMF(tt.x, tt.r, tt.p); err != nil || !within(got, tt.want, 1e-14) {
				t.Errorf("NegativeBinomialTrialsPMF(%d, %d, %v) = (%v, %v), want %v", tt.x, tt.r, tt.p, got, err, tt.want)
			}
		})
	}
}

func TestNegativeBinomialTrialsPMF_errors(t *testing.T) {
	_, err := NegativeBinomialTrialsPMF(4, 0, 0.5)
	requireErr(t, "r < 1", err, ErrOutOfRange)
	_, err = NegativeBinomialTrialsPMF(4, 2, 0)
	requireErr(t, "p = 0", err, ErrOutOfRange)
	_, err = NegativeBinomialTrialsPMF(4, 2, nan)
	requireErr(t, "NaN", err, ErrNotFinite)
	_, err = NegativeBinomialTrialsPMF(4, 2, -1)
	requireErr(t, "p < 0", err, ErrProbability)
}

// Oracle: exact big.Rat for random r ≤ 30, x ≤ 300. Tolerance 1e-12 (as BinomialPMF's oracle).
func TestNegativeBinomialTrialsPMF_oracle(t *testing.T) {
	rng := newTestRand()
	for trial := range 300 {
		r := 1 + rng.IntN(30)
		x := r + rng.IntN(270)
		p := 0.01 + 0.99*rng.Float64()
		want := exactNegativeBinomialPMF(x, r, p)
		if want < 1e-290 {
			continue
		}
		if got, err := NegativeBinomialTrialsPMF(x, r, p); err != nil || !within(got, want, 1e-12) {
			t.Fatalf("trial %d: NegativeBinomialTrialsPMF(%d, %d, %v) = (%v, %v), want %v", trial, x, r, p, got, err, want)
		}
	}
}

func TestNegativeBinomialTrialsPMF_noAlloc(t *testing.T) {
	requireNoAlloc(t, "NegativeBinomialTrialsPMF", func() { sinkFloat, _ = NegativeBinomialTrialsPMF(40, 10, 0.3) })
}

func BenchmarkNegativeBinomialTrialsPMF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = NegativeBinomialTrialsPMF(40, 10, 0.3)
	}
}
