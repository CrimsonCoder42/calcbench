package anthropometry

import "testing"

// Worked examples, HAND-DERIVED from the Tinning & Acworth 2007 abstract's "<12 months" formula,
// (age in months + 9) / 2: 6 → 15 / 2 = 7.5; 0 → 9 / 2 = 4.5; 11 → 20 / 2 = 10; 11.9999 → 20.9999 / 2 = 10.49995.
// Tolerance 1e-15 relative: one add, one exact halving.
func TestEstimatedWeightInfantBestGuess_example(t *testing.T) {
	tests := []struct{ ageMonths, want float64 }{
		{6, 7.5},
		{0, 4.5},
		{11, 10},
		{11.9999, 10.49995},
	}
	for _, tc := range tests {
		if got, err := EstimatedWeightInfantBestGuess(tc.ageMonths); err != nil || !relNear(got, tc.want, 1e-15) {
			t.Errorf("EstimatedWeightInfantBestGuess(%v) = (%v, %v), want %v", tc.ageMonths, got, err, tc.want)
		}
	}
}

func TestEstimatedWeightInfantBestGuess_errors(t *testing.T) {
	call := func(m float64) func() error {
		return func() error { _, err := EstimatedWeightInfantBestGuess(m); return err }
	}
	runErrCases(t, []errCase{
		{"NaN", call(nan), ErrNotFinite},
		{"+Inf", call(pInf), ErrNotFinite},
		{"-Inf", call(nInf), ErrNotFinite},
		{"just below 0", call(-0.0001), ErrOutOfRange},
		{"12 months", call(12), ErrOutOfRange},
		{"above 12", call(30), ErrOutOfRange},
	})
}

func TestEstimatedWeightInfantBestGuess_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = EstimatedWeightInfantBestGuess(6) }); n != 0 {
		t.Errorf("EstimatedWeightInfantBestGuess allocates %v times, want 0", n)
	}
}

func BenchmarkEstimatedWeightInfantBestGuess(b *testing.B) {
	b.ReportAllocs()
	m := 6.0
	for b.Loop() {
		sink, _ = EstimatedWeightInfantBestGuess(m)
	}
}
