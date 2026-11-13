package trig

import (
	"math"
	"testing"
)

// Worked example (hand-derived from the POSIX remainder() definition): from 350° to 10° is +20°, from
// 10° to 350° is −20°, a half turn is +π. The inputs are not exact in radians: tolerance a few ulps.
func TestAngleDifference_example(t *testing.T) {
	tests := []struct{ from, to, want float64 }{
		{radOf(350), radOf(10), radOf(20)},
		{radOf(10), radOf(350), -radOf(20)},
		{0, math.Pi, math.Pi},
		{math.Pi, 0, math.Pi},
	}
	for _, tc := range tests {
		if got, err := AngleDifference(tc.from, tc.to); err != nil || !near(got, tc.want, 1e-15) {
			t.Errorf("AngleDifference(%v, %v) = (%v, %v), want %v", tc.from, tc.to, got, err, tc.want)
		}
	}
}

func TestAngleDifference_cases(t *testing.T) {
	tests := []struct {
		name           string
		from, to, want float64
	}{
		{"same angle", 1, 1, 0},
		{"zero to zero", 0, 0, 0},
		{"quarter turn ahead", 0, math.Pi / 2, math.Pi / 2},
		{"quarter turn behind", math.Pi / 2, 0, -math.Pi / 2},
		{"several turns apart", 0, 3 * twoPi, 0},
		{"huge inputs that fit", -1e307, 1e307, remainderTwoPi(2e307)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := AngleDifference(tc.from, tc.to); err != nil || got != tc.want {
				t.Errorf("AngleDifference(%v, %v) = (%v, %v), want %v", tc.from, tc.to, got, err, tc.want)
			}
		})
	}
}

func TestAngleDifference_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"from NaN", func() error { _, err := AngleDifference(math.NaN(), 0); return err }, ErrNotFinite},
		{"to +Inf", func() error { _, err := AngleDifference(0, math.Inf(1)); return err }, ErrNotFinite},
		{"from -Inf", func() error { _, err := AngleDifference(math.Inf(-1), 0); return err }, ErrNotFinite},
		{"difference overflows", func() error { _, err := AngleDifference(-math.MaxFloat64, math.MaxFloat64); return err }, ErrOverflow},
	})
}

// Oracle: from + difference is congruent to to, and the difference is in (−π, π]. Inputs within ±50 rad,
// so the congruence holds within a few ulps of 50 (1e-13).
func TestAngleDifference_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		from, to := (rng.Float64()-0.5)*100, (rng.Float64()-0.5)*100
		d, err := AngleDifference(from, to)
		if err != nil || d <= -math.Pi || d > math.Pi {
			t.Fatalf("AngleDifference(%v, %v) = (%v, %v)", from, to, d, err)
		}
		if r := remainderTwoPi(from + d - to); math.Abs(r) > 1e-13 {
			t.Fatalf("AngleDifference(%v, %v) = %v: from + d − to = %v (mod 2π)", from, to, d, r)
		}
	}
}

func TestAngleDifference_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = AngleDifference(6.1, 0.2) }); n != 0 {
		t.Errorf("AngleDifference allocates %v times, want 0", n)
	}
}

func BenchmarkAngleDifference(b *testing.B) {
	b.ReportAllocs()
	from, to := 6.1, 0.2
	for b.Loop() {
		sink, _ = AngleDifference(from, to)
	}
}
