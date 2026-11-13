package trig

import (
	"math"
	"testing"
)

// Worked example: DLMF Table 4.23.1, arccsc 1 = π/2 and arccsc(−1) = −π/2; arccsc 2 = π/6 (hand-derived:
// sin(π/6) = 1/2).
func TestAcsc_dlmfTable(t *testing.T) {
	for _, tc := range []struct{ x, want, tol float64 }{{1, math.Pi / 2, 0}, {-1, -math.Pi / 2, 0}, {2, math.Pi / 6, ulps4}} {
		if got, err := Acsc(tc.x); err != nil || !near(got, tc.want, tc.tol) {
			t.Errorf("Acsc(%v) = (%v, %v), want %v", tc.x, got, err, tc.want)
		}
	}
}

func TestAcsc_cases(t *testing.T) {
	tests := []struct {
		name    string
		x, want float64
	}{
		{"large positive tends to +0", 1e300, 1e-300},
		{"large negative tends to −0", -1e300, -1e-300},
		{"largest float64", math.MaxFloat64, 1 / math.MaxFloat64},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := Acsc(tc.x); err != nil || !relNear(got, tc.want, ulps4) {
				t.Errorf("Acsc(%v) = (%v, %v), want %v", tc.x, got, err, tc.want)
			}
		})
	}
	// Odd (DLMF 4.23.13).
	p, _ := Acsc(3)
	if m, _ := Acsc(-3); m != -p {
		t.Errorf("Acsc(−3) = %v, want −Acsc(3) = %v", m, -p)
	}
}

func TestAcsc_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { _, err := Acsc(math.NaN()); return err }, ErrNotFinite},
		{"+Inf", func() error { _, err := Acsc(math.Inf(1)); return err }, ErrNotFinite},
		{"-Inf", func() error { _, err := Acsc(math.Inf(-1)); return err }, ErrNotFinite},
		{"zero", func() error { _, err := Acsc(0); return err }, ErrDomain},
		{"inside (−1, 1)", func() error { _, err := Acsc(-0.999); return err }, ErrDomain},
	})
}

// Oracle: csc(Acsc x) = x for 1.5 ≤ |x| ≤ 100, to 1e-12 relative (as for Asec).
func TestAcsc_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		x := 1.5 + rng.Float64()*98.5
		if rng.IntN(2) == 0 {
			x = -x
		}
		y, err := Acsc(x)
		if err != nil || math.Abs(y) > math.Pi/2 || !relNear(1/math.Sin(y), x, 1e-12) {
			t.Fatalf("Acsc(%v) = (%v, %v)", x, y, err)
		}
	}
}

func TestAcsc_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = Acsc(2.5) }); n != 0 {
		t.Errorf("Acsc allocates %v times, want 0", n)
	}
}

func BenchmarkAcsc(b *testing.B) {
	b.ReportAllocs()
	x := 2.5
	for b.Loop() {
		sink, _ = Acsc(x)
	}
}
