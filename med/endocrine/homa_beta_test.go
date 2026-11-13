package endocrine

import (
	"math"
	"testing"
)

// HAND-DERIVED from the formula as Singh 2010 (PMC3083884) and Lee 2014 (PMC3938726) state it:
// 20 × insulin / (glucose − 3.5). Tolerance 1e-12 relative.
func TestHOMABeta_example(t *testing.T) {
	tests := []struct {
		name     string
		ins, glu float64
		want     float64
	}{
		{"worked", 10, 5, 200 / 1.5},
		{"100 %", 5, 4.5, 100},
		{"just above the domain edge", 1, 3.5 + 1.0/1024, 20 * 1024},
		{"high glucose", 10, 13.5, 20},
	}
	for _, tc := range tests {
		if got, err := HOMABeta(tc.ins, tc.glu); err != nil || !relNear(got, tc.want, 1e-12) {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// parity: ledger bs_homa_b. The app converts glucose mg/dL / 18, then 20 × insulin / (G − 3.5), and returns
// nothing when G ≤ 3.5. Cases: (100 mg/dL, 10) → 97.297…; (180, 20) → 61.538…; (63 mg/dL = 3.5 mmol/L) → refused.
func TestHOMABeta_parity(t *testing.T) {
	cases := []struct{ gluMgDL, ins float64 }{{100, 10}, {180, 20}, {300, 100}}
	for i, c := range cases {
		want := 20 * c.ins / (c.gluMgDL/18 - 3.5) // parity: ledger bs_homa_b, source case i+1
		if got, err := HOMABeta(c.ins, c.gluMgDL/18); err != nil || !relNear(got, want, 1e-12) {
			t.Errorf("case %d: got (%v, %v), want %v", i+1, got, err, want)
		}
	}
	if _, err := HOMABeta(10, 3.5); err == nil { // parity: ledger bs_homa_b, source case 4 (app returns null)
		t.Error("glucose 3.5 mmol/L: want an error")
	}
}

func TestHOMABeta_errors(t *testing.T) {
	call := func(i, g float64) func() error {
		return func() error { _, err := HOMABeta(i, g); return err }
	}
	runErrCases(t, []errCase{
		{"insulin NaN", call(nan, 5), ErrNotFinite},
		{"glucose +Inf", call(10, pInf), ErrNotFinite},
		{"glucose -Inf", call(10, nInf), ErrNotFinite},
		{"NaN before range", call(0, nan), ErrNotFinite},
		{"insulin 0", call(0, 5), ErrOutOfRange},
		{"glucose 0", call(10, 0), ErrOutOfRange},
		{"glucose 3.5", call(10, 3.5), ErrDomain},
		{"glucose 3.0", call(10, 3), ErrDomain},
		{"overflow", call(math.MaxFloat64, 5), ErrDomain},
	})
}

func TestHOMABeta_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = HOMABeta(10, 5) }); n != 0 {
		t.Errorf("HOMABeta allocates %v times, want 0", n)
	}
}

func BenchmarkHOMABeta(b *testing.B) {
	b.ReportAllocs()
	ins, g := 10.0, 5.0
	for b.Loop() {
		sink, _ = HOMABeta(ins, g)
	}
}
