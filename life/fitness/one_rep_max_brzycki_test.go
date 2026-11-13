package fitness

import (
	"math"
	"testing"
)

// HAND-DERIVED from the equation as printed in Lin 2021 Table 1 and Beia 2024 Table 1:
// 1RM = load / (1.0278 − 0.0278 × reps). Tolerance 1e-12 relative: one subtraction and one division.
func TestOneRepMaxBrzycki_example(t *testing.T) {
	cases := []struct {
		load float64
		reps int
		want float64
	}{
		{100, 10, 133.36889837289942}, // 100 / 0.7498
		{225, 5, 253.15031503150314},  // 225 / 0.8888
		{100, 1, 100},                 // denominator 1.0000: the load itself
		{100, 36, 3703.7037037036853}, // 100 / 0.0270, the largest whole rep count
	}
	for _, c := range cases {
		got, err := OneRepMaxBrzycki(c.load, c.reps)
		if err != nil || !relNear(got, c.want, 1e-12) {
			t.Errorf("OneRepMaxBrzycki(%v, %d) = (%v, %v), want %v", c.load, c.reps, got, err, c.want)
		}
	}
}

// Lin 2021 prints the same equation scaled by 100: 100 × weight ÷ (102.78 − 2.78 × reps). The two forms agree
// to 1e-12 relative for every rep count the function accepts.
func TestOneRepMaxBrzycki_linForm(t *testing.T) {
	for reps := 1; reps <= 36; reps++ {
		got, err := OneRepMaxBrzycki(80, reps)
		want := 100 * 80 / (102.78 - 2.78*float64(reps))
		if err != nil || !relNear(got, want, 1e-12) {
			t.Errorf("reps %d: got (%v, %v), want %v", reps, got, err, want)
		}
	}
}

// parity: ledger brzycki_1rm — the source entry is cite-only (no formula, no test file, no example); no
// source cases exist to reproduce.
func TestOneRepMaxBrzycki_parity(t *testing.T) {
	if _, err := OneRepMaxBrzycki(100, 10); err != nil {
		t.Fatalf("unexpected error %v", err)
	}
}

// The estimate rises with reps and is homogeneous in the load (kg in → kg out, lb in → lb out; D-063).
func TestOneRepMaxBrzycki_properties(t *testing.T) {
	prev := 0.0
	for reps := 1; reps <= 36; reps++ {
		v, err := OneRepMaxBrzycki(50, reps)
		if err != nil || v <= prev {
			t.Fatalf("reps %d: %v (%v), want > %v", reps, v, err, prev)
		}
		prev = v
		v2, _ := OneRepMaxBrzycki(100, reps)
		if !relNear(v2, 2*v, 1e-15) {
			t.Errorf("reps %d: not homogeneous: %v vs 2 × %v", reps, v2, v)
		}
	}
}

func TestOneRepMaxBrzycki_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"load NaN", func() error { _, err := OneRepMaxBrzycki(nan, 5); return err }, ErrNotFinite},
		{"load +Inf", func() error { _, err := OneRepMaxBrzycki(pInf, 5); return err }, ErrNotFinite},
		{"load -Inf", func() error { _, err := OneRepMaxBrzycki(nInf, 5); return err }, ErrNotFinite},
		{"load 0", func() error { _, err := OneRepMaxBrzycki(0, 5); return err }, ErrOutOfRange},
		{"load negative", func() error { _, err := OneRepMaxBrzycki(-1, 5); return err }, ErrOutOfRange},
		{"reps 0", func() error { _, err := OneRepMaxBrzycki(100, 0); return err }, ErrOutOfRange},
		{"reps negative", func() error { _, err := OneRepMaxBrzycki(100, -3); return err }, ErrOutOfRange},
		{"reps 37", func() error { _, err := OneRepMaxBrzycki(100, 37); return err }, ErrOutOfRange},
		{"reps max int", func() error { _, err := OneRepMaxBrzycki(100, math.MaxInt); return err }, ErrOutOfRange},
		{"result overflows", func() error { _, err := OneRepMaxBrzycki(math.MaxFloat64, 36); return err }, ErrDomain},
	})
}

func TestOneRepMaxBrzycki_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = OneRepMaxBrzycki(100, 10) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkOneRepMaxBrzycki(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink, _ = OneRepMaxBrzycki(100, 10)
	}
}
