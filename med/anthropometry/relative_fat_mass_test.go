package anthropometry

import (
	"math"
	"testing"
)

// HAND-DERIVED from Woolcott & Bergman 2018 (abstract). Tolerance 1e-13 relative (a divide, a multiply and two
// adds; the result can be small after cancellation, so the absolute error is about 1e-14).
func TestRelativeFatMass_example(t *testing.T) {
	tests := []struct {
		h, w   float64
		female bool
		want   float64
	}{
		{175, 88, false, 24.227272727272727},
		{165, 80, true, 34.75},
		{160, 100, false, 32},
	}
	for _, tc := range tests {
		if got, err := RelativeFatMass(tc.h, tc.w, tc.female); err != nil || !relNear(got, tc.want, 1e-13) {
			t.Errorf("RelativeFatMass(%v, %v, %v) = (%v, %v), want %v", tc.h, tc.w, tc.female, got, err, tc.want)
		}
	}
}

// parity: ledger relative_fat_mass, source cases 1–6 (inches; the ratio is unit-free, so inches × 2.54 give the
// same result). The source prints one decimal: 22.8, 33.3, 33.4, 45.0, 12.6, 25.2 — each within 0.05 of the
// exact value asserted here.
func TestRelativeFatMass_parity(t *testing.T) {
	tests := []struct {
		hIn, wIn      float64
		female        bool
		want, printed float64
	}{
		{70, 34, false, 22.82352941176471, 22.8},
		{64, 30, true, 33.333333333333336, 33.3},
		{68, 44.5, false, 33.438202247191015, 33.4},
		{62, 40, true, 45.0, 45.0},
		{72, 28, false, 12.57142857142857, 12.6},
		{66, 26, true, 25.230769230769234, 25.2},
	}
	for i, tc := range tests {
		got, err := RelativeFatMass(tc.hIn*2.54, tc.wIn*2.54, tc.female)
		if err != nil || !relNear(got, tc.want, 1e-13) || math.Abs(got-tc.printed) > 0.05 {
			t.Errorf("source case %d: got (%v, %v), want %v (printed %v)", i+1, got, err, tc.want, tc.printed)
		}
	}
}

func TestRelativeFatMass_errors(t *testing.T) {
	call := func(h, w float64) func() error {
		return func() error { _, err := RelativeFatMass(h, w, false); return err }
	}
	runErrCases(t, []errCase{
		{"height NaN", call(nan, 80), ErrNotFinite},
		{"height +Inf", call(pInf, 80), ErrNotFinite},
		{"waist -Inf", call(170, nInf), ErrNotFinite},
		{"NaN before range", call(0, nan), ErrNotFinite},
		{"height 0", call(0, 80), ErrOutOfRange},
		{"waist negative", call(170, -1), ErrOutOfRange},
		// 1e308 / 1e-10 overflows to +Inf, so the result is −Inf.
		{"ratio overflows", call(1e308, 1e-10), ErrDomain},
	})
}

func TestRelativeFatMass_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = RelativeFatMass(175, 88, false) }); n != 0 {
		t.Errorf("RelativeFatMass allocates %v times, want 0", n)
	}
}

func BenchmarkRelativeFatMass(b *testing.B) {
	b.ReportAllocs()
	h, w := 175.0, 88.0
	for b.Loop() {
		sink, _ = RelativeFatMass(h, w, false)
	}
}
