package cardio

import (
	"math"
	"testing"
)

// O'Brien 2015, Tables 2 and 3.
func TestORBITBleedingScore_items(t *testing.T) {
	tests := []struct {
		age  float64
		f    ORBITBleedingFactors
		want int
	}{
		{74.9, ORBITBleedingFactors{}, 0},
		{75, ORBITBleedingFactors{}, 1},
		{78, ORBITBleedingFactors{AntiplateletTreatment: true}, 2},
		{80, ORBITBleedingFactors{true, true, true, true}, 7},
	}
	for _, tc := range tests {
		if got, err := ORBITBleedingScore(tc.age, tc.f); err != nil || got != tc.want {
			t.Errorf("(%v, %+v) = (%d, %v), want %d", tc.age, tc.f, got, err, tc.want)
		}
	}
	_, err := ORBITBleedingScore(math.NaN(), ORBITBleedingFactors{})
	checkErr(t, "NaN", err, ErrNotFinite)
	_, err = ORBITBleedingScore(121, ORBITBleedingFactors{})
	checkErr(t, "121", err, ErrOutOfRange)
}

func TestORBITBleedingGroup_edges(t *testing.T) {
	tests := []struct {
		s    int
		want ORBITBleedingCategory
	}{{0, ORBITBleedingLow}, {2, ORBITBleedingLow}, {3, ORBITBleedingMedium}, {4, ORBITBleedingHigh}, {7, ORBITBleedingHigh}}
	for _, tc := range tests {
		if got, err := ORBITBleedingGroup(tc.s); err != nil || got != tc.want {
			t.Errorf("%d: (%q, %v), want %q", tc.s, got, err, tc.want)
		}
	}
	for _, s := range []int{-1, 8} {
		_, err := ORBITBleedingGroup(s)
		checkErr(t, "range", err, ErrOutOfRange)
	}
}

func TestORBITBleedingScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		sinkInt, _ = ORBITBleedingScore(78, ORBITBleedingFactors{})
		g, _ := ORBITBleedingGroup(3)
		sinkStr = string(g)
	}); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkORBITBleedingScore(b *testing.B) {
	b.ReportAllocs()
	f := ORBITBleedingFactors{AntiplateletTreatment: true}
	for b.Loop() {
		sinkInt, _ = ORBITBleedingScore(78, f)
	}
}

func BenchmarkORBITBleedingGroup(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		g, _ := ORBITBleedingGroup(3)
		sinkStr = string(g)
	}
}
