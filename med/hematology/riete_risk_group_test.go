package hematology

import (
	"errors"
	"math"
	"testing"
)

// Ruíz-Giménez 2008 abstract: 0 "low risk"; 1–4 "intermediate"; > 4 "high risk".
func TestRIETERiskGroup_example(t *testing.T) {
	tests := []struct {
		s    float64
		want RIETEGroup
	}{
		{0, RIETELow}, {1, RIETEIntermediate}, {1.5, RIETEIntermediate}, {4, RIETEIntermediate},
		{4.5, RIETEHigh}, {8, RIETEHigh},
	}
	for _, tc := range tests {
		if got, err := RIETERiskGroup(tc.s); err != nil || got != tc.want {
			t.Errorf("RIETERiskGroup(%v) = (%q, %v), want %q", tc.s, got, err, tc.want)
		}
	}
}

func TestRIETERiskGroup_errors(t *testing.T) {
	cases := []struct {
		s    float64
		want error
	}{
		{math.NaN(), ErrNotFinite}, {math.Inf(1), ErrNotFinite}, {math.Inf(-1), ErrNotFinite},
		{-0.5, ErrOutOfRange}, {8.5, ErrOutOfRange}, {0.5, ErrOutOfRange}, {1.25, ErrOutOfRange},
	}
	for _, c := range cases {
		if _, err := RIETERiskGroup(c.s); !errors.Is(err, c.want) {
			t.Errorf("RIETERiskGroup(%v) error = %v, want %v", c.s, err, c.want)
		}
	}
}

func TestRIETERiskGroup_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { g, _ := RIETERiskGroup(4.5); sinkString = string(g) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkRIETERiskGroup(b *testing.B) {
	b.ReportAllocs()
	s := 4.5
	for b.Loop() {
		g, _ := RIETERiskGroup(s)
		sinkString = string(g)
	}
}
