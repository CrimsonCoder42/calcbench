package cardio

import (
	"math"
	"testing"
)

// Tisdale 2013 table: maximum 21 with both drug items; HAND-DERIVED cases.
func TestTisdaleQTRiskScore_items(t *testing.T) {
	all := TisdaleQTFactors{true, true, true, true, true, true}
	tests := []struct {
		name   string
		age    float64
		female bool
		drugs  int
		f      TisdaleQTFactors
		want   int
	}{
		{"none, 67.9", 67.9, false, 0, TisdaleQTFactors{}, 0},
		{"age edge 68", 68, false, 0, TisdaleQTFactors{}, 1},
		{"one drug", 50, false, 1, TisdaleQTFactors{}, 3},
		{"example", 70, true, 2, TisdaleQTFactors{HeartFailure: true}, 11},
		{"maximum 21", 80, true, 5, all, 21},
	}
	for _, tc := range tests {
		got, err := TisdaleQTRiskScore(tc.age, tc.female, tc.drugs, tc.f)
		if err != nil || got != tc.want {
			t.Errorf("%s: got (%d, %v), want %d", tc.name, got, err, tc.want)
		}
	}
	_, err := TisdaleQTRiskScore(math.NaN(), false, 0, TisdaleQTFactors{})
	checkErr(t, "NaN", err, ErrNotFinite)
	_, err = TisdaleQTRiskScore(17, false, 0, TisdaleQTFactors{})
	checkErr(t, "age 17", err, ErrOutOfRange)
	_, err = TisdaleQTRiskScore(50, false, -1, TisdaleQTFactors{})
	checkErr(t, "drugs -1", err, ErrOutOfRange)
}

// Tisdale 2013: low < 7, moderate 7–10, high ≥ 11; both sides of each edge.
func TestTisdaleQTRiskGroup_edges(t *testing.T) {
	tests := []struct {
		s    int
		want TisdaleQTRiskCategory
	}{{0, TisdaleQTLow}, {6, TisdaleQTLow}, {7, TisdaleQTModerate}, {10, TisdaleQTModerate}, {11, TisdaleQTHigh}, {21, TisdaleQTHigh}}
	for _, tc := range tests {
		if got, err := TisdaleQTRiskGroup(tc.s); err != nil || got != tc.want {
			t.Errorf("%d: (%q, %v), want %q", tc.s, got, err, tc.want)
		}
	}
	for _, s := range []int{-1, 22} {
		_, err := TisdaleQTRiskGroup(s)
		checkErr(t, "range", err, ErrOutOfRange)
	}
}

func TestTisdaleQTRiskScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		sinkInt, _ = TisdaleQTRiskScore(70, true, 1, TisdaleQTFactors{})
		g, _ := TisdaleQTRiskGroup(9)
		sinkStr = string(g)
	}); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkTisdaleQTRiskScore(b *testing.B) {
	b.ReportAllocs()
	f := TisdaleQTFactors{HeartFailure: true}
	for b.Loop() {
		sinkInt, _ = TisdaleQTRiskScore(70, true, 2, f)
	}
}

func BenchmarkTisdaleQTRiskGroup(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		g, _ := TisdaleQTRiskGroup(9)
		sinkStr = string(g)
	}
}
