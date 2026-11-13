package cardio

import (
	"math"
	"testing"
)

// Liu 2012: MSI = HR / MAP. Tolerance 1e-12 relative: one division.
func TestModifiedShockIndex(t *testing.T) {
	tests := []struct{ hr, mapp, want float64 }{
		{110, 70, 110.0 / 70}, // doc Example, 1.5714…
		{70, 100, 0.7},        // the paper's lower threshold
		{130, 100, 1.3},       // the paper's upper threshold
	}
	for _, tc := range tests {
		got, err := ModifiedShockIndex(tc.hr, tc.mapp)
		if err != nil || math.Abs(got-tc.want) > 1e-12*tc.want {
			t.Errorf("(%v, %v) = (%v, %v), want %v", tc.hr, tc.mapp, got, err, tc.want)
		}
	}
}

// parity: ledger modified_shock_index — no source test file or JSON example. The app derives MAP as
// (SBP + 2·DBP)/3 itself; with SBP 120, DBP 60 (MAP 80) and HR 100 it shows 1.25 (source case 1).
func TestModifiedShockIndex_parity(t *testing.T) {
	if got, _ := ModifiedShockIndex(100, 80); math.Abs(got-1.25) > 1e-12 { // parity: ledger modified_shock_index, source case 1
		t.Errorf("got %v", got)
	}
}

func TestModifiedShockIndex_errors(t *testing.T) {
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		_, err := ModifiedShockIndex(v, 80)
		checkErr(t, "hr", err, ErrNotFinite)
		_, err = ModifiedShockIndex(100, v)
		checkErr(t, "map", err, ErrNotFinite)
	}
	for _, in := range [][2]float64{{0, 80}, {100, 0}, {-1, 80}} {
		_, err := ModifiedShockIndex(in[0], in[1])
		checkErr(t, "range", err, ErrOutOfRange)
	}
}

func TestModifiedShockIndex_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = ModifiedShockIndex(110, 70) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkModifiedShockIndex(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = ModifiedShockIndex(110, 70)
	}
}
