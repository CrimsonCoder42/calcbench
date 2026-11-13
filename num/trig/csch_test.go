package trig

import (
	"math"
	"testing"
)

// Worked example, hand-derived from DLMF 4.28.5 (csch z = 1/sinh z): csch 1 = 2/(e − 1/e) =
// 0.8509181282393216, and csch(−1) is its negative.
func TestCsch_dlmfExample(t *testing.T) {
	for _, tc := range []struct{ x, want float64 }{{1, 0.8509181282393216}, {-1, -0.8509181282393216}} {
		if got, err := Csch(tc.x); err != nil || !relNear(got, tc.want, ulps4) {
			t.Errorf("Csch(%v) = (%v, %v), want %v", tc.x, got, err, tc.want)
		}
	}
}

func TestCsch_cases(t *testing.T) {
	tests := []struct {
		name    string
		x, want float64
	}{
		{"sinh overflows, result +0", 800, 0},
		{"negative large, result −0", -800, math.Copysign(0, -1)},
		{"tiny", 1e-300, 1e300},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Csch(tc.x)
			if err != nil || !relNear(got, tc.want, ulps4) || math.Signbit(got) != math.Signbit(tc.want) {
				t.Errorf("Csch(%v) = (%v, %v), want %v", tc.x, got, err, tc.want)
			}
		})
	}
}

func TestCsch_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { _, err := Csch(math.NaN()); return err }, ErrNotFinite},
		{"+Inf", func() error { _, err := Csch(math.Inf(1)); return err }, ErrNotFinite},
		{"+0 pole", func() error { _, err := Csch(0); return err }, ErrDomain},
		{"-0 pole", func() error { _, err := Csch(math.Copysign(0, -1)); return err }, ErrDomain},
		{"subnormal overflows", func() error { _, err := Csch(5e-324); return err }, ErrOverflow},
	})
}

// Oracle: coth² x − csch² x = 1 (A&S 4.5.18) for 0.1 ≤ |x| ≤ 20, where both terms are at most about
// 100: tolerance 1e-13.
func TestCsch_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		x := 0.1 + rng.Float64()*19.9
		if rng.IntN(2) == 0 {
			x = -x
		}
		c, err := Csch(x)
		ct := 1 / math.Tanh(x)
		if err != nil || math.Abs(ct*ct-c*c-1) > 1e-13 {
			t.Fatalf("Csch(%v) = (%v, %v)", x, c, err)
		}
	}
}

func TestCsch_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = Csch(0.7) }); n != 0 {
		t.Errorf("Csch allocates %v times, want 0", n)
	}
}

func BenchmarkCsch(b *testing.B) {
	b.ReportAllocs()
	x := 0.7
	for b.Loop() {
		sink, _ = Csch(x)
	}
}
