package endocrine

import (
	"errors"
	"math"
	"testing"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Shared by the package's tests.

// sink, sinkInt and sinkBool keep benchmark results observable.
var (
	sink     float64
	sinkInt  int
	sinkBool bool
)

// relNear reports |got − want| ≤ tol·|want| (tol·1e-300 when want is 0).
func relNear(got, want, tol float64) bool {
	return math.Abs(got-want) <= tol*math.Max(math.Abs(want), 1e-300)
}

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

// Bergenstal 2018, Table 1 (PMC6196826): mean glucose 150, 175, 200, 250 mg/dL → GMI 6.9, 7.5, 8.1, 9.3 %
// as printed (one decimal). Exact values 6.898, 7.496, 8.094, 9.29; tolerance 0.05 (the table's rounding) and
// 1e-12 relative against the exact value.
func TestGlucoseManagementIndicatorPercent_table1(t *testing.T) {
	tests := []struct{ mean, printed, exact float64 }{
		{150, 6.9, 6.898},
		{175, 7.5, 7.496},
		{200, 8.1, 8.094},
		{250, 9.3, 9.29},
	}
	for _, tc := range tests {
		got, err := GlucoseManagementIndicatorPercent(tc.mean)
		if err != nil || math.Abs(got-tc.printed) > 0.05 || !relNear(got, tc.exact, 1e-12) {
			t.Errorf("GMI%%(%v) = (%v, %v), want %v (printed %v)", tc.mean, got, err, tc.exact, tc.printed)
		}
	}
}

// parity: ledger bs_gmi, JSON documented formula "GMI (%) = 3.31 + 0.02392 × mean glucose (mg/dL)"; the
// app's mean-glucose default 154 mg/dL → 3.31 + 3.68368 = 6.99368 (HAND-DERIVED).
func TestGlucoseManagementIndicatorPercent_parity(t *testing.T) {
	if got, err := GlucoseManagementIndicatorPercent(154); err != nil || !relNear(got, 6.99368, 1e-12) {
		t.Errorf("got (%v, %v), want 6.99368", got, err)
	}
}

// HAND-DERIVED from the source's mmol/L equation: 7 → 12.71 + 32.94109 = 45.65109; 10 → 59.7687.
func TestGlucoseManagementIndicatorMmolMol_example(t *testing.T) {
	tests := []struct{ mean, want float64 }{{7, 45.65109}, {10, 59.7687}}
	for _, tc := range tests {
		if got, err := GlucoseManagementIndicatorMmolMol(tc.mean); err != nil || !relNear(got, tc.want, 1e-12) {
			t.Errorf("GMI mmol/mol(%v) = (%v, %v), want %v", tc.mean, got, err, tc.want)
		}
	}
}

func TestGlucoseManagementIndicator_range(t *testing.T) {
	if got, err := GlucoseManagementIndicatorPercent(math.MaxFloat64); err != nil || !reduce.AllFinite(got) {
		t.Errorf("MaxFloat64: got (%v, %v), want a finite value", got, err)
	}
	if _, err := GlucoseManagementIndicatorPercent(5e-324); err != nil {
		t.Errorf("smallest positive refused: %v", err)
	}
	if _, err := GlucoseManagementIndicatorMmolMol(5e-324); err != nil {
		t.Errorf("smallest positive refused: %v", err)
	}
}

func TestGlucoseManagementIndicator_errors(t *testing.T) {
	pct := func(v float64) func() error {
		return func() error { _, err := GlucoseManagementIndicatorPercent(v); return err }
	}
	si := func(v float64) func() error {
		return func() error { _, err := GlucoseManagementIndicatorMmolMol(v); return err }
	}
	runErrCases(t, []errCase{
		{"pct NaN", pct(nan), ErrNotFinite},
		{"pct +Inf", pct(pInf), ErrNotFinite},
		{"pct -Inf", pct(nInf), ErrNotFinite},
		{"pct 0", pct(0), ErrOutOfRange},
		{"pct negative", pct(-1), ErrOutOfRange},
		{"si NaN", si(nan), ErrNotFinite},
		{"si +Inf", si(pInf), ErrNotFinite},
		{"si -Inf", si(nInf), ErrNotFinite},
		{"si 0", si(0), ErrOutOfRange},
		{"si negative", si(-1), ErrOutOfRange},
		{"si overflow", si(math.MaxFloat64), ErrDomain},
	})
}

func TestGlucoseManagementIndicator_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = GlucoseManagementIndicatorPercent(150) }); n != 0 {
		t.Errorf("GlucoseManagementIndicatorPercent allocates %v times, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sink, _ = GlucoseManagementIndicatorMmolMol(8) }); n != 0 {
		t.Errorf("GlucoseManagementIndicatorMmolMol allocates %v times, want 0", n)
	}
}

func BenchmarkGlucoseManagementIndicatorPercent(b *testing.B) {
	b.ReportAllocs()
	v := 150.0
	for b.Loop() {
		sink, _ = GlucoseManagementIndicatorPercent(v)
	}
}

func BenchmarkGlucoseManagementIndicatorMmolMol(b *testing.B) {
	b.ReportAllocs()
	v := 8.0
	for b.Loop() {
		sink, _ = GlucoseManagementIndicatorMmolMol(v)
	}
}
