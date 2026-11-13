package cardio

import (
	"errors"
	"math"
	"testing"
)

// Tanaka 2001 (full text, comparing 220 − age): at 70 years the equation gives 208 − 0.7 × 70 = 159 beats/min.
func TestMaxHeartRateTanaka_sourceExample(t *testing.T) {
	got, err := MaxHeartRateTanaka(70)
	// 0.7 × 70 is 49.0000000000000071 in float64; 1e-12 absolute covers that rounding.
	if err != nil || math.Abs(got-159) > 1e-12 {
		t.Errorf("MaxHeartRateTanaka(70) = (%v, %v), want (159, nil)", got, err)
	}
}

func TestMaxHeartRateTanaka_rangeEdges(t *testing.T) {
	tests := []struct {
		name    string
		age     float64
		want    float64 // 208 − 0.7 × age, hand-derived
		wantErr error
	}{
		{"lower edge 18", 18, 195.4, nil},
		{"upper edge 81", 81, 151.3, nil},
		{"just below 18", 17.999, 0, ErrOutOfRange},
		{"just above 81", 81.001, 0, ErrOutOfRange},
		{"NaN", math.NaN(), 0, ErrNotFinite},
		{"+Inf", math.Inf(1), 0, ErrNotFinite},
		{"-Inf", math.Inf(-1), 0, ErrNotFinite},
	}
	for _, tc := range tests {
		got, err := MaxHeartRateTanaka(tc.age)
		if !errors.Is(err, tc.wantErr) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.wantErr)
			continue
		}
		if tc.wantErr == nil && math.Abs(got-tc.want) > 1e-12 {
			t.Errorf("%s: got %v, want %v", tc.name, got, tc.want)
		}
	}
}

var sinkFloat float64

func TestMaxHeartRateTanaka_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = MaxHeartRateTanaka(45) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkMaxHeartRateTanaka(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = MaxHeartRateTanaka(45)
	}
}
