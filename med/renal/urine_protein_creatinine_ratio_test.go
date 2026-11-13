package renal

import (
	"math"
	"testing"
)

// Ginsberg 1983: ratio in mg/mg. Tolerance 1e-15 relative: one division.
func TestUrineProteinCreatinineRatio_source(t *testing.T) {
	tests := []struct{ p, c, want float64 }{{350, 100, 3.5}, {20, 100, 0.2}, {0, 50, 0}, {150, 75, 2}}
	for _, tc := range tests {
		got, err := UrineProteinCreatinineRatio(tc.p, tc.c)
		if err != nil || !relNear(got, tc.want, 1e-15) {
			t.Errorf("(%v, %v) → (%v, %v), want %v", tc.p, tc.c, got, err, tc.want)
		}
	}
}

// parity: ledger upee, source case 1 (app ratio = urine_protein / urine_creatinine; null when creatinine ≤ 0).
func TestUrineProteinCreatinineRatio_parity(t *testing.T) {
	if got, _ := UrineProteinCreatinineRatio(300, 120); !relNear(got, 2.5, 1e-15) {
		t.Errorf("got %v, want 2.5", got)
	}
}

func TestUrineProteinCreatinineRatio_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"protein NaN", func() error { _, e := UrineProteinCreatinineRatio(nan, 1); return e }, ErrNotFinite},
		{"protein +Inf", func() error { _, e := UrineProteinCreatinineRatio(pInf, 1); return e }, ErrNotFinite},
		{"creatinine -Inf", func() error { _, e := UrineProteinCreatinineRatio(1, nInf); return e }, ErrNotFinite},
		{"creatinine 0", func() error { _, e := UrineProteinCreatinineRatio(1, 0); return e }, ErrOutOfRange},
		{"protein < 0", func() error { _, e := UrineProteinCreatinineRatio(-1, 1); return e }, ErrOutOfRange},
		{"overflow", func() error { _, e := UrineProteinCreatinineRatio(math.MaxFloat64, 1e-300); return e }, ErrDomain},
	})
}

func TestUrineProteinCreatinineRatio_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = UrineProteinCreatinineRatio(300, 120) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkUrineProteinCreatinineRatio(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sink, _ = UrineProteinCreatinineRatio(300, 120)
	}
}
