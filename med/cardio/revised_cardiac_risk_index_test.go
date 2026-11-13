package cardio

import "testing"

// Lee 1999 abstract: one point per factor, 0–6.
func TestRevisedCardiacRiskIndex_items(t *testing.T) {
	if got := RevisedCardiacRiskIndex(RCRIFactors{}); got != 0 {
		t.Errorf("none = %d, want 0", got)
	}
	if got := RevisedCardiacRiskIndex(RCRIFactors{InsulinTreatment: true, CreatinineAbove2MgDL: true}); got != 2 {
		t.Errorf("insulin + creatinine = %d, want 2", got)
	}
	if got := RevisedCardiacRiskIndex(RCRIFactors{true, true, true, true, true, true}); got != 6 {
		t.Errorf("all = %d, want 6", got)
	}
}

func TestRevisedCardiacRiskIndex_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt = RevisedCardiacRiskIndex(RCRIFactors{HighRiskSurgery: true}) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkRevisedCardiacRiskIndex(b *testing.B) {
	b.ReportAllocs()
	f := RCRIFactors{HighRiskSurgery: true}
	for b.Loop() {
		sinkInt = RevisedCardiacRiskIndex(f)
	}
}
