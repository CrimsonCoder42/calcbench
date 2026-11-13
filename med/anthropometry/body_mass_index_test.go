package anthropometry

import (
	"errors"
	"math"
	"testing"
)

// Shared by the package's tests.

// sink keeps benchmark results observable.
var (
	sink         float64
	sinkCategory BMICategory
)

// relNear reports |got − want| ≤ tol·|want| (tol·1e-300 when want is 0).
func relNear(got, want, tol float64) bool {
	return math.Abs(got-want) <= tol*math.Max(math.Abs(want), 1e-300)
}

// errCase is one call that must fail with want.
type errCase struct {
	name string
	call func() error
	want error
}

func runErrCases(t *testing.T, cases []errCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.call(); !errors.Is(err, tc.want) {
				t.Errorf("err = %v, want %v", err, tc.want)
			}
		})
	}
}

var (
	nan  = math.NaN()
	pInf = math.Inf(1)
	nInf = math.Inf(-1)
)

// Worked example, HAND-DERIVED from the CDC formula "weight (kg)/[height (m)]2" (CDC BMI FAQ, read 2026-11-12):
// 70 kg, 1.75 m → 1.75² = 3.0625, 70 / 3.0625 = 22.857142857142858. 80 kg, 2 m → 80 / 4 = 20 exactly.
// Tolerance 1e-15 relative: one multiply and one divide, each correctly rounded.
func TestBodyMassIndex_example(t *testing.T) {
	tests := []struct{ weightKg, heightM, want float64 }{
		{70, 1.75, 22.857142857142858},
		{80, 2, 20},
		{50, 1, 50},
	}
	for _, tc := range tests {
		if got, err := BodyMassIndex(tc.weightKg, tc.heightM); err != nil || !relNear(got, tc.want, 1e-15) {
			t.Errorf("BodyMassIndex(%v, %v) = (%v, %v), want %v", tc.weightKg, tc.heightM, got, err, tc.want)
		}
	}
}

func TestBodyMassIndex_range(t *testing.T) {
	// Both sides of the > 0 edges: the smallest positive weight and a small height are accepted.
	if got, err := BodyMassIndex(5e-324, 1); err != nil || got != 5e-324 {
		t.Errorf("smallest weight: got (%v, %v), want 5e-324", got, err)
	}
	if got, err := BodyMassIndex(1, 1e-100); err != nil || !relNear(got, 1e200, 1e-15) {
		t.Errorf("height 1e-100: got (%v, %v), want 1e200", got, err)
	}
}

func TestBodyMassIndex_errors(t *testing.T) {
	call := func(w, h float64) func() error {
		return func() error { _, err := BodyMassIndex(w, h); return err }
	}
	runErrCases(t, []errCase{
		{"weight NaN", call(nan, 1.7), ErrNotFinite},
		{"weight +Inf", call(pInf, 1.7), ErrNotFinite},
		{"weight -Inf", call(nInf, 1.7), ErrNotFinite},
		{"height NaN", call(70, nan), ErrNotFinite},
		{"height +Inf", call(70, pInf), ErrNotFinite},
		{"height -Inf", call(70, nInf), ErrNotFinite},
		{"NaN before range", call(-1, nan), ErrNotFinite},
		{"weight 0", call(0, 1.7), ErrOutOfRange},
		{"weight negative", call(-70, 1.7), ErrOutOfRange},
		{"height 0", call(70, 0), ErrOutOfRange},
		{"height negative", call(70, -1.7), ErrOutOfRange},
		// 1e-200² underflows to 0, so the quotient is +Inf.
		{"height square underflows", call(70, 1e-200), ErrDomain},
		// 1e308 / 0.01 = 1e310, beyond float64.
		{"quotient overflows", call(1e308, 0.1), ErrDomain},
	})
}

func TestBodyMassIndex_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = BodyMassIndex(70, 1.75) }); n != 0 {
		t.Errorf("BodyMassIndex allocates %v times, want 0", n)
	}
}

func BenchmarkBodyMassIndex(b *testing.B) {
	b.ReportAllocs()
	w, h := 70.0, 1.75
	for b.Loop() {
		sink, _ = BodyMassIndex(w, h)
	}
}
