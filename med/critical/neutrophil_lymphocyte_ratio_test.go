package critical

import (
	"math"
	"testing"
)

// Definition (Zahorec 2001 title). Tolerance 1e-15 relative: one division.
func TestNeutrophilLymphocyteRatio_source(t *testing.T) {
	tests := []struct{ n, l, want float64 }{{8, 1, 8}, {4.5, 1.5, 3}, {0, 2, 0}, {1, 4, 0.25}}
	for _, tc := range tests {
		got, err := NeutrophilLymphocyteRatio(tc.n, tc.l)
		if err != nil || math.Abs(got-tc.want) > 1e-15*math.Max(1, tc.want) {
			t.Errorf("(%v, %v) → (%v, %v), want %v", tc.n, tc.l, got, err, tc.want)
		}
	}
}

// parity: ledger nlr_calculator, source case 1 (app nlr = neut / lymph; null when lymph ≤ 0 → ErrOutOfRange here).
func TestNeutrophilLymphocyteRatio_parity(t *testing.T) {
	if got, _ := NeutrophilLymphocyteRatio(10, 2); got != 5 {
		t.Errorf("got %v, want 5", got)
	}
	_, err := NeutrophilLymphocyteRatio(10, 0)
	wantErrB6(t, "lymph 0", err, ErrOutOfRange)
}

func TestNeutrophilLymphocyteRatio_errors(t *testing.T) {
	for _, v := range nonFiniteB6 {
		_, err := NeutrophilLymphocyteRatio(v, 1)
		wantErrB6(t, "neut", err, ErrNotFinite)
		_, err = NeutrophilLymphocyteRatio(1, v)
		wantErrB6(t, "lymph", err, ErrNotFinite)
	}
	_, err := NeutrophilLymphocyteRatio(-1, 1)
	wantErrB6(t, "neut < 0", err, ErrOutOfRange)
	_, err = NeutrophilLymphocyteRatio(1, -1)
	wantErrB6(t, "lymph < 0", err, ErrOutOfRange)
	_, err = NeutrophilLymphocyteRatio(math.MaxFloat64, 1e-300)
	wantErrB6(t, "overflow", err, ErrOutOfRange)
}

func TestNeutrophilLymphocyteRatio_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkB6Float, _ = NeutrophilLymphocyteRatio(8, 1.2) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkNeutrophilLymphocyteRatio(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkB6Float, _ = NeutrophilLymphocyteRatio(8, 1.2)
	}
}
