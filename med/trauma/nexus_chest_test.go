package trauma

import "testing"

// Rodriguez 2013 abstract: seven criteria, age ">60 years". HAND-DERIVED.
func TestNEXUSChest_source(t *testing.T) {
	tests := []struct {
		age  float64
		f    NEXUSChestFindings
		want int
	}{
		{60, NEXUSChestFindings{}, 0},
		{60.01, NEXUSChestFindings{}, 1},
		{61, NEXUSChestFindings{true, true, true, true, true, true}, 7},
		{30, NEXUSChestFindings{ChestPain: true, Intoxication: true}, 2},
	}
	for i, tc := range tests {
		if got, err := NEXUSChest(tc.age, tc.f); err != nil || got != tc.want {
			t.Errorf("case %d: (%d, %v), want %d", i, got, err, tc.want)
		}
	}
}

// parity: ledger nexus_chest, source case 1 (app sums seven 0/1 items; age_60 = 1 is age 70 here).
func TestNEXUSChest_parity(t *testing.T) {
	if got, _ := NEXUSChest(70, NEXUSChestFindings{RapidDeceleration: true}); got != 2 {
		t.Errorf("got %d, want 2", got)
	}
}

func TestNEXUSChest_errors(t *testing.T) {
	for _, v := range nonFinite {
		_, err := NEXUSChest(v, NEXUSChestFindings{})
		checkErr(t, "non-finite", err, ErrNotFinite)
	}
	for _, v := range []float64{-1, 131} {
		_, err := NEXUSChest(v, NEXUSChestFindings{})
		checkErr(t, "range", err, ErrOutOfRange)
	}
}

func TestNEXUSChest_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = NEXUSChest(61, NEXUSChestFindings{ChestPain: true}) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkNEXUSChest(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = NEXUSChest(61, NEXUSChestFindings{ChestPain: true})
	}
}
