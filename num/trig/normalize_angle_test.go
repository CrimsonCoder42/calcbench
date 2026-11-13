package trig

import (
	"math"
	"math/big"
	"testing"
)

// Worked example (hand-derived from the POSIX fmod() definition, IEEE Std 1003.1-2024): 7π/2 → 3π/2,
// −π/2 → 3π/2, 2π → 0. 7π/2 is not exact in float64, so its result is within a few ulps of 3π/2.
func TestNormalizeAngle_fmodExample(t *testing.T) {
	for _, tc := range []struct{ in, want, tol float64 }{{7 * math.Pi / 2, 3 * math.Pi / 2, ulps4}, {-math.Pi / 2, 3 * math.Pi / 2, ulps4}, {2 * math.Pi, 0, 0}} {
		if got, err := NormalizeAngle(tc.in); err != nil || !near(got, tc.want, tc.tol) {
			t.Errorf("NormalizeAngle(%v) = (%v, %v), want %v", tc.in, got, err, tc.want)
		}
	}
}

func TestNormalizeAngle_cases(t *testing.T) {
	tests := []struct {
		name     string
		in, want float64
	}{
		{"zero", 0, 0},
		{"negative zero becomes +0", math.Copysign(0, -1), 0},
		{"already in range", 1, 1},
		{"just below 2π stays", math.Nextafter(twoPi, 0), math.Nextafter(twoPi, 0)},
		{"tiny negative rounds up to 2π, so 0", -1e-300, 0},
		{"negative full turn", -twoPi, 0},
		{"π", math.Pi, math.Pi},
		{"−π", -math.Pi, math.Pi},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NormalizeAngle(tc.in)
			if err != nil || got != tc.want || math.Signbit(got) {
				t.Errorf("NormalizeAngle(%v) = (%v, %v), want +%v", tc.in, got, err, tc.want)
			}
		})
	}
	// The largest float64 still lands in [0, 2π).
	if got, err := NormalizeAngle(math.MaxFloat64); err != nil || got < 0 || got >= twoPi {
		t.Errorf("NormalizeAngle(MaxFloat64) = (%v, %v), want in [0, 2π)", got, err)
	}
}

func TestNormalizeAngle_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { _, err := NormalizeAngle(math.NaN()); return err }, ErrNotFinite},
		{"+Inf", func() error { _, err := NormalizeAngle(math.Inf(1)); return err }, ErrNotFinite},
		{"-Inf", func() error { _, err := NormalizeAngle(math.Inf(-1)); return err }, ErrNotFinite},
	})
}

// Oracle: θ − 2π_f·⌊θ/2π_f⌋ in exact big.Float arithmetic (2π_f = float64 2π), rounded once. math.Mod is
// exact, and the only extra rounding is the +2π for negative θ: tolerance one ulp of 2π, 2^-50.
func TestNormalizeAngle_oracle(t *testing.T) {
	rng := newRNG()
	tp := new(big.Float).SetPrec(2000).SetFloat64(twoPi)
	for range 500 {
		theta := (rng.Float64() - 0.5) * math.Pow(10, float64(rng.IntN(12)))
		x := new(big.Float).SetPrec(2000).SetFloat64(theta)
		q, _ := new(big.Float).SetPrec(2000).Quo(x, tp).Int(nil)
		if x.Sign() < 0 && new(big.Float).SetPrec(2000).Mul(new(big.Float).SetInt(q), tp).Cmp(x) != 0 {
			q.Sub(q, bigOne) // floor for negatives
		}
		r := new(big.Float).SetPrec(2000).Sub(x, new(big.Float).SetPrec(2000).Mul(new(big.Float).SetPrec(2000).SetInt(q), tp))
		want, _ := r.Float64()
		got, err := NormalizeAngle(theta)
		if err != nil || got < 0 || got >= twoPi || math.Abs(got-want) > 0x1p-50 {
			t.Fatalf("NormalizeAngle(%v) = (%v, %v), want %v", theta, got, err, want)
		}
	}
}

func TestNormalizeAngle_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = NormalizeAngle(-12.5) }); n != 0 {
		t.Errorf("NormalizeAngle allocates %v times, want 0", n)
	}
}

func BenchmarkNormalizeAngle(b *testing.B) {
	b.ReportAllocs()
	x := -12.5
	for b.Loop() {
		sink, _ = NormalizeAngle(x)
	}
}
