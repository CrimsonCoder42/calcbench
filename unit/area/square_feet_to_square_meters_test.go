package area

import (
	"math"
	"testing"
)

// NIST SP 811 (2008) B.8: 1 ft² = 0.09290304 m² exactly. Hand computation: 100 / 0.09290304 = 1076.3910416709722.
// Tolerance 1e-12 relative: one division or multiplication, a few ulps.
func TestSquareFeetToSquareMeters_cases(t *testing.T) {
	cases := []struct{ in, want float64 }{
		{1, 0.09290304},
		{1076.3910416709722, 100},
		{0, 0},
	}
	for _, c := range cases {
		got, err := SquareFeetToSquareMeters(c.in)
		if err != nil || math.Abs(got-c.want) > 1e-12*math.Max(1, c.want) {
			t.Errorf("%v ft²: got %v %v, want %v", c.in, got, err, c.want)
		}
	}
}

func BenchmarkSquareFeetToSquareMeters(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = SquareFeetToSquareMeters(1076.39)
	}
}
