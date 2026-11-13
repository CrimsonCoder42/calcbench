package cardio

import (
	"math"
	"testing"
)

// Rodger 2008: four characteristics, thresholds ≥ 250, ≥ 30, ≥ 65; both sides of each.
func TestHERDOO2Count_items(t *testing.T) {
	tests := []struct {
		age, bmi, dd float64
		legs         bool
		want         int
	}{
		{64.9, 29.9, 249.9, false, 0},
		{65, 29.9, 249.9, false, 1},
		{40, 30, 249.9, false, 1},
		{40, 20, 250, false, 1},
		{70, 31, 200, false, 2},
		{70, 31, 300, true, 4},
	}
	for _, tc := range tests {
		if got, err := HERDOO2Count(tc.age, tc.bmi, tc.dd, tc.legs); err != nil || got != tc.want {
			t.Errorf("(%v, %v, %v, %v) = (%d, %v), want %d", tc.age, tc.bmi, tc.dd, tc.legs, got, err, tc.want)
		}
	}
	_, err := HERDOO2Count(math.NaN(), 25, 100, false)
	checkErr(t, "NaN", err, ErrNotFinite)
	_, err = HERDOO2Count(40, -1, 100, false)
	checkErr(t, "bmi < 0", err, ErrOutOfRange)
	_, err = HERDOO2Count(40, 25, math.Inf(1), false)
	checkErr(t, "dd +Inf", err, ErrNotFinite)
}

func TestHERDOO2Count_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = HERDOO2Count(70, 31, 200, false) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkHERDOO2Count(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = HERDOO2Count(70, 31, 200, false)
	}
}
