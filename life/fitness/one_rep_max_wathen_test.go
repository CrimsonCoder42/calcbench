package fitness

import (
	"math"
	"testing"
)

// Lin 2021 (PMC8080630), Methods: "if a participant can perform an 80 kg squat for 7 repetitions at most, the
// 1RM=100×80 kg÷(48.8+53.8e^−0.075×reps)= 99 kg". The source prints whole kilograms; the unrounded value
// 99.22397105515728 is HAND-DERIVED (8000 / 80.6257…), tolerance 1e-12 relative (one exp, one division).
func TestOneRepMaxWathen_example(t *testing.T) {
	got, err := OneRepMaxWathen(80, 7)
	if err != nil || math.Round(got) != 99 || !relNear(got, 99.22397105515728, 1e-12) {
		t.Errorf("OneRepMaxWathen(80, 7) = (%v, %v), want 99.22397105515728 (Lin 2021: 99 kg)", got, err)
	}
}

// HAND-DERIVED from the equation 100 × load / (48.8 + 53.8 × e^(−0.075 × reps)); tolerance 1e-12 relative.
func TestOneRepMaxWathen_cases(t *testing.T) {
	cases := []struct {
		load float64
		reps int
		want float64
	}{
		{100, 1, 101.30419059087079},
		{100, 10, 134.74669948168537},
		{225, 5, 262.3106369051758},
	}
	for _, c := range cases {
		got, err := OneRepMaxWathen(c.load, c.reps)
		if err != nil || !relNear(got, c.want, 1e-12) {
			t.Errorf("OneRepMaxWathen(%v, %d) = (%v, %v), want %v", c.load, c.reps, got, err, c.want)
		}
	}
	// Very large reps: e^(−0.075 × reps) underflows to 0, the estimate tends to 100 × load / 48.8.
	got, err := OneRepMaxWathen(48.8, math.MaxInt)
	if err != nil || !relNear(got, 100, 1e-12) {
		t.Errorf("OneRepMaxWathen(48.8, MaxInt) = (%v, %v), want 100", got, err)
	}
}

// parity: ledger wathen_1rm — the source entry is cite-only (no formula, no test file, no example); no
// source cases exist to reproduce.
func TestOneRepMaxWathen_parity(t *testing.T) {
	if _, err := OneRepMaxWathen(100, 10); err != nil {
		t.Fatalf("unexpected error %v", err)
	}
}

// The estimate rises with reps and is homogeneous in the load (D-063).
func TestOneRepMaxWathen_properties(t *testing.T) {
	prev := 0.0
	for reps := 1; reps <= 100; reps++ {
		v, err := OneRepMaxWathen(50, reps)
		if err != nil || v <= prev {
			t.Fatalf("reps %d: %v (%v), want > %v", reps, v, err, prev)
		}
		prev = v
		v2, _ := OneRepMaxWathen(100, reps)
		if !relNear(v2, 2*v, 1e-15) {
			t.Errorf("reps %d: not homogeneous: %v vs 2 × %v", reps, v2, v)
		}
	}
}

func TestOneRepMaxWathen_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"load NaN", func() error { _, err := OneRepMaxWathen(nan, 5); return err }, ErrNotFinite},
		{"load +Inf", func() error { _, err := OneRepMaxWathen(pInf, 5); return err }, ErrNotFinite},
		{"load -Inf", func() error { _, err := OneRepMaxWathen(nInf, 5); return err }, ErrNotFinite},
		{"load 0", func() error { _, err := OneRepMaxWathen(0, 5); return err }, ErrOutOfRange},
		{"load negative", func() error { _, err := OneRepMaxWathen(-1, 5); return err }, ErrOutOfRange},
		{"reps 0", func() error { _, err := OneRepMaxWathen(100, 0); return err }, ErrOutOfRange},
		{"reps negative", func() error { _, err := OneRepMaxWathen(100, -1); return err }, ErrOutOfRange},
		{"result overflows", func() error { _, err := OneRepMaxWathen(math.MaxFloat64, 10); return err }, ErrDomain},
	})
}

func TestOneRepMaxWathen_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = OneRepMaxWathen(100, 10) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkOneRepMaxWathen(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink, _ = OneRepMaxWathen(100, 10)
	}
}
