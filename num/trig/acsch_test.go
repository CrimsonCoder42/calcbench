package trig

import (
	"math"
	"testing"
)

// Worked example, hand-derived from DLMF 4.37.7 (arccsch z = arcsinh(1/z)): arccsch 1 = arcsinh 1 =
// ln(1 + √2) = 0.881373587019543.
func TestAcsch_dlmfExample(t *testing.T) {
	for _, tc := range []struct{ x, want float64 }{{1, 0.881373587019543}, {-1, -0.881373587019543}} {
		if got, err := Acsch(tc.x); err != nil || !relNear(got, tc.want, ulps4) {
			t.Errorf("Acsch(%v) = (%v, %v), want %v", tc.x, got, err, tc.want)
		}
	}
}

func TestAcsch_cases(t *testing.T) {
	tests := []struct {
		name    string
		x, want float64
	}{
		{"¼ is arcsinh 4", 0.25, math.Asinh(4)},
		{"4 is arcsinh ¼", 4, math.Asinh(0.25)},
		// A&S 4.6.23 with x² below an ulp of 1: ln(2/|x|), with the sign of x.
		{"smallest subnormal", 5e-324, math.Ln2 - math.Log(5e-324)},
		{"negative subnormal", -5e-324, -(math.Ln2 - math.Log(5e-324))},
		{"largest float64", math.MaxFloat64, 1 / math.MaxFloat64},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := Acsch(tc.x); err != nil || !relNear(got, tc.want, ulps4) {
				t.Errorf("Acsch(%v) = (%v, %v), want %v", tc.x, got, err, tc.want)
			}
		})
	}
}

func TestAcsch_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { _, err := Acsch(math.NaN()); return err }, ErrNotFinite},
		{"+Inf", func() error { _, err := Acsch(math.Inf(1)); return err }, ErrNotFinite},
		{"-Inf", func() error { _, err := Acsch(math.Inf(-1)); return err }, ErrNotFinite},
		{"+0 pole", func() error { _, err := Acsch(0); return err }, ErrDomain},
		{"-0 pole", func() error { _, err := Acsch(math.Copysign(0, -1)); return err }, ErrDomain},
	})
}

// Oracle: csch(Acsch x) = x and the sign is kept, for 0.01 ≤ |x| ≤ 100 (both branches): 1e-13 relative.
func TestAcsch_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		x := math.Pow(10, rng.Float64()*4-2)
		if rng.IntN(2) == 0 {
			x = -x
		}
		y, err := Acsch(x)
		if err != nil || math.Signbit(y) != math.Signbit(x) || !relNear(1/math.Sinh(y), x, 1e-13) {
			t.Fatalf("Acsch(%v) = (%v, %v)", x, y, err)
		}
	}
}

func TestAcsch_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = Acsch(0.7) }); n != 0 {
		t.Errorf("Acsch allocates %v times, want 0", n)
	}
}

func BenchmarkAcsch(b *testing.B) {
	b.ReportAllocs()
	x := 0.7
	for b.Loop() {
		sink, _ = Acsch(x)
	}
}
