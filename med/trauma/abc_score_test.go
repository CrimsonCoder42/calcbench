package trauma

import (
	"errors"
	"math"
	"testing"
)

var (
	sinkInt  int
	sinkBool bool
)

// checkErr fails unless err wraps want.
func checkErr(t *testing.T, name string, err, want error) {
	t.Helper()
	if !errors.Is(err, want) {
		t.Errorf("%s: err = %v, want %v", name, err, want)
	}
}

var nonFinite = []float64{math.NaN(), math.Inf(1), math.Inf(-1)}

// Nunez 2009 abstract: four nonweighted parameters; SBP "90 mm Hg or less", HR "≥120 bpm". HAND-DERIVED.
func TestABCScore_source(t *testing.T) {
	tests := []struct {
		pen, fast bool
		sbp, hr   float64
		want      int
	}{
		{false, false, 120, 80, 0},
		{true, false, 85, 125, 3},
		{true, true, 90, 120, 4},         // both edges included
		{false, false, 90.01, 119.99, 0}, // just outside both edges
		{false, true, 91, 119, 1},
	}
	for i, tc := range tests {
		if got, err := ABCScore(tc.pen, tc.fast, tc.sbp, tc.hr); err != nil || got != tc.want {
			t.Errorf("case %d: (%d, %v), want %d", i, got, err, tc.want)
		}
	}
}

// parity: ledger abc_massive_transfusion, source case 1 (the app sums four 0/1 inputs; JSON carries no example).
// Here the two vital-sign items are measured values, so the app's sbpLow = hrHigh = 1 is SBP 80, HR 130.
func TestABCScore_parity(t *testing.T) {
	if got, _ := ABCScore(true, true, 80, 130); got != 4 {
		t.Errorf("all four → %d, want 4", got)
	}
}

func TestABCScore_errors(t *testing.T) {
	for _, v := range nonFinite {
		_, err := ABCScore(false, false, v, 80)
		checkErr(t, "sbp non-finite", err, ErrNotFinite)
		_, err = ABCScore(false, false, 120, v)
		checkErr(t, "hr non-finite", err, ErrNotFinite)
	}
	_, err := ABCScore(false, false, -1, 80)
	checkErr(t, "sbp < 0", err, ErrOutOfRange)
	_, err = ABCScore(false, false, 120, -0.1)
	checkErr(t, "hr < 0", err, ErrOutOfRange)
	if _, err := ABCScore(false, false, 0, 0); err != nil {
		t.Errorf("zero vitals: %v", err)
	}
}

func TestABCScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = ABCScore(true, false, 85, 125) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkABCScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = ABCScore(true, false, 85, 125)
	}
}
