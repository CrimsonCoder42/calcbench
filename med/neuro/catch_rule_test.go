package neuro

import "testing"

// Osmond 2010 abstract: four high-risk, three medium-risk factors; each alone, HAND-DERIVED.
func TestCATCHRule_example(t *testing.T) {
	for _, f := range []CATCHFindings{{GCSBelow15At2Hours: true}, {SuspectedOpenSkullFx: true}, {WorseningHeadache: true}, {Irritability: true}} {
		if r := CATCHRule(f); !r.HighRisk || r.MediumRisk {
			t.Errorf("CATCHRule(%+v) = %+v, want HighRisk only", f, r)
		}
	}
	for _, f := range []CATCHFindings{{LargeBoggyScalpHematoma: true}, {BasalSkullFractureSign: true}, {DangerousMechanism: true}} {
		if r := CATCHRule(f); r.HighRisk || !r.MediumRisk {
			t.Errorf("CATCHRule(%+v) = %+v, want MediumRisk only", f, r)
		}
	}
	if r := CATCHRule(CATCHFindings{}); r.HighRisk || r.MediumRisk {
		t.Errorf("none: %+v", r)
	}
}

// parity: ledger catch_head_injury, source cases 1–4 (consumer: high+medium, high, medium, none; the
// consumer's signs-of-basal-skull-fracture item is medium-risk as in the source).
func TestCATCHRule_parity(t *testing.T) {
	cases := []struct {
		f         CATCHFindings
		high, med bool
	}{
		{CATCHFindings{Irritability: true, DangerousMechanism: true}, true, true}, // parity: ledger catch_head_injury, source case 1
		{CATCHFindings{WorseningHeadache: true}, true, false},                     // parity: ledger catch_head_injury, source case 2
		{CATCHFindings{BasalSkullFractureSign: true}, false, true},                // parity: ledger catch_head_injury, source case 3
		{CATCHFindings{}, false, false},                                           // parity: ledger catch_head_injury, source case 4
	}
	for i, c := range cases {
		if r := CATCHRule(c.f); r.HighRisk != c.high || r.MediumRisk != c.med {
			t.Errorf("case %d: %+v", i+1, r)
		}
	}
}

func TestCATCHRule_noAlloc(t *testing.T) {
	f := CATCHFindings{Irritability: true}
	if n := testing.AllocsPerRun(100, func() { sinkBool = CATCHRule(f).HighRisk }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkCATCHRule(b *testing.B) {
	b.ReportAllocs()
	f := CATCHFindings{Irritability: true}
	for b.Loop() {
		sinkBool = CATCHRule(f).HighRisk
	}
}
