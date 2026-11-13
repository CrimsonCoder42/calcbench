package cardio

import (
	"math"
	"testing"
)

// Subherwal 2009 Table 4, one predictor varied at a time from the 1-point baseline (Hct 45, CrCl 130, HR 60,
// SBP 150 → 0 + 0 + 0 + 1).
func TestCRUSADEBleedingScore_table4(t *testing.T) {
	type in struct{ hct, crcl, hr, sbp float64 }
	tests := []struct {
		in   in
		want int
	}{
		{in{45, 130, 60, 150}, 1},
		{in{30.9, 130, 60, 150}, 10}, {in{31, 130, 60, 150}, 8}, {in{33.9, 130, 60, 150}, 8}, {in{34, 130, 60, 150}, 4},
		{in{36.9, 130, 60, 150}, 4}, {in{37, 130, 60, 150}, 3}, {in{39.9, 130, 60, 150}, 3}, {in{40, 130, 60, 150}, 1},
		{in{45, 15, 60, 150}, 40}, {in{45, 15.1, 60, 150}, 36}, {in{45, 30, 60, 150}, 36}, {in{45, 30.1, 60, 150}, 29},
		{in{45, 60, 60, 150}, 29}, {in{45, 60.1, 60, 150}, 18}, {in{45, 90, 60, 150}, 18}, {in{45, 90.1, 60, 150}, 8},
		{in{45, 120, 60, 150}, 8}, {in{45, 120.1, 60, 150}, 1},
		{in{45, 130, 70, 150}, 1}, {in{45, 130, 71, 150}, 2}, {in{45, 130, 80, 150}, 2}, {in{45, 130, 81, 150}, 4},
		{in{45, 130, 90, 150}, 4}, {in{45, 130, 91, 150}, 7}, {in{45, 130, 100, 150}, 7}, {in{45, 130, 101, 150}, 9},
		{in{45, 130, 110, 150}, 9}, {in{45, 130, 111, 150}, 11}, {in{45, 130, 120, 150}, 11}, {in{45, 130, 121, 150}, 12},
		{in{45, 130, 60, 90}, 10}, {in{45, 130, 60, 91}, 8}, {in{45, 130, 60, 100}, 8}, {in{45, 130, 60, 101}, 5},
		{in{45, 130, 60, 120}, 5}, {in{45, 130, 60, 121}, 1}, {in{45, 130, 60, 180}, 1}, {in{45, 130, 60, 181}, 3},
		{in{45, 130, 60, 200}, 3}, {in{45, 130, 60, 201}, 5},
	}
	for _, tc := range tests {
		got, err := CRUSADEBleedingScore(tc.in.hct, tc.in.crcl, tc.in.hr, tc.in.sbp, CRUSADEFactors{})
		if err != nil || got != tc.want {
			t.Errorf("%+v = (%d, %v), want %d", tc.in, got, err, tc.want)
		}
	}
	// yes/no items: 8, 7, 6, 6
	for _, c := range []struct {
		f    CRUSADEFactors
		want int
	}{
		{CRUSADEFactors{Female: true}, 9}, {CRUSADEFactors{SignsOfCHF: true}, 8},
		{CRUSADEFactors{PriorVascularDisease: true}, 7}, {CRUSADEFactors{DiabetesMellitus: true}, 7},
	} {
		if got, _ := CRUSADEBleedingScore(45, 130, 60, 150, c.f); got != c.want {
			t.Errorf("%+v = %d, want %d", c.f, got, c.want)
		}
	}
	// doc Example and the maximum 96
	if got, _ := CRUSADEBleedingScore(32, 45, 95, 110, CRUSADEFactors{Female: true, DiabetesMellitus: true}); got != 60 {
		t.Errorf("example = %d, want 60", got)
	}
	if got, _ := CRUSADEBleedingScore(25, 10, 130, 80, CRUSADEFactors{true, true, true, true}); got != 96 {
		t.Errorf("max = %d, want 96", got)
	}
}

// Subherwal 2009 Results: the five groups, both sides of each edge.
func TestCRUSADEBleedingGroupOf(t *testing.T) {
	tests := []struct {
		s    int
		want CRUSADEBleedingGroup
	}{
		{1, CRUSADEVeryLow}, {20, CRUSADEVeryLow}, {21, CRUSADELow}, {30, CRUSADELow}, {31, CRUSADEModerate},
		{40, CRUSADEModerate}, {41, CRUSADEHigh}, {50, CRUSADEHigh}, {51, CRUSADEVeryHigh}, {96, CRUSADEVeryHigh},
	}
	for _, tc := range tests {
		if got, err := CRUSADEBleedingGroupOf(tc.s); err != nil || got != tc.want {
			t.Errorf("%d = (%q, %v), want %q", tc.s, got, err, tc.want)
		}
	}
	for _, s := range []int{0, 97} {
		_, err := CRUSADEBleedingGroupOf(s)
		checkErr(t, "group range", err, ErrOutOfRange)
	}
}

// parity: ledger crusade_score — no source test file or JSON example; the app's bands agree with Table 4 on
// integer inputs (source case 1 is the app's ≤ 20 / > 50 group edges).
func TestCRUSADEBleedingScore_parity(t *testing.T) {
	g, _ := CRUSADEBleedingGroupOf(20) // parity: ledger crusade_score, source case 1
	if g != CRUSADEVeryLow {
		t.Errorf("20 → %q", g)
	}
}

func TestCRUSADEBleedingScore_errors(t *testing.T) {
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		_, err := CRUSADEBleedingScore(v, 130, 60, 150, CRUSADEFactors{})
		checkErr(t, "non-finite", err, ErrNotFinite)
	}
	for _, in := range [][4]float64{{9.9, 130, 60, 150}, {70.1, 130, 60, 150}, {45, -0.1, 60, 150}, {45, 300.1, 60, 150}, {45, 130, 19.9, 150}, {45, 130, 250.1, 150}, {45, 130, 60, 39.9}, {45, 130, 60, 300.1}} {
		_, err := CRUSADEBleedingScore(in[0], in[1], in[2], in[3], CRUSADEFactors{})
		checkErr(t, "range", err, ErrOutOfRange)
	}
}

func TestCRUSADEBleedingScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		sinkInt, _ = CRUSADEBleedingScore(32, 45, 95, 110, CRUSADEFactors{Female: true})
		g, _ := CRUSADEBleedingGroupOf(60)
		sinkStr = string(g)
	}); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkCRUSADEBleedingScore(b *testing.B) {
	b.ReportAllocs()
	f := CRUSADEFactors{Female: true}
	for b.Loop() {
		sinkInt, _ = CRUSADEBleedingScore(32, 45, 95, 110, f)
	}
}

func BenchmarkCRUSADEBleedingGroupOf(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		g, _ := CRUSADEBleedingGroupOf(60)
		sinkStr = string(g)
	}
}
