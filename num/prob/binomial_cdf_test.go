package prob

import (
	"errors"
	"math/big"
	"testing"
)

// exactBinomialCDF returns Σ_{i≤k} C(n,i)pⁱqⁿ⁻ⁱ in exact rationals for the float64 p, rounded once.
func exactBinomialCDF(k, n int, p float64) float64 {
	pr := new(big.Rat).SetFloat64(p)
	qr := new(big.Rat).Sub(big.NewRat(1, 1), pr)
	sum := new(big.Rat)
	for i := 0; i <= k; i++ {
		term := new(big.Rat).SetInt(new(big.Int).Binomial(int64(n), int64(i)))
		term.Mul(term, ratPow(pr, i))
		term.Mul(term, ratPow(qr, n-i))
		sum.Add(sum, term)
	}
	f, _ := sum.Float64()
	return f
}

// NIST §1.3.6.6.18 (the CDF is the sum of the PMF): n = 10, p = 0.5, k = 5 → 638/1024 = 0.623046875.
// Also n = 100, p = 0.3, k = 30 → 0.5491236007687904 (exact rational sum, rounded). Tolerance 1e-14.
func TestBinomialCDF_sourceExample(t *testing.T) {
	if got, err := BinomialCDF(5, 10, 0.5); err != nil || !within(got, 0.623046875, 1e-14) {
		t.Errorf("BinomialCDF(5, 10, 0.5) = (%v, %v), want (0.623046875, nil)", got, err)
	}
	if got, err := BinomialCDF(30, 100, 0.3); err != nil || !within(got, 0.5491236007687904, 1e-14) {
		t.Errorf("BinomialCDF(30, 100, 0.3) = (%v, %v), want (0.5491236007687904, nil)", got, err)
	}
}

func TestBinomialCDF_cases(t *testing.T) {
	tests := []struct {
		name string
		k, n int
		p    float64
		want float64
	}{
		{"k < 0", -1, 10, 0.5, 0},
		{"k = n", 10, 10, 0.5, 1},
		{"k > n", 11, 10, 0.5, 1},
		{"p = 0", 0, 10, 0, 1},
		{"p = 1", 9, 10, 1, 0},
		{"k = 0", 0, 10, 0.5, 1.0 / 1024},
		{"k = n − 1", 9, 10, 0.5, 1023.0 / 1024},
		{"n = 0", 0, 0, 0.5, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := BinomialCDF(tt.k, tt.n, tt.p); err != nil || !within(got, tt.want, 1e-14) {
				t.Errorf("BinomialCDF(%d, %d, %v) = (%v, %v), want %v", tt.k, tt.n, tt.p, got, err, tt.want)
			}
		})
	}
	// A far lower tail keeps relative accuracy: P(X ≤ 0) for n = 2000, p = 0.5 is 2⁻²⁰⁰⁰ → underflows to 0,
	// and P(X ≤ 400) for n = 2000, p = 0.5 is tiny but representable.
	if got, err := BinomialCDF(400, 2000, 0.5); err != nil || !within(got, exactBinomialCDF(400, 2000, 0.5), 1e-11) {
		t.Errorf("BinomialCDF(400, 2000, 0.5) = (%v, %v)", got, err)
	}
}

func TestBinomialCDF_errors(t *testing.T) {
	_, err := BinomialCDF(0, -1, 0.5)
	requireErr(t, "n < 0", err, ErrOutOfRange)
	_, err = BinomialCDF(0, 5, nan)
	requireErr(t, "NaN", err, ErrNotFinite)
	_, err = BinomialCDF(0, 5, -0.5)
	requireErr(t, "p < 0", err, ErrProbability)
	// σ = 2³⁰ near the mean: more than 2²⁸ terms.
	if _, err = BinomialCDF(1<<61, 1<<62, 0.5); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("BinomialCDF(huge σ) err = %v, want ErrOutOfRange", err)
	}
}

// Oracle: exact big.Rat CDF for random n ≤ 200. Tolerance 2e-12 relative (as BinomialPMF's oracle, plus the recurrence) plus 1e-15 absolute (the
// complement side is 1 − tail, accurate to an ulp of 1).
func TestBinomialCDF_oracle(t *testing.T) {
	rng := newTestRand()
	for trial := range 300 {
		n := rng.IntN(201)
		k := rng.IntN(n + 1)
		p := rng.Float64()
		want := exactBinomialCDF(k, n, p)
		got, err := BinomialCDF(k, n, p)
		if err != nil || absDiff(got, want) > 2e-12*want+1e-15 {
			t.Fatalf("trial %d: BinomialCDF(%d, %d, %v) = (%v, %v), want %v", trial, k, n, p, got, err, want)
		}
	}
}

func TestBinomialCDF_noAlloc(t *testing.T) {
	requireNoAlloc(t, "BinomialCDF", func() { sinkFloat, _ = BinomialCDF(500, 1000, 0.5) })
}

func BenchmarkBinomialCDF(b *testing.B) {
	for _, size := range []struct {
		name string
		n    int
	}{{"n=1e3", 1_000}, {"n=1e5", 100_000}, {"n=1e7", 10_000_000}} {
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, _ = BinomialCDF(size.n/2, size.n, 0.5)
			}
		})
	}
}
