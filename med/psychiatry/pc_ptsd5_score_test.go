package psychiatry

import "testing"

// Prins 2016: exposure "No" → 0; otherwise the yes count. HAND-DERIVED edges.
func TestPCPTSD5Score_example(t *testing.T) {
	all := [5]bool{true, true, true, true, true}
	if got := PCPTSD5Score(false, all); got != 0 {
		t.Errorf("no exposure, all yes = %d, want 0", got)
	}
	if got := PCPTSD5Score(true, all); got != 5 {
		t.Errorf("exposure, all yes = %d, want 5", got)
	}
	if got := PCPTSD5Score(true, [5]bool{true, false, true}); got != 2 {
		t.Errorf("exposure, two yes = %d, want 2", got)
	}
}

// parity: ledger pc_ptsd5, source cases 1–6. Items in source order: nightmares, avoidance, onGuard, numb, guilt.
func TestPCPTSD5Score_parity(t *testing.T) {
	tests := []struct {
		exposure bool
		items    [5]bool
		want     int
	}{
		{false, [5]bool{true, true, true, true, true}, 0}, // source case 1 "no trauma, symptoms ignored"
		{true, [5]bool{}, 0},                              // source case 2 "trauma, none"
		{true, [5]bool{true, true}, 2},                    // source case 3 "two"
		{true, [5]bool{true, true, true}, 3},              // source case 4 "three"
		{true, [5]bool{true, true, true, true}, 4},        // source case 5 "four"
		{true, [5]bool{true, true, true, true, true}, 5},  // source case 6 "five"
	}
	for i, tc := range tests {
		if got := PCPTSD5Score(tc.exposure, tc.items); got != tc.want {
			t.Errorf("source case %d: PCPTSD5Score = %d, want %d", i+1, got, tc.want)
		}
	}
}

func TestPCPTSD5Score_noAlloc(t *testing.T) {
	in := [5]bool{true, true, false, true}
	if n := testing.AllocsPerRun(100, func() { sinkInt = PCPTSD5Score(true, in) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkPCPTSD5Score(b *testing.B) {
	b.ReportAllocs()
	in := [5]bool{true, true, false, true}
	for b.Loop() {
		sinkInt = PCPTSD5Score(true, in)
	}
}
