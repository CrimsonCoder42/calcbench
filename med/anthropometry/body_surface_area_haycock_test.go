package anthropometry

import (
	"math"
	"testing"
)

// Worked examples, HAND-DERIVED from the formula in the Haycock 1978 abstract, written out here:
// SA (m²) = weight (kg)^0.5378 × height (cm)^0.3964 × 0.024265. 1 kg and 1 cm give the constant itself.
// Tolerance 1e-14 relative: two pow calls and two multiplies, in any order.
func TestBodySurfaceAreaHaycock_example(t *testing.T) {
	tests := []struct{ weightKg, heightCm, want float64 }{
		{70, 175, math.Pow(70, 0.5378) * math.Pow(175, 0.3964) * 0.024265},
		{3.5, 50, math.Pow(3.5, 0.5378) * math.Pow(50, 0.3964) * 0.024265},
		{1, 1, 0.024265},
	}
	for _, tc := range tests {
		if got, err := BodySurfaceAreaHaycock(tc.weightKg, tc.heightCm); err != nil || !relNear(got, tc.want, 1e-14) {
			t.Errorf("BodySurfaceAreaHaycock(%v, %v) = (%v, %v), want %v", tc.weightKg, tc.heightCm, got, err, tc.want)
		}
	}
	// The source reports measured SA from less than 0.2 m² to greater than 2.0 m²; an adult of 70 kg, 175 cm
	// falls inside that span (a sanity check on the exponents, not a source value).
	if got, _ := BodySurfaceAreaHaycock(70, 175); got < 0.2 || got > 2.5 {
		t.Errorf("70 kg, 175 cm: %v m², want within the source's reported span", got)
	}
}

func TestBodySurfaceAreaHaycock_range(t *testing.T) {
	// Smallest positive inputs: finite and positive (about 2e-304). Largest: finite (about 1e288).
	if got, err := BodySurfaceAreaHaycock(5e-324, 5e-324); err != nil || !(got > 0) {
		t.Errorf("smallest inputs: got (%v, %v), want > 0", got, err)
	}
	if got, err := BodySurfaceAreaHaycock(math.MaxFloat64, math.MaxFloat64); err != nil || math.IsInf(got, 0) {
		t.Errorf("largest inputs: got (%v, %v), want finite", got, err)
	}
}

func TestBodySurfaceAreaHaycock_errors(t *testing.T) {
	call := func(w, h float64) func() error {
		return func() error { _, err := BodySurfaceAreaHaycock(w, h); return err }
	}
	runErrCases(t, []errCase{
		{"weight NaN", call(nan, 175), ErrNotFinite},
		{"weight +Inf", call(pInf, 175), ErrNotFinite},
		{"weight -Inf", call(nInf, 175), ErrNotFinite},
		{"height NaN", call(70, nan), ErrNotFinite},
		{"height +Inf", call(70, pInf), ErrNotFinite},
		{"height -Inf", call(70, nInf), ErrNotFinite},
		{"weight 0", call(0, 175), ErrOutOfRange},
		{"weight negative", call(-70, 175), ErrOutOfRange},
		{"height 0", call(70, 0), ErrOutOfRange},
		{"height negative", call(70, -175), ErrOutOfRange},
	})
}

func TestBodySurfaceAreaHaycock_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = BodySurfaceAreaHaycock(70, 175) }); n != 0 {
		t.Errorf("BodySurfaceAreaHaycock allocates %v times, want 0", n)
	}
}

func BenchmarkBodySurfaceAreaHaycock(b *testing.B) {
	b.ReportAllocs()
	w, h := 70.0, 175.0
	for b.Loop() {
		sink, _ = BodySurfaceAreaHaycock(w, h)
	}
}
