package trauma

import "testing"

// Hoffman 2000 abstract: low probability only when all five criteria are met.
func TestNEXUSCSpineLowProbability_source(t *testing.T) {
	if !NEXUSCSpineLowProbability(NEXUSCSpineFindings{}) {
		t.Error("no finding → want true")
	}
	for i, f := range []NEXUSCSpineFindings{
		{MidlineCervicalTenderness: true}, {FocalNeurologicDeficit: true}, {AbnormalAlertness: true},
		{Intoxication: true}, {PainfulDistractingInjury: true},
	} {
		if NEXUSCSpineLowProbability(f) {
			t.Errorf("finding %d alone → want false", i)
		}
	}
}

// parity: ledger nexus_cspine, source case 1 (app canClear = criteriaCount === 0).
func TestNEXUSCSpineLowProbability_parity(t *testing.T) {
	if NEXUSCSpineLowProbability(NEXUSCSpineFindings{Intoxication: true, PainfulDistractingInjury: true}) {
		t.Error("two findings → want false")
	}
}

func TestNEXUSCSpineLowProbability_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkBool = NEXUSCSpineLowProbability(NEXUSCSpineFindings{}) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkNEXUSCSpineLowProbability(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkBool = NEXUSCSpineLowProbability(NEXUSCSpineFindings{Intoxication: true})
	}
}
