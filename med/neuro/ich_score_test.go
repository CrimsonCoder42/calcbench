package neuro

import (
	"errors"
	"math"
	"testing"
)

// HAND-DERIVED from Hemphill 2001 abstract points; each edge on both sides.
func TestICHScore_example(t *testing.T) {
	tests := []struct {
		gcs   int
		age   float64
		infra bool
		vol   float64
		ivh   bool
		want  int
	}{
		{15, 79.9, false, 29.9, false, 0},
		{13, 50, false, 0, false, 0},
		{12, 50, false, 0, false, 1},
		{5, 50, false, 0, false, 1},
		{4, 50, false, 0, false, 2},
		{3, 50, false, 0, false, 2},
		{15, 80, false, 0, false, 1},
		{15, 50, true, 0, false, 1},
		{15, 50, false, 30, false, 1},
		{15, 50, false, 0, true, 1},
		{10, 82, false, 35, true, 4}, // doc example
		{3, 90, true, 60, true, 6},
	}
	for _, tc := range tests {
		if got, err := ICHScore(tc.gcs, tc.age, tc.infra, tc.vol, tc.ivh); err != nil || got != tc.want {
			t.Errorf("ICHScore(%d, %v, %v, %v, %v) = (%d, %v), want %d", tc.gcs, tc.age, tc.infra, tc.vol, tc.ivh, got, err, tc.want)
		}
	}
}

func TestICHScore_errors(t *testing.T) {
	tests := []struct {
		gcs      int
		age, vol float64
		want     error
	}{
		{2, 50, 10, ErrOutOfRange}, {16, 50, 10, ErrOutOfRange},
		{10, math.NaN(), 10, ErrNotFinite}, {10, 50, math.Inf(1), ErrNotFinite},
		{10, -1, 10, ErrOutOfRange}, {10, 50, -1, ErrOutOfRange},
	}
	for _, tc := range tests {
		if _, err := ICHScore(tc.gcs, tc.age, false, tc.vol, false); !errors.Is(err, tc.want) {
			t.Errorf("ICHScore(%d, %v, _, %v, _) error = %v, want %v", tc.gcs, tc.age, tc.vol, err, tc.want)
		}
	}
}

func TestICHScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = ICHScore(10, 82, false, 35, true) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkICHScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = ICHScore(10, 82, false, 35, true)
	}
}
