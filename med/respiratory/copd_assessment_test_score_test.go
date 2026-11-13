package respiratory

import (
	"errors"
	"testing"
)

// HAND-DERIVED: the sum of eight 0–5 items (Jones 2011, PMID 21835018; Houben-Wilke 2018, PMID 30376861);
// the extremes 0 and 40.
func TestCOPDAssessmentTestScore_example(t *testing.T) {
	cases := []struct {
		items [8]int
		want  int
	}{
		{[8]int{2, 3, 1, 4, 2, 1, 3, 2}, 18},
		{[8]int{}, 0},
		{[8]int{5, 5, 5, 5, 5, 5, 5, 5}, 40},
	}
	for _, c := range cases {
		if got, err := COPDAssessmentTestScore(c.items); err != nil || got != c.want {
			t.Errorf("%v: got (%d, %v), want %d", c.items, got, err, c.want)
		}
	}
}

// parity: ledger copd_cat, source case 1: the source rounds and sums the eight answers; 1,1,1,1,1,1,1,1 → 8.
func TestCOPDAssessmentTestScore_parity(t *testing.T) {
	if got, err := COPDAssessmentTestScore([8]int{1, 1, 1, 1, 1, 1, 1, 1}); err != nil || got != 8 {
		t.Errorf("got (%d, %v), want 8", got, err)
	}
}

func TestCOPDAssessmentTestScore_errors(t *testing.T) {
	for _, items := range [][8]int{{-1}, {6}, {0, 0, 0, 0, 0, 0, 0, 6}, {0, 0, 0, -1}} {
		if _, err := COPDAssessmentTestScore(items); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("%v: err = %v, want ErrOutOfRange", items, err)
		}
	}
}

func TestCOPDAssessmentTestScore_noAlloc(t *testing.T) {
	items := [8]int{2, 3, 1, 4, 2, 1, 3, 2}
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = COPDAssessmentTestScore(items) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkCOPDAssessmentTestScore(b *testing.B) {
	items := [8]int{2, 3, 1, 4, 2, 1, 3, 2}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt, _ = COPDAssessmentTestScore(items)
	}
}
