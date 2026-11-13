package hepatic

import (
	"math"
	"testing"
)

// HAND-DERIVED from Harrison 2008 (abstract): BMI ≥ 28 → 1, AST/ALT ≥ 0.8 → 2, diabetes → 1. Both sides of
// each edge; 32/40 is exactly the double nearest 0.8.
func TestBARDScore_example(t *testing.T) {
	tests := []struct {
		name          string
		bmi, ast, alt float64
		dm            bool
		want          int
	}{
		{"worked", 30, 40, 45, false, 3},
		{"none", 25, 20, 40, false, 0},
		{"BMI 27.99", 27.99, 20, 40, false, 0},
		{"BMI 28", 28, 20, 40, false, 1},
		{"AAR 0.8", 25, 32, 40, false, 2},
		{"AAR 0.7975", 25, 31.9, 40, false, 0},
		{"diabetes", 25, 20, 40, true, 1},
		{"maximum", 35, 60, 30, true, 4},
		{"ratio overflows to +Inf", 25, math.MaxFloat64, 1e-10, false, 2},
	}
	for _, tc := range tests {
		if got, err := BARDScore(tc.bmi, tc.ast, tc.alt, tc.dm); err != nil || got != tc.want {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// parity: ledger bard_score. The app takes BMI as a yes/no item (bmi_elevated) and computes ratio = ast / alt,
// ≥ 0.8 → 2. Its JSON ranges AST/ALT 1–500 U/L: AST 50, ALT 50, BMI elevated, diabetes → 1 + 2 + 1 = 4.
func TestBARDScore_parity(t *testing.T) {
	if got, err := BARDScore(30, 50, 50, true); err != nil || got != 4 {
		t.Errorf("got (%v, %v), want 4", got, err)
	}
}

func TestBARDScore_errors(t *testing.T) {
	call := func(b, a, l float64) func() error {
		return func() error { _, err := BARDScore(b, a, l, false); return err }
	}
	runErrCases(t, []errCase{
		{"BMI NaN", call(nan, 40, 40), ErrNotFinite},
		{"AST +Inf", call(30, pInf, 40), ErrNotFinite},
		{"ALT -Inf", call(30, 40, nInf), ErrNotFinite},
		{"NaN before range", call(0, 40, nan), ErrNotFinite},
		{"BMI 0", call(0, 40, 40), ErrOutOfRange},
		{"AST 0", call(30, 0, 40), ErrOutOfRange},
		{"ALT negative", call(30, 40, -1), ErrOutOfRange},
	})
}

func TestBARDScore_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = BARDScore(30, 40, 45, true) }); n != 0 {
		t.Errorf("BARDScore allocates %v times, want 0", n)
	}
}

func BenchmarkBARDScore(b *testing.B) {
	b.ReportAllocs()
	bmi, ast, alt := 30.0, 40.0, 45.0
	for b.Loop() {
		sinkInt, _ = BARDScore(bmi, ast, alt, true)
	}
}
