package cardio

import "testing"

// Fowler 2023 Table 1: Definite (2 major; 1 major + 3 minor; 5 minor; pathologic), Possible (1 major + 1 minor;
// 3 minor), Rejected otherwise; each edge from both sides.
func TestDukeISCVIDCategoryOf(t *testing.T) {
	tests := []struct {
		f    DukeISCVIDFindings
		want DukeISCVIDCategory
	}{
		{DukeISCVIDFindings{}, DukeISCVIDRejected},
		{DukeISCVIDFindings{PathologicCriterion: true}, DukeISCVIDDefinite},
		{DukeISCVIDFindings{MajorCriteria: 2}, DukeISCVIDDefinite},
		{DukeISCVIDFindings{MajorCriteria: 1, MinorCriteria: 3}, DukeISCVIDDefinite},
		{DukeISCVIDFindings{MajorCriteria: 1, MinorCriteria: 2}, DukeISCVIDPossible},
		{DukeISCVIDFindings{MajorCriteria: 1, MinorCriteria: 1}, DukeISCVIDPossible},
		{DukeISCVIDFindings{MajorCriteria: 1}, DukeISCVIDRejected},
		{DukeISCVIDFindings{MinorCriteria: 5}, DukeISCVIDDefinite},
		{DukeISCVIDFindings{MinorCriteria: 4}, DukeISCVIDPossible},
		{DukeISCVIDFindings{MinorCriteria: 3}, DukeISCVIDPossible},
		{DukeISCVIDFindings{MinorCriteria: 2}, DukeISCVIDRejected},
		{DukeISCVIDFindings{MinorCriteria: 3, RejectionFinding: true}, DukeISCVIDRejected},
		{DukeISCVIDFindings{MajorCriteria: 2, RejectionFinding: true}, DukeISCVIDDefinite},
		{DukeISCVIDFindings{MajorCriteria: 3, MinorCriteria: 7}, DukeISCVIDDefinite},
	}
	for _, tc := range tests {
		if got, err := DukeISCVIDCategoryOf(tc.f); err != nil || got != tc.want {
			t.Errorf("%+v = (%q, %v), want %q", tc.f, got, err, tc.want)
		}
	}
	for _, f := range []DukeISCVIDFindings{{MajorCriteria: -1}, {MajorCriteria: 4}, {MinorCriteria: -1}, {MinorCriteria: 8}} {
		_, err := DukeISCVIDCategoryOf(f)
		checkErr(t, "range", err, ErrOutOfRange)
	}
}

// parity: ledger duke_iscvid — no source test file or JSON example; the app's count rules match Table 1.
func TestDukeISCVIDCategoryOf_parity(t *testing.T) {
	got, _ := DukeISCVIDCategoryOf(DukeISCVIDFindings{MajorCriteria: 1, MinorCriteria: 1}) // parity: ledger duke_iscvid, source case 1
	if got != DukeISCVIDPossible {
		t.Errorf("got %q", got)
	}
}

func TestDukeISCVIDCategoryOf_noAlloc(t *testing.T) {
	f := DukeISCVIDFindings{MajorCriteria: 1, MinorCriteria: 3}
	if n := testing.AllocsPerRun(100, func() {
		c, _ := DukeISCVIDCategoryOf(f)
		sinkStr = string(c)
	}); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkDukeISCVIDCategoryOf(b *testing.B) {
	b.ReportAllocs()
	f := DukeISCVIDFindings{MajorCriteria: 1, MinorCriteria: 3}
	for b.Loop() {
		c, _ := DukeISCVIDCategoryOf(f)
		sinkStr = string(c)
	}
}
