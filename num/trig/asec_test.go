package trig

import (
	"math"
	"testing"
)

// Worked example: DLMF Table 4.23.1, arcsec 1 = 0 and arcsec(−1) = π; arcsec 2 = π/3 (hand-derived:
// cos(π/3) = 1/2).
func TestAsec_dlmfTable(t *testing.T) {
	for _, tc := range []struct{ x, want, tol float64 }{{1, 0, 0}, {-1, math.Pi, 0}, {2, math.Pi / 3, ulps4}} {
		if got, err := Asec(tc.x); err != nil || !near(got, tc.want, tc.tol) {
			t.Errorf("Asec(%v) = (%v, %v), want %v", tc.x, got, err, tc.want)
		}
	}
}

func TestAsec_cases(t *testing.T) {
	tests := []struct {
		name   string
		x      float64
		lo, hi float64 // the result must be in [lo, hi]
	}{
		{"large positive tends to π/2 from below", 1e300, 0, math.Pi / 2},
		{"large negative tends to π/2 from above", -1e300, math.Pi / 2, math.Pi},
		{"largest float64", math.MaxFloat64, 0, math.Pi / 2},
		{"just above 1", math.Nextafter(1, 2), 0, 1e-7},
		{"just below −1", math.Nextafter(-1, -2), math.Pi - 1e-7, math.Pi},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := Asec(tc.x); err != nil || got < tc.lo || got > tc.hi {
				t.Errorf("Asec(%v) = (%v, %v), want in [%v, %v]", tc.x, got, err, tc.lo, tc.hi)
			}
		})
	}
}

func TestAsec_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { _, err := Asec(math.NaN()); return err }, ErrNotFinite},
		{"+Inf", func() error { _, err := Asec(math.Inf(1)); return err }, ErrNotFinite},
		{"-Inf", func() error { _, err := Asec(math.Inf(-1)); return err }, ErrNotFinite},
		{"zero", func() error { _, err := Asec(0); return err }, ErrDomain},
		{"inside (−1, 1)", func() error { _, err := Asec(0.5); return err }, ErrDomain},
		{"just inside 1", func() error { _, err := Asec(math.Nextafter(1, 0)); return err }, ErrDomain},
	})
}

// Oracle: sec(Asec x) = x for 1.5 ≤ |x| ≤ 100, where arccos is well conditioned; the round trip
// amplifies rounding by at most |x|·tan, so the tolerance is 1e-12 relative.
func TestAsec_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		x := 1.5 + rng.Float64()*98.5
		if rng.IntN(2) == 0 {
			x = -x
		}
		y, err := Asec(x)
		if err != nil || y < 0 || y > math.Pi || !relNear(1/math.Cos(y), x, 1e-12) {
			t.Fatalf("Asec(%v) = (%v, %v)", x, y, err)
		}
	}
}

func TestAsec_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = Asec(2.5) }); n != 0 {
		t.Errorf("Asec allocates %v times, want 0", n)
	}
}

func BenchmarkAsec(b *testing.B) {
	b.ReportAllocs()
	x := 2.5
	for b.Loop() {
		sink, _ = Asec(x)
	}
}
