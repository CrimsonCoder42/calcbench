package neuro

import (
	"errors"
	"testing"
)

// HAND-DERIVED from Delgado Almandoz 2010 Table 2; age edges on both sides.
func TestSICHScore_example(t *testing.T) {
	tests := []struct {
		ncct    SICHNCCT
		age     int
		female  bool
		neither bool
		want    int
	}{
		{SICHNCCTLowProbability, 71, false, false, 0},
		{SICHNCCTLowProbability, 70, false, false, 1},
		{SICHNCCTLowProbability, 46, false, false, 1},
		{SICHNCCTLowProbability, 45, false, false, 2},
		{SICHNCCTLowProbability, 18, false, false, 2},
		{SICHNCCTIndeterminate, 80, false, false, 1},
		{SICHNCCTHighProbability, 80, false, false, 2},
		{SICHNCCTLowProbability, 80, true, false, 1},
		{SICHNCCTLowProbability, 80, false, true, 1},
		{SICHNCCTHighProbability, 40, true, true, 6}, // doc example, the maximum
	}
	for _, tc := range tests {
		if got, err := SICHScore(tc.ncct, tc.age, tc.female, tc.neither); err != nil || got != tc.want {
			t.Errorf("SICHScore(%d, %d, %v, %v) = (%d, %v), want %d", tc.ncct, tc.age, tc.female, tc.neither, got, err, tc.want)
		}
	}
}

func TestSICHScore_errors(t *testing.T) {
	for _, tc := range []struct {
		ncct SICHNCCT
		age  int
	}{{-1, 50}, {3, 50}, {0, 17}} {
		if _, err := SICHScore(tc.ncct, tc.age, false, false); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("SICHScore(%d, %d) error = %v, want ErrOutOfRange", tc.ncct, tc.age, err)
		}
	}
}

func TestSICHScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = SICHScore(SICHNCCTHighProbability, 40, true, true) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkSICHScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = SICHScore(SICHNCCTHighProbability, 40, true, true)
	}
}
