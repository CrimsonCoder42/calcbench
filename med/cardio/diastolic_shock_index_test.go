package cardio

import (
	"math"
	"testing"
)

// HAND-DERIVED from the printed definition DSI = HR / DAP (Ospina-Tascón 2020). Tolerance 1e-12 relative.
func TestDiastolicShockIndex_values(t *testing.T) {
	tests := []struct {
		call func() (float64, error)
		want float64
	}{
		{func() (float64, error) { return DiastolicShockIndex(110, 50) }, 2.2},
		{func() (float64, error) { return DiastolicShockIndex(60, 60) }, 1},
	}
	for i, tc := range tests {
		got, err := tc.call()
		if err != nil || math.Abs(got-tc.want) > 1e-12*math.Abs(tc.want) {
			t.Errorf("case %d: (%v, %v), want %v", i, got, err, tc.want)
		}
	}
}

func TestDiastolicShockIndex_errors(t *testing.T) {
	var err error
	_, err = DiastolicShockIndex(math.Inf(1), 50)
	checkErr(t, "math.Inf(1), 50", err, ErrNotFinite)
	_, err = DiastolicShockIndex(110, math.NaN())
	checkErr(t, "110, math.NaN()", err, ErrNotFinite)
	_, err = DiastolicShockIndex(110, 0)
	checkErr(t, "110, 0", err, ErrOutOfRange)
	_, err = DiastolicShockIndex(0, 50)
	checkErr(t, "0, 50", err, ErrOutOfRange)
}

func TestDiastolicShockIndex_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = DiastolicShockIndex(110, 50) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkDiastolicShockIndex(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = DiastolicShockIndex(110, 50)
	}
}
