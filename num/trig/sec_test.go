package trig

import (
	"math"
	"testing"
)

// Worked example, hand-derived from DLMF 4.14.6 (sec z = 1/cos z): sec 0 = 1, sec π = −1. At math.Pi/2,
// which is π/2 − 6.123233995736766e-17, the true secant is 1/6.123233995736766e-17 = 1.633123935319537e16.
// The 4e-15 tolerance at the multiples of π/2: Go's math.Sincos reduces the argument with a three-part
// π/4, so the tiny cos or sin there carries about 1.5e-15 relative error (measured on darwin/arm64).
func TestSec_dlmfExample(t *testing.T) {
	for _, tc := range []struct{ x, want, tol float64 }{{0, 1, 0}, {math.Pi, -1, 0}, {math.Pi / 3, 2, ulps4}, {math.Pi / 2, 1.633123935319537e16, 4e-15}} {
		if got, err := Sec(tc.x); err != nil || !relNear(got, tc.want, tc.tol) {
			t.Errorf("Sec(%v) = (%v, %v), want %v", tc.x, got, err, tc.want)
		}
	}
}

func TestSec_cases(t *testing.T) {
	for _, tc := range []struct {
		name string
		x    float64
	}{{"negative", -1}, {"large argument", 1e300}, {"tiny", 1e-300}, {"largest float64", math.MaxFloat64}} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Sec(tc.x)
			if err != nil || math.Abs(got) < 1 || !relNear(got, 1/math.Cos(tc.x), 0) {
				t.Errorf("Sec(%v) = (%v, %v)", tc.x, got, err)
			}
		})
	}
}

func TestSec_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { _, err := Sec(math.NaN()); return err }, ErrNotFinite},
		{"+Inf", func() error { _, err := Sec(math.Inf(1)); return err }, ErrNotFinite},
		{"-Inf", func() error { _, err := Sec(math.Inf(-1)); return err }, ErrNotFinite},
	})
}

// Oracle: the identity sec² x − tan² x = 1 (A&S 4.3.11), on x where |cos x| ≥ 0.1 so both terms are at
// most 100 and the check holds to about 100 ulps (1e-13).
func TestSec_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		x := (rng.Float64() - 0.5) * 20
		if math.Abs(math.Cos(x)) < 0.1 {
			continue
		}
		s, err := Sec(x)
		tn := math.Tan(x)
		if err != nil || math.Abs(s*s-tn*tn-1) > 1e-13 {
			t.Fatalf("Sec(%v) = (%v, %v): sec² − tan² = %v", x, s, err, s*s-tn*tn)
		}
	}
}

func TestSec_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = Sec(0.7) }); n != 0 {
		t.Errorf("Sec allocates %v times, want 0", n)
	}
}

func BenchmarkSec(b *testing.B) {
	b.ReportAllocs()
	x := 0.7
	for b.Loop() {
		sink, _ = Sec(x)
	}
}
