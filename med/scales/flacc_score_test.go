package scales

import (
	"errors"
	"testing"
)

var sinkFLACC int

// Uitti 2018 Table 1 and Curtis 2007 (both crediting Merkel 1997): five categories each 0–2, total 0–10. Sums, the bounds, and each item at 0, 1 and 2. HAND-DERIVED from the sum.
func TestFLACCScore_source(t *testing.T) {
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
		got, err := FLACCScore(tc.in[0], tc.in[1], tc.in[2], tc.in[3], tc.in[4])
		if err != nil || got != tc.want {
			t.Errorf("%v → %d (%v), want %d", tc.in, got, err, tc.want)
		}
	}
	for pos := 0; pos < 5; pos++ {
		for v := 0; v <= 2; v++ {
			var in [5]int
			in[pos] = v
			if got, err := FLACCScore(in[0], in[1], in[2], in[3], in[4]); err != nil || got != v {
				t.Errorf("item %d = %d → %d (%v), want %d", pos, v, got, err, v)
			}
		}
	}
}

// parity: ledger flacc_scale, source case 1. No _test.go or JSON example; the entry sums the same five 0–2
// inputs (face, legs, activity, cry, consolability), so its maximum is 10.
func TestFLACCScore_parity(t *testing.T) {
	if got, _ := FLACCScore(2, 2, 2, 2, 2); got != 10 {
		t.Errorf("got %d, want 10", got)
	}
}

func TestFLACCScore_errors(t *testing.T) {
	for pos := 0; pos < 5; pos++ {
		for _, v := range []int{-1, 3} {
			var in [5]int
			in[pos] = v
			if _, err := FLACCScore(in[0], in[1], in[2], in[3], in[4]); !errors.Is(err, ErrOutOfRange) {
				t.Errorf("item %d = %d: err = %v, want ErrOutOfRange", pos, v, err)
			}
		}
	}
}

func TestFLACCScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFLACC, _ = FLACCScore(1, 0, 1, 2, 1) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkFLACCScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFLACC, _ = FLACCScore(1, 0, 1, 2, 1)
	}
}
