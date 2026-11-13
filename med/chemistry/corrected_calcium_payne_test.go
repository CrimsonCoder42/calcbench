package chemistry

import (
	"errors"
	"math"
	"testing"
)

// Shared by the package's tests.

// sink keeps benchmark results observable.
var sink float64

// relNear reports |got − want| ≤ tol·|want| (tol·1e-300 when want is 0).
func relNear(got, want, tol float64) bool {
	return math.Abs(got-want) <= tol*math.Max(math.Abs(want), 1e-300)
}

// near reports |got − want| ≤ tol, an absolute check.
func near(got, want, tol float64) bool { return math.Abs(got-want) <= tol }

// errCase is one call that must fail with want.
type errCase struct {
	name string
	call func() error
	want error
}

func runErrCases(t *testing.T, cases []errCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.call(); !errors.Is(err, tc.want) {
				t.Errorf("err = %v, want %v", err, tc.want)
			}
		})
	}
}

var (
	nan  = math.NaN()
	pInf = math.Inf(1)
	nInf = math.Inf(-1)
)

// Worked examples, HAND-DERIVED from Payne et al. 1973 (abstract: calcium − albumin + 4.0):
// 8.0 − 2.0 + 4.0 = 10.0; 9.5 − 4.0 + 4.0 = 9.5; 7.2 − 1.8 + 4.0 = 9.4. Tolerance 1e-14 absolute (two
// correctly rounded additions on values near 10).
func TestCorrectedCalciumPayne_example(t *testing.T) {
	tests := []struct{ ca, alb, want float64 }{
		{8.0, 2.0, 10.0},
		{9.5, 4.0, 9.5},
		{7.2, 1.8, 9.4},
		// No limits in the source: the formula's value is returned even when ≤ 0 (1 − 10 + 4 = −5).
		{1, 10, -5},
	}
	for _, tc := range tests {
		if got, err := CorrectedCalciumPayne(tc.ca, tc.alb); err != nil || !near(got, tc.want, 1e-14) {
			t.Errorf("CorrectedCalciumPayne(%v, %v) = (%v, %v), want %v", tc.ca, tc.alb, got, err, tc.want)
		}
	}
}

func TestCorrectedCalciumPayne_errors(t *testing.T) {
	call := func(ca, alb float64) func() error {
		return func() error { _, err := CorrectedCalciumPayne(ca, alb); return err }
	}
	runErrCases(t, []errCase{
		{"calcium NaN", call(nan, 3), ErrNotFinite},
		{"calcium +Inf", call(pInf, 3), ErrNotFinite},
		{"calcium -Inf", call(nInf, 3), ErrNotFinite},
		{"albumin NaN", call(9, nan), ErrNotFinite},
		{"albumin +Inf", call(9, pInf), ErrNotFinite},
		{"albumin -Inf", call(9, nInf), ErrNotFinite},
		{"NaN before range", call(0, nan), ErrNotFinite},
		{"calcium 0", call(0, 3), ErrOutOfRange},
		{"calcium negative", call(-1, 3), ErrOutOfRange},
		{"albumin 0", call(9, 0), ErrOutOfRange},
		{"albumin negative", call(9, -3), ErrOutOfRange},
	})
	if _, err := CorrectedCalciumPayne(1e-300, 1e-300); err != nil {
		t.Errorf("smallest positive inputs: err = %v, want nil", err)
	}
}

// The SI twin. HAND-DERIVED: 2.0 mmol/L ÷ 0.2495 = 8.016032… mg/dL, 20 g/L = 2.0 g/dL;
// (8.016032… − 2.0 + 4.0) × 0.2495 = 2.0 + 2.0 × 0.2495 = 2.499 mmol/L. Agreement with the mg/dL form at
// calcium × 0.2495 and albumin × 10; tolerance 1e-13 relative (a division and a multiplication by 0.2495).
func TestCorrectedCalciumPayneSI_agrees(t *testing.T) {
	if got, err := CorrectedCalciumPayneSI(2.0, 20); err != nil || !relNear(got, 2.0+2.0*0.2495, 1e-13) {
		t.Errorf("2.0 mmol/L, 20 g/L: got (%v, %v), want 2.499", got, err)
	}
	for _, ca := range []float64{6.5, 8.0, 10.2, 13} {
		for _, alb := range []float64{1.5, 2.8, 4.4} {
			mg, _ := CorrectedCalciumPayne(ca, alb)
			got, err := CorrectedCalciumPayneSI(ca*0.2495, alb*10)
			if err != nil || !relNear(got, mg*0.2495, 1e-13) {
				t.Errorf("Ca %v alb %v: SI (%v, %v), mg/dL × 0.2495 = %v", ca, alb, got, err, mg*0.2495)
			}
		}
	}
}

func TestCorrectedCalciumPayneSI_errors(t *testing.T) {
	call := func(ca, alb float64) func() error {
		return func() error { _, err := CorrectedCalciumPayneSI(ca, alb); return err }
	}
	runErrCases(t, []errCase{
		{"calcium NaN", call(nan, 30), ErrNotFinite},
		{"calcium +Inf", call(pInf, 30), ErrNotFinite},
		{"albumin -Inf", call(2, nInf), ErrNotFinite},
		{"albumin NaN", call(2, nan), ErrNotFinite},
		{"calcium 0", call(0, 30), ErrOutOfRange},
		{"albumin 0", call(2, 0), ErrOutOfRange},
		// 1e308 ÷ 0.2495 is beyond float64.
		{"result beyond float64", call(1e308, 30), ErrDomain},
	})
}

func TestCorrectedCalciumPayne_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = CorrectedCalciumPayne(8, 2) }); n != 0 {
		t.Errorf("CorrectedCalciumPayne allocates %v times, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sink, _ = CorrectedCalciumPayneSI(2, 20) }); n != 0 {
		t.Errorf("CorrectedCalciumPayneSI allocates %v times, want 0", n)
	}
}

func BenchmarkCorrectedCalciumPayne(b *testing.B) {
	b.ReportAllocs()
	ca, alb := 8.0, 2.0
	for b.Loop() {
		sink, _ = CorrectedCalciumPayne(ca, alb)
	}
}

func BenchmarkCorrectedCalciumPayneSI(b *testing.B) {
	b.ReportAllocs()
	ca, alb := 2.0, 20.0
	for b.Loop() {
		sink, _ = CorrectedCalciumPayneSI(ca, alb)
	}
}
