package trig

import (
	"math"
	"testing"
)

// Worked example, hand-derived from DLMF 4.28.7 (coth z = 1/tanh z): coth 1 = (e² + 1)/(e² − 1) =
// 1.3130352854993312, to rounding.
func TestCoth_dlmfExample(t *testing.T) {
	for _, tc := range []struct{ x, want float64 }{{1, 1.3130352854993312}, {-1, -1.3130352854993312}} {
		if got, err := Coth(tc.x); err != nil || !relNear(got, tc.want, ulps4) {
			t.Errorf("Coth(%v) = (%v, %v), want %v", tc.x, got, err, tc.want)
		}
	}
}

func TestCoth_cases(t *testing.T) {
	tests := []struct {
		name    string
		x, want float64
	}{
		{"large rounds to 1", 50, 1},
		{"large negative rounds to −1", -50, -1},
		{"largest float64", math.MaxFloat64, 1},
		{"tiny", 1e-300, 1e300},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := Coth(tc.x); err != nil || !relNear(got, tc.want, ulps4) {
				t.Errorf("Coth(%v) = (%v, %v), want %v", tc.x, got, err, tc.want)
			}
		})
	}
}

func TestCoth_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { _, err := Coth(math.NaN()); return err }, ErrNotFinite},
		{"-Inf", func() error { _, err := Coth(math.Inf(-1)); return err }, ErrNotFinite},
		{"+0 pole", func() error { _, err := Coth(0); return err }, ErrDomain},
		{"-0 pole", func() error { _, err := Coth(math.Copysign(0, -1)); return err }, ErrDomain},
		{"subnormal overflows", func() error { _, err := Coth(-5e-324); return err }, ErrOverflow},
	})
}

// Oracle: coth x · tanh x = 1 and |coth x| ≥ 1, for 0.01 ≤ |x| ≤ 20, to a few ulps (1e-15).
func TestCoth_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		x := 0.01 + rng.Float64()*19.99
		if rng.IntN(2) == 0 {
			x = -x
		}
		c, err := Coth(x)
		if err != nil || math.Abs(c) < 1 || math.Abs(c*math.Tanh(x)-1) > 1e-15 {
			t.Fatalf("Coth(%v) = (%v, %v)", x, c, err)
		}
	}
}

func TestCoth_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = Coth(0.7) }); n != 0 {
		t.Errorf("Coth allocates %v times, want 0", n)
	}
}

func BenchmarkCoth(b *testing.B) {
	b.ReportAllocs()
	x := 0.7
	for b.Loop() {
		sink, _ = Coth(x)
	}
}
