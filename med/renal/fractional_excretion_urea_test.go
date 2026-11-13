package renal

import "testing"

// Worked examples, HAND-DERIVED from the FEUrea definition (Carvounis et al. 2002):
// UUrea 300, PUrea 30, UCr 100, PCr 2 → (300 × 2) / (30 × 100) × 100 = 600/3000 × 100 = 20.
// UUrea 500, PUrea 40, UCr 120, PCr 3 → 1500/4800 × 100 = 31.25.
// Tolerance 1e-14 relative (four correctly rounded operations).
func TestFractionalExcretionUrea_example(t *testing.T) {
	tests := []struct{ uU, pU, uCr, pCr, want float64 }{
		{300, 30, 100, 2, 20},
		{500, 40, 120, 3, 31.25},
		// Creatinine in µmol/L instead (both × 88.4): the ratio cancels the unit.
		{300, 30, 8840, 176.8, 20},
	}
	for _, tc := range tests {
		if got, err := FractionalExcretionUrea(tc.uU, tc.pU, tc.uCr, tc.pCr); err != nil || !relNear(got, tc.want, 1e-14) {
			t.Errorf("FractionalExcretionUrea(%v, %v, %v, %v) = (%v, %v), want %v", tc.uU, tc.pU, tc.uCr, tc.pCr, got, err, tc.want)
		}
	}
}

func TestFractionalExcretionUrea_errors(t *testing.T) {
	call := func(uU, pU, uCr, pCr float64) func() error {
		return func() error { _, err := FractionalExcretionUrea(uU, pU, uCr, pCr); return err }
	}
	runErrCases(t, []errCase{
		{"UUrea NaN", call(nan, 30, 100, 2), ErrNotFinite},
		{"UUrea -Inf", call(nInf, 30, 100, 2), ErrNotFinite},
		{"PUrea +Inf", call(300, pInf, 100, 2), ErrNotFinite},
		{"PUrea NaN", call(300, nan, 100, 2), ErrNotFinite},
		{"UCr NaN", call(300, 30, nan, 2), ErrNotFinite},
		{"UCr -Inf", call(300, 30, nInf, 2), ErrNotFinite},
		{"PCr +Inf", call(300, 30, 100, pInf), ErrNotFinite},
		{"PCr NaN", call(300, 30, 100, nan), ErrNotFinite},
		{"UUrea 0", call(0, 30, 100, 2), ErrOutOfRange},
		{"PUrea 0", call(300, 0, 100, 2), ErrOutOfRange},
		{"UCr negative", call(300, 30, -100, 2), ErrOutOfRange},
		{"PCr 0", call(300, 30, 100, 0), ErrOutOfRange},
		{"result beyond float64", call(1e300, 1e-300, 1e-300, 1e300), ErrDomain},
	})
}

func TestFractionalExcretionUrea_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = FractionalExcretionUrea(300, 30, 100, 2) }); n != 0 {
		t.Errorf("FractionalExcretionUrea allocates %v times, want 0", n)
	}
}

func BenchmarkFractionalExcretionUrea(b *testing.B) {
	b.ReportAllocs()
	uU, pU, uCr, pCr := 300.0, 30.0, 100.0, 2.0
	for b.Loop() {
		sink, _ = FractionalExcretionUrea(uU, pU, uCr, pCr)
	}
}
