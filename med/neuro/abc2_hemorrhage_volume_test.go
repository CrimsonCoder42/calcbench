package neuro

import (
	"errors"
	"math"
	"testing"
)

var (
	sinkInt    int
	sinkFloat  float64
	sinkString string
)

// HAND-DERIVED from Kothari 1996 abstract: A × B × C / 2.
func TestABC2HemorrhageVolumeML_example(t *testing.T) {
	tests := []struct{ a, b, c, want float64 }{
		{4, 3, 2.5, 15}, // 4 × 3 × 2.5 / 2
		{1, 1, 1, 0.5},  // unit cube halves
		{6, 5, 4, 60},   // 120 / 2
		{0.1, 0.1, 0.1, 0.0005},
	}
	for _, tc := range tests {
		got, err := ABC2HemorrhageVolumeML(tc.a, tc.b, tc.c)
		// tolerance 1e-12 relative: one multiply chain and a halving.
		if err != nil || math.Abs(got-tc.want) > 1e-12*tc.want {
			t.Errorf("ABC2HemorrhageVolumeML(%v, %v, %v) = (%v, %v), want %v", tc.a, tc.b, tc.c, got, err, tc.want)
		}
	}
}

func TestABC2HemorrhageVolumeML_errors(t *testing.T) {
	tests := []struct {
		a, b, c float64
		want    error
	}{
		{math.NaN(), 1, 1, ErrNotFinite},
		{1, math.Inf(1), 1, ErrNotFinite},
		{1, 1, math.Inf(-1), ErrNotFinite},
		{0, 1, 1, ErrOutOfRange},
		{1, 0, 1, ErrOutOfRange},
		{1, 1, 0, ErrOutOfRange},
		{-1, 1, 1, ErrOutOfRange},
		{math.MaxFloat64, math.MaxFloat64, 2, ErrOutOfRange}, // overflow
	}
	for _, tc := range tests {
		if _, err := ABC2HemorrhageVolumeML(tc.a, tc.b, tc.c); !errors.Is(err, tc.want) {
			t.Errorf("ABC2HemorrhageVolumeML(%v, %v, %v) error = %v, want %v", tc.a, tc.b, tc.c, err, tc.want)
		}
	}
}

func TestABC2HemorrhageVolumeML_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = ABC2HemorrhageVolumeML(4, 3, 2.5) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkABC2HemorrhageVolumeML(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = ABC2HemorrhageVolumeML(4, 3, 2.5)
	}
}
