package cardio

import "testing"

// Nguyen 2018 Table 2: each item's points, from 0 decades.
func TestREADMITSScore_table2(t *testing.T) {
	tests := []struct {
		dec  int
		f    READMITSFactors
		want int
	}{
		{0, READMITSFactors{}, 0},
		{1, READMITSFactors{}, 1},
		{9, READMITSFactors{}, 9},
		{0, READMITSFactors{CreatinineAbove2MgDL: true}, 6},
		{0, READMITSFactors{ElevatedBNP: true}, 8},
		{0, READMITSFactors{DiabetesMellitus: true}, 4},
		{0, READMITSFactors{Female: true}, 2},
		{0, READMITSFactors{NoTimelyPCI: true}, 1},
		{0, READMITSFactors{SystolicBelow100: true}, 3},
		{4, READMITSFactors{DiabetesMellitus: true, Female: true}, 10}, // doc Example
		{9, READMITSFactors{true, true, true, true, true, true}, 33},   // maximum
	}
	for _, tc := range tests {
		if got, err := READMITSScore(tc.dec, tc.f); err != nil || got != tc.want {
			t.Errorf("%d %+v = (%d, %v), want %d", tc.dec, tc.f, got, err, tc.want)
		}
	}
	for _, d := range []int{-1, 10} {
		_, err := READMITSScore(d, READMITSFactors{})
		checkErr(t, "decades range", err, ErrOutOfRange)
	}
}

// Nguyen 2018 Table 4: both sides of every group edge.
func TestREADMITSRiskGroupOf(t *testing.T) {
	tests := []struct {
		s    int
		want READMITSRiskGroup
	}{
		{0, READMITSExtremelyLow}, {13, READMITSExtremelyLow}, {14, READMITSLow}, {15, READMITSLow},
		{16, READMITSModerate}, {17, READMITSModerate}, {18, READMITSHigh}, {19, READMITSHigh},
		{20, READMITSExtremelyHigh}, {33, READMITSExtremelyHigh},
	}
	for _, tc := range tests {
		if got, err := READMITSRiskGroupOf(tc.s); err != nil || got != tc.want {
			t.Errorf("%d = (%q, %v), want %q", tc.s, got, err, tc.want)
		}
	}
	for _, s := range []int{-1, 34} {
		_, err := READMITSRiskGroupOf(s)
		checkErr(t, "group range", err, ErrOutOfRange)
	}
}

// parity: ledger readmits_ami — no source test file or JSON example. The app computes the age points as
// floor((age − 18) / 10); for age 70 that is 5 decades. Source case 1: age 70, BNP, renal → 5 + 8 + 6 = 19, "High".
func TestREADMITSScore_parity(t *testing.T) {
	s, err := READMITSScore(5, READMITSFactors{CreatinineAbove2MgDL: true, ElevatedBNP: true}) // parity: ledger readmits_ami, source case 1
	g, _ := READMITSRiskGroupOf(s)
	if err != nil || s != 19 || g != READMITSHigh {
		t.Errorf("(%d, %q, %v), want (19, High)", s, g, err)
	}
}

func TestREADMITSScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		sinkInt, _ = READMITSScore(4, READMITSFactors{Female: true})
		g, _ := READMITSRiskGroupOf(18)
		sinkStr = string(g)
	}); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkREADMITSScore(b *testing.B) {
	b.ReportAllocs()
	f := READMITSFactors{Female: true, DiabetesMellitus: true}
	for b.Loop() {
		sinkInt, _ = READMITSScore(4, f)
	}
}

func BenchmarkREADMITSRiskGroupOf(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		g, _ := READMITSRiskGroupOf(18)
		sinkStr = string(g)
	}
}
