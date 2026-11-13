package endocrine

import (
	"math"
	"testing"
)

// HAND-DERIVED from QUICKI = 1/(log10 I + log10 G) (Katz 2000 abstract; units per Singh 2010 and Gutch 2015).
// Tolerance 1e-12 relative (log10 is exact on powers of ten in Go's implementation within 1 ulp).
func TestQUICKI_example(t *testing.T) {
	tests := []struct {
		name     string
		ins, glu float64
		want     float64
	}{
		{"worked", 10, 100, 1.0 / 3},
		{"insulin 1", 1, 100, 0.5},
		{"product 1000", 100, 10, 1.0 / 3},
		{"typical", 5, 90, 1 / (math.Log10(5) + math.Log10(90))},
	}
	for _, tc := range tests {
		if got, err := QUICKI(tc.ins, tc.glu); err != nil || !relNear(got, tc.want, 1e-12) {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// parity: ledger bs_quicki. The app computes 1/(log10 insulin + log10 glucose) with glucose in mg/dL and
// returns nothing when the denominator ≤ 0.
func TestQUICKI_parity(t *testing.T) {
	cases := []struct{ glu, ins float64 }{{100, 10}, {90, 5}, {300, 100}, {40, 1}}
	for i, c := range cases {
		want := 1 / (math.Log10(c.ins) + math.Log10(c.glu)) // parity: ledger bs_quicki, source case i+1
		if got, err := QUICKI(c.ins, c.glu); err != nil || !relNear(got, want, 1e-12) {
			t.Errorf("case %d: got (%v, %v), want %v", i+1, got, err, want)
		}
	}
}

func TestQUICKI_errors(t *testing.T) {
	call := func(i, g float64) func() error {
		return func() error { _, err := QUICKI(i, g); return err }
	}
	runErrCases(t, []errCase{
		{"insulin NaN", call(nan, 100), ErrNotFinite},
		{"glucose +Inf", call(10, pInf), ErrNotFinite},
		{"glucose -Inf", call(10, nInf), ErrNotFinite},
		{"NaN before range", call(0, nan), ErrNotFinite},
		{"insulin 0", call(0, 100), ErrOutOfRange},
		{"glucose negative", call(10, -5), ErrOutOfRange},
		{"product 1", call(0.5, 2), ErrDomain},
		{"product below 1", call(0.1, 5), ErrDomain},
	})
	if got, err := QUICKI(math.MaxFloat64, math.MaxFloat64); err != nil || !(got > 0) {
		t.Errorf("huge inputs: got (%v, %v), want a small positive value", got, err)
	}
}

func TestQUICKI_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = QUICKI(10, 100) }); n != 0 {
		t.Errorf("QUICKI allocates %v times, want 0", n)
	}
}

func BenchmarkQUICKI(b *testing.B) {
	b.ReportAllocs()
	ins, g := 10.0, 100.0
	for b.Loop() {
		sink, _ = QUICKI(ins, g)
	}
}
