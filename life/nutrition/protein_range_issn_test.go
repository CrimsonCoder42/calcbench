package nutrition

import (
	"math"
	"testing"
)

// Tolerance 1e-12 relative: one product with an exact-decimal factor.

// ISSN 2017 point 2 (PMC5477153); HAND-DERIVED: 70 kg → 1.4 × 70 = 98, 2.0 × 70 = 140 g/day.
func TestProteinRangeISSN_example(t *testing.T) {
	lo, hi, err := ProteinRangeISSN(70)
	if err != nil || !relNear(lo, 98, 1e-12) || !relNear(hi, 140, 1e-12) {
		t.Errorf("got (%v, %v, %v), want 98, 140", lo, hi, err)
	}
}

// parity: ledger protein_per_kg, source cases 1–6.
func TestProteinRangeISSN_parity(t *testing.T) {
	tests := []struct{ kg, lo, hi float64 }{
		{50, 70, 100}, {70, 98, 140}, {80, 112, 160}, {62.5, 87.5, 125}, {100, 140, 200}, {120, 168, 240},
	}
	for i, tc := range tests {
		lo, hi, err := ProteinRangeISSN(tc.kg)
		if err != nil || !relNear(lo, tc.lo, 1e-12) || !relNear(hi, tc.hi, 1e-12) {
			t.Errorf("source case %d (%v kg): got (%v, %v, %v), want %v, %v", i+1, tc.kg, lo, hi, err, tc.lo, tc.hi)
		}
	}
}

// parity: ledger protein_per_kg, source refusals: 10 kg (the source's physical bound 20–300 kg) and age 16
// (the source's adult default). The stand states no weight or age limits and the function takes no age, so
// 10 kg is accepted: 14–20 g/day.
func TestProteinRangeISSN_parityTenKg(t *testing.T) {
	lo, hi, err := ProteinRangeISSN(10)
	if err != nil || !relNear(lo, 14, 1e-12) || !relNear(hi, 20, 1e-12) {
		t.Errorf("got (%v, %v, %v), want 14, 20", lo, hi, err)
	}
}

func TestProteinRangeISSN_errors(t *testing.T) {
	f := func(kg float64) func() error {
		return func() error { _, _, err := ProteinRangeISSN(kg); return err }
	}
	runErrCases(t, []errCase{
		{"NaN", f(nan), ErrNotFinite},
		{"+Inf", f(pInf), ErrNotFinite},
		{"-Inf", f(nInf), ErrNotFinite},
		{"0", f(0), ErrOutOfRange},
		{"negative", f(-70), ErrOutOfRange},
		{"overflow", f(math.MaxFloat64), ErrDomain},
	})
}

func TestProteinRangeISSN_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _, _ = ProteinRangeISSN(70) }); n != 0 {
		t.Errorf("ProteinRangeISSN allocates %v times, want 0", n)
	}
}

func BenchmarkProteinRangeISSN(b *testing.B) {
	b.ReportAllocs()
	kg := 70.0
	for b.Loop() {
		sink, _, _ = ProteinRangeISSN(kg)
	}
}
