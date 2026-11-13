package hematology

import (
	"errors"
	"math"
	"testing"
)

// HAND-DERIVED from Passamonti 2017 abstract: 2 each for Hb < 11, blasts ≥ 3%, CALR-unmutated; 1 each for
// platelets < 150 and symptoms; 0.15 per year of age. Tolerance 1e-9 (0.15 is inexact in binary).
func TestMYSECPMScore_example(t *testing.T) {
	tests := []struct {
		age, hb, blasts float64
		calr            bool
		plt             float64
		sym             bool
		want            float64
	}{
		{60, 10.9, 3, true, 149, true, 17},
		{60, 11, 2.9, false, 150, false, 9},
		{60, 10.99, 2.9, false, 150, false, 11},
		{60, 11, 3, false, 150, false, 11},
		{60, 11, 2.9, true, 150, false, 11},
		{60, 11, 2.9, false, 149.9, false, 10},
		{60, 11, 2.9, false, 150, true, 10},
		{0, 12, 0, false, 300, false, 0},
		{45.5, 12, 0, false, 300, false, 6.825},
	}
	for _, tc := range tests {
		got, err := MYSECPMScore(tc.age, tc.hb, tc.blasts, tc.calr, tc.plt, tc.sym)
		if err != nil || math.Abs(got-tc.want) > 1e-9 {
			t.Errorf("MYSECPMScore(%+v) = (%v, %v), want %v", tc, got, err, tc.want)
		}
	}
}

func TestMYSECPMScore_errors(t *testing.T) {
	cases := []struct {
		age, hb, bl, plt float64
		want             error
	}{
		{math.NaN(), 12, 0, 200, ErrNotFinite}, {60, math.Inf(1), 0, 200, ErrNotFinite},
		{60, 12, math.Inf(-1), 200, ErrNotFinite}, {60, 12, 0, math.NaN(), ErrNotFinite},
		{-1, 12, 0, 200, ErrOutOfRange}, {131, 12, 0, 200, ErrOutOfRange}, {60, -1, 0, 200, ErrOutOfRange},
		{60, 31, 0, 200, ErrOutOfRange}, {60, 12, -1, 200, ErrOutOfRange}, {60, 12, 101, 200, ErrOutOfRange},
		{60, 12, 0, -1, ErrOutOfRange}, {60, 12, 0, 5001, ErrOutOfRange},
	}
	for _, c := range cases {
		if _, err := MYSECPMScore(c.age, c.hb, c.bl, false, c.plt, false); !errors.Is(err, c.want) {
			t.Errorf("MYSECPMScore(%+v) error = %v, want %v", c, err, c.want)
		}
	}
}

func TestMYSECPMScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = MYSECPMScore(60, 10, 4, true, 140, true) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkMYSECPMScore(b *testing.B) {
	b.ReportAllocs()
	age, hb, bl, plt := 60.0, 10.0, 4.0, 140.0
	for b.Loop() {
		sinkFloat, _ = MYSECPMScore(age, hb, bl, true, plt, true)
	}
}
