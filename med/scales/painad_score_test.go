package scales

import (
	"errors"
	"testing"
)

var sinkPAINAD int

// Shaw 2023 and Pinto 2015 Chart 1 (both citing Warden 2003): five items each 0–2, total 0–10. Sums, the bounds, and each item at 0, 1 and 2. HAND-DERIVED from the sum.
func TestPAINADScore_source(t *testing.T) {
	tests := []struct {
		in   [5]int
		want int
	}{
		{[5]int{0, 0, 0, 0, 0}, 0},
		{[5]int{2, 2, 2, 2, 2}, 10},
		{[5]int{1, 0, 1, 2, 1}, 5},
		{[5]int{1, 1, 2, 1, 0}, 5},
		{[5]int{0, 0, 0, 0, 1}, 1},
		{[5]int{2, 0, 0, 0, 0}, 2},
	}
	for _, tc := range tests {
		got, err := PAINADScore(tc.in[0], tc.in[1], tc.in[2], tc.in[3], tc.in[4])
		if err != nil || got != tc.want {
			t.Errorf("%v → %d (%v), want %d", tc.in, got, err, tc.want)
		}
	}
	for pos := 0; pos < 5; pos++ {
		for v := 0; v <= 2; v++ {
			var in [5]int
			in[pos] = v
			if got, err := PAINADScore(in[0], in[1], in[2], in[3], in[4]); err != nil || got != v {
				t.Errorf("item %d = %d → %d (%v), want %d", pos, v, got, err, v)
			}
		}
	}
}

// parity: ledger painad, source case 1. No _test.go or JSON example; the entry sums the same five 0–2
// inputs (breathing, negativeVocalization, facialExpression, bodyLanguage, consolability), so its maximum is 10.
func TestPAINADScore_parity(t *testing.T) {
	if got, _ := PAINADScore(2, 2, 2, 2, 2); got != 10 {
		t.Errorf("got %d, want 10", got)
	}
}

func TestPAINADScore_errors(t *testing.T) {
	for pos := 0; pos < 5; pos++ {
		for _, v := range []int{-1, 3} {
			var in [5]int
			in[pos] = v
			if _, err := PAINADScore(in[0], in[1], in[2], in[3], in[4]); !errors.Is(err, ErrOutOfRange) {
				t.Errorf("item %d = %d: err = %v, want ErrOutOfRange", pos, v, err)
			}
		}
	}
}

func TestPAINADScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkPAINAD, _ = PAINADScore(1, 0, 1, 2, 1) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkPAINADScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkPAINAD, _ = PAINADScore(1, 0, 1, 2, 1)
	}
}
