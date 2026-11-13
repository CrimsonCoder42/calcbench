package respiratory

import "testing"

// HAND-DERIVED from Webb 2016 Table 3 (PMID 26856838): the maximum 4 × 2 + 6 × 1 = 14 ("Total possible
// points: 14") and the three ways to reach the ≥ 4 threshold: 2 major, 1 major + 2 minor, 4 minor.
func TestDRIPScore_tableExample(t *testing.T) {
	cases := []struct {
		name string
		f    DRIPFindings
		want int
	}{
		{"none", DRIPFindings{}, 0},
		{"all", DRIPFindings{true, true, true, true, true, true, true, true, true, true}, 14},
		{"2 major", DRIPFindings{AntibioticsWithin60Days: true, TubeFeeding: true}, 4},
		{"1 major + 2 minor", DRIPFindings{LongTermCareResidence: true, WoundCare: true, MRSAColonizationWithin1: true}, 4},
		{"4 minor", DRIPFindings{HospitalizationWithin60Days: true, ChronicPulmonaryDisease: true, PoorFunctionalStatus: true, GastricAcidSuppression: true}, 4},
		{"1 major + 1 minor", DRIPFindings{PriorDRPInfectionWithin1: true, ChronicPulmonaryDisease: true}, 3},
	}
	for _, c := range cases {
		if got := DRIPScore(c.f); got != c.want {
			t.Errorf("%s: got %d, want %d", c.name, got, c.want)
		}
	}
}

// parity: ledger drip_pneumonia, source cases 1–2 (its worked text): antibiotics + chronic pulmonary disease +
// PPI = 4; chronic pulmonary disease + recent hospitalization = 2. The source's weights agree with Table 3.
func TestDRIPScore_parity(t *testing.T) {
	if got := DRIPScore(DRIPFindings{AntibioticsWithin60Days: true, ChronicPulmonaryDisease: true, GastricAcidSuppression: true}); got != 4 {
		t.Errorf("case 1: got %d, want 4", got)
	}
	if got := DRIPScore(DRIPFindings{ChronicPulmonaryDisease: true, HospitalizationWithin60Days: true}); got != 2 {
		t.Errorf("case 2: got %d, want 2", got)
	}
}

func TestDRIPScore_noAlloc(t *testing.T) {
	f := DRIPFindings{AntibioticsWithin60Days: true, WoundCare: true}
	if n := testing.AllocsPerRun(100, func() { sinkInt = DRIPScore(f) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkDRIPScore(b *testing.B) {
	f := DRIPFindings{AntibioticsWithin60Days: true, WoundCare: true}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt = DRIPScore(f)
	}
}
