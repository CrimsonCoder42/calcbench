package trig

import (
	"math"
	"math/big"
	"testing"
)

// Worked example: the inverse of NIST SP 811 (2008), §5.1.1, Table 6, 1° = (π/180) rad.
func TestRadToDeg_sp811Example(t *testing.T) {
	for _, tc := range []struct{ rad, want float64 }{{math.Pi, 180}, {math.Pi / 2, 90}, {-math.Pi / 4, -45}, {2 * math.Pi, 360}} {
		if got, err := RadToDeg(tc.rad); err != nil || got != tc.want {
			t.Errorf("RadToDeg(%v) = (%v, %v), want (%v, nil)", tc.rad, got, err, tc.want)
		}
	}
}

func TestRadToDeg_cases(t *testing.T) {
	tests := []struct {
		name      string
		rad, want float64
		tol       float64
	}{
		{"zero", 0, 0, 0},
		{"one radian", 1, 57.29577951308232, ulps4},
		{"negative", -1, -57.29577951308232, ulps4},
		{"largest that fits", 3e306, 3e306 / math.Pi * 180, ulps4},
		{"smallest subnormal", 5e-324, 0, 1e-320},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := RadToDeg(tc.rad); err != nil || !near(got, tc.want, tc.tol) {
				t.Errorf("RadToDeg(%v) = (%v, %v), want %v", tc.rad, got, err, tc.want)
			}
		})
	}
}

func TestRadToDeg_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { _, err := RadToDeg(math.NaN()); return err }, ErrNotFinite},
		{"+Inf", func() error { _, err := RadToDeg(math.Inf(1)); return err }, ErrNotFinite},
		{"-Inf", func() error { _, err := RadToDeg(math.Inf(-1)); return err }, ErrNotFinite},
		{"overflow", func() error { _, err := RadToDeg(1e307); return err }, ErrOverflow},
		{"overflow, largest float64", func() error { _, err := RadToDeg(-math.MaxFloat64); return err }, ErrOverflow},
	})
}

// Oracle: rad·180/π in 256-bit big.Float; two roundings in the code, so one ulp (2^-52 relative).
// Also the round trip through DegToRad, within four roundings.
func TestRadToDeg_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		rad := (rng.Float64() - 0.5) * 1e4
		want, _ := new(big.Float).SetPrec(256).Quo(new(big.Float).SetPrec(256).Mul(big.NewFloat(rad), big.NewFloat(180)), bigPi).Float64()
		got, err := RadToDeg(rad)
		if err != nil || !relNear(got, want, 0x1p-52) {
			t.Fatalf("RadToDeg(%v) = (%v, %v), want %v", rad, got, err, want)
		}
		if back, _ := DegToRad(got); !relNear(back, rad, ulps4) {
			t.Fatalf("DegToRad(RadToDeg(%v)) = %v", rad, back)
		}
	}
}

func TestRadToDeg_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = RadToDeg(0.75) }); n != 0 {
		t.Errorf("RadToDeg allocates %v times, want 0", n)
	}
}

func BenchmarkRadToDeg(b *testing.B) {
	b.ReportAllocs()
	x := 0.75
	for b.Loop() {
		sink, _ = RadToDeg(x)
	}
}
