package neuro

import (
	"errors"
	"testing"
)

// HAND-DERIVED from Katz 2015 abstract: gaze ≥ 1 → 2; arm ≥ 2 → 1; LOC questions ≥ 1 and commands ≥ 1 → 1.
func TestCPSSSScore_example(t *testing.T) {
	tests := []struct{ gaze, arm, q, c, want int }{
		{0, 1, 0, 1, 0}, // arm 1 below its edge; commands alone do not count
		{1, 0, 0, 0, 2},
		{0, 2, 0, 0, 1},
		{0, 0, 1, 1, 1},
		{0, 0, 1, 0, 0},
		{1, 3, 1, 0, 3}, // doc example
		{2, 4, 2, 2, 4},
	}
	for _, tc := range tests {
		if got, err := CPSSSScore(tc.gaze, tc.arm, tc.q, tc.c); err != nil || got != tc.want {
			t.Errorf("CPSSSScore(%d, %d, %d, %d) = (%d, %v), want %d", tc.gaze, tc.arm, tc.q, tc.c, got, err, tc.want)
		}
	}
}

func TestCPSSSScore_errors(t *testing.T) {
	for _, in := range [][4]int{{-1, 0, 0, 0}, {0, -1, 0, 0}, {0, 0, -1, 0}, {0, 0, 0, -1}} {
		if _, err := CPSSSScore(in[0], in[1], in[2], in[3]); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("CPSSSScore(%v) error = %v, want ErrOutOfRange", in, err)
		}
	}
}

func TestCPSSSScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = CPSSSScore(1, 3, 1, 0) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkCPSSSScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = CPSSSScore(1, 3, 1, 0)
	}
}
