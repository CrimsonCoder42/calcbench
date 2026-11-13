package chemistry

import (
	"math"
	"testing"
)

// HAND-DERIVED from Kitabchi 2009 Table 1 footnote: AG = Na − (Cl + HCO3). Integer-valued inputs give exact
// float64 results, so the comparison is exact.
func TestAnionGap_example(t *testing.T) {
	tests := []struct {
		name         string
		na, cl, hco3 float64
		want         float64
	}{
		{"worked", 140, 104, 24, 12},
		{"wide gap (Table 1 moderate DKA > 12)", 135, 100, 10, 25},
		{"zero gap", 130, 110, 20, 0},
		{"negative gap returned as computed", 120, 110, 20, -10},
		{"fractional", 138.5, 101.25, 22.25, 15},
	}
	for _, tc := range tests {
		if got, err := AnionGap(tc.na, tc.cl, tc.hco3); err != nil || got != tc.want {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// parity: ledger serum_anion_gap, source case 1 (the app's default inputs Na 140, Cl 104, HCO3 24 → AG 12).
// The app's albumin correction (+2.5 × (4.0 − albumin)) is MED-CHEM-008 with normalAlbumin 4.0, not repeated here.
func TestAnionGap_parity(t *testing.T) {
	if got, err := AnionGap(140, 104, 24); err != nil || got != 12 {
		t.Errorf("got (%v, %v), want 12", got, err)
	}
	ag, _ := AnionGap(140, 104, 24)
	if got, err := AdjustedAnionGapFiggeGDL(ag, 2.0, 4.0); err != nil || !near(got, 17, 1e-12) {
		t.Errorf("albumin 2.0 g/dL: got (%v, %v), want 17 (parity: ledger serum_anion_gap, source case 2)", got, err)
	}
}

func TestAnionGap_errors(t *testing.T) {
	call := func(na, cl, h float64) func() error {
		return func() error { _, err := AnionGap(na, cl, h); return err }
	}
	runErrCases(t, []errCase{
		{"Na NaN", call(nan, 104, 24), ErrNotFinite},
		{"Cl +Inf", call(140, pInf, 24), ErrNotFinite},
		{"HCO3 -Inf", call(140, 104, nInf), ErrNotFinite},
		{"NaN before range", call(0, 104, nan), ErrNotFinite},
		{"Na 0", call(0, 104, 24), ErrOutOfRange},
		{"Cl negative", call(140, -1, 24), ErrOutOfRange},
		{"HCO3 0", call(140, 104, 0), ErrOutOfRange},
		{"overflow", call(math.MaxFloat64, math.MaxFloat64, math.MaxFloat64), ErrDomain},
	})
}

func TestAnionGap_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = AnionGap(140, 104, 24) }); n != 0 {
		t.Errorf("AnionGap allocates %v times, want 0", n)
	}
}

func BenchmarkAnionGap(b *testing.B) {
	b.ReportAllocs()
	na, cl, h := 140.0, 104.0, 24.0
	for b.Loop() {
		sink, _ = AnionGap(na, cl, h)
	}
}
