package trig

import (
	"math"
	"math/big"
	"testing"
)

// bigOne is the integer 1, for floor adjustments in the oracles.
var bigOne = big.NewInt(1)

// Worked example (hand-derived from the POSIX remainder() definition, IEEE Std 1003.1-2024):
// 7π/2 → −π/2, 3π/2 → −π/2, −π → π.
func TestNormalizeAngleSigned_remainderExample(t *testing.T) {
	for _, tc := range []struct{ in, want, tol float64 }{{7 * math.Pi / 2, -math.Pi / 2, ulps4}, {3 * math.Pi / 2, -math.Pi / 2, ulps4}, {-math.Pi, math.Pi, 0}} {
		if got, err := NormalizeAngleSigned(tc.in); err != nil || !near(got, tc.want, tc.tol) {
			t.Errorf("NormalizeAngleSigned(%v) = (%v, %v), want %v", tc.in, got, err, tc.want)
		}
	}
}

func TestNormalizeAngleSigned_cases(t *testing.T) {
	tests := []struct {
		name     string
		in, want float64
	}{
		{"zero", 0, 0},
		{"π stays", math.Pi, math.Pi},
		{"−π tie goes to +π", -math.Pi, math.Pi},
		{"3π tie goes to +π", 3 * math.Pi, math.Pi},
		{"−3π tie goes to +π", -3 * math.Pi, math.Pi},
		{"just above −π stays", math.Nextafter(-math.Pi, 0), math.Nextafter(-math.Pi, 0)},
		{"full turn", twoPi, 0},
		{"small negative", -0.5, -0.5},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := NormalizeAngleSigned(tc.in); err != nil || got != tc.want {
				t.Errorf("NormalizeAngleSigned(%v) = (%v, %v), want %v", tc.in, got, err, tc.want)
			}
		})
	}
	if got, err := NormalizeAngleSigned(-math.MaxFloat64); err != nil || got <= -math.Pi || got > math.Pi {
		t.Errorf("NormalizeAngleSigned(-MaxFloat64) = (%v, %v), want in (−π, π]", got, err)
	}
}

func TestNormalizeAngleSigned_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { _, err := NormalizeAngleSigned(math.NaN()); return err }, ErrNotFinite},
		{"+Inf", func() error { _, err := NormalizeAngleSigned(math.Inf(1)); return err }, ErrNotFinite},
		{"-Inf", func() error { _, err := NormalizeAngleSigned(math.Inf(-1)); return err }, ErrNotFinite},
	})
}

// Oracle: the result is in (−π, π] and differs from NormalizeAngle's by 0 or exactly 2π_f, within the
// one rounding NormalizeAngle may add (2^-50).
func TestNormalizeAngleSigned_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		theta := (rng.Float64() - 0.5) * math.Pow(10, float64(rng.IntN(12)))
		got, err := NormalizeAngleSigned(theta)
		u, _ := NormalizeAngle(theta)
		if err != nil || got <= -math.Pi || got > math.Pi {
			t.Fatalf("NormalizeAngleSigned(%v) = (%v, %v), out of (−π, π]", theta, got, err)
		}
		if d := u - got; math.Abs(d) > 0x1p-50 && math.Abs(d-twoPi) > 0x1p-50 {
			t.Fatalf("NormalizeAngleSigned(%v) = %v, NormalizeAngle = %v: not congruent", theta, got, u)
		}
	}
}

func TestNormalizeAngleSigned_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = NormalizeAngleSigned(-12.5) }); n != 0 {
		t.Errorf("NormalizeAngleSigned allocates %v times, want 0", n)
	}
}

func BenchmarkNormalizeAngleSigned(b *testing.B) {
	b.ReportAllocs()
	x := -12.5
	for b.Loop() {
		sink, _ = NormalizeAngleSigned(x)
	}
}
