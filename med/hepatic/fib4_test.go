package hepatic

import (
	"errors"
	"math"
	"testing"
)

// Shared by the package's tests.

// sink and sinkInt keep benchmark results observable.
var (
	sink    float64
	sinkInt int
)

// relNear reports |got − want| ≤ tol·|want|.
func relNear(got, want, tol float64) bool { return math.Abs(got-want) <= tol*math.Abs(want) }

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

// Worked examples, HAND-DERIVED from Sterling et al. 2006 (abstract formula):
// age 50, AST 40, ALT 25, PLT 200 → 50 × 40 / (200 × 5) = 2.0.
// age 64, AST 90, ALT 36, PLT 120 → 5760 / (120 × 6) = 8.
// age 30, AST 20, ALT 16, PLT 300 → 600 / 1200 = 0.5.
// Tolerance 1e-14 relative (exact square roots, three correctly rounded operations).
func TestFIB4_example(t *testing.T) {
	tests := []struct{ age, ast, alt, plt, want float64 }{
		{50, 40, 25, 200, 2.0},
		{64, 90, 36, 120, 8},
		{30, 20, 16, 300, 0.5},
	}
	for _, tc := range tests {
		if got, err := FIB4(tc.age, tc.ast, tc.alt, tc.plt); err != nil || !relNear(got, tc.want, 1e-14) {
			t.Errorf("FIB4(%v, %v, %v, %v) = (%v, %v), want %v", tc.age, tc.ast, tc.alt, tc.plt, got, err, tc.want)
		}
	}
}

// Extreme in-range inputs whose products alone overflow or underflow: the ratio form keeps the digits.
// 1e300 × 1e300 / (1e300 × √1) = 1e300; 1e-300 × 1e-300 / (1e-300 × √1) = 1e-300.
func TestFIB4_extremes(t *testing.T) {
	if got, err := FIB4(1e300, 1e300, 1, 1e300); err != nil || !relNear(got, 1e300, 1e-14) {
		t.Errorf("large: got (%v, %v), want 1e300", got, err)
	}
	if got, err := FIB4(1e-300, 1e-300, 1, 1e-300); err != nil || !relNear(got, 1e-300, 1e-14) {
		t.Errorf("small: got (%v, %v), want 1e-300", got, err)
	}
}

func TestFIB4_errors(t *testing.T) {
	call := func(age, ast, alt, plt float64) func() error {
		return func() error { _, err := FIB4(age, ast, alt, plt); return err }
	}
	runErrCases(t, []errCase{
		{"age NaN", call(nan, 40, 25, 200), ErrNotFinite},
		{"age +Inf", call(pInf, 40, 25, 200), ErrNotFinite},
		{"age -Inf", call(nInf, 40, 25, 200), ErrNotFinite},
		{"AST NaN", call(50, nan, 25, 200), ErrNotFinite},
		{"AST +Inf", call(50, pInf, 25, 200), ErrNotFinite},
		{"AST -Inf", call(50, nInf, 25, 200), ErrNotFinite},
		{"ALT NaN", call(50, 40, nan, 200), ErrNotFinite},
		{"ALT +Inf", call(50, 40, pInf, 200), ErrNotFinite},
		{"ALT -Inf", call(50, 40, nInf, 200), ErrNotFinite},
		{"PLT NaN", call(50, 40, 25, nan), ErrNotFinite},
		{"PLT +Inf", call(50, 40, 25, pInf), ErrNotFinite},
		{"PLT -Inf", call(50, 40, 25, nInf), ErrNotFinite},
		{"NaN before range", call(0, 40, 25, nan), ErrNotFinite},
		{"age 0", call(0, 40, 25, 200), ErrOutOfRange},
		{"AST 0", call(50, 0, 25, 200), ErrOutOfRange},
		{"ALT 0", call(50, 40, 0, 200), ErrOutOfRange},
		{"ALT negative", call(50, 40, -25, 200), ErrOutOfRange},
		{"PLT 0", call(50, 40, 25, 0), ErrOutOfRange},
		{"result beyond float64", call(1e308, 1e308, 1, 1e-10), ErrDomain},
	})
	if _, err := FIB4(1e-9, 40, 25, 200); err != nil {
		t.Errorf("age just above 0: err = %v, want nil", err)
	}
}

func TestFIB4_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = FIB4(50, 40, 25, 200) }); n != 0 {
		t.Errorf("FIB4 allocates %v times, want 0", n)
	}
}

func BenchmarkFIB4(b *testing.B) {
	b.ReportAllocs()
	age, ast, alt, plt := 50.0, 40.0, 25.0, 200.0
	for b.Loop() {
		sink, _ = FIB4(age, ast, alt, plt)
	}
}
