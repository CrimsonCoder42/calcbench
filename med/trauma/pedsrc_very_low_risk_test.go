package trauma

import "testing"

// Streck 2017 abstract: five variables, "aspartate aminotransferase >200 U/L". HAND-DERIVED.
func TestPedSRCVeryLowRisk_source(t *testing.T) {
	tests := []struct {
		ast  float64
		f    PedSRCFindings
		n    int
		very bool
	}{
		{200, PedSRCFindings{}, 0, true},
		{200.1, PedSRCFindings{}, 1, false},
		{50, PedSRCFindings{AbnormalAbdominalExam: true, AbdominalPain: true}, 2, false},
		{300, PedSRCFindings{true, true, true, true}, 5, false},
		{0, PedSRCFindings{AbnormalChestXRay: true}, 1, false},
		{0, PedSRCFindings{AbnormalPancreaticEnzymes: true}, 1, false},
	}
	for i, tc := range tests {
		n, very, err := PedSRCVeryLowRisk(tc.ast, tc.f)
		if err != nil || n != tc.n || very != tc.very {
			t.Errorf("case %d: (%d, %v, %v), want (%d, %v)", i, n, very, err, tc.n, tc.very)
		}
	}
}

// parity: ledger pedsrc, source case 1 (app sums five 0/1 items; ast = 1 is AST 250 here).
func TestPedSRCVeryLowRisk_parity(t *testing.T) {
	if n, _, _ := PedSRCVeryLowRisk(250, PedSRCFindings{AbdominalPain: true}); n != 2 {
		t.Errorf("n = %d, want 2", n)
	}
}

func TestPedSRCVeryLowRisk_errors(t *testing.T) {
	for _, v := range nonFinite {
		_, _, err := PedSRCVeryLowRisk(v, PedSRCFindings{})
		checkErr(t, "non-finite", err, ErrNotFinite)
	}
	_, _, err := PedSRCVeryLowRisk(-1, PedSRCFindings{})
	checkErr(t, "range", err, ErrOutOfRange)
}

func TestPedSRCVeryLowRisk_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, sinkBool, _ = PedSRCVeryLowRisk(250, PedSRCFindings{}) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkPedSRCVeryLowRisk(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, sinkBool, _ = PedSRCVeryLowRisk(250, PedSRCFindings{})
	}
}
