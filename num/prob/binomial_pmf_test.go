package prob

import (
	"math"
	"math/big"
	"testing"
)

// exactBinomialPMF returns C(n,k)·pᵏ·(1−p)ⁿ⁻ᵏ in exact rationals for the float64 p, rounded once.
func exactBinomialPMF(k, n int, p float64) float64 {
	pr := new(big.Rat).SetFloat64(p)
	qr := new(big.Rat).Sub(big.NewRat(1, 1), pr)
	r := new(big.Rat).SetInt(new(big.Int).Binomial(int64(n), int64(k)))
	r.Mul(r, ratPow(pr, k))
	r.Mul(r, ratPow(qr, n-k))
	f, _ := r.Float64()
	return f
}

// ratPow returns xᵉ for e ≥ 0 by square and multiply.
func ratPow(x *big.Rat, e int) *big.Rat {
	result, base := big.NewRat(1, 1), new(big.Rat).Set(x)
	for ; e > 0; e >>= 1 {
		if e&1 == 1 {
			result.Mul(result, base)
		}
		base.Mul(base, base)
	}
	return result
}

// NIST §1.3.6.6.18 / Ross §4.6: n = 10, p = 0.5, k = 5 → 252/1024 = 0.24609375. Large n where the naive
// form overflows: n = 1000, p = 0.5, k = 500 → C(1000,500)/2¹⁰⁰⁰ = 0.0252250181783608 (exact rational,
// rounded). Tolerance 1e-14 relative (Loader's method is accurate to a few ulps).
func TestBinomialPMF_sourceExample(t *testing.T) {
	for _, c := range []struct {
		k, n int
		want float64
	}{{5, 10, 0.24609375}, {500, 1000, 0.0252250181783608}} {
		if got, err := BinomialPMF(c.k, c.n, 0.5); err != nil || !within(got, c.want, 1e-14) {
			t.Errorf("BinomialPMF(%d, %d, 0.5) = (%v, %v), want (%v, nil)", c.k, c.n, got, err, c.want)
		}
	}
}

func TestBinomialPMF_cases(t *testing.T) {
	tests := []struct {
		name string
		k, n int
		p    float64
		want float64
	}{
		{"n = 0", 0, 0, 0.3, 1},
		{"k below support", -1, 5, 0.3, 0},
		{"k above support", 6, 5, 0.3, 0},
		{"p = 0 at 0", 0, 5, 0, 1},
		{"p = 0 elsewhere", 1, 5, 0, 0},
		{"p = 1 at n", 5, 5, 1, 1},
		{"p = 1 elsewhere", 4, 5, 1, 0},
		{"k = 0 edge", 0, 3, 0.5, 0.125},
		{"k = n edge", 3, 3, 0.5, 0.125},
		{"one trial", 1, 1, 0.3, 0.3},
		{"tiny p, k = 0", 0, 1_000_000, 1e-9, math.Exp(1e6 * math.Log1p(-1e-9))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := BinomialPMF(tt.k, tt.n, tt.p); err != nil || !within(got, tt.want, 1e-14) {
				t.Errorf("BinomialPMF(%d, %d, %v) = (%v, %v), want %v", tt.k, tt.n, tt.p, got, err, tt.want)
			}
		})
	}
	// Huge n: finite, in [0, 1], and the mode term is near 1/√(2πnpq).
	got, err := BinomialPMF(math.MaxInt/2, math.MaxInt, 0.5)
	if want := 1 / math.Sqrt(2*math.Pi*float64(math.MaxInt)/4); err != nil || !within(got, want, 1e-9) {
		t.Errorf("BinomialPMF(MaxInt/2, MaxInt, 0.5) = (%v, %v), want ≈ %v", got, err, want)
	}
}

func TestBinomialPMF_errors(t *testing.T) {
	_, err := BinomialPMF(0, -1, 0.5)
	requireErr(t, "n < 0", err, ErrOutOfRange)
	_, err = BinomialPMF(0, 5, nan)
	requireErr(t, "NaN", err, ErrNotFinite)
	_, err = BinomialPMF(0, 5, 1.1)
	requireErr(t, "p > 1", err, ErrProbability)
}

// Oracle: exact big.Rat binomial for random n ≤ 400, k, p. Tolerance 1e-12 relative: Loader's error
// grows with the deviance term, about 2⁻⁵³·n·(a cancellation factor ≤ 10 at the series switch), so up to
// a few thousand ulps here; values below 1e-290 are skipped (subnormal range).
func TestBinomialPMF_oracle(t *testing.T) {
	rng := newTestRand()
	for trial := range 400 {
		n := rng.IntN(401)
		k := rng.IntN(n + 1)
		p := rng.Float64()
		want := exactBinomialPMF(k, n, p)
		if want < 1e-290 {
			continue
		}
		if got, err := BinomialPMF(k, n, p); err != nil || !within(got, want, 1e-12) {
			t.Fatalf("trial %d: BinomialPMF(%d, %d, %v) = (%v, %v), want %v", trial, k, n, p, got, err, want)
		}
	}
}

func TestBinomialPMF_noAlloc(t *testing.T) {
	requireNoAlloc(t, "BinomialPMF", func() { sinkFloat, _ = BinomialPMF(500, 1000, 0.5) })
}

func BenchmarkBinomialPMF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = BinomialPMF(500, 1000, 0.5)
	}
}
