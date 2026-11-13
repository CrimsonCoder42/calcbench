package neuro

import "testing"

// Perry 2021 Table 1: thirteen weights, total −3 to 23.
func TestCanadianTIAScore_example(t *testing.T) {
	all := CanadianTIAFindings{true, true, true, true, true, true, false, true, true, true, true, true, true}
	if got := CanadianTIAScore(all); got != 23 {
		t.Errorf("all but vertigo = %d, want 23", got)
	}
	if got := CanadianTIAScore(CanadianTIAFindings{HistoryOfVertigo: true}); got != -3 {
		t.Errorf("vertigo = %d, want -3", got)
	}
	if got := CanadianTIAScore(CanadianTIAFindings{FirstTIAInLifetime: true, SymptomsTenMinutesOrMore: true}); got != 4 {
		t.Errorf("first TIA + ≥10 min = %d, want 4", got)
	}
	single := []struct {
		f    CanadianTIAFindings
		want int
	}{
		{CanadianTIAFindings{OnAntiplateletTherapy: true}, 3}, {CanadianTIAFindings{TriageDiastolicAtLeast110: true}, 3},
		{CanadianTIAFindings{GlucoseAtLeast15MmolL: true}, 3}, {CanadianTIAFindings{PlateletsAtLeast400: true}, 2},
		{CanadianTIAFindings{AtrialFibrillationOnECG: true}, 2}, {CanadianTIAFindings{HistoryOfCarotidStenosis: true}, 2},
		{CanadianTIAFindings{InfarctionOnCT: true}, 1}, {CanadianTIAFindings{DysarthriaOrAphasia: true}, 1},
		{CanadianTIAFindings{HistoryOfGaitDisturbance: true}, 1}, {CanadianTIAFindings{HistoryOfUnilateralWeakness: true}, 1},
	}
	for _, c := range single {
		if got := CanadianTIAScore(c.f); got != c.want {
			t.Errorf("CanadianTIAScore(%+v) = %d, want %d", c.f, got, c.want)
		}
	}
}

// parity: ledger canadian_tia_score. The consumer computes a different item set (unilateral weakness 3, prior TIA,
// carotid bruit, "non-cerebrovascular" −3, …) that matches no table of the source; no consumer case is shared.
// The one comparable item, atrial fibrillation on ECG (+2), agrees.
func TestCanadianTIAScore_parity(t *testing.T) {
	if got := CanadianTIAScore(CanadianTIAFindings{AtrialFibrillationOnECG: true}); got != 2 { // parity: ledger canadian_tia_score, source case 1
		t.Errorf("case 1 = %d, want 2", got)
	}
}

func TestCanadianTIAScore_noAlloc(t *testing.T) {
	f := CanadianTIAFindings{FirstTIAInLifetime: true}
	if n := testing.AllocsPerRun(100, func() { sinkInt = CanadianTIAScore(f) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkCanadianTIAScore(b *testing.B) {
	b.ReportAllocs()
	f := CanadianTIAFindings{FirstTIAInLifetime: true}
	for b.Loop() {
		sinkInt = CanadianTIAScore(f)
	}
}
