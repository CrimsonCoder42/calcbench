package anthropometry

import (
	"math"
	"testing"
)

// Worked examples, HAND-DERIVED from the formula in the Traub & Kichen 1983 abstract, written out here:
// IBM (kg) = 2.396 × e^(0.01863 × height cm). 100 cm → 2.396 × e^1.863. Tolerance 1e-14 relative: one exp,
// two multiplies.
func TestIdealBodyMassTraubKichen_example(t *testing.T) {
	tests := []struct{ heightCm, ageYears, want float64 }{
		{100, 4, 2.396 * math.Exp(1.863)},
		{150, 12, 2.396 * math.Exp(0.01863*150)},
		{75, 1, 2.396 * math.Exp(0.01863*75)},
	}
	for _, tc := range tests {
		if got, err := IdealBodyMassTraubKichen(tc.heightCm, tc.ageYears); err != nil || !relNear(got, tc.want, 1e-14) {
			t.Errorf("IdealBodyMassTraubKichen(%v, %v) = (%v, %v), want %v", tc.heightCm, tc.ageYears, got, err, tc.want)
		}
	}
}

// Age edges: 1 and 17 inside, 0.9999 and 17.0001 outside. The age does not change the result.
func TestIdealBodyMassTraubKichen_range(t *testing.T) {
	want := 2.396 * math.Exp(1.863)
	for _, age := range []float64{1, 17} {
		if got, err := IdealBodyMassTraubKichen(100, age); err != nil || !relNear(got, want, 1e-14) {
			t.Errorf("age %v: got (%v, %v), want %v", age, got, err, want)
		}
	}
	// Just inside the heightCm > 0 edge: e^0 → 2.396.
	if got, err := IdealBodyMassTraubKichen(5e-324, 5); err != nil || got != 2.396 {
		t.Errorf("smallest height: got (%v, %v), want 2.396", got, err)
	}
}

func TestIdealBodyMassTraubKichen_errors(t *testing.T) {
	call := func(h, age float64) func() error {
		return func() error { _, err := IdealBodyMassTraubKichen(h, age); return err }
	}
	runErrCases(t, []errCase{
		{"height NaN", call(nan, 5), ErrNotFinite},
		{"height +Inf", call(pInf, 5), ErrNotFinite},
		{"height -Inf", call(nInf, 5), ErrNotFinite},
		{"age NaN", call(100, nan), ErrNotFinite},
		{"age +Inf", call(100, pInf), ErrNotFinite},
		{"age -Inf", call(100, nInf), ErrNotFinite},
		{"height 0", call(0, 5), ErrOutOfRange},
		{"height negative", call(-100, 5), ErrOutOfRange},
		{"age just below 1", call(100, 0.9999), ErrOutOfRange},
		{"age just above 17", call(100, 17.0001), ErrOutOfRange},
		// 0.01863 × 40000 = 745.2 > ln(MaxFloat64) ≈ 709.78.
		{"result overflows", call(40000, 5), ErrDomain},
	})
}

func TestIdealBodyMassTraubKichen_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = IdealBodyMassTraubKichen(100, 4) }); n != 0 {
		t.Errorf("IdealBodyMassTraubKichen allocates %v times, want 0", n)
	}
}

func BenchmarkIdealBodyMassTraubKichen(b *testing.B) {
	b.ReportAllocs()
	h, age := 100.0, 4.0
	for b.Loop() {
		sink, _ = IdealBodyMassTraubKichen(h, age)
	}
}
