package anthropometry

import (
	"math"
	"testing"
)

// Worked examples, HAND-DERIVED from Schwartz et al. 1978 (abstract: "normalized body mass plus 40% of the
// adipose mass"): ideal + 0.4 × (actual − ideal).
// ideal 70, actual 120 → 70 + 0.4 × 50 = 90. ideal 50, actual 100 → 50 + 20 = 70. ideal 60, actual 61 → 60.4.
// Tolerance 1e-15 relative: one subtract, one multiply, one add.
func TestAdjustedBodyWeight_example(t *testing.T) {
	tests := []struct{ ideal, actual, want float64 }{
		{70, 120, 90},
		{50, 100, 70},
		{60, 61, 60.4},
	}
	for _, tc := range tests {
		if got, err := AdjustedBodyWeight(tc.ideal, tc.actual); err != nil || !relNear(got, tc.want, 1e-15) {
			t.Errorf("AdjustedBodyWeight(%v, %v) = (%v, %v), want %v", tc.ideal, tc.actual, got, err, tc.want)
		}
	}
}

func TestAdjustedBodyWeight_range(t *testing.T) {
	// Just inside the actual > ideal edge.
	if got, err := AdjustedBodyWeight(70, math.Nextafter(70, 100)); err != nil || got != 70 && got != math.Nextafter(70, 100) {
		t.Errorf("actual one ulp above ideal: got (%v, %v), want 70 or the next float", got, err)
	}
	// Largest finite actual: the result lies between the inputs, so it stays finite.
	if got, err := AdjustedBodyWeight(1, math.MaxFloat64); err != nil || math.IsInf(got, 0) {
		t.Errorf("largest actual: got (%v, %v), want finite", got, err)
	}
}

func TestAdjustedBodyWeight_errors(t *testing.T) {
	call := func(ideal, actual float64) func() error {
		return func() error { _, err := AdjustedBodyWeight(ideal, actual); return err }
	}
	runErrCases(t, []errCase{
		{"ideal NaN", call(nan, 120), ErrNotFinite},
		{"ideal +Inf", call(pInf, 120), ErrNotFinite},
		{"ideal -Inf", call(nInf, 120), ErrNotFinite},
		{"actual NaN", call(70, nan), ErrNotFinite},
		{"actual +Inf", call(70, pInf), ErrNotFinite},
		{"actual -Inf", call(70, nInf), ErrNotFinite},
		{"ideal 0", call(0, 120), ErrOutOfRange},
		{"ideal negative", call(-70, 120), ErrOutOfRange},
		{"actual equals ideal", call(70, 70), ErrOutOfRange},
		{"actual below ideal", call(70, 60), ErrOutOfRange},
	})
}

func TestAdjustedBodyWeight_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = AdjustedBodyWeight(70, 120) }); n != 0 {
		t.Errorf("AdjustedBodyWeight allocates %v times, want 0", n)
	}
}

func BenchmarkAdjustedBodyWeight(b *testing.B) {
	b.ReportAllocs()
	ideal, actual := 70.0, 120.0
	for b.Loop() {
		sink, _ = AdjustedBodyWeight(ideal, actual)
	}
}
