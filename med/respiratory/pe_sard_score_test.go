package respiratory

import (
	"errors"
	"math"
	"testing"
)

var sinkPESARD PESARDRisk

// Chopard 2021 abstract (PMID 34217683): syncope +1.5, Hb < 12 g/dL +2.5, GFR < 60 mL/min +1. Every combination;
// sums of these halves are exact, so the comparison is exact.
func TestPESARDScore_allCombinations(t *testing.T) {
	cases := []struct {
		syncope  bool
		hb, gfr  float64
		want     float64
		wantRisk PESARDRisk
	}{
		{false, 13, 90, 0, PESARDLow},
		{false, 13, 59, 1, PESARDIntermediate},
		{true, 13, 90, 1.5, PESARDIntermediate},
		{false, 11, 90, 2.5, PESARDIntermediate},
		{true, 13, 59, 2.5, PESARDIntermediate},
		{false, 11, 59, 3.5, PESARDHigh},
		{true, 11, 90, 4, PESARDHigh},
		{true, 11, 59, 5, PESARDHigh},
	}
	for _, c := range cases {
		got, err := PESARDScore(c.syncope, c.hb, c.gfr)
		if err != nil || got != c.want {
			t.Errorf("%v: got (%v, %v), want %v", c, got, err, c.want)
			continue
		}
		if r, err := PESARDRiskGroup(got); err != nil || r != c.wantRisk {
			t.Errorf("%v: group (%q, %v), want %q", got, r, err, c.wantRisk)
		}
	}
}

// Edges: Hb 12 does not score, 11.99 does; GFR 60 does not score, 59.99 does.
func TestPESARDScore_edges(t *testing.T) {
	cases := []struct{ hb, gfr, want float64 }{{12, 90, 0}, {11.99, 90, 2.5}, {13, 60, 0}, {13, 59.99, 1}}
	for _, c := range cases {
		if got, err := PESARDScore(false, c.hb, c.gfr); err != nil || got != c.want {
			t.Errorf("hb %v gfr %v: got (%v, %v), want %v", c.hb, c.gfr, got, err, c.want)
		}
	}
}

func TestPESARDScore_errors(t *testing.T) {
	tests := []struct {
		name    string
		hb, gfr float64
		wantErr error
	}{
		{"zero Hb", 0, 90, ErrOutOfRange},
		{"negative GFR", 13, -1, ErrOutOfRange},
		{"NaN Hb", math.NaN(), 90, ErrNotFinite},
		{"+Inf GFR", 13, math.Inf(1), ErrNotFinite},
		{"-Inf Hb", math.Inf(-1), 90, ErrNotFinite},
	}
	for _, tc := range tests {
		if _, err := PESARDScore(false, tc.hb, tc.gfr); !errors.Is(err, tc.wantErr) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.wantErr)
		}
	}
	for _, bad := range []struct {
		s       float64
		wantErr error
	}{{0.5, ErrOutOfRange}, {-1, ErrOutOfRange}, {6, ErrOutOfRange}, {3, ErrOutOfRange}, {math.NaN(), ErrNotFinite}, {math.Inf(1), ErrNotFinite}} {
		if _, err := PESARDRiskGroup(bad.s); !errors.Is(err, bad.wantErr) {
			t.Errorf("group %v: err = %v, want %v", bad.s, err, bad.wantErr)
		}
	}
}

func TestPESARDScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = PESARDScore(true, 11, 59) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkPESARD, _ = PESARDRiskGroup(2.5) }); n != 0 {
		t.Errorf("group allocs = %v, want 0", n)
	}
}

func BenchmarkPESARDScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = PESARDScore(true, 11, 59)
	}
}

func BenchmarkPESARDRiskGroup(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkPESARD, _ = PESARDRiskGroup(2.5)
	}
}
