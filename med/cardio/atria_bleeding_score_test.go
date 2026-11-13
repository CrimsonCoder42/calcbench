package cardio

import (
	"math"
	"testing"
)

// Fang 2011 abstract points; HAND-DERIVED totals.
func TestATRIABleedingScore_items(t *testing.T) {
	tests := []struct {
		age  float64
		f    ATRIABleedingFactors
		want int
	}{
		{74.9, ATRIABleedingFactors{}, 0},
		{75, ATRIABleedingFactors{}, 2},
		{80, ATRIABleedingFactors{Anemia: true}, 5},
		{90, ATRIABleedingFactors{true, true, true, true}, 10},
	}
	for _, tc := range tests {
		if got, err := ATRIABleedingScore(tc.age, tc.f); err != nil || got != tc.want {
			t.Errorf("(%v, %+v) = (%d, %v), want %d", tc.age, tc.f, got, err, tc.want)
		}
	}
	_, err := ATRIABleedingScore(math.NaN(), ATRIABleedingFactors{})
	checkErr(t, "NaN", err, ErrNotFinite)
	_, err = ATRIABleedingScore(10, ATRIABleedingFactors{})
	checkErr(t, "10", err, ErrOutOfRange)
}

// Fang 2011 abstract: low 0–3, intermediate 4, high 5–10.
func TestATRIABleedingGroup_edges(t *testing.T) {
	tests := []struct {
		s    int
		want ATRIABleedingCategory
	}{{0, ATRIABleedingLow}, {3, ATRIABleedingLow}, {4, ATRIABleedingIntermediate}, {5, ATRIABleedingHigh}, {10, ATRIABleedingHigh}}
	for _, tc := range tests {
		if got, err := ATRIABleedingGroup(tc.s); err != nil || got != tc.want {
			t.Errorf("%d: (%q, %v), want %q", tc.s, got, err, tc.want)
		}
	}
	for _, s := range []int{-1, 11} {
		_, err := ATRIABleedingGroup(s)
		checkErr(t, "range", err, ErrOutOfRange)
	}
}

func TestATRIABleedingScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		sinkInt, _ = ATRIABleedingScore(80, ATRIABleedingFactors{Anemia: true})
		g, _ := ATRIABleedingGroup(4)
		sinkStr = string(g)
	}); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkATRIABleedingScore(b *testing.B) {
	b.ReportAllocs()
	f := ATRIABleedingFactors{Anemia: true}
	for b.Loop() {
		sinkInt, _ = ATRIABleedingScore(80, f)
	}
}

func BenchmarkATRIABleedingGroup(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		g, _ := ATRIABleedingGroup(4)
		sinkStr = string(g)
	}
}
