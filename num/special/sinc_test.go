package special

import (
	"math"
	"testing"
)

func TestSinc_Values(t *testing.T) {
	cases := []struct {
		name    string
		x, want float64
	}{
		{"zero", 0, 1},
		{"half (Oppenheim & Schafer §2.9 normalised form)", 0.5, 2 / math.Pi},
		{"minus half (even)", -0.5, 2 / math.Pi},
		{"one and a half", 1.5, -1 / (1.5 * math.Pi)},
		{"tiny (Taylor branch)", 1e-9, 1},
		{"large non-integer", 1e6 + 0.5, 1 / ((1e6 + 0.5) * math.Pi)},
		{"+Inf limit", math.Inf(1), 0},
		{"-Inf limit", math.Inf(-1), 0},
	}
	for _, c := range cases {
		if got := Sinc(c.x); math.Abs(got-c.want) > 1e-15 {
			t.Errorf("%s: Sinc(%v) = %v; want %v", c.name, c.x, got, c.want)
		}
	}
	for _, n := range []float64{1, -1, 2, 3, -7, 1e6} {
		if got := Sinc(n); got != 0 {
			t.Errorf("Sinc(%v) = %v; want exactly 0", n, got)
		}
	}
	if got := Sinc(math.NaN()); !math.IsNaN(got) {
		t.Errorf("Sinc(NaN) = %v; want NaN", got)
	}
	// Taylor and direct branches meet at |πx| = 1e-4. Two different points differ by the true slope
	// (≈ −π²x/3 · Δx ≈ 6.6e-15 here), so each side is instead checked against sin(πx)/(πx) at the
	// SAME point: near 1e-4 that quotient is accurate to ~1 ulp, and the Taylor remainder y⁶/5040 is ~2e-28.
	x := 1e-4 / math.Pi
	for _, p := range []float64{x * 0.999999, x * 1.000001} {
		want := math.Sin(math.Pi*p) / (math.Pi * p)
		if d := Sinc(p) - want; math.Abs(d) > 1e-15 {
			t.Errorf("branch mismatch at %v: Sinc = %v, sin(πx)/(πx) = %v", p, Sinc(p), want)
		}
	}
}

func BenchmarkSinc(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkF = Sinc(0.37)
	}
}
