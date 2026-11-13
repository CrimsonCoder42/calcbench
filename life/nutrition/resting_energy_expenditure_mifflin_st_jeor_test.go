package nutrition

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

// Tolerance 1e-12 relative: the equations are a few exact-decimal products and sums.

// Mifflin 1990 abstract (PMID 2305711), the derived equation; HAND-DERIVED: man 70 kg, 175 cm, 30 y →
// 699.3 + 1093.75 − 147.6 + 166 − 161 = 1650.45.
func TestRestingEnergyExpenditureMifflinStJeor_example(t *testing.T) {
	if got, err := RestingEnergyExpenditureMifflinStJeor(70, 175, 30, false); err != nil || !relNear(got, 1650.45, 1e-12) {
		t.Errorf("got (%v, %v), want 1650.45", got, err)
	}
}

// parity: ledger mifflin_st_jeor_bmr, source cases 1–6 (sex 1 = male, 0 = female).
func TestRestingEnergyExpenditureMifflinStJeor_parity(t *testing.T) {
	tests := []struct {
		name              string
		female            bool
		kg, cm, age, want float64
	}{
		{"source case 1", false, 70, 175, 30, 1650.45},
		{"source case 2", true, 60, 165, 40, 1272.85},
		{"source case 3", false, 100, 180, 50, 1883},
		{"source case 4", true, 80, 160, 19, 1544.72},
		{"source case 5", false, 55, 170, 78, 1233.19},
		{"source case 6", true, 120, 155, 65, 1686.75},
	}
	for _, tc := range tests {
		got, err := RestingEnergyExpenditureMifflinStJeor(tc.kg, tc.cm, tc.age, tc.female)
		if err != nil || !relNear(got, tc.want, 1e-12) {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// Mifflin 1990 abstract, simplified forms; HAND-DERIVED: man 70/175/30 → 1648.75; woman 60/165/40 → 1270.25.
func TestRestingEnergyExpenditureMifflinStJeorSimplified_example(t *testing.T) {
	tests := []struct {
		female            bool
		kg, cm, age, want float64
	}{
		{false, 70, 175, 30, 1648.75},
		{true, 60, 165, 40, 1270.25},
	}
	for _, tc := range tests {
		got, err := RestingEnergyExpenditureMifflinStJeorSimplified(tc.kg, tc.cm, tc.age, tc.female)
		if err != nil || !relNear(got, tc.want, 1e-12) {
			t.Errorf("(%v, %v, %v, %v) = (%v, %v), want %v", tc.kg, tc.cm, tc.age, tc.female, got, err, tc.want)
		}
	}
}

// The derived and simplified forms differ by a few kcal/day (the abstract: the simplification "did not affect
// its predictive value"); the 70/175/30 man differs by 1.7.
func TestRestingEnergyExpenditureMifflinStJeor_formsAgree(t *testing.T) {
	a, _ := RestingEnergyExpenditureMifflinStJeor(70, 175, 30, false)
	b, _ := RestingEnergyExpenditureMifflinStJeorSimplified(70, 175, 30, false)
	if math.Abs(a-b-1.7) > 1e-9 {
		t.Errorf("difference %v, want 1.7", a-b)
	}
}

func TestRestingEnergyExpenditureMifflinStJeor_rangeEdges(t *testing.T) {
	for _, age := range []float64{19, 78, 78.999} {
		if _, err := RestingEnergyExpenditureMifflinStJeor(70, 175, age, true); err != nil {
			t.Errorf("age %v refused: %v", age, err)
		}
		if _, err := RestingEnergyExpenditureMifflinStJeorSimplified(70, 175, age, true); err != nil {
			t.Errorf("simplified age %v refused: %v", age, err)
		}
	}
}

func TestRestingEnergyExpenditureMifflinStJeor_errors(t *testing.T) {
	d := func(kg, cm, age float64) func() error {
		return func() error { _, err := RestingEnergyExpenditureMifflinStJeor(kg, cm, age, false); return err }
	}
	s := func(kg, cm, age float64) func() error {
		return func() error { _, err := RestingEnergyExpenditureMifflinStJeorSimplified(kg, cm, age, true); return err }
	}
	runErrCases(t, []errCase{
		// parity: ledger mifflin_st_jeor_bmr, source refusals: age 18 and age 79.
		{"age 18", d(70, 175, 18), ErrOutOfRange},
		{"age 79", d(70, 175, 79), ErrOutOfRange},
		{"age 18.999", d(70, 175, 18.999), ErrOutOfRange},
		{"weight NaN", d(nan, 175, 30), ErrNotFinite},
		{"height +Inf", d(70, pInf, 30), ErrNotFinite},
		{"age -Inf", d(70, 175, nInf), ErrNotFinite},
		{"age NaN", d(70, 175, nan), ErrNotFinite},
		{"weight 0", d(0, 175, 30), ErrOutOfRange},
		{"height negative", d(70, -1, 30), ErrOutOfRange},
		{"overflow", d(math.MaxFloat64, 175, 30), ErrDomain},
		{"simplified age 79", s(70, 175, 79), ErrOutOfRange},
		{"simplified NaN", s(nan, 175, 30), ErrNotFinite},
		{"simplified weight 0", s(0, 175, 30), ErrOutOfRange},
		{"simplified overflow", s(math.MaxFloat64, 175, 30), ErrDomain},
	})
}

func TestRestingEnergyExpenditureMifflinStJeor_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = RestingEnergyExpenditureMifflinStJeor(70, 175, 30, false) }); n != 0 {
		t.Errorf("RestingEnergyExpenditureMifflinStJeor allocates %v times, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sink, _ = RestingEnergyExpenditureMifflinStJeorSimplified(70, 175, 30, false) }); n != 0 {
		t.Errorf("RestingEnergyExpenditureMifflinStJeorSimplified allocates %v times, want 0", n)
	}
}

func BenchmarkRestingEnergyExpenditureMifflinStJeor(b *testing.B) {
	b.ReportAllocs()
	kg, cm, age := 70.0, 175.0, 30.0
	for b.Loop() {
		sink, _ = RestingEnergyExpenditureMifflinStJeor(kg, cm, age, false)
	}
}

func BenchmarkRestingEnergyExpenditureMifflinStJeorSimplified(b *testing.B) {
	b.ReportAllocs()
	kg, cm, age := 70.0, 175.0, 30.0
	for b.Loop() {
		sink, _ = RestingEnergyExpenditureMifflinStJeorSimplified(kg, cm, age, false)
	}
}
