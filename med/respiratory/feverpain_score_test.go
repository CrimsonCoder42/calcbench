package respiratory

import "testing"

// HAND-DERIVED from Little 2013 BMJ Open Table 2 (PMID 24163209), 1 point each: none 0, all five 5, three 3.
func TestFeverPAINScore_tableExample(t *testing.T) {
	cases := []struct {
		f, p, a, i, n bool
		want          int
	}{
		{false, false, false, false, false, 0},
		{true, true, true, true, true, 5},
		{true, true, true, false, false, 3},
		{false, false, false, false, true, 1},
		{false, true, false, true, false, 2},
	}
	for _, c := range cases {
		if got := FeverPAINScore(c.f, c.p, c.a, c.i, c.n); got != c.want {
			t.Errorf("%v: got %d, want %d", c, got, c.want)
		}
	}
}

// parity: ledger feverpain, source case 1: the source sums the five yes/no items (fever 1, purulence 1,
// attendance ≤ 3 days 1, inflamed 1, no cough/coryza 1); 4 items → 4.
func TestFeverPAINScore_parity(t *testing.T) {
	if got := FeverPAINScore(true, true, true, true, false); got != 4 {
		t.Errorf("got %d, want 4", got)
	}
}

func TestFeverPAINScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt = FeverPAINScore(true, false, true, false, true) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkFeverPAINScore(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt = FeverPAINScore(true, false, true, false, true)
	}
}
