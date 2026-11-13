package area

import (
	"errors"
	"math"
	"testing"
)

// NIST SP 811 (2008) B.8: 1 ft² = 0.09290304 m² exactly. Hand computation: 100 / 0.09290304 = 1076.3910416709722.
// Tolerance 1e-12 relative: one division or multiplication, a few ulps.
func TestSquareMetersToSquareFeet_cases(t *testing.T) {
	cases := []struct{ in, want float64 }{
		{100, 1076.3910416709722},
		{0.09290304, 1},
		{0, 0},
	}
	for _, c := range cases {
		got, err := SquareMetersToSquareFeet(c.in)
		if err != nil || math.Abs(got-c.want) > 1e-12*math.Max(1, c.want) {
			t.Errorf("%v m²: got %v %v, want %v", c.in, got, err, c.want)
		}
	}
	if _, err := SquareMetersToSquareFeet(math.MaxFloat64); !errors.Is(err, ErrOverflow) {
		t.Errorf("overflow: %v", err)
	}
}

func TestArea_errors(t *testing.T) {
	for _, f := range []func(float64) (float64, error){SquareFeetToSquareMeters, SquareMetersToSquareFeet} {
		if _, err := f(math.NaN()); !errors.Is(err, ErrNotFinite) {
			t.Errorf("NaN: %v", err)
		}
		if _, err := f(math.Inf(1)); !errors.Is(err, ErrNotFinite) {
			t.Errorf("Inf: %v", err)
		}
		if _, err := f(-1); !errors.Is(err, ErrNegative) {
			t.Errorf("negative: %v", err)
		}
	}
}

func BenchmarkSquareMetersToSquareFeet(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = SquareMetersToSquareFeet(100)
	}
}
