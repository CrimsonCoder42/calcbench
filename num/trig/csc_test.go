package trig

import (
	"math"
	"testing"
)

// Worked example, hand-derived from DLMF 4.14.5 (csc z = 1/sin z): csc(π/2) = 1, csc(π/6) = 2 (to
// rounding: math.Pi/6 is not π/6), and csc(math.Pi) = 1/1.2246467991473532e-16 = 8.165619676597685e15.
// The 4e-15 tolerance at the multiples of π/2: Go's math.Sincos reduces the argument with a three-part
// π/4, so the tiny cos or sin there carries about 1.5e-15 relative error (measured on darwin/arm64).
func TestCsc_dlmfExample(t *testing.T) {
	for _, tc := range []struct{ x, want, tol float64 }{{math.Pi / 2, 1, 0}, {math.Pi / 6, 2, ulps4}, {math.Pi, 8.165619676597685e15, 4e-15}} {
		if got, err := Csc(tc.x); err != nil || !relNear(got, tc.want, tc.tol) {
			t.Errorf("Csc(%v) = (%v, %v), want %v", tc.x, got, err, tc.want)
		}
	}
}

func TestCsc_cases(t *testing.T) {
	for _, tc := range []struct {
		name string
		x    float64
	}{{"negative", -1}, {"odd", 1}, {"tiny but finite result", 1e-300}, {"large argument", 1e300}} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Csc(tc.x)
			if err != nil || math.Abs(got) < 1 || got != 1/math.Sin(tc.x) {
				t.Errorf("Csc(%v) = (%v, %v)", tc.x, got, err)
			}
		})
	}
	if p, _ := Csc(0.5); func() bool { m, _ := Csc(-0.5); return m != -p }() {
		t.Error("Csc is not odd at 0.5")
	}
}

func TestCsc_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { _, err := Csc(math.NaN()); return err }, ErrNotFinite},
		{"+Inf", func() error { _, err := Csc(math.Inf(1)); return err }, ErrNotFinite},
		{"+0 pole", func() error { _, err := Csc(0); return err }, ErrDomain},
		{"-0 pole", func() error { _, err := Csc(math.Copysign(0, -1)); return err }, ErrDomain},
		{"subnormal overflows", func() error { _, err := Csc(5e-324); return err }, ErrOverflow},
		{"negative subnormal overflows", func() error { _, err := Csc(-1e-310); return err }, ErrOverflow},
	})
}

// Oracle: csc² x − cot² x = 1 (A&S 4.3.12) on x with |sin x| ≥ 0.1, to about 100 ulps (1e-13).
func TestCsc_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		x := (rng.Float64() - 0.5) * 20
		if math.Abs(math.Sin(x)) < 0.1 {
			continue
		}
		c, err := Csc(x)
		ct := math.Cos(x) / math.Sin(x)
		if err != nil || math.Abs(c*c-ct*ct-1) > 1e-13 {
			t.Fatalf("Csc(%v) = (%v, %v): csc² − cot² = %v", x, c, err, c*c-ct*ct)
		}
	}
}

func TestCsc_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = Csc(0.7) }); n != 0 {
		t.Errorf("Csc allocates %v times, want 0", n)
	}
}

func BenchmarkCsc(b *testing.B) {
	b.ReportAllocs()
	x := 0.7
	for b.Loop() {
		sink, _ = Csc(x)
	}
}
