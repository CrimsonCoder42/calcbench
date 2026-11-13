package renal

import "testing"

var sinkSTONEGroup STONEGroup

// Moore 2014 Results: groups 0–5, 6–9, 10–13; both sides of each edge.
func TestSTONEProbabilityGroup_edges(t *testing.T) {
	tests := []struct {
		s    int
		want STONEGroup
	}{
		{0, STONELowProbability}, {5, STONELowProbability}, {6, STONEModerateProbability},
		{9, STONEModerateProbability}, {10, STONEHighProbability}, {13, STONEHighProbability},
	}
	for _, tc := range tests {
		if got, err := STONEProbabilityGroup(tc.s); err != nil || got != tc.want {
			t.Errorf("%d → (%q, %v), want %q", tc.s, got, err, tc.want)
		}
	}
	runErrCases(t, []errCase{
		{"-1", func() error { _, e := STONEProbabilityGroup(-1); return e }, ErrOutOfRange},
		{"14", func() error { _, e := STONEProbabilityGroup(14); return e }, ErrOutOfRange},
	})
}

// parity: ledger stone_score, source case 2 (app bands are its own "~%" likelihoods; the primary's group labels are
// asserted: 12 → "high probability").
func TestSTONEProbabilityGroup_parity(t *testing.T) {
	if got, _ := STONEProbabilityGroup(12); got != "high probability" {
		t.Errorf("got %q", got)
	}
}

func TestSTONEProbabilityGroup_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkSTONEGroup, _ = STONEProbabilityGroup(7) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkSTONEProbabilityGroup(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkSTONEGroup, _ = STONEProbabilityGroup(7)
	}
}
