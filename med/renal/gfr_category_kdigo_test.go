package renal

import (
	"errors"
	"math"
	"testing"
)

// KDIGO 2024 p. S125, GFR categories: G1 ≥ 90 "Normal or high", G2 60–89 "Mildly decreased",
// G3a 45–59 "Mildly to moderately decreased", G3b 30–44 "Moderately to severely decreased",
// G4 15–29 "Severely decreased", G5 < 15 "Kidney failure". Both sides of every edge (ASSUMPTION A-M3:
// the integer table read as ≥ the lower printed bound).
func TestGFRCategoryKDIGO_edges(t *testing.T) {
	tests := []struct {
		egfr float64
		want GFRCategory
		term string
	}{
		{1e9, G1, "Normal or high"},
		{90, G1, "Normal or high"},
		{89.999, G2, "Mildly decreased"},
		{60, G2, "Mildly decreased"},
		{59.999, G3a, "Mildly to moderately decreased"},
		{45, G3a, "Mildly to moderately decreased"},
		{44.999, G3b, "Moderately to severely decreased"},
		{30, G3b, "Moderately to severely decreased"},
		{29.999, G4, "Severely decreased"},
		{15, G4, "Severely decreased"},
		{14.999, G5, "Kidney failure"},
		{0, G5, "Kidney failure"},
		{math.Copysign(0, -1), G5, "Kidney failure"},
	}
	for _, tc := range tests {
		got, err := GFRCategoryKDIGO(tc.egfr)
		if err != nil || got != tc.want || got.Term() != tc.term {
			t.Errorf("GFRCategoryKDIGO(%v) = (%q %q, %v), want %q %q", tc.egfr, got, got.Term(), err, tc.want, tc.term)
		}
	}
}

func TestGFRCategory_Term_unknown(t *testing.T) {
	if got := GFRCategory("G6").Term(); got != "" {
		t.Errorf(`GFRCategory("G6").Term() = %q, want ""`, got)
	}
}

func TestGFRCategoryKDIGO_errors(t *testing.T) {
	for _, v := range []float64{nan, pInf, nInf} {
		if _, err := GFRCategoryKDIGO(v); !errors.Is(err, ErrNotFinite) {
			t.Errorf("GFRCategoryKDIGO(%v): err = %v, want ErrNotFinite", v, err)
		}
	}
	for _, v := range []float64{-0.001, -90} {
		if _, err := GFRCategoryKDIGO(v); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("GFRCategoryKDIGO(%v): err = %v, want ErrOutOfRange", v, err)
		}
	}
}

func TestGFRCategoryKDIGO_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkCategory, _ = GFRCategoryKDIGO(52) }); n != 0 {
		t.Errorf("GFRCategoryKDIGO allocates %v times, want 0", n)
	}
}

func BenchmarkGFRCategoryKDIGO(b *testing.B) {
	b.ReportAllocs()
	egfr := 52.0
	for b.Loop() {
		sinkCategory, _ = GFRCategoryKDIGO(egfr)
	}
}
