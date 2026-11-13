package hepatic

import (
	"math"
	"testing"
)

// HAND-DERIVED from 4.6 × (PT − control) + bilirubin (Mitri 2023 PMC10494564; Dugum 2015 PMC4548353).
// Tolerance 1e-12 relative (4.6 is not exact in binary).
func TestMaddreyDiscriminantFunction_example(t *testing.T) {
	tests := []struct {
		name          string
		pt, ctl, bili float64
		want          float64
	}{
		{"worked", 20, 12, 10, 46.8},
		{"equal PT", 12, 12, 3.5, 3.5},
		{"near 32", 17, 12, 9, 32},
		{"PT shorter than control", 11, 12, 1, -3.6},
	}
	for _, tc := range tests {
		got, err := MaddreyDiscriminantFunction(tc.pt, tc.ctl, tc.bili)
		if err != nil || math.Abs(got-tc.want) > 1e-12*math.Max(1, math.Abs(tc.want)) {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// parity: ledger maddrey_discriminant. The app computes 4.6 × (pt − pt_control) + bilirubin; case 1 is its
// default (14, 12, 1.0) → 10.2.
func TestMaddreyDiscriminantFunction_parity(t *testing.T) {
	cases := []struct{ pt, ctl, bili, want float64 }{{14, 12, 1.0, 10.2}, {30, 12, 15, 97.8}, {8, 20, 0.1, -55.1}}
	for i, c := range cases {
		got, err := MaddreyDiscriminantFunction(c.pt, c.ctl, c.bili) // parity: ledger maddrey_discriminant, source case i+1
		if err != nil || math.Abs(got-c.want) > 1e-12*math.Max(1, math.Abs(c.want)) {
			t.Errorf("case %d: got (%v, %v), want %v", i+1, got, err, c.want)
		}
	}
}

func TestMaddreyDiscriminantFunction_errors(t *testing.T) {
	call := func(p, c, b float64) func() error {
		return func() error { _, err := MaddreyDiscriminantFunction(p, c, b); return err }
	}
	runErrCases(t, []errCase{
		{"PT NaN", call(nan, 12, 5), ErrNotFinite},
		{"control +Inf", call(14, pInf, 5), ErrNotFinite},
		{"bilirubin -Inf", call(14, 12, nInf), ErrNotFinite},
		{"NaN before range", call(0, 12, nan), ErrNotFinite},
		{"PT 0", call(0, 12, 5), ErrOutOfRange},
		{"control negative", call(14, -1, 5), ErrOutOfRange},
		{"bilirubin 0", call(14, 12, 0), ErrOutOfRange},
		{"overflow", call(math.MaxFloat64, 1, 5), ErrDomain},
	})
}

func TestMaddreyDiscriminantFunction_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = MaddreyDiscriminantFunction(20, 12, 10) }); n != 0 {
		t.Errorf("MaddreyDiscriminantFunction allocates %v times, want 0", n)
	}
}

func BenchmarkMaddreyDiscriminantFunction(b *testing.B) {
	b.ReportAllocs()
	pt, ctl, bili := 20.0, 12.0, 10.0
	for b.Loop() {
		sink, _ = MaddreyDiscriminantFunction(pt, ctl, bili)
	}
}
