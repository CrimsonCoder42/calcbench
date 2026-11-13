package cardio

import (
	"math"
	"testing"
)

// Quinn 2004 abstract: five items; hematocrit < 30% and SBP < 90 mm Hg on both sides of each edge.
func TestSanFranciscoSyncopeRuleCount(t *testing.T) {
	tests := []struct {
		ecg, sob, chf bool
		hct, sbp      float64
		want          int
	}{
		{false, false, false, 40, 120, 0},
		{false, false, false, 30, 90, 0},
		{false, false, false, 29.9, 90, 1},
		{false, false, false, 30, 89.9, 1},
		{true, false, false, 40, 120, 1},
		{false, true, false, 40, 120, 1},
		{false, false, true, 40, 120, 1},
		{true, true, true, 29, 80, 5},
	}
	for _, tc := range tests {
		c, any, err := SanFranciscoSyncopeRuleCount(tc.ecg, tc.sob, tc.chf, tc.hct, tc.sbp)
		if err != nil || c != tc.want || any != (tc.want > 0) {
			t.Errorf("%+v = (%d, %v, %v), want %d", tc, c, any, err, tc.want)
		}
	}
}

// parity: ledger sf_syncope_rule — no source test file or JSON example; source case 1 is the doc Example.
func TestSanFranciscoSyncopeRuleCount_parity(t *testing.T) {
	c, any, _ := SanFranciscoSyncopeRuleCount(false, false, false, 29, 120) // parity: ledger sf_syncope_rule, source case 1
	if c != 1 || !any {
		t.Errorf("(%d, %v)", c, any)
	}
}

func TestSanFranciscoSyncopeRuleCount_errors(t *testing.T) {
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		_, _, err := SanFranciscoSyncopeRuleCount(false, false, false, v, 120)
		checkErr(t, "hct", err, ErrNotFinite)
		_, _, err = SanFranciscoSyncopeRuleCount(false, false, false, 40, v)
		checkErr(t, "sbp", err, ErrNotFinite)
	}
	for _, in := range [][2]float64{{4.9, 120}, {75.1, 120}, {40, 29.9}, {40, 300.1}} {
		_, _, err := SanFranciscoSyncopeRuleCount(false, false, false, in[0], in[1])
		checkErr(t, "range", err, ErrOutOfRange)
	}
}

func TestSanFranciscoSyncopeRuleCount_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		sinkInt, sinkBool, _ = SanFranciscoSyncopeRuleCount(true, false, false, 29, 120)
	}); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkSanFranciscoSyncopeRuleCount(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, sinkBool, _ = SanFranciscoSyncopeRuleCount(true, false, false, 29, 120)
	}
}
