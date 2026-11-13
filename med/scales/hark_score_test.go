package scales

import "testing"

// Sohal 2007 Table 1: one point per yes, 0–4. All 16 combinations.
func TestHARKScore_example(t *testing.T) {
	for mask := 0; mask < 16; mask++ {
		h, a, r, k := mask&1 != 0, mask&2 != 0, mask&4 != 0, mask&8 != 0
		want := 0
		for m := mask; m != 0; m &= m - 1 {
			want++
		}
		if got := HARKScore(h, a, r, k); got != want {
			t.Errorf("HARKScore(%v, %v, %v, %v) = %d, want %d", h, a, r, k, got, want)
		}
	}
}

// parity: ledger hark_ipv, source cases 1–5 (score only; the source's band labels are not the paper's).
func TestHARKScore_parity(t *testing.T) {
	tests := []struct {
		h, a, r, k bool
		want       int
	}{
		{false, false, false, false, 0}, // source case 1 "negative"
		{true, false, false, false, 1},  // source case 2 "one"
		{true, true, false, false, 2},   // source case 3 "two"
		{true, true, true, false, 3},    // source case 4 "three"
		{true, true, true, true, 4},     // source case 5 "four"
	}
	for i, tc := range tests {
		if got := HARKScore(tc.h, tc.a, tc.r, tc.k); got != tc.want {
			t.Errorf("source case %d: HARKScore = %d, want %d", i+1, got, tc.want)
		}
	}
}

func TestHARKScore_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt = HARKScore(true, false, true, false) }); n != 0 {
		t.Errorf("HARKScore allocates %v times, want 0", n)
	}
}

func BenchmarkHARKScore(b *testing.B) {
	b.ReportAllocs()
	x := true
	for b.Loop() {
		sinkInt = HARKScore(x, !x, x, !x)
	}
}
