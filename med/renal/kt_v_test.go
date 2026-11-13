package renal

import (
	"math"
	"testing"
)

// Gotch & Sargent 1985 definition, HAND-DERIVED: K × t / V. 250 × 240 / 42 000 = 1.428571…
func TestKtV_definition(t *testing.T) {
	tests := []struct{ k, tm, v, want float64 }{
		{250, 240, 42000, 250.0 * 240 / 42000},
		{0, 240, 42000, 0},
		{300, 0, 35000, 0},
		{200, 180, 30000, 1.2},
	}
	for _, tc := range tests {
		if got, err := KtV(tc.k, tc.tm, tc.v); err != nil || !relNear(got, tc.want, 1e-12) {
			t.Errorf("%+v: (%v, %v)", tc, got, err)
		}
	}
}

// parity: ledger kt_v_dialysis, source case 1. The app takes V = 0.6 × weight (no source found); with that V
// passed explicitly (70 kg → 42 000 mL) the ratio agrees: K 250, t 240 → 1.4286.
func TestKtV_parity(t *testing.T) {
	if got, err := KtV(250, 240, 42000); err != nil || math.Abs(got-1.4286) > 1e-4 {
		t.Errorf("got (%v, %v), want 1.4286", got, err)
	}
}

func TestKtV_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"K NaN", func() error { _, e := KtV(nan, 240, 42000); return e }, ErrNotFinite},
		{"t +Inf", func() error { _, e := KtV(250, pInf, 42000); return e }, ErrNotFinite},
		{"V -Inf", func() error { _, e := KtV(250, 240, nInf); return e }, ErrNotFinite},
		{"K -1", func() error { _, e := KtV(-1, 240, 42000); return e }, ErrOutOfRange},
		{"t -1", func() error { _, e := KtV(250, -1, 42000); return e }, ErrOutOfRange},
		{"V 0", func() error { _, e := KtV(250, 240, 0); return e }, ErrOutOfRange},
		{"overflow", func() error { _, e := KtV(math.MaxFloat64, 2, 1); return e }, ErrDomain},
	})
}

func TestKtV_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = KtV(250, 240, 42000) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkKtV(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sink, _ = KtV(250, 240, 42000)
	}
}
