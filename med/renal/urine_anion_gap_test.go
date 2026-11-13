package renal

import (
	"math"
	"testing"
)

// Goldstein 1986: UAG = Na + K − Cl. Exact for these integers.
func TestUrineAnionGap_source(t *testing.T) {
	tests := []struct{ na, k, cl, want float64 }{{50, 30, 100, -20}, {40, 20, 40, 20}, {0, 0, 0, 0}, {60.5, 25.5, 86, 0}}
	for _, tc := range tests {
		if got, err := UrineAnionGap(tc.na, tc.k, tc.cl); err != nil || got != tc.want {
			t.Errorf("%+v → (%v, %v), want %v", tc, got, err, tc.want)
		}
	}
}

// parity: ledger urine_anion_gap, source case 1 (app uag = na + k − cl, displayed to 0.1; unrounded here).
func TestUrineAnionGap_parity(t *testing.T) {
	if got, _ := UrineAnionGap(45, 35, 60); got != 20 {
		t.Errorf("got %v, want 20", got)
	}
}

func TestUrineAnionGap_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"Na NaN", func() error { _, e := UrineAnionGap(nan, 1, 1); return e }, ErrNotFinite},
		{"K +Inf", func() error { _, e := UrineAnionGap(1, pInf, 1); return e }, ErrNotFinite},
		{"Cl -Inf", func() error { _, e := UrineAnionGap(1, 1, nInf); return e }, ErrNotFinite},
		{"Na < 0", func() error { _, e := UrineAnionGap(-1, 1, 1); return e }, ErrOutOfRange},
		{"K < 0", func() error { _, e := UrineAnionGap(1, -1, 1); return e }, ErrOutOfRange},
		{"Cl < 0", func() error { _, e := UrineAnionGap(1, 1, -1); return e }, ErrOutOfRange},
		{"overflow", func() error { _, e := UrineAnionGap(math.MaxFloat64, math.MaxFloat64, 0); return e }, ErrDomain},
	})
}

func TestUrineAnionGap_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = UrineAnionGap(45, 35, 60) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkUrineAnionGap(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sink, _ = UrineAnionGap(45, 35, 60)
	}
}
