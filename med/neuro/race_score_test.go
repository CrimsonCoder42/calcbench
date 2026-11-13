package neuro

import (
	"errors"
	"testing"
)

// HAND-DERIVED from Pérez de la Ossa 2014 abstract item ranges: the sum.
func TestRACEScore_example(t *testing.T) {
	tests := []struct{ f, a, l, g, ap, want int }{
		{0, 0, 0, 0, 0, 0},
		{1, 2, 1, 1, 0, 5},
		{2, 2, 2, 1, 2, 9},
		{1, 1, 1, 1, 0, 4},
	}
	for _, tc := range tests {
		if got, err := RACEScore(tc.f, tc.a, tc.l, tc.g, tc.ap); err != nil || got != tc.want {
			t.Errorf("RACEScore(%d, %d, %d, %d, %d) = (%d, %v), want %d", tc.f, tc.a, tc.l, tc.g, tc.ap, got, err, tc.want)
		}
	}
}

func TestRACEScore_errors(t *testing.T) {
	for _, in := range [][5]int{
		{-1, 0, 0, 0, 0}, {3, 0, 0, 0, 0}, {0, 3, 0, 0, 0}, {0, 0, 3, 0, 0},
		{0, 0, 0, 2, 0}, {0, 0, 0, -1, 0}, {0, 0, 0, 0, 3},
	} {
		if _, err := RACEScore(in[0], in[1], in[2], in[3], in[4]); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("RACEScore(%v) error = %v, want ErrOutOfRange", in, err)
		}
	}
}

func TestRACEScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = RACEScore(1, 2, 1, 1, 0) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkRACEScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = RACEScore(1, 2, 1, 1, 0)
	}
}
