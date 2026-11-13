package endocrine

import (
	"math"
	"testing"
)

// HAND-DERIVED from the formula as Singh 2010 (PMC3083884) and Lee 2014 (PMC3938726) state it: insulin ×
// glucose (mmol/L) / 22.5. Tolerance 1e-12 relative (one multiplication and one division round).
func TestHOMAIR_example(t *testing.T) {
	tests := []struct {
		name     string
		ins, glu float64
		want     float64
	}{
		{"worked", 10, 5, 50.0 / 22.5},
		{"exactly 1", 4.5, 5, 1},
		{"exactly 10", 22.5, 10, 10},
		{"small", 1, 0.1, 0.1 / 22.5},
	}
	for _, tc := range tests {
		if got, err := HOMAIR(tc.ins, tc.glu); err != nil || !relNear(got, tc.want, 1e-12) {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// parity: ledger bs_homa_ir / homa_ir. The app computes glucose (mg/dL) × insulin / 405, i.e. glucose/18 in
// mmol/L over 22.5; the caller's mg/dL → mmol/L step uses the app's own 18. Cases: (100 mg/dL, 10) → 2.4691…;
// (90, 5) → 1.1111…; (300, 100) → 74.074…
func TestHOMAIR_parity(t *testing.T) {
	cases := []struct{ gluMgDL, ins float64 }{{100, 10}, {90, 5}, {300, 100}, {40, 1}}
	for i, c := range cases {
		want := c.gluMgDL * c.ins / 405 // parity: ledger bs_homa_ir, source case i+1
		if got, err := HOMAIR(c.ins, c.gluMgDL/18); err != nil || !relNear(got, want, 1e-12) {
			t.Errorf("case %d: got (%v, %v), want %v", i+1, got, err, want)
		}
	}
}

func TestHOMAIR_errors(t *testing.T) {
	call := func(i, g float64) func() error {
		return func() error { _, err := HOMAIR(i, g); return err }
	}
	runErrCases(t, []errCase{
		{"insulin NaN", call(nan, 5), ErrNotFinite},
		{"glucose +Inf", call(10, pInf), ErrNotFinite},
		{"glucose -Inf", call(10, nInf), ErrNotFinite},
		{"NaN before range", call(0, nan), ErrNotFinite},
		{"insulin 0", call(0, 5), ErrOutOfRange},
		{"glucose negative", call(10, -1), ErrOutOfRange},
		{"overflow", call(math.MaxFloat64, 100), ErrDomain},
	})
}

func TestHOMAIR_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = HOMAIR(10, 5) }); n != 0 {
		t.Errorf("HOMAIR allocates %v times, want 0", n)
	}
}

func BenchmarkHOMAIR(b *testing.B) {
	b.ReportAllocs()
	ins, g := 10.0, 5.0
	for b.Loop() {
		sink, _ = HOMAIR(ins, g)
	}
}
