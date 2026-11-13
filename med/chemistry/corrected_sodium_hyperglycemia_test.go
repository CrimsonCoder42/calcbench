package chemistry

import (
	"math"
	"testing"
)

// HAND-DERIVED from Kitabchi 2009 (Diagnosis): add 1.6 mEq/L per 100 mg/dL of glucose above 100 mg/dL.
// Tolerance 1e-12 relative.
func TestCorrectedSodiumHyperglycemia_example(t *testing.T) {
	tests := []struct {
		name    string
		na, glu float64
		want    float64
	}{
		{"worked", 130, 600, 138},
		{"glucose 100: nothing added", 135, 100, 135},
		{"glucose 100.5: just above", 135, 100.5, 135.008},
		{"glucose 99: nothing added", 135, 99, 135},
		{"glucose 200", 140, 200, 141.6},
		{"glucose 1200", 125, 1200, 142.6},
	}
	for _, tc := range tests {
		if got, err := CorrectedSodiumHyperglycemia(tc.na, tc.glu); err != nil || !relNear(got, tc.want, 1e-12) {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// parity: ledger bs_corrected_na. Case 1: app defaults Na 135, glucose 300 → 135 + 1.6 × 2 = 138.2 (agrees).
// Case 2: glucose 500 — the app switches to 2.4 per 100 mg/dL above 400 and shows 144.6; the primary states
// 1.6 only, so this asserts 135 + 1.6 × 4 = 141.4 (finding: Hillier 1999's abstract recommends 4.0, not 2.4,
// above 400 mg/dL).
func TestCorrectedSodiumHyperglycemia_parity(t *testing.T) {
	cases := []struct{ na, glu, want float64 }{
		{135, 300, 138.2}, // parity: ledger bs_corrected_na, source case 1
		{135, 500, 141.4}, // parity: ledger bs_corrected_na, source case 2 (app 144.6)
	}
	for _, c := range cases {
		if got, err := CorrectedSodiumHyperglycemia(c.na, c.glu); err != nil || !relNear(got, c.want, 1e-12) {
			t.Errorf("Na %v glucose %v: got (%v, %v), want %v", c.na, c.glu, got, err, c.want)
		}
	}
}

func TestCorrectedSodiumHyperglycemia_errors(t *testing.T) {
	call := func(na, g float64) func() error {
		return func() error { _, err := CorrectedSodiumHyperglycemia(na, g); return err }
	}
	runErrCases(t, []errCase{
		{"Na NaN", call(nan, 300), ErrNotFinite},
		{"glucose +Inf", call(135, pInf), ErrNotFinite},
		{"glucose -Inf", call(135, nInf), ErrNotFinite},
		{"NaN before range", call(0, nan), ErrNotFinite},
		{"Na 0", call(0, 300), ErrOutOfRange},
		{"glucose 0", call(135, 0), ErrOutOfRange},
		{"overflow", call(math.MaxFloat64, math.MaxFloat64), ErrDomain},
	})
}

func TestCorrectedSodiumHyperglycemia_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = CorrectedSodiumHyperglycemia(130, 600) }); n != 0 {
		t.Errorf("CorrectedSodiumHyperglycemia allocates %v times, want 0", n)
	}
}

func BenchmarkCorrectedSodiumHyperglycemia(b *testing.B) {
	b.ReportAllocs()
	na, g := 130.0, 600.0
	for b.Loop() {
		sink, _ = CorrectedSodiumHyperglycemia(na, g)
	}
}
