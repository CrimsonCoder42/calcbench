package electrical

import (
	"errors"
	"math"
	"testing"
)

// Hand computation from R = ρL/A (OpenStax University Physics Vol. 2, §9.3): 30.48 m, d = 0.00127 m,
// σ = 5.96e7 S/m: A = π·0.00127²/4 = 1.26677e-6 m², R = 30.48 / (5.96e7 · A) = 0.4037117406548103 Ω.
// Tolerance 1e-12 relative.
func TestWireResistance_cases(t *testing.T) {
	got, err := WireResistance(30.48, 0.00127, 5.96e7)
	if err != nil || !near(got, 0.4037117406548103) {
		t.Errorf("got %v %v", got, err)
	}
	if got, err := WireResistance(0, 0.001, 5.96e7); err != nil || got != 0 {
		t.Errorf("zero length: %v %v", got, err)
	}
	for _, c := range []struct {
		l, d, s float64
		want    error
	}{
		{1, 0, 1, ErrZero},
		{1, 1, 0, ErrZero},
		{-1, 1, 1, ErrNegative},
		{1, math.NaN(), 1, ErrNotFinite},
		{1, 1, math.Inf(1), ErrNotFinite},
		{1, 1e-200, 1e-200, ErrOverflow},
		{math.MaxFloat64, 1e-3, 1, ErrOverflow},
	} {
		if _, err := WireResistance(c.l, c.d, c.s); !errors.Is(err, c.want) {
			t.Errorf("(%v, %v, %v): %v, want %v", c.l, c.d, c.s, err, c.want)
		}
	}
}

func BenchmarkWireResistance(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = WireResistance(30.48, 0.00127, 5.96e7)
	}
}
