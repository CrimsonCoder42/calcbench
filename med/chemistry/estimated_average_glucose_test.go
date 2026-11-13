package chemistry

import (
	"errors"
	"math"
	"testing"
)

// Nathan et al. 2008, the eAG table (whole A1C 5–12, mg/dL printed as integers): 5→97, 6→126, 7→154, 8→183,
// 9→212, 10→240, 11→269, 12→298. The equation's exact values (28.7 × A1C − 46.7) are 96.8, 125.5, 154.2,
// 182.9, 211.6, 240.3, 269.0, 297.7; math.Round of each gives the printed integer (125.5 rounds up).
// Exact values compared to 1e-13 relative (a few correctly rounded operations).
func TestEstimatedAverageGlucose_table(t *testing.T) {
	tests := []struct{ a1c, exact, printed float64 }{
		{5, 96.8, 97}, {6, 125.5, 126}, {7, 154.2, 154}, {8, 182.9, 183},
		{9, 211.6, 212}, {10, 240.3, 240}, {11, 269.0, 269}, {12, 297.7, 298},
	}
	for _, tc := range tests {
		got, err := EstimatedAverageGlucose(tc.a1c)
		if err != nil || !relNear(got, tc.exact, 1e-13) || math.Round(got) != tc.printed {
			t.Errorf("EstimatedAverageGlucose(%v) = (%v, %v), want %v (table %v)", tc.a1c, got, err, tc.exact, tc.printed)
		}
	}
}

// Nathan et al. 2008, the eAG table in mmol/L (one decimal): 5→5.4, 6→7.0, 7→8.6, 8→10.2, 9→11.8, 10→13.4,
// 11→14.9, 12→16.5. The equation's exact values (1.5944 × A1C − 2.594) are 5.378, 6.9724, 8.5668, 10.1612,
// 11.7556, 13.35, 14.9444, 16.5388; each is within 0.05 of the printed value (13.35 is the tie the table
// prints as 13.4). Exact values compared to 1e-13 relative.
func TestEstimatedAverageGlucoseMmolL_table(t *testing.T) {
	tests := []struct{ a1c, exact, printed float64 }{
		{5, 5.378, 5.4}, {6, 6.9724, 7.0}, {7, 8.5668, 8.6}, {8, 10.1612, 10.2},
		{9, 11.7556, 11.8}, {10, 13.35, 13.4}, {11, 14.9444, 14.9}, {12, 16.5388, 16.5},
	}
	for _, tc := range tests {
		got, err := EstimatedAverageGlucoseMmolL(tc.a1c)
		if err != nil || !relNear(got, tc.exact, 1e-13) || !near(got, tc.printed, 0.05+1e-12) {
			t.Errorf("EstimatedAverageGlucoseMmolL(%v) = (%v, %v), want %v (table %v)", tc.a1c, got, err, tc.exact, tc.printed)
		}
	}
}

// The two forms are separate regressions printed by the source, not a unit conversion of each other; they
// agree to the table's precision: mg/dL ÷ 18 is within 0.1 mmol/L of the mmol/L form over 4–12 (HAND-DERIVED:
// at 4, 68.1/18 = 3.783 vs 3.7836; at 12, 297.7/18 = 16.539 vs 16.5388).
func TestEstimatedAverageGlucose_formsAgree(t *testing.T) {
	for a1c := 4.0; a1c <= 12; a1c += 0.5 {
		mg, _ := EstimatedAverageGlucose(a1c)
		mmol, _ := EstimatedAverageGlucoseMmolL(a1c)
		if !near(mg/18, mmol, 0.1) {
			t.Errorf("A1C %v: %v mg/dL ÷ 18 = %v, mmol/L form %v", a1c, mg, mg/18, mmol)
		}
	}
}

func TestEstimatedAverageGlucose_range(t *testing.T) {
	for _, f := range []func(float64) (float64, error){EstimatedAverageGlucose, EstimatedAverageGlucoseMmolL} {
		for _, a1c := range []float64{4, 12} {
			if _, err := f(a1c); err != nil {
				t.Errorf("A1C %v: err = %v, want nil", a1c, err)
			}
		}
		for _, a1c := range []float64{3.999, 12.001, 0, -5} {
			if _, err := f(a1c); !errors.Is(err, ErrOutOfRange) {
				t.Errorf("A1C %v: err = %v, want ErrOutOfRange", a1c, err)
			}
		}
		for _, a1c := range []float64{nan, pInf, nInf} {
			if _, err := f(a1c); !errors.Is(err, ErrNotFinite) {
				t.Errorf("A1C %v: err = %v, want ErrNotFinite", a1c, err)
			}
		}
	}
	// HAND-DERIVED lower edge: 28.7 × 4 − 46.7 = 68.1; 1.5944 × 4 − 2.594 = 3.7836.
	if got, _ := EstimatedAverageGlucose(4); !relNear(got, 68.1, 1e-13) {
		t.Errorf("EstimatedAverageGlucose(4) = %v, want 68.1", got)
	}
	if got, _ := EstimatedAverageGlucoseMmolL(4); !relNear(got, 3.7836, 1e-13) {
		t.Errorf("EstimatedAverageGlucoseMmolL(4) = %v, want 3.7836", got)
	}
}

func TestEstimatedAverageGlucose_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = EstimatedAverageGlucose(7) }); n != 0 {
		t.Errorf("EstimatedAverageGlucose allocates %v times, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sink, _ = EstimatedAverageGlucoseMmolL(7) }); n != 0 {
		t.Errorf("EstimatedAverageGlucoseMmolL allocates %v times, want 0", n)
	}
}

func BenchmarkEstimatedAverageGlucose(b *testing.B) {
	b.ReportAllocs()
	a1c := 7.0
	for b.Loop() {
		sink, _ = EstimatedAverageGlucose(a1c)
	}
}

func BenchmarkEstimatedAverageGlucoseMmolL(b *testing.B) {
	b.ReportAllocs()
	a1c := 7.0
	for b.Loop() {
		sink, _ = EstimatedAverageGlucoseMmolL(a1c)
	}
}
