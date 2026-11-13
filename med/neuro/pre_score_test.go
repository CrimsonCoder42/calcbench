package neuro

import (
	"errors"
	"math"
	"testing"
)

// HAND-DERIVED from Rangaraju 2015 abstract formula: age + 2 × NIHSS − 10 × ASPECTS.
func TestPREScore_example(t *testing.T) {
	tests := []struct {
		age            float64
		nihss, aspects int
		want           float64
	}{
		{70, 20, 8, 30},  // 70 + 40 − 80
		{40, 2, 10, -56}, // 40 + 4 − 100
		{80, 20, 4, 80},  // 80 + 40 − 40
		{0, 0, 0, 0},
		{59.5, 10, 6, 19.5}, // 59.5 + 20 − 60
	}
	for _, tc := range tests {
		got, err := PREScore(tc.age, tc.nihss, tc.aspects)
		// exact: small integers and one half are exact in float64.
		if err != nil || got != tc.want {
			t.Errorf("PREScore(%v, %d, %d) = (%v, %v), want %v", tc.age, tc.nihss, tc.aspects, got, err, tc.want)
		}
	}
}

func TestPREScore_errors(t *testing.T) {
	tests := []struct {
		age            float64
		nihss, aspects int
		want           error
	}{
		{math.NaN(), 10, 5, ErrNotFinite}, {math.Inf(1), 10, 5, ErrNotFinite},
		{-1, 10, 5, ErrOutOfRange}, {50, -1, 5, ErrOutOfRange}, {50, 10, -1, ErrOutOfRange}, {50, 10, 11, ErrOutOfRange},
	}
	for _, tc := range tests {
		if _, err := PREScore(tc.age, tc.nihss, tc.aspects); !errors.Is(err, tc.want) {
			t.Errorf("PREScore(%v, %d, %d) error = %v, want %v", tc.age, tc.nihss, tc.aspects, err, tc.want)
		}
	}
}

func TestPREScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = PREScore(70, 20, 8) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkPREScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = PREScore(70, 20, 8)
	}
}
