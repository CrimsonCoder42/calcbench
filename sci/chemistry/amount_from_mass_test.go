package chemistry

import (
	"errors"
	"math"
	"testing"
)

// Hand computation: 100 g / 58.44 g/mol = 1.7111567419575633 mol. Tolerance 1e-12 relative: one division.
func TestAmountFromMass_cases(t *testing.T) {
	got, err := AmountFromMass(100, 58.44)
	if err != nil || math.Abs(got-1.7111567419575633) > 1e-12 {
		t.Errorf("got %v %v", got, err)
	}
	if got, err := AmountFromMass(0, 58.44); err != nil || got != 0 {
		t.Errorf("zero mass: %v %v", got, err)
	}
	for _, c := range []struct {
		m, mm float64
		want  error
	}{
		{1, 0, ErrZero},
		{-1, 1, ErrNegative},
		{1, -1, ErrNegative},
		{math.NaN(), 1, ErrNotFinite},
		{1, math.Inf(1), ErrNotFinite},
		{math.MaxFloat64, 0.5, ErrOverflow},
	} {
		if _, err := AmountFromMass(c.m, c.mm); !errors.Is(err, c.want) {
			t.Errorf("AmountFromMass(%v, %v): %v, want %v", c.m, c.mm, err, c.want)
		}
	}
}

func BenchmarkAmountFromMass(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = AmountFromMass(100, 58.44)
	}
}
