package hematology

import (
	"errors"
	"math"
	"testing"
)

// Thieblemont 2017 abstract: age ≥ 70, stage III/IV, elevated LDH, one each.
func TestMALTIPIScore_example(t *testing.T) {
	tests := []struct {
		age        float64
		stage, ldh bool
		want       int
	}{
		{69.9, false, false, 0}, {70, false, false, 1}, {50, true, false, 1}, {50, false, true, 1},
		{70, true, false, 2}, {80, true, true, 3},
	}
	for _, tc := range tests {
		if got, err := MALTIPIScore(tc.age, tc.stage, tc.ldh); err != nil || got != tc.want {
			t.Errorf("MALTIPIScore(%v, %v, %v) = (%d, %v), want %d", tc.age, tc.stage, tc.ldh, got, err, tc.want)
		}
	}
	for _, c := range []struct {
		age  float64
		want error
	}{{math.NaN(), ErrNotFinite}, {math.Inf(1), ErrNotFinite}, {math.Inf(-1), ErrNotFinite}, {-1, ErrOutOfRange}, {131, ErrOutOfRange}} {
		if _, err := MALTIPIScore(c.age, false, false); !errors.Is(err, c.want) {
			t.Errorf("MALTIPIScore(%v) error = %v, want %v", c.age, err, c.want)
		}
	}
}

func TestMALTIPIScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = MALTIPIScore(72, true, false) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkMALTIPIScore(b *testing.B) {
	b.ReportAllocs()
	age := 72.0
	for b.Loop() {
		sinkInt, _ = MALTIPIScore(age, true, false)
	}
}
