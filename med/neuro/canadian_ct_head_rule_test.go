package neuro

import "testing"

var sinkBool bool

// Stiell 2001 abstract: five high-risk factors, two medium-risk factors; each factor alone, HAND-DERIVED.
func TestCanadianCTHeadRule_example(t *testing.T) {
	high := []CanadianCTHeadFindings{
		{GCSBelow15At2Hours: true}, {SuspectedOpenSkullFx: true}, {BasalSkullFractureSign: true},
		{VomitingMoreThan2: true}, {AgeOver65: true},
	}
	for _, f := range high {
		if r := CanadianCTHeadRule(f); !r.HighRisk || r.MediumRisk {
			t.Errorf("CanadianCTHeadRule(%+v) = %+v, want HighRisk only", f, r)
		}
	}
	for _, f := range []CanadianCTHeadFindings{{AmnesiaBeforeImpactOver30: true}, {DangerousMechanism: true}} {
		if r := CanadianCTHeadRule(f); r.HighRisk || !r.MediumRisk {
			t.Errorf("CanadianCTHeadRule(%+v) = %+v, want MediumRisk only", f, r)
		}
	}
	if r := CanadianCTHeadRule(CanadianCTHeadFindings{}); r.HighRisk || r.MediumRisk {
		t.Errorf("no factor: %+v, want neither", r)
	}
	if r := CanadianCTHeadRule(CanadianCTHeadFindings{AgeOver65: true, DangerousMechanism: true}); !r.HighRisk || !r.MediumRisk {
		t.Errorf("both groups: %+v, want both", r)
	}
}

// parity: ledger canadian_ct_head, source case 1–3 (consumer switch: any high-risk; else any medium-risk; else none).
func TestCanadianCTHeadRule_parity(t *testing.T) {
	r := CanadianCTHeadRule(CanadianCTHeadFindings{VomitingMoreThan2: true, AmnesiaBeforeImpactOver30: true})
	if !r.HighRisk { // parity: ledger canadian_ct_head, source case 1 (high-risk branch)
		t.Errorf("case 1: %+v", r)
	}
	r = CanadianCTHeadRule(CanadianCTHeadFindings{DangerousMechanism: true})
	if r.HighRisk || !r.MediumRisk { // parity: ledger canadian_ct_head, source case 2 (medium-risk branch)
		t.Errorf("case 2: %+v", r)
	}
	r = CanadianCTHeadRule(CanadianCTHeadFindings{})
	if r.HighRisk || r.MediumRisk { // parity: ledger canadian_ct_head, source case 3 (default branch)
		t.Errorf("case 3: %+v", r)
	}
}

func TestCanadianCTHeadRule_noAlloc(t *testing.T) {
	f := CanadianCTHeadFindings{AgeOver65: true}
	if n := testing.AllocsPerRun(100, func() { sinkBool = CanadianCTHeadRule(f).HighRisk }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkCanadianCTHeadRule(b *testing.B) {
	b.ReportAllocs()
	f := CanadianCTHeadFindings{AgeOver65: true}
	for b.Loop() {
		sinkBool = CanadianCTHeadRule(f).HighRisk
	}
}
