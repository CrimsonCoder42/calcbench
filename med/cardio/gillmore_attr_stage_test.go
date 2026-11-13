package cardio

import (
	"math"
	"testing"
)

// Gillmore 2018 abstract: stage I NT-proBNP ≤ 3000 and eGFR ≥ 45; stage III NT-proBNP > 3000 and eGFR < 45.
func TestGillmoreATTRStage_edges(t *testing.T) {
	tests := []struct {
		bnp, egfr float64
		want      GillmoreStage
	}{
		{3000, 45, 1},
		{3000.01, 45, 2},
		{3000, 44.99, 2},
		{3000.01, 44.99, 3},
		{500, 90, 1},
		{9000, 20, 3},
	}
	for _, tc := range tests {
		if got, err := GillmoreATTRStage(tc.bnp, tc.egfr); err != nil || got != tc.want {
			t.Errorf("(%v, %v) = (%d, %v), want %d", tc.bnp, tc.egfr, got, err, tc.want)
		}
	}
	_, err := GillmoreATTRStage(math.NaN(), 50)
	checkErr(t, "NaN", err, ErrNotFinite)
	_, err = GillmoreATTRStage(100, math.Inf(1))
	checkErr(t, "+Inf", err, ErrNotFinite)
	_, err = GillmoreATTRStage(-1, 50)
	checkErr(t, "negative", err, ErrOutOfRange)
}

var sinkStage GillmoreStage

func TestGillmoreATTRStage_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkStage, _ = GillmoreATTRStage(3000, 45) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkGillmoreATTRStage(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkStage, _ = GillmoreATTRStage(3000, 45)
	}
}
