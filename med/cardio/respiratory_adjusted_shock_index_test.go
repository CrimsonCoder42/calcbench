package cardio

import (
	"math"
	"testing"
)

// HAND-DERIVED from the printed definition RASI = HR / SBP × RR / 10 (Caputo 2018). Tolerance 1e-12 relative.
func TestRespiratoryAdjustedShockIndex_values(t *testing.T) {
	tests := []struct {
		call func() (float64, error)
		want float64
	}{
		{func() (float64, error) { return RespiratoryAdjustedShockIndex(100, 100, 20) }, 2},
		{func() (float64, error) { return RespiratoryAdjustedShockIndex(90, 120, 19) }, 0.75 * 1.9},
	}
	for i, tc := range tests {
		got, err := tc.call()
		if err != nil || math.Abs(got-tc.want) > 1e-12*math.Abs(tc.want) {
			t.Errorf("case %d: (%v, %v), want %v", i, got, err, tc.want)
		}
	}
}

func TestRespiratoryAdjustedShockIndex_errors(t *testing.T) {
	var err error
	_, err = RespiratoryAdjustedShockIndex(100, 100, math.NaN())
	checkErr(t, "100, 100, math.NaN()", err, ErrNotFinite)
	_, err = RespiratoryAdjustedShockIndex(math.Inf(1), 100, 20)
	checkErr(t, "math.Inf(1), 100, 20", err, ErrNotFinite)
	_, err = RespiratoryAdjustedShockIndex(100, 100, 0)
	checkErr(t, "100, 100, 0", err, ErrOutOfRange)
	_, err = RespiratoryAdjustedShockIndex(100, -5, 20)
	checkErr(t, "100, -5, 20", err, ErrOutOfRange)
}

func TestRespiratoryAdjustedShockIndex_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = RespiratoryAdjustedShockIndex(100, 100, 20) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkRespiratoryAdjustedShockIndex(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = RespiratoryAdjustedShockIndex(100, 100, 20)
	}
}
