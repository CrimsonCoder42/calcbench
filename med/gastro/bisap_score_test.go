package gastro

import (
	"errors"
	"math"
	"testing"
)

// Shared by the package's tests.

// sinkInt and sinkFloat keep benchmark results observable.
var (
	sinkInt   int
	sinkFloat float64
)

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

// HAND-DERIVED from Wu 2008 (abstract). The edges are strict: BUN > 25 mg/dL, age > 60 years.
func TestBISAPScore_example(t *testing.T) {
	tests := []struct {
		name string
		f    BISAPFindings
		want int
	}{
		{"BUN and age", BISAPFindings{BUNMgDL: 30, AgeYears: 65}, 2},
		{"BUN 25, age 60", BISAPFindings{BUNMgDL: 25, AgeYears: 60}, 0},
		{"BUN 25.1", BISAPFindings{BUNMgDL: 25.1, AgeYears: 40}, 1},
		{"age 60.01", BISAPFindings{BUNMgDL: 10, AgeYears: 60.01}, 1},
		{"all five", BISAPFindings{40, true, true, 70, true}, 5},
		{"mental only", BISAPFindings{BUNMgDL: 10, ImpairedMentalStatus: true, AgeYears: 30}, 1},
		{"SIRS only", BISAPFindings{BUNMgDL: 10, SIRS: true, AgeYears: 30}, 1},
		{"effusion only", BISAPFindings{BUNMgDL: 10, AgeYears: 30, PleuralEffusion: true}, 1},
		{"zeros accepted", BISAPFindings{}, 0},
	}
	for _, tc := range tests {
		if got, err := BISAPScore(tc.f); err != nil || got != tc.want {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// parity: ledger bisap_pancreatitis, the app's five yes/no items, one point each (0–5): all yes → 5, all no → 0,
// three yes → 3.
func TestBISAPScore_parity(t *testing.T) {
	if got, err := BISAPScore(BISAPFindings{BUNMgDL: 26, SIRS: true, AgeYears: 61}); err != nil || got != 3 {
		t.Errorf("three items: got (%v, %v), want 3", got, err)
	}
}

func TestBISAPScore_errors(t *testing.T) {
	call := func(f BISAPFindings) func() error {
		return func() error { _, err := BISAPScore(f); return err }
	}
	runErrCases(t, []errCase{
		{"BUN NaN", call(BISAPFindings{BUNMgDL: nan, AgeYears: 50}), ErrNotFinite},
		{"BUN +Inf", call(BISAPFindings{BUNMgDL: pInf, AgeYears: 50}), ErrNotFinite},
		{"age -Inf", call(BISAPFindings{BUNMgDL: 10, AgeYears: nInf}), ErrNotFinite},
		{"BUN negative", call(BISAPFindings{BUNMgDL: -1, AgeYears: 50}), ErrOutOfRange},
		{"age negative", call(BISAPFindings{BUNMgDL: 10, AgeYears: -1}), ErrOutOfRange},
	})
}

func TestBISAPScore_noAllocation(t *testing.T) {
	f := BISAPFindings{30, true, false, 65, false}
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = BISAPScore(f) }); n != 0 {
		t.Errorf("BISAPScore allocates %v times, want 0", n)
	}
}

func BenchmarkBISAPScore(b *testing.B) {
	b.ReportAllocs()
	f := BISAPFindings{30, true, false, 65, false}
	for b.Loop() {
		sinkInt, _ = BISAPScore(f)
	}
}
