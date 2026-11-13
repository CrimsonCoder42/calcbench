package critical

import "testing"

// Klastersky 2000 abstract weights; the ≥ 21 threshold on both sides. HAND-DERIVED.
func TestMASCCRiskIndex_source(t *testing.T) {
	all := MASCCFactors{true, true, true, true, true}
	tests := []struct {
		b    MASCCBurden
		f    MASCCFactors
		age  float64
		want int
		low  bool
	}{
		{MASCCBurdenNoOrMildSymptoms, all, 45, 26, true},
		{MASCCBurdenNeither, MASCCFactors{}, 70, 0, false},
		{MASCCBurdenModerateSymptoms, MASCCFactors{NoHypotension: true, NoCOPD: true, SolidTumorOrNoPriorFungal: true, NoDehydration: true}, 70, 19, false},
		{MASCCBurdenModerateSymptoms, all, 70, 22, true},
		{MASCCBurdenNoOrMildSymptoms, MASCCFactors{NoHypotension: true, NoCOPD: true, SolidTumorOrNoPriorFungal: true, NoDehydration: true}, 60, 21, true},
		{MASCCBurdenNoOrMildSymptoms, MASCCFactors{NoHypotension: true, NoCOPD: true, SolidTumorOrNoPriorFungal: true, Outpatient: true}, 60, 21, true},
		{MASCCBurdenNoOrMildSymptoms, MASCCFactors{NoHypotension: true, NoCOPD: true, SolidTumorOrNoPriorFungal: true}, 59.9, 20, false},
	}
	for i, tc := range tests {
		r, err := MASCCRiskIndex(tc.b, tc.f, tc.age)
		if err != nil || r.Score != tc.want || r.LowRisk != tc.low {
			t.Errorf("case %d: (%+v, %v), want %d %v", i, r, err, tc.want, tc.low)
		}
	}
}

// parity: ledger mascc_febrile_neutropenia, source case 1 (app: burden index 2 = mild/none 5, every 1 = favourable,
// age_under_60 = 1 → 26).
func TestMASCCRiskIndex_parity(t *testing.T) {
	r, _ := MASCCRiskIndex(MASCCBurdenNoOrMildSymptoms, MASCCFactors{true, true, true, true, true}, 50)
	if r.Score != 26 {
		t.Errorf("got %d, want 26", r.Score)
	}
}

func TestMASCCRiskIndex_errors(t *testing.T) {
	for _, v := range nonFiniteB6 {
		_, err := MASCCRiskIndex(MASCCBurdenNeither, MASCCFactors{}, v)
		wantErrB6(t, "non-finite", err, ErrNotFinite)
	}
	_, err := MASCCRiskIndex(MASCCBurdenNeither, MASCCFactors{}, -1)
	wantErrB6(t, "age", err, ErrOutOfRange)
	for _, b := range []MASCCBurden{-1, 3} {
		_, err := MASCCRiskIndex(b, MASCCFactors{}, 40)
		wantErrB6(t, "burden", err, ErrOutOfRange)
	}
}

func TestMASCCRiskIndex_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		r, _ := MASCCRiskIndex(MASCCBurdenModerateSymptoms, MASCCFactors{NoCOPD: true}, 50)
		sinkB6Int = r.Score
	}); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkMASCCRiskIndex(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		r, _ := MASCCRiskIndex(MASCCBurdenModerateSymptoms, MASCCFactors{NoCOPD: true}, 50)
		sinkB6Int = r.Score
	}
}
