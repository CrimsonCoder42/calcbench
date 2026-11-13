package infectious

import "testing"

// Menza 2009 Table 2 weights: STI 4, methamphetamine/nitrites 11, ≥ 10 partners 3, non-concordant UAS 1.
func TestMenzaScore_example(t *testing.T) {
	tests := []struct {
		sti, meth, p10, uas bool
		want                int
	}{
		{false, false, false, false, 0}, {true, false, false, false, 4}, {false, true, false, false, 11},
		{false, false, true, false, 3}, {false, false, false, true, 1}, {true, true, false, false, 15},
		{true, true, true, true, 19},
	}
	for _, tc := range tests {
		if got := MenzaScore(tc.sti, tc.meth, tc.p10, tc.uas); got != tc.want {
			t.Errorf("MenzaScore(%+v) = %d, want %d", tc, got, tc.want)
		}
	}
}

func TestMenzaScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt = MenzaScore(true, false, true, true) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkMenzaScore(b *testing.B) {
	b.ReportAllocs()
	x := true
	for b.Loop() {
		sinkInt = MenzaScore(x, !x, x, x)
	}
}
