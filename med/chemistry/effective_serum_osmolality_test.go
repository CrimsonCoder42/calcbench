package chemistry

import (
	"math"
	"testing"
)

// HAND-DERIVED from Kitabchi 2009 Table 1 footnote: 2 × Na + glucose/18. Tolerance 1e-12 relative (one
// division rounds).
func TestEffectiveSerumOsmolality_example(t *testing.T) {
	tests := []struct {
		name    string
		na, glu float64
		want    float64
	}{
		{"worked", 140, 180, 290},
		{"HHS threshold side (Table 1: > 320)", 150, 360, 320},
		{"HHS range", 150, 900, 350},
		{"low glucose", 135, 72, 274},
	}
	for _, tc := range tests {
		if got, err := EffectiveSerumOsmolality(tc.na, tc.glu); err != nil || !relNear(got, tc.want, 1e-12) {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// parity: ledger bs_osmolality, source case 1 (app defaults Na 140, glucose 200 → effective 280 + 200/18).
// The app's "total" adds BUN/2.8, which Kitabchi 2009 does not state (consumer defect in the ledger).
func TestEffectiveSerumOsmolality_parity(t *testing.T) {
	want := 280 + 200.0/18
	if got, err := EffectiveSerumOsmolality(140, 200); err != nil || !relNear(got, want, 1e-12) {
		t.Errorf("got (%v, %v), want %v", got, err, want)
	}
}

func TestEffectiveSerumOsmolality_errors(t *testing.T) {
	call := func(na, g float64) func() error {
		return func() error { _, err := EffectiveSerumOsmolality(na, g); return err }
	}
	runErrCases(t, []errCase{
		{"Na NaN", call(nan, 100), ErrNotFinite},
		{"glucose +Inf", call(140, pInf), ErrNotFinite},
		{"glucose -Inf", call(140, nInf), ErrNotFinite},
		{"NaN before range", call(0, nan), ErrNotFinite},
		{"Na 0", call(0, 100), ErrOutOfRange},
		{"glucose negative", call(140, -1), ErrOutOfRange},
		{"overflow", call(math.MaxFloat64, 100), ErrDomain},
	})
}

func TestEffectiveSerumOsmolality_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = EffectiveSerumOsmolality(140, 180) }); n != 0 {
		t.Errorf("EffectiveSerumOsmolality allocates %v times, want 0", n)
	}
}

func BenchmarkEffectiveSerumOsmolality(b *testing.B) {
	b.ReportAllocs()
	na, g := 140.0, 180.0
	for b.Loop() {
		sink, _ = EffectiveSerumOsmolality(na, g)
	}
}
