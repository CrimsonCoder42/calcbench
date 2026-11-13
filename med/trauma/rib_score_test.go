package trauma

import "testing"

// Chapman 2016 abstract: one point per variable.
func TestRibScore_source(t *testing.T) {
	if got := RibScore(RibScoreFindings{}); got != 0 {
		t.Errorf("none → %d", got)
	}
	if got := RibScore(RibScoreFindings{BilateralFractures: true, FlailChest: true}); got != 2 {
		t.Errorf("two → %d", got)
	}
	if got := RibScore(RibScoreFindings{true, true, true, true, true, true}); got != 6 {
		t.Errorf("all → %d", got)
	}
}

// parity: ledger rib_score, source case 1 (app sums six 0/1 items).
func TestRibScore_parity(t *testing.T) {
	if got := RibScore(RibScoreFindings{SixOrMoreRibFractures: true, FirstRibFracture: true, FractureInAllThreeAnatomicAreas: true}); got != 3 {
		t.Errorf("got %d, want 3", got)
	}
}

func TestRibScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt = RibScore(RibScoreFindings{FlailChest: true}) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkRibScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt = RibScore(RibScoreFindings{FlailChest: true})
	}
}
