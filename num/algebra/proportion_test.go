package algebra

import (
	"errors"
	"math"
	"testing"
)

// Hand computations by cross products (Prealgebra 2e §6.5): 3 : 4 = 600 : 800; 250 : 280 = 100 : 112;
// 1 : 3 = 1 : 3; −2 : 5 = 4 : −10; a = 0 has no fourth term (ErrDegenerate).
func TestProportionFourthTerm_cases(t *testing.T) {
	tests := []struct{ a, b, c, want float64 }{
		{3, 4, 600, 800},
		{250, 280, 100, 112},
		{1, 3, 1, 3},
		{-2, 5, 4, -10},
		{5, 0, 7, 0},
		{1e-300, 1e300, 1e-300, 1e300},
		{4, math.MaxFloat64, 2, math.MaxFloat64 / 2},
	}
	for _, tc := range tests {
		if got, err := ProportionFourthTerm(tc.a, tc.b, tc.c); err != nil || !closePct(got, tc.want) {
			t.Errorf("ProportionFourthTerm(%v, %v, %v) = %v, %v; want %v", tc.a, tc.b, tc.c, got, err, tc.want)
		}
	}
}

func TestProportionFourthTerm_errors(t *testing.T) {
	if _, err := ProportionFourthTerm(0, 1, 1); !errors.Is(err, ErrDegenerate) {
		t.Errorf("a = 0: %v", err)
	}
	if _, err := ProportionFourthTerm(1, math.NaN(), 1); !errors.Is(err, ErrNotFinite) {
		t.Errorf("NaN: %v", err)
	}
	if _, err := ProportionFourthTerm(0.5, math.MaxFloat64, 4); !errors.Is(err, ErrOverflow) {
		t.Errorf("overflow: %v", err)
	}
	if a := testing.AllocsPerRun(100, func() { _, _ = ProportionFourthTerm(3, 4, 600) }); a != 0 {
		t.Errorf("allocates %v", a)
	}
}

func BenchmarkProportionFourthTerm(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = ProportionFourthTerm(3, 4, 600)
	}
}
