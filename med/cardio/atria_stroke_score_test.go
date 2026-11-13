package cardio

import (
	"math"
	"testing"
)

// Singer 2013, Table 3: age rows in both columns, edges on both sides; "0 to 12" and "7 to 15" extremes.
func TestATRIAStrokeScore_items(t *testing.T) {
	all := ATRIAStrokeFactors{false, true, true, true, true, true}
	allPrior := all
	allPrior.PriorStroke = true
	tests := []struct {
		name   string
		age    float64
		female bool
		f      ATRIAStrokeFactors
		want   int
	}{
		{"64.9 none", 64.9, false, ATRIAStrokeFactors{}, 0},
		{"65", 65, false, ATRIAStrokeFactors{}, 3},
		{"74.9", 74.9, false, ATRIAStrokeFactors{}, 3},
		{"75", 75, false, ATRIAStrokeFactors{}, 5},
		{"84.9", 84.9, false, ATRIAStrokeFactors{}, 5},
		{"85", 85, false, ATRIAStrokeFactors{}, 6},
		{"prior <65", 50, false, ATRIAStrokeFactors{PriorStroke: true}, 8},
		{"prior 65", 65, false, ATRIAStrokeFactors{PriorStroke: true}, 7},
		{"prior 75", 75, false, ATRIAStrokeFactors{PriorStroke: true}, 7},
		{"prior 85", 85, false, ATRIAStrokeFactors{PriorStroke: true}, 9},
		{"example", 80, true, ATRIAStrokeFactors{Hypertension: true}, 7},
		{"max no prior 12", 90, true, all, 12},
		{"max prior 15", 90, true, allPrior, 15},
	}
	for _, tc := range tests {
		if got, err := ATRIAStrokeScore(tc.age, tc.female, tc.f); err != nil || got != tc.want {
			t.Errorf("%s: (%d, %v), want %d", tc.name, got, err, tc.want)
		}
	}
	_, err := ATRIAStrokeScore(math.Inf(1), false, ATRIAStrokeFactors{})
	checkErr(t, "+Inf", err, ErrNotFinite)
	_, err = ATRIAStrokeScore(17.9, false, ATRIAStrokeFactors{})
	checkErr(t, "17.9", err, ErrOutOfRange)
}

// Singer 2013: low 0–5, moderate 6, high 7–15.
func TestATRIAStrokeGroup_edges(t *testing.T) {
	tests := []struct {
		s    int
		want ATRIAStrokeCategory
	}{{0, ATRIAStrokeLow}, {5, ATRIAStrokeLow}, {6, ATRIAStrokeModerate}, {7, ATRIAStrokeHigh}, {15, ATRIAStrokeHigh}}
	for _, tc := range tests {
		if got, err := ATRIAStrokeGroup(tc.s); err != nil || got != tc.want {
			t.Errorf("%d: (%q, %v), want %q", tc.s, got, err, tc.want)
		}
	}
	for _, s := range []int{-1, 16} {
		_, err := ATRIAStrokeGroup(s)
		checkErr(t, "range", err, ErrOutOfRange)
	}
}

func TestATRIAStrokeScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		sinkInt, _ = ATRIAStrokeScore(80, true, ATRIAStrokeFactors{})
		g, _ := ATRIAStrokeGroup(6)
		sinkStr = string(g)
	}); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkATRIAStrokeScore(b *testing.B) {
	b.ReportAllocs()
	f := ATRIAStrokeFactors{Hypertension: true}
	for b.Loop() {
		sinkInt, _ = ATRIAStrokeScore(80, true, f)
	}
}

func BenchmarkATRIAStrokeGroup(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		g, _ := ATRIAStrokeGroup(6)
		sinkStr = string(g)
	}
}
