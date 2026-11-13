package neuro

import "testing"

// HAND-DERIVED from Struck 2017: BIRDs 2, five findings 1 each.
func TestTwoHELPS2BScore_example(t *testing.T) {
	tests := []struct {
		in   [6]bool
		want int
	}{
		{[6]bool{}, 0},
		{[6]bool{true}, 2},
		{[6]bool{false, true}, 1},
		{[6]bool{false, false, true}, 1},
		{[6]bool{false, false, false, true}, 1},
		{[6]bool{false, false, false, false, true}, 1},
		{[6]bool{false, false, false, false, false, true}, 1},
		{[6]bool{true, true, false, false, true, false}, 4}, // doc example
		{[6]bool{true, true, true, true, true, true}, 7},
	}
	for _, tc := range tests {
		b := tc.in
		if got := TwoHELPS2BScore(b[0], b[1], b[2], b[3], b[4], b[5]); got != tc.want {
			t.Errorf("TwoHELPS2BScore(%v) = %d, want %d", b, got, tc.want)
		}
	}
}

func TestTwoHELPS2BScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt = TwoHELPS2BScore(true, true, false, false, true, false) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkTwoHELPS2BScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt = TwoHELPS2BScore(true, true, false, false, true, false)
	}
}
