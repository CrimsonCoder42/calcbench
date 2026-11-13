package trauma

import "testing"

// Holmes 2013 abstract: every item must be met; "Glasgow Coma Scale score greater than 13".
func TestPECARNIntraAbdominalInjuryVeryLowRisk_source(t *testing.T) {
	if ok, err := PECARNIntraAbdominalInjuryVeryLowRisk(14, PECARNIAIFindings{}); err != nil || !ok {
		t.Errorf("GCS 14, none → (%v, %v), want true", ok, err)
	}
	if ok, _ := PECARNIntraAbdominalInjuryVeryLowRisk(13, PECARNIAIFindings{}); ok {
		t.Error("GCS 13 → want false")
	}
	for i, f := range []PECARNIAIFindings{
		{AbdominalWallTraumaOrSeatBeltSign: true}, {AbdominalTenderness: true}, {ThoracicWallTrauma: true},
		{AbdominalPain: true}, {DecreasedBreathSounds: true}, {Vomiting: true},
	} {
		if ok, _ := PECARNIntraAbdominalInjuryVeryLowRisk(15, f); ok {
			t.Errorf("finding %d → want false", i)
		}
	}
}

// parity: ledger pecarn_iai, source case 1 (app isVeryLowRisk = positiveCount === 0; gcs_low = 1 is GCS 12 here).
func TestPECARNIntraAbdominalInjuryVeryLowRisk_parity(t *testing.T) {
	if ok, _ := PECARNIntraAbdominalInjuryVeryLowRisk(12, PECARNIAIFindings{}); ok {
		t.Error("gcs low → want false")
	}
}

func TestPECARNIntraAbdominalInjuryVeryLowRisk_errors(t *testing.T) {
	for _, g := range []int{2, 16} {
		_, err := PECARNIntraAbdominalInjuryVeryLowRisk(g, PECARNIAIFindings{})
		checkErr(t, "gcs", err, ErrOutOfRange)
	}
}

func TestPECARNIntraAbdominalInjuryVeryLowRisk_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkBool, _ = PECARNIntraAbdominalInjuryVeryLowRisk(15, PECARNIAIFindings{}) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkPECARNIntraAbdominalInjuryVeryLowRisk(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkBool, _ = PECARNIntraAbdominalInjuryVeryLowRisk(15, PECARNIAIFindings{})
	}
}
