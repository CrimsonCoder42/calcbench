package anthropometry

import "testing"

// Worked example, HAND-DERIVED from the CDC formula "weight (pounds)/[height (in)]2 x 703" (CDC BMI FAQ, read
// 2026-11-12): 150 lb, 65 in → 65² = 4225, 150 / 4225 × 703 = 24.958579881656806.
// Tolerance 1e-15 relative: three correctly rounded operations.
func TestBodyMassIndexImperial_example(t *testing.T) {
	tests := []struct{ weightLb, heightIn, want float64 }{
		{150, 65, 24.958579881656806},
		{150, 65, 150.0 / (65 * 65) * 703},
		{200, 70, 200.0 / (70 * 70) * 703},
	}
	for _, tc := range tests {
		if got, err := BodyMassIndexImperial(tc.weightLb, tc.heightIn); err != nil || !relNear(got, tc.want, 1e-15) {
			t.Errorf("BodyMassIndexImperial(%v, %v) = (%v, %v), want %v", tc.weightLb, tc.heightIn, got, err, tc.want)
		}
	}
}

// Unit-system twins: the same person in lb/in and kg/m. The test converts with the exact constants
// (1 lb = 0.45359237 kg, 1 in = 0.0254 m; NIST SP 811 (2008) §B.8), written here, not taken from the code.
// The CDC's 703 is the rounding of 0.45359237 / 0.0254² = 703.0695796…, so the twins differ by
// 0.0695796 / 703.0696 ≈ 9.9e-5 relative; tolerance 1e-3 (0.1%, the spec's bound) and a check that the
// difference is the factor ratio itself, to 1e-12.
func TestBodyMassIndexImperial_agreesWithMetric(t *testing.T) {
	const kgPerLb, mPerIn = 0.45359237, 0.0254
	exactFactor := kgPerLb / (mPerIn * mPerIn)
	for _, p := range []struct{ lb, in float64 }{{150, 65}, {100, 60}, {250, 74}, {8, 20}, {400, 80}} {
		imp, err1 := BodyMassIndexImperial(p.lb, p.in)
		met, err2 := BodyMassIndex(p.lb*kgPerLb, p.in*mPerIn)
		if err1 != nil || err2 != nil || !relNear(imp, met, 1e-3) {
			t.Errorf("%v lb, %v in: imperial (%v, %v), metric (%v, %v), want within 0.1%%", p.lb, p.in, imp, err1, met, err2)
		}
		if !relNear(imp/met, 703/exactFactor, 1e-12) {
			t.Errorf("%v lb, %v in: ratio %v, want 703/%v", p.lb, p.in, imp/met, exactFactor)
		}
	}
}

func TestBodyMassIndexImperial_errors(t *testing.T) {
	call := func(w, h float64) func() error {
		return func() error { _, err := BodyMassIndexImperial(w, h); return err }
	}
	runErrCases(t, []errCase{
		{"weight NaN", call(nan, 65), ErrNotFinite},
		{"weight +Inf", call(pInf, 65), ErrNotFinite},
		{"weight -Inf", call(nInf, 65), ErrNotFinite},
		{"height NaN", call(150, nan), ErrNotFinite},
		{"height +Inf", call(150, pInf), ErrNotFinite},
		{"height -Inf", call(150, nInf), ErrNotFinite},
		{"weight 0", call(0, 65), ErrOutOfRange},
		{"weight negative", call(-1, 65), ErrOutOfRange},
		{"height 0", call(150, 0), ErrOutOfRange},
		{"height negative", call(150, -65), ErrOutOfRange},
		{"height square underflows", call(150, 1e-200), ErrDomain},
		// 1e308 / 1 × 703 overflows.
		{"product overflows", call(1e308, 1), ErrDomain},
	})
	// Just inside: the smallest positive weight gives a finite, positive result.
	if got, err := BodyMassIndexImperial(5e-324, 1); err != nil || !(got > 0) {
		t.Errorf("smallest weight: got (%v, %v), want > 0", got, err)
	}
}

func TestBodyMassIndexImperial_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = BodyMassIndexImperial(150, 65) }); n != 0 {
		t.Errorf("BodyMassIndexImperial allocates %v times, want 0", n)
	}
}

func BenchmarkBodyMassIndexImperial(b *testing.B) {
	b.ReportAllocs()
	w, h := 150.0, 65.0
	for b.Loop() {
		sink, _ = BodyMassIndexImperial(w, h)
	}
}
