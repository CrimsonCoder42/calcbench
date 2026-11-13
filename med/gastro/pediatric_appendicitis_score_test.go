package gastro

import "testing"

// HAND-DERIVED from Samuel 2002 (abstract): items 1 and 5 score 2, the others 1; maximum 10.
func TestPediatricAppendicitisScore_example(t *testing.T) {
	tests := []struct {
		name string
		f    PediatricAppendicitisFindings
		want int
	}{
		{"none", PediatricAppendicitisFindings{}, 0},
		{"all", PediatricAppendicitisFindings{true, true, true, true, true, true, true, true}, 10},
		{"worked", PediatricAppendicitisFindings{CoughPercussionHoppingTendernessRLQ: true, Anorexia: true, Leukocytosis: true}, 4},
		{"item 1 alone", PediatricAppendicitisFindings{CoughPercussionHoppingTendernessRLQ: true}, 2},
		{"item 5 alone", PediatricAppendicitisFindings{TendernessRightIliacFossa: true}, 2},
		{"items 2,3,4,6,7,8", PediatricAppendicitisFindings{false, true, true, true, false, true, true, true}, 6},
	}
	for _, tc := range tests {
		if got := PediatricAppendicitisScore(tc.f); got != tc.want {
			t.Errorf("%s: got %v, want %v", tc.name, got, tc.want)
		}
	}
}

// parity: ledger pediatric_appendicitis_score, the app's sum rlq_cough(2) + anorexia + fever + nausea_vomiting +
// rif_tenderness(2) + leukocytosis + neutrophilia + migration; e.g. rlq_cough, rif_tenderness, fever → 5.
func TestPediatricAppendicitisScore_parity(t *testing.T) {
	f := PediatricAppendicitisFindings{CoughPercussionHoppingTendernessRLQ: true, TendernessRightIliacFossa: true, Pyrexia: true}
	if got := PediatricAppendicitisScore(f); got != 5 {
		t.Errorf("got %v, want 5", got)
	}
}

func TestPediatricAppendicitisScore_noAllocation(t *testing.T) {
	f := PediatricAppendicitisFindings{true, true, false, true, true, false, true, false}
	if n := testing.AllocsPerRun(100, func() { sinkInt = PediatricAppendicitisScore(f) }); n != 0 {
		t.Errorf("PediatricAppendicitisScore allocates %v times, want 0", n)
	}
}

func BenchmarkPediatricAppendicitisScore(b *testing.B) {
	b.ReportAllocs()
	f := PediatricAppendicitisFindings{true, true, false, true, true, false, true, false}
	for b.Loop() {
		sinkInt = PediatricAppendicitisScore(f)
	}
}
