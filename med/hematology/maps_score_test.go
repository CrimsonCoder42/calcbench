package hematology

import (
	"errors"
	"math"
	"testing"
)

// Pardanani 2018 Figure 3: advanced SM 2; age > 60, platelets < 150, ALP high, adverse mutation 1 each.
func TestMAPSScore_example(t *testing.T) {
	tests := []struct {
		adv      bool
		age, plt float64
		alp, mut bool
		want     int
	}{
		{false, 60, 150, false, false, 0}, {false, 60.1, 150, false, false, 1}, {false, 60, 149.9, false, false, 1},
		{true, 60, 150, false, false, 2}, {false, 60, 150, true, false, 1}, {false, 60, 150, false, true, 1},
		{true, 61, 149, true, true, 6},
	}
	for _, tc := range tests {
		if got, err := MAPSScore(tc.adv, tc.age, tc.plt, tc.alp, tc.mut); err != nil || got != tc.want {
			t.Errorf("MAPSScore(%+v) = (%d, %v), want %d", tc, got, err, tc.want)
		}
	}
}

func TestMAPSScore_errors(t *testing.T) {
	cases := []struct {
		age, plt float64
		want     error
	}{
		{math.NaN(), 100, ErrNotFinite}, {50, math.Inf(1), ErrNotFinite}, {math.Inf(-1), 100, ErrNotFinite},
		{-1, 100, ErrOutOfRange}, {131, 100, ErrOutOfRange}, {50, -1, ErrOutOfRange}, {50, 5001, ErrOutOfRange},
	}
	for _, c := range cases {
		if _, err := MAPSScore(false, c.age, c.plt, false, false); !errors.Is(err, c.want) {
			t.Errorf("MAPSScore(age %v, plt %v) error = %v, want %v", c.age, c.plt, err, c.want)
		}
	}
}

func TestMAPSScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = MAPSScore(true, 65, 120, true, false) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkMAPSScore(b *testing.B) {
	b.ReportAllocs()
	age, plt := 65.0, 120.0
	for b.Loop() {
		sinkInt, _ = MAPSScore(true, age, plt, true, false)
	}
}
