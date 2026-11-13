package special

import (
	"errors"
	"math"
	"testing"
)

func TestRiemannZeta_Values(t *testing.T) {
	cases := []struct{ s, want, tol float64 }{
		{2, math.Pi * math.Pi / 6, 1e-15},                // DLMF 25.6.1
		{4, math.Pow(math.Pi, 4) / 90, 1e-15},            // DLMF 25.6.1
		{3, 1.2020569031595942, 1e-15},                   // Apéry's constant, A&S Table 23.3
		{1.5, 2.612375348685488, 1e-14},                  // A&S Table 23.3 (ζ(3/2))
		{20, 1.0000009539620338, 1e-15},                  // A&S Table 23.3
		{1.01, 100.57794333849687, 1e-12},                // near the pole: 1/(s−1) + γ + O(s−1)
		{60, 1 + math.Pow(2, -60) + math.Pow(3, -60), 0}, // series branch
	}
	for _, c := range cases {
		got, err := RiemannZeta(c.s)
		if err != nil || math.Abs(got-c.want) > c.tol*c.want+1e-16 {
			t.Errorf("zeta(%v) = %v, %v; want %v", c.s, got, err, c.want)
		}
	}
	// The two branches agree at s = 20.
	lo, _ := RiemannZeta(math.Nextafter(20, 0))
	hi, _ := RiemannZeta(20)
	if math.Abs(lo-hi) > 1e-15 {
		t.Errorf("branch jump at 20: %v vs %v", lo, hi)
	}
}

func TestRiemannZeta_Errors(t *testing.T) {
	for _, s := range []float64{1, 0.5, -2} {
		if _, err := RiemannZeta(s); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("zeta(%v): err = %v; want ErrOutOfRange", s, err)
		}
	}
	for _, s := range []float64{math.NaN(), math.Inf(1)} {
		if _, err := RiemannZeta(s); !errors.Is(err, ErrNotFinite) {
			t.Errorf("zeta(%v): err = %v; want ErrNotFinite", s, err)
		}
	}
}

func BenchmarkRiemannZeta(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkF, _ = RiemannZeta(2.5)
	}
}
