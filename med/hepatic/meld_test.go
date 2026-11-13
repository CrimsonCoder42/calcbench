package hepatic

import (
	"math"
	"testing"
)

// meldOf is the test's own copy of the formula quoted in Kim et al. 2021 (9.57 ln Cr + 3.78 ln bili +
// 11.20 ln INR + 6.43), applied to already-bounded values.
func meldOf(cr, bili, inr float64) float64 {
	return 9.57*math.Log(cr) + 3.78*math.Log(bili) + 11.20*math.Log(inr) + 6.43
}

// Worked examples, HAND-DERIVED from the quoted formula. Every input 1 → 6.43 (every log is 0).
// Bounds: values below 1 are set to 1; creatinine is capped at 4.0 and set to 4.0 with dialysis.
// Tolerance 1e-13 relative (three logs and a short sum).
func TestMELD_example(t *testing.T) {
	tests := []struct {
		name          string
		bili, inr, cr float64
		dialysis      bool
		want          float64
	}{
		{"all 1", 1, 1, 1, false, 6.43},
		{"all below 1 set to 1", 0.3, 0.8, 0.5, false, 6.43},
		{"bili 2, INR 1.5, Cr 2", 2, 1.5, 2, false, meldOf(2, 2, 1.5)},
		{"creatinine at the cap 4.0", 3, 2, 4.0, false, meldOf(4, 3, 2)},
		{"creatinine above the cap", 3, 2, 7.5, false, meldOf(4, 3, 2)},
		{"creatinine just below the cap", 3, 2, 3.999, false, meldOf(3.999, 3, 2)},
		{"dialysis sets creatinine to 4.0", 3, 2, 1.1, true, meldOf(4, 3, 2)},
		{"dialysis with low creatinine", 1, 1, 0.5, true, 9.57*math.Log(4) + 6.43},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := MELD(tc.bili, tc.inr, tc.cr, tc.dialysis)
			if err != nil || !relNear(got, tc.want, 1e-13) {
				t.Errorf("got (%v, %v), want %v", got, err, tc.want)
			}
		})
	}
}

func TestMELD_errors(t *testing.T) {
	call := func(bili, inr, cr float64) func() error {
		return func() error { _, err := MELD(bili, inr, cr, false); return err }
	}
	runErrCases(t, []errCase{
		{"bilirubin NaN", call(nan, 1, 1), ErrNotFinite},
		{"bilirubin +Inf", call(pInf, 1, 1), ErrNotFinite},
		{"bilirubin -Inf", call(nInf, 1, 1), ErrNotFinite},
		{"INR NaN", call(1, nan, 1), ErrNotFinite},
		{"INR +Inf", call(1, pInf, 1), ErrNotFinite},
		{"INR -Inf", call(1, nInf, 1), ErrNotFinite},
		{"creatinine NaN", call(1, 1, nan), ErrNotFinite},
		{"creatinine +Inf", call(1, 1, pInf), ErrNotFinite},
		{"creatinine -Inf", call(1, 1, nInf), ErrNotFinite},
		{"NaN before range", call(0, 1, nan), ErrNotFinite},
		{"bilirubin 0", call(0, 1, 1), ErrOutOfRange},
		{"INR 0", call(1, 0, 1), ErrOutOfRange},
		{"creatinine 0", call(1, 1, 0), ErrOutOfRange},
		{"creatinine negative", call(1, 1, -1), ErrOutOfRange},
	})
	// Largest finite inputs: logs keep the result finite (no overflow path exists).
	if got, err := MELD(math.MaxFloat64, math.MaxFloat64, math.MaxFloat64, false); err != nil || math.IsInf(got, 0) {
		t.Errorf("MaxFloat64 inputs: got (%v, %v), want finite", got, err)
	}
}

func TestMELD_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = MELD(2, 1.5, 2, false) }); n != 0 {
		t.Errorf("MELD allocates %v times, want 0", n)
	}
}

func BenchmarkMELD(b *testing.B) {
	b.ReportAllocs()
	bili, inr, cr := 2.0, 1.5, 2.0
	for b.Loop() {
		sink, _ = MELD(bili, inr, cr, false)
	}
}
