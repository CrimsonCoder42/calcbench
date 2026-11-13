package fitness

import (
	"math"
	"testing"
)

// HAND-DERIVED: 1 RM 200 lb at 75 % → 200 × (75 / 100) = 150 (McBride 2009's protocol load "75% of a 1RM").
// Tolerance 1e-15 relative: two correctly rounded operations (0.9 is inexact in binary).
func TestLoadAtPercentOneRepMax_example(t *testing.T) {
	cases := []struct{ orm, pct, want float64 }{
		{200, 75, 150},
		{100, 60, 60},
		{140, 90, 126},
		{100, 100, 100},
	}
	for _, c := range cases {
		if got, err := LoadAtPercentOneRepMax(c.orm, c.pct); err != nil || !relNear(got, c.want, 1e-15) {
			t.Errorf("LoadAtPercentOneRepMax(%v, %v) = (%v, %v), want %v", c.orm, c.pct, got, err, c.want)
		}
	}
}

// parity: ledger training_load_from_1rm. The source is not runnable (no formula, no test case); its JSON
// shape is one_rep_max (1–1000) × percent_1rm (1–100) → training_load. Its two extreme cases, HAND-DERIVED.
func TestLoadAtPercentOneRepMax_parity(t *testing.T) {
	if got, err := LoadAtPercentOneRepMax(1000, 100); err != nil || got != 1000 { // parity: ledger training_load_from_1rm, JSON bounds case 1
		t.Errorf("(1000, 100) = (%v, %v), want 1000", got, err)
	}
	if got, err := LoadAtPercentOneRepMax(1, 1); err != nil || !relNear(got, 0.01, 1e-15) { // parity: ledger training_load_from_1rm, JSON bounds case 2
		t.Errorf("(1, 1) = (%v, %v), want 0.01", got, err)
	}
}

// Edges: percent 100 accepted, just above refused; percent just above 0 accepted, 0 refused.
func TestLoadAtPercentOneRepMax_errors(t *testing.T) {
	if _, err := LoadAtPercentOneRepMax(100, math.SmallestNonzeroFloat64); err != nil {
		t.Errorf("smallest percent: %v", err)
	}
	call := func(o, p float64) func() error {
		return func() error { _, err := LoadAtPercentOneRepMax(o, p); return err }
	}
	runErrCases(t, []errCase{
		{"percent just above 100", call(100, math.Nextafter(100, 200)), ErrOutOfRange},
		{"percent 0", call(100, 0), ErrOutOfRange},
		{"percent negative", call(100, -5), ErrOutOfRange},
		{"1RM 0", call(0, 50), ErrOutOfRange},
		{"1RM negative", call(-1, 50), ErrOutOfRange},
		{"1RM NaN", call(nan, 50), ErrNotFinite},
		{"1RM +Inf", call(pInf, 50), ErrNotFinite},
		{"1RM -Inf", call(nInf, 50), ErrNotFinite},
		{"percent NaN", call(100, nan), ErrNotFinite},
		{"percent +Inf", call(100, pInf), ErrNotFinite},
		{"percent -Inf", call(100, nInf), ErrNotFinite},
	})
	// Largest 1RM at 100 %: percent / 100 is applied first, so MaxFloat64 does not overflow.
	if got, err := LoadAtPercentOneRepMax(math.MaxFloat64, 100); err != nil || got != math.MaxFloat64 {
		t.Errorf("MaxFloat64 at 100 %%: (%v, %v), want MaxFloat64", got, err)
	}
}

func TestLoadAtPercentOneRepMax_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = LoadAtPercentOneRepMax(200, 75) }); n != 0 {
		t.Errorf("LoadAtPercentOneRepMax allocates %v times, want 0", n)
	}
}

func BenchmarkLoadAtPercentOneRepMax(b *testing.B) {
	b.ReportAllocs()
	o, p := 200.0, 75.0
	for b.Loop() {
		sink, _ = LoadAtPercentOneRepMax(o, p)
	}
}
