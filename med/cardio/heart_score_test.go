package cardio

import (
	"math"
	"testing"
)

// Mahler 2015 Table 2: age < 45 → 0, 45–65 → 1, > 65 → 2; items 0–2 summed.
func TestHEARTScore_age(t *testing.T) {
	tests := []struct {
		age  float64
		want int
	}{{44.9, 0}, {45, 1}, {65, 1}, {65.1, 2}, {18, 0}, {110, 2}}
	for _, tc := range tests {
		if got, err := HEARTScore(0, 0, tc.age, 0, 0); err != nil || got != tc.want {
			t.Errorf("age %v = (%d, %v), want %d", tc.age, got, err, tc.want)
		}
	}
	if got, _ := HEARTScore(1, 0, 50, 1, 0); got != 3 { // doc Example
		t.Errorf("example = %d, want 3", got)
	}
	if got, _ := HEARTScore(2, 2, 70, 2, 2); got != 10 {
		t.Errorf("max = %d, want 10", got)
	}
}

// parity: ledger heart_pathway — no source test file or JSON example. The app scores age 65 as 2 (≥ 65); Mahler
// 2015 prints "> 65", so 65 scores 1 here (source case 1 asserts the PRIMARY; the app would give 2).
func TestHEARTScore_parity(t *testing.T) {
	if got, _ := HEARTScore(0, 0, 65, 0, 0); got != 1 { // parity: ledger heart_pathway, source case 1 (app: 2)
		t.Errorf("age 65 = %d, want 1 (primary)", got)
	}
}

func TestHEARTScore_errors(t *testing.T) {
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		_, err := HEARTScore(0, 0, v, 0, 0)
		checkErr(t, "age", err, ErrNotFinite)
	}
	for _, in := range [][4]int{{-1, 0, 0, 0}, {3, 0, 0, 0}, {0, -1, 0, 0}, {0, 3, 0, 0}, {0, 0, -1, 0}, {0, 0, 3, 0}, {0, 0, 0, -1}, {0, 0, 0, 3}} {
		_, err := HEARTScore(in[0], in[1], 50, in[2], in[3])
		checkErr(t, "item", err, ErrOutOfRange)
	}
	for _, a := range []float64{17.9, 110.1} {
		_, err := HEARTScore(0, 0, a, 0, 0)
		checkErr(t, "age range", err, ErrOutOfRange)
	}
}

func TestHEARTScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = HEARTScore(1, 0, 50, 1, 0) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkHEARTScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = HEARTScore(1, 0, 50, 1, 0)
	}
}
