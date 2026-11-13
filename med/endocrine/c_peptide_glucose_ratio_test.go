package endocrine

import (
	"math"
	"testing"
)

// HAND-DERIVED from Saisho 2016 section 4 (C-peptide ng/mL / glucose mg/dL × 100). Tolerance 1e-15 relative
// (one divide, one multiply).
// parity: ledger c_peptide_glucose_ratio, the app's postprandial branch "(cPeptideNgml / glucose) * 100".
// The app's fasting branch (C-peptide pmol/L / glucose mg/dL, Fritsche) is not ported (ledger row).
func TestCPeptideGlucoseRatio_example(t *testing.T) {
	tests := []struct{ cp, glu, want float64 }{
		{3.0, 150, 2.0},
		{1.0, 200, 0.5},
		{4.04, 200, 2.02},
	}
	for _, tc := range tests {
		if got, err := CPeptideGlucoseRatio(tc.cp, tc.glu); err != nil || !relNear(got, tc.want, 1e-15) {
			t.Errorf("CPeptideGlucoseRatio(%v, %v) = (%v, %v), want %v", tc.cp, tc.glu, got, err, tc.want)
		}
	}
}

func TestCPeptideGlucoseRatio_errors(t *testing.T) {
	call := func(c, g float64) func() error {
		return func() error { _, err := CPeptideGlucoseRatio(c, g); return err }
	}
	runErrCases(t, []errCase{
		{"cpep NaN", call(nan, 100), ErrNotFinite},
		{"cpep +Inf", call(pInf, 100), ErrNotFinite},
		{"cpep -Inf", call(nInf, 100), ErrNotFinite},
		{"glucose NaN", call(1, nan), ErrNotFinite},
		{"glucose +Inf", call(1, pInf), ErrNotFinite},
		{"glucose -Inf", call(1, nInf), ErrNotFinite},
		{"cpep 0", call(0, 100), ErrOutOfRange},
		{"glucose 0", call(1, 0), ErrOutOfRange},
		{"glucose negative", call(1, -5), ErrOutOfRange},
		{"overflow", call(math.MaxFloat64, 1e-10), ErrDomain},
	})
}

func TestCPeptideGlucoseRatio_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = CPeptideGlucoseRatio(3, 150) }); n != 0 {
		t.Errorf("CPeptideGlucoseRatio allocates %v times, want 0", n)
	}
}

func BenchmarkCPeptideGlucoseRatio(b *testing.B) {
	b.ReportAllocs()
	c, g := 3.0, 150.0
	for b.Loop() {
		sink, _ = CPeptideGlucoseRatio(c, g)
	}
}
