package neuro

import (
	"errors"
	"testing"
)

// HAND-DERIVED from Lima 2016 Table 1: the sum of the item points.
func TestFASTEDScore_example(t *testing.T) {
	tests := []struct{ f, a, s, e, n, want int }{
		{0, 0, 0, 0, 0, 0},
		{1, 2, 1, 1, 0, 5},
		{1, 2, 2, 2, 2, 9},
		{0, 1, 1, 1, 0, 3},
		{1, 1, 1, 1, 0, 4},
	}
	for _, tc := range tests {
		if got, err := FASTEDScore(tc.f, tc.a, tc.s, tc.e, tc.n); err != nil || got != tc.want {
			t.Errorf("FASTEDScore(%d, %d, %d, %d, %d) = (%d, %v), want %d", tc.f, tc.a, tc.s, tc.e, tc.n, got, err, tc.want)
		}
	}
}

func TestFASTEDScore_errors(t *testing.T) {
	for _, in := range [][5]int{
		{-1, 0, 0, 0, 0}, {2, 0, 0, 0, 0}, {0, 3, 0, 0, 0}, {0, -1, 0, 0, 0},
		{0, 0, 3, 0, 0}, {0, 0, 0, 3, 0}, {0, 0, 0, 0, 3}, {0, 0, 0, 0, -1},
	} {
		if _, err := FASTEDScore(in[0], in[1], in[2], in[3], in[4]); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("FASTEDScore(%v) error = %v, want ErrOutOfRange", in, err)
		}
	}
}

func TestFASTEDScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = FASTEDScore(1, 2, 1, 1, 0) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkFASTEDScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = FASTEDScore(1, 2, 1, 1, 0)
	}
}
