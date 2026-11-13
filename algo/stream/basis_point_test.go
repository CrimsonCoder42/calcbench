package stream

import (
	"errors"
	"math"
	"testing"
)

// 1 bp = 0.01% = 1e-4 (Hull, Glossary, "Basis Point"); each quotient and product below is the float64
// nearest the exact decimal, so the comparisons are exact.
func TestBasisPointsToFraction(t *testing.T) {
	for _, c := range [][2]float64{{25, 0.0025}, {1, 1e-4}, {-50, -0.005}, {0, 0}, {10000, 1}} {
		if got, err := BasisPointsToFraction(c[0]); err != nil || got != c[1] {
			t.Errorf("BasisPointsToFraction(%v) = (%v, %v), want %v", c[0], got, err, c[1])
		}
	}
	for _, bp := range []float64{math.NaN(), math.Inf(1)} {
		if _, err := BasisPointsToFraction(bp); !errors.Is(err, ErrNotFinite) {
			t.Errorf("BasisPointsToFraction(%v) error = %v, want ErrNotFinite", bp, err)
		}
	}
}

func TestFractionToBasisPoints(t *testing.T) {
	for _, c := range [][2]float64{{0.0025, 25}, {1e-4, 1}, {-0.005, -50}, {0, 0}, {1, 10000}} {
		if got, err := FractionToBasisPoints(c[0]); err != nil || got != c[1] {
			t.Errorf("FractionToBasisPoints(%v) = (%v, %v), want %v", c[0], got, err, c[1])
		}
	}
	if _, err := FractionToBasisPoints(math.Inf(-1)); !errors.Is(err, ErrNotFinite) {
		t.Errorf("FractionToBasisPoints(-Inf) error = %v, want ErrNotFinite", err)
	}
	if _, err := FractionToBasisPoints(1e305); !errors.Is(err, ErrOverflow) {
		t.Errorf("FractionToBasisPoints(1e305) error = %v, want ErrOverflow", err)
	}
}

func BenchmarkBasisPointsToFraction(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = BasisPointsToFraction(25)
	}
}

func BenchmarkFractionToBasisPoints(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = FractionToBasisPoints(0.0025)
	}
}
