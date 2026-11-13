package trauma

import "testing"

// Rodriguez 2015: CT-All seven criteria; CT-Major "the same criteria without rapid deceleration mechanism".
func TestNEXUSChestCT_source(t *testing.T) {
	tests := []struct {
		f          NEXUSChestCTFindings
		all, major bool
	}{
		{NEXUSChestCTFindings{}, false, false},
		{NEXUSChestCTFindings{RapidDeceleration: true}, true, false},
		{NEXUSChestCTFindings{AbnormalCXR: true}, true, true},
		{NEXUSChestCTFindings{ScapularTenderness: true}, true, true},
		{NEXUSChestCTFindings{SternalTenderness: true, RapidDeceleration: true}, true, true},
		{NEXUSChestCTFindings{ThoracicSpineTenderness: true}, true, true},
		{NEXUSChestCTFindings{DistractingInjury: true}, true, true},
		{NEXUSChestCTFindings{ChestWallTenderness: true}, true, true},
	}
	for i, tc := range tests {
		if r := NEXUSChestCT(tc.f); r.AnyChestCTAll != tc.all || r.AnyChestCTMajor != tc.major {
			t.Errorf("case %d: %+v, want all %v major %v", i, r, tc.all, tc.major)
		}
	}
}

// parity: ledger nexus_chest_ct, source case 1 (app: anyMajor over six items, rapid deceleration separate).
func TestNEXUSChestCT_parity(t *testing.T) {
	if r := NEXUSChestCT(NEXUSChestCTFindings{}); r.AnyChestCTAll {
		t.Errorf("none → %+v", r)
	}
}

func TestNEXUSChestCT_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkBool = NEXUSChestCT(NEXUSChestCTFindings{AbnormalCXR: true}).AnyChestCTAll }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkNEXUSChestCT(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkBool = NEXUSChestCT(NEXUSChestCTFindings{AbnormalCXR: true}).AnyChestCTAll
	}
}
