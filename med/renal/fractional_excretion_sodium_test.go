package renal

import "testing"

// Worked examples, HAND-DERIVED from the FENa definition (Espinel 1976):
// UNa 20, PNa 140, UCr 100, PCr 2 → (20 × 2) / (140 × 100) × 100 = 40/14000 × 100 = 0.2857142857142857.
// UNa 40, PNa 130, UCr 50, PCr 5 → 200/6500 × 100 = 3.076923076923077.
// Tolerance 1e-14 relative (four correctly rounded operations).
func TestFractionalExcretionSodium_example(t *testing.T) {
	tests := []struct{ uNa, pNa, uCr, pCr, want float64 }{
		{20, 140, 100, 2, 2.0 / 7},
		{40, 130, 50, 5, 200.0 / 65},
		// Creatinine in µmol/L instead (both × 88.4): the ratio cancels the unit.
		{20, 140, 8840, 176.8, 2.0 / 7},
	}
	for _, tc := range tests {
		if got, err := FractionalExcretionSodium(tc.uNa, tc.pNa, tc.uCr, tc.pCr); err != nil || !relNear(got, tc.want, 1e-14) {
			t.Errorf("FractionalExcretionSodium(%v, %v, %v, %v) = (%v, %v), want %v", tc.uNa, tc.pNa, tc.uCr, tc.pCr, got, err, tc.want)
		}
	}
}

// Extreme in-range inputs: products that would overflow (1e200 × 1e200) are avoided by the ratio form;
// 1e200 × 1e200 / (1e200 × 1e200) × 100 = 100.
func TestFractionalExcretionSodium_extremes(t *testing.T) {
	if got, err := FractionalExcretionSodium(1e200, 1e200, 1e200, 1e200); err != nil || !relNear(got, 100, 1e-14) {
		t.Errorf("all 1e200: got (%v, %v), want 100", got, err)
	}
	if got, err := FractionalExcretionSodium(1e-200, 1e-200, 1e-200, 1e-200); err != nil || !relNear(got, 100, 1e-14) {
		t.Errorf("all 1e-200: got (%v, %v), want 100", got, err)
	}
}

func TestFractionalExcretionSodium_errors(t *testing.T) {
	call := func(uNa, pNa, uCr, pCr float64) func() error {
		return func() error { _, err := FractionalExcretionSodium(uNa, pNa, uCr, pCr); return err }
	}
	runErrCases(t, []errCase{
		{"UNa NaN", call(nan, 140, 100, 2), ErrNotFinite},
		{"UNa +Inf", call(pInf, 140, 100, 2), ErrNotFinite},
		{"PNa -Inf", call(20, nInf, 100, 2), ErrNotFinite},
		{"PNa NaN", call(20, nan, 100, 2), ErrNotFinite},
		{"UCr +Inf", call(20, 140, pInf, 2), ErrNotFinite},
		{"UCr NaN", call(20, 140, nan, 2), ErrNotFinite},
		{"PCr -Inf", call(20, 140, 100, nInf), ErrNotFinite},
		{"PCr NaN", call(20, 140, 100, nan), ErrNotFinite},
		{"NaN before range", call(0, 140, 100, nan), ErrNotFinite},
		{"UNa 0", call(0, 140, 100, 2), ErrOutOfRange},
		{"PNa 0", call(20, 0, 100, 2), ErrOutOfRange},
		{"UCr 0", call(20, 140, 0, 2), ErrOutOfRange},
		{"PCr 0", call(20, 140, 100, 0), ErrOutOfRange},
		{"PCr negative", call(20, 140, 100, -2), ErrOutOfRange},
		{"result beyond float64", call(1e300, 1e-300, 1e-300, 1e300), ErrDomain},
	})
}

func TestFractionalExcretionSodium_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = FractionalExcretionSodium(20, 140, 100, 2) }); n != 0 {
		t.Errorf("FractionalExcretionSodium allocates %v times, want 0", n)
	}
}

func BenchmarkFractionalExcretionSodium(b *testing.B) {
	b.ReportAllocs()
	uNa, pNa, uCr, pCr := 20.0, 140.0, 100.0, 2.0
	for b.Loop() {
		sink, _ = FractionalExcretionSodium(uNa, pNa, uCr, pCr)
	}
}
