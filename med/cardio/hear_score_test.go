package cardio

import (
	"math"
	"testing"
)

// Smith 2020: HEAR = HEART without troponin; 0–1 is the paper's very-low-risk group.
func TestHEARScore(t *testing.T) {
	tests := []struct {
		h, e int
		age  float64
		r    int
		want int
	}{
		{0, 0, 40, 1, 1}, // doc Example
		{0, 0, 40, 0, 0},
		{0, 0, 45, 0, 1},
		{0, 0, 65.1, 0, 2},
		{2, 2, 80, 2, 8},
	}
	for _, tc := range tests {
		if got, err := HEARScore(tc.h, tc.e, tc.age, tc.r); err != nil || got != tc.want {
			t.Errorf("%+v = (%d, %v), want %d", tc, got, err, tc.want)
		}
	}
	// HEAR equals HEART with troponin 0
	for _, age := range []float64{30, 50, 70} {
		a, _ := HEARScore(1, 1, age, 1)
		b, _ := HEARTScore(1, 1, age, 1, 0)
		if a != b {
			t.Errorf("age %v: HEAR %d ≠ HEART(trop 0) %d", age, a, b)
		}
	}
}

// parity: ledger hear_score — no source test file or JSON example; the app sums four caller-chosen 0–2 items
// (age as a 0–2 group), source case 1: 1 + 0 + 1 + 0 = 2.
func TestHEARScore_parity(t *testing.T) {
	if got, _ := HEARScore(1, 0, 50, 0); got != 2 { // parity: ledger hear_score, source case 1
		t.Errorf("got %d", got)
	}
}

func TestHEARScore_errors(t *testing.T) {
	_, err := HEARScore(0, 0, math.NaN(), 0)
	checkErr(t, "NaN", err, ErrNotFinite)
	_, err = HEARScore(3, 0, 50, 0)
	checkErr(t, "item", err, ErrOutOfRange)
	_, err = HEARScore(0, 0, 17, 0)
	checkErr(t, "age", err, ErrOutOfRange)
}

func TestHEARScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = HEARScore(0, 0, 40, 1) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkHEARScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = HEARScore(0, 0, 40, 1)
	}
}
