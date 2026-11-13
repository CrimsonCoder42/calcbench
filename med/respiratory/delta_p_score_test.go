package respiratory

import "testing"

// parity: ledger delta_p, source cases 1–5 (source order: bulbar, erectile, weight loss, tobacco, age at onset,
// KPS). The source's band labels (Low/Intermediate/High Risk of SCLC) were not seen in the primary and are not ported.
func TestDeltaPScore_parity(t *testing.T) {
	cases := []struct {
		f    DeltaPFindings
		want int
	}{
		{DeltaPFindings{}, 0},                                                                                // source case 1
		{DeltaPFindings{BulbarInvolvement: true}, 1},                                                         // source case 2
		{DeltaPFindings{BulbarInvolvement: true, ErectileDysfunction: true}, 2},                              // source case 3
		{DeltaPFindings{BulbarInvolvement: true, ErectileDysfunction: true, WeightLossAtLeast5Pct: true}, 3}, // source case 4
		{DeltaPFindings{true, true, true, true, true, true}, 6},                                              // source case 5
	}
	for i, c := range cases {
		if got := DeltaPScore(c.f); got != c.want {
			t.Errorf("case %d: got %d, want %d", i+1, got, c.want)
		}
	}
}

// Each item alone scores exactly one point (Titulaer 2011 abstract: six items, scores 0–6).
func TestDeltaPScore_eachItemOnePoint(t *testing.T) {
	items := []DeltaPFindings{
		{AgeAtOnsetAtLeast50: true}, {SmokingAtDiagnosis: true}, {WeightLossAtLeast5Pct: true},
		{BulbarInvolvement: true}, {ErectileDysfunction: true}, {KarnofskyBelow70: true},
	}
	for i, f := range items {
		if got := DeltaPScore(f); got != 1 {
			t.Errorf("item %d: got %d, want 1", i, got)
		}
	}
}

func TestDeltaPScore_noAlloc(t *testing.T) {
	f := DeltaPFindings{BulbarInvolvement: true}
	if n := testing.AllocsPerRun(100, func() { sinkInt = DeltaPScore(f) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkDeltaPScore(b *testing.B) {
	b.ReportAllocs()
	f := DeltaPFindings{BulbarInvolvement: true, KarnofskyBelow70: true}
	for b.Loop() {
		sinkInt = DeltaPScore(f)
	}
}
