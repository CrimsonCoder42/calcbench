package renal

import (
	"errors"
	"testing"
)

var sinkB6Int int

// Mehta 2007 Table 2: every creatinine and urine-output edge on both sides. HAND-DERIVED. Normal urine output
// (1.0 mL/kg/h for 24 h) isolates the creatinine stage; baseline = current isolates the urine-output stage.
func TestAKINStage_table2(t *testing.T) {
	tests := []struct {
		name               string
		base, cur, uo, hrs float64
		rrt                bool
		want               int
	}{
		{"no change", 1.0, 1.0, 1.0, 24, false, 0},
		{"rise 0.29", 1.0, 1.29, 1.0, 24, false, 0},
		{"rise 0.3", 1.0, 1.3, 1.0, 24, false, 1},
		{"rise 0.3 from 0.9", 0.9, 1.2, 1.0, 24, false, 1},
		{"149% of small baseline", 0.4, 0.596, 1.0, 24, false, 0},
		{"150% of small baseline", 0.4, 0.6, 1.0, 24, false, 1},
		{"200%", 1.0, 2.0, 1.0, 24, false, 1},
		{"201%", 1.0, 2.01, 1.0, 24, false, 2},
		{"300%", 1.0, 3.0, 1.0, 24, false, 2},
		{"301%", 1.0, 3.01, 1.0, 24, false, 3},
		{"≥4.0 with rise 0.5", 3.5, 4.0, 1.0, 24, false, 3},
		{"≥4.0 with rise 0.49", 3.51, 4.0, 1.0, 24, false, 1},
		{"3.99 with rise 0.5", 3.49, 3.99, 1.0, 24, false, 1},
		{"uo 0.49 for 6 h", 1, 1, 0.49, 6, false, 0},
		{"uo 0.49 for 6.1 h", 1, 1, 0.49, 6.1, false, 1},
		{"uo 0.5 for 13 h", 1, 1, 0.5, 13, false, 0},
		{"uo 0.49 for 12 h", 1, 1, 0.49, 12, false, 1},
		{"uo 0.49 for 12.1 h", 1, 1, 0.49, 12.1, false, 2},
		{"uo 0.29 for 23.9 h", 1, 1, 0.29, 23.9, false, 2},
		{"uo 0.29 for 24 h", 1, 1, 0.29, 24, false, 3},
		{"uo 0.3 for 24 h", 1, 1, 0.3, 24, false, 2},
		{"anuria 11.9 h", 1, 1, 0, 11.9, false, 1},
		{"anuria 12 h", 1, 1, 0, 12, false, 3},
		{"higher of the two", 1.0, 1.3, 0.4, 13, false, 2},
		{"RRT", 1.0, 1.0, 1.0, 24, true, 3},
	}
	for _, tc := range tests {
		if got, err := AKINStage(tc.base, tc.cur, tc.uo, tc.hrs, tc.rrt); err != nil || got != tc.want {
			t.Errorf("%s: (%d, %v), want %d", tc.name, got, err, tc.want)
		}
	}
}

// parity: ledger akin_classification, source cases. The app takes the creatinine stage and urine-output stage as
// 0–3 pickers and returns their maximum, RRT → 3; here the stages are derived from the measured values.
// creatinine picker 2 with urine picker 1 → 2.
func TestAKINStage_parity(t *testing.T) {
	if got, _ := AKINStage(1.0, 2.5, 0.4, 7, false); got != 2 {
		t.Errorf("got %d, want 2", got)
	}
	if got, _ := AKINStage(1.0, 1.0, 1.0, 24, true); got != 3 {
		t.Errorf("rrt: got %d, want 3", got)
	}
}

func TestAKINStage_errors(t *testing.T) {
	var cases []errCase
	for _, v := range []float64{nan, pInf, nInf} {
		v := v
		cases = append(cases,
			errCase{"base non-finite", func() error { _, e := AKINStage(v, 1, 1, 1, false); return e }, ErrNotFinite},
			errCase{"cur non-finite", func() error { _, e := AKINStage(1, v, 1, 1, false); return e }, ErrNotFinite},
			errCase{"uo non-finite", func() error { _, e := AKINStage(1, 1, v, 1, false); return e }, ErrNotFinite},
			errCase{"hours non-finite", func() error { _, e := AKINStage(1, 1, 1, v, false); return e }, ErrNotFinite},
		)
	}
	cases = append(cases,
		errCase{"base 0", func() error { _, e := AKINStage(0, 1, 1, 1, false); return e }, ErrOutOfRange},
		errCase{"cur 0", func() error { _, e := AKINStage(1, 0, 1, 1, false); return e }, ErrOutOfRange},
		errCase{"uo < 0", func() error { _, e := AKINStage(1, 1, -0.1, 1, false); return e }, ErrOutOfRange},
		errCase{"hours < 0", func() error { _, e := AKINStage(1, 1, 1, -1, false); return e }, ErrOutOfRange},
	)
	runErrCases(t, cases)
}

func TestAKINStage_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkB6Int, _ = AKINStage(1.0, 2.5, 0.4, 7, false) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkAKINStage(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkB6Int, _ = AKINStage(1.0, 2.5, 0.4, 7, false)
	}
}

// Creatinines whose thousandths overflow float64 are ErrDomain, never a stage (review finding: MaxFloat64 twice
// returned stage 1 because Inf ≥ Inf).
func TestAKINStage_overflow(t *testing.T) {
	if _, err := AKINStage(1.7976931348623157e308, 1.7976931348623157e308, 0, 0, false); !errors.Is(err, ErrDomain) {
		t.Errorf("AKINStage(MaxFloat64, MaxFloat64, …) err = %v, want ErrDomain", err)
	}
}
