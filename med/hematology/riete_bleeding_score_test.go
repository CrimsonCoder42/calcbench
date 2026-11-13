package hematology

import (
	"errors"
	"math"
	"testing"
)

// Ruíz-Giménez 2008 abstract: recent bleeding 2; creatinine > 1.2 mg/dL 1.5; anemia 1.5; age > 75, cancer,
// PE 1 each. Exact comparison: every result is a multiple of 0.5.
func TestRIETEBleedingScore_example(t *testing.T) {
	tests := []struct {
		age           float64
		bleed, cancer bool
		cr            float64
		anemia, pe    bool
		want          float64
	}{
		{75, false, false, 1.2, false, false, 0},
		{75.1, false, false, 1.2, false, false, 1},
		{60, true, false, 1.0, false, false, 2},
		{60, false, false, 1.21, false, false, 1.5},
		{60, false, false, 1.0, true, false, 1.5},
		{60, false, true, 1.0, false, false, 1},
		{60, false, false, 1.0, false, true, 1},
		{60, true, false, 1.3, true, false, 5},
		{76, false, true, 1.0, false, true, 3},
		{80, true, true, 2.0, true, true, 8},
	}
	for _, tc := range tests {
		if got, err := RIETEBleedingScore(tc.age, tc.bleed, tc.cancer, tc.cr, tc.anemia, tc.pe); err != nil || got != tc.want {
			t.Errorf("RIETEBleedingScore(%+v) = (%v, %v), want %v", tc, got, err, tc.want)
		}
	}
}

func TestRIETEBleedingScore_errors(t *testing.T) {
	cases := []struct {
		age, cr float64
		want    error
	}{
		{math.NaN(), 1, ErrNotFinite}, {60, math.Inf(1), ErrNotFinite}, {math.Inf(-1), 1, ErrNotFinite},
		{-1, 1, ErrOutOfRange}, {131, 1, ErrOutOfRange}, {60, -0.1, ErrOutOfRange}, {60, 30.1, ErrOutOfRange},
	}
	for _, c := range cases {
		if _, err := RIETEBleedingScore(c.age, false, false, c.cr, false, false); !errors.Is(err, c.want) {
			t.Errorf("RIETEBleedingScore(age %v, cr %v) error = %v, want %v", c.age, c.cr, err, c.want)
		}
	}
}

func TestRIETEBleedingScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = RIETEBleedingScore(80, true, false, 1.5, true, false) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkRIETEBleedingScore(b *testing.B) {
	b.ReportAllocs()
	age, cr := 80.0, 1.5
	for b.Loop() {
		sinkFloat, _ = RIETEBleedingScore(age, true, false, cr, true, false)
	}
}
