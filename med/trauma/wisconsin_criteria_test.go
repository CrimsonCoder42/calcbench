package trauma

import "testing"

// Sitzman 2011 abstract: five criteria, "Glasgow Coma Scale score less than 14". HAND-DERIVED.
func TestWisconsinCriteria_source(t *testing.T) {
	tests := []struct {
		gcs  int
		f    WisconsinFindings
		want int
	}{
		{14, WisconsinFindings{}, 0},
		{13, WisconsinFindings{}, 1},
		{15, WisconsinFindings{Malocclusion: true, ToothAbsence: true}, 2},
		{3, WisconsinFindings{true, true, true, true}, 5},
	}
	for i, tc := range tests {
		if got, err := WisconsinCriteria(tc.gcs, tc.f); err != nil || got != tc.want {
			t.Errorf("case %d: (%d, %v), want %d", i, got, err, tc.want)
		}
	}
}

// parity: ledger wisconsin_criteria, source case 1 (app counts five 0/1 items; gcs_low = 1 is GCS 10 here).
func TestWisconsinCriteria_parity(t *testing.T) {
	if got, _ := WisconsinCriteria(10, WisconsinFindings{PeriorbitalSwellingContusion: true}); got != 2 {
		t.Errorf("got %d, want 2", got)
	}
}

func TestWisconsinCriteria_errors(t *testing.T) {
	for _, g := range []int{2, 16} {
		_, err := WisconsinCriteria(g, WisconsinFindings{})
		checkErr(t, "gcs", err, ErrOutOfRange)
	}
}

func TestWisconsinCriteria_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = WisconsinCriteria(13, WisconsinFindings{}) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkWisconsinCriteria(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = WisconsinCriteria(13, WisconsinFindings{})
	}
}
