package trig

import (
	"math"
	"testing"
)

// Worked example, hand-derived from DLMF 4.14.7 (cot z = cos z / sin z): cot(π/4) = 1 to rounding
// (1.0000000000000002), and cot(math.Pi/2) = 6.123233995736766e-17, not 0, because math.Pi/2 is not π/2.
// The 4e-15 tolerance at the multiples of π/2: Go's math.Sincos reduces the argument with a three-part
// π/4, so the tiny cos or sin there carries about 1.5e-15 relative error (measured on darwin/arm64).
func TestCot_dlmfExample(t *testing.T) {
	for _, tc := range []struct{ x, want, tol float64 }{{math.Pi / 4, 1, ulps4}, {math.Pi / 2, 6.123233995736766e-17, 4e-15}} {
		if got, err := Cot(tc.x); err != nil || !relNear(got, tc.want, tc.tol) {
			t.Errorf("Cot(%v) = (%v, %v), want %v", tc.x, got, err, tc.want)
		}
	}
}

func TestCot_cases(t *testing.T) {
	for _, tc := range []struct {
		name    string
		x, want float64
		tol     float64
	}{{"negative quarter", -math.Pi / 4, -1, ulps4}, {"tiny", 1e-300, 1e300, ulps4}, {"π/3", math.Pi / 3, 1 / math.Sqrt(3), ulps4}} {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := Cot(tc.x); err != nil || !relNear(got, tc.want, tc.tol) {
				t.Errorf("Cot(%v) = (%v, %v), want %v", tc.x, got, err, tc.want)
			}
		})
	}
}

func TestCot_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { _, err := Cot(math.NaN()); return err }, ErrNotFinite},
		{"-Inf", func() error { _, err := Cot(math.Inf(-1)); return err }, ErrNotFinite},
		{"+0 pole", func() error { _, err := Cot(0); return err }, ErrDomain},
		{"-0 pole", func() error { _, err := Cot(math.Copysign(0, -1)); return err }, ErrDomain},
		{"subnormal overflows", func() error { _, err := Cot(5e-324); return err }, ErrOverflow},
	})
}

// Oracle: cot x · tan x = 1 on x with |sin x| and |cos x| ≥ 0.1, to a few ulps (1e-14).
func TestCot_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		x := (rng.Float64() - 0.5) * 20
		if math.Abs(math.Sin(x)) < 0.1 || math.Abs(math.Cos(x)) < 0.1 {
			continue
		}
		c, err := Cot(x)
		if err != nil || math.Abs(c*math.Tan(x)-1) > 1e-14 {
			t.Fatalf("Cot(%v) = (%v, %v)", x, c, err)
		}
	}
}

func TestCot_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = Cot(0.7) }); n != 0 {
		t.Errorf("Cot allocates %v times, want 0", n)
	}
}

func BenchmarkCot(b *testing.B) {
	b.ReportAllocs()
	x := 0.7
	for b.Loop() {
		sink, _ = Cot(x)
	}
}
