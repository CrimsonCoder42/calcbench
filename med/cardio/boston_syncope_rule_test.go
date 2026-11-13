package cardio

import "testing"

// Grossman 2007: any of the eight groups present makes the rule positive.
func TestBostonSyncopeRuleCount(t *testing.T) {
	if c, any := BostonSyncopeRuleCount(BostonSyncopeFindings{}); c != 0 || any {
		t.Errorf("none = (%d, %v)", c, any)
	}
	if c, any := BostonSyncopeRuleCount(BostonSyncopeFindings{VolumeDepletion: true}); c != 1 || !any {
		t.Errorf("one = (%d, %v)", c, any)
	}
	all := BostonSyncopeFindings{true, true, true, true, true, true, true, true}
	if c, any := BostonSyncopeRuleCount(all); c != 8 || !any {
		t.Errorf("all = (%d, %v)", c, any)
	}
}

// parity: ledger boston_syncope — no source test file or JSON example; the app counts the same eight groups.
func TestBostonSyncopeRuleCount_parity(t *testing.T) {
	c, any := BostonSyncopeRuleCount(BostonSyncopeFindings{AcuteCoronarySyndrome: true, PrimaryCNSEvent: true}) // parity: ledger boston_syncope, source case 1
	if c != 2 || !any {
		t.Errorf("(%d, %v)", c, any)
	}
}

func TestBostonSyncopeRuleCount_noAlloc(t *testing.T) {
	f := BostonSyncopeFindings{VolumeDepletion: true}
	if n := testing.AllocsPerRun(100, func() { sinkInt, sinkBool = BostonSyncopeRuleCount(f) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkBostonSyncopeRuleCount(b *testing.B) {
	b.ReportAllocs()
	f := BostonSyncopeFindings{VolumeDepletion: true}
	for b.Loop() {
		sinkInt, sinkBool = BostonSyncopeRuleCount(f)
	}
}
