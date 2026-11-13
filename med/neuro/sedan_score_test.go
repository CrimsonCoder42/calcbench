package neuro

import (
	"errors"
	"math"
	"testing"
)

// HAND-DERIVED from Strbian 2012 (Ann Neurol) abstract points; each edge on both sides.
func TestSEDANScore_example(t *testing.T) {
	tests := []struct {
		glu          float64
		early, dense bool
		age          float64
		nihss        int
		want         int
	}{
		{8.0, false, false, 75, 9, 0},
		{8.05, false, false, 50, 0, 1}, // between printed bands: counted in 8.1–12.0
		{8.1, false, false, 50, 0, 1},
		{12.0, false, false, 50, 0, 1},
		{12.01, false, false, 50, 0, 2},
		{5, true, false, 50, 0, 1},
		{5, false, true, 50, 0, 1},
		{5, false, false, 75.1, 0, 1},
		{5, false, false, 50, 10, 1},
		{13, true, false, 80, 12, 5}, // doc example
		{20, true, true, 90, 25, 6},
	}
	for _, tc := range tests {
		if got, err := SEDANScore(tc.glu, tc.early, tc.dense, tc.age, tc.nihss); err != nil || got != tc.want {
			t.Errorf("SEDANScore(%v, %v, %v, %v, %d) = (%d, %v), want %d", tc.glu, tc.early, tc.dense, tc.age, tc.nihss, got, err, tc.want)
		}
	}
}

func TestSEDANScore_errors(t *testing.T) {
	tests := []struct {
		glu, age float64
		nihss    int
		want     error
	}{
		{math.NaN(), 50, 0, ErrNotFinite}, {5, math.Inf(1), 0, ErrNotFinite},
		{-1, 50, 0, ErrOutOfRange}, {5, -1, 0, ErrOutOfRange}, {5, 50, -1, ErrOutOfRange},
	}
	for _, tc := range tests {
		if _, err := SEDANScore(tc.glu, false, false, tc.age, tc.nihss); !errors.Is(err, tc.want) {
			t.Errorf("SEDANScore(%v, _, _, %v, %d) error = %v, want %v", tc.glu, tc.age, tc.nihss, err, tc.want)
		}
	}
}

func TestSEDANScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = SEDANScore(13, true, false, 80, 12) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkSEDANScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = SEDANScore(13, true, false, 80, 12)
	}
}
