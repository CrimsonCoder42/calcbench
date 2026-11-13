package prob

import (
	"math/big"
	"testing"
)

// Wheatcroft 2021 eq. (1) (π = Σ 1/Oᵢ − 1): [1.9, 1.9] → 2/1.9 − 1 ≈ 0.0526315789473684. Hand-checked; tolerance
// 4 ulps relative to the overround (1/1.9 is rounded once).
func TestOverround_sourceExample(t *testing.T) {
	if got, err := Overround([]float64{1.9, 1.9}); err != nil || !within(got, 1.0/19, 1e-14) {
		t.Errorf("Overround([1.9 1.9]) = (%v, %v), want (≈ 1/19, nil)", got, err)
	}
}

// Wheatcroft 2021, preprint §2.1: the two-outcome market O₁ = 3, O₂ = 1.4 sums to more than one; by eq. (1)
// π = 1/3 + 1/1.4 − 1 = 1/21 exactly in rationals. Hand-checked; tolerance 4 ulps relative.
func TestOverround_wheatcroftExample(t *testing.T) {
	if got, err := Overround([]float64{3, 1.4}); err != nil || !within(got, 1.0/21, 1e-14) {
		t.Errorf("Overround([3 1.4]) = (%v, %v), want (≈ 1/21, nil)", got, err)
	}
}

func TestOverround_cases(t *testing.T) {
	tests := []struct {
		name string
		odds []float64
		want float64
	}{
		{"fair coin", []float64{2, 2}, 0},
		{"fair three-way", []float64{3, 3, 3}, 0},
		{"arbitrage is negative", []float64{2.5, 2.5}, -0.2},
		{"certain favourite", []float64{1, 1e300}, 1e-300},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Overround(tt.odds); err != nil || !within(got, tt.want, 1e-15) {
				t.Errorf("Overround(%v) = (%v, %v), want %v", tt.odds, got, err, tt.want)
			}
		})
	}
}

func TestOverround_errors(t *testing.T) {
	_, err := Overround([]float64{2})
	requireErr(t, "one outcome", err, ErrOutOfRange)
	_, err = Overround([]float64{2, 0.5})
	requireErr(t, "d < 1", err, ErrOutOfRange)
	_, err = Overround([]float64{2, nan})
	requireErr(t, "NaN", err, ErrNotFinite)
}

// Oracle: the exact rational Σ fl(1/dᵢ) − 1 for random markets. Tolerance: 2⁻⁵² absolute (the sum is
// compensated; each 1/dᵢ is rounded once, which the oracle reproduces).
func TestOverround_oracle(t *testing.T) {
	rng := newTestRand()
	for trial := range 300 {
		n := 2 + rng.IntN(20)
		odds := make([]float64, n)
		exact := big.NewRat(-1, 1)
		for i := range n {
			odds[i] = 1 + 20*rng.Float64()
			exact.Add(exact, new(big.Rat).SetFloat64(1/odds[i]))
		}
		want, _ := exact.Float64()
		if got, err := Overround(odds); err != nil || absDiff(got, want) > 0x1p-52 {
			t.Fatalf("trial %d: Overround = (%v, %v), want %v", trial, got, err, want)
		}
	}
}

// absDiff returns |a − b|.
func absDiff(a, b float64) float64 {
	if a > b {
		return a - b
	}
	return b - a
}

func TestOverround_noAlloc(t *testing.T) {
	odds := []float64{1.9, 1.9}
	requireNoAlloc(t, "Overround", func() { sinkFloat, _ = Overround(odds) })
}

func BenchmarkOverround(b *testing.B) {
	for _, size := range []struct {
		name string
		n    int
	}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}} {
		odds := make([]float64, size.n)
		for i := range odds {
			odds[i] = float64(size.n)
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, _ = Overround(odds)
			}
		})
	}
}
