package logs

import (
	"errors"
	"math"
	"testing"
)

// relClose compares with relative tolerance 1e-12: each function is one exp2 or log plus two roundings
// (a few ulps); the tolerance only absorbs that.
func relClose(got, want float64) bool {
	return math.Abs(got-want) <= 1e-12*math.Max(1, math.Abs(want))
}

// Abramson, College Algebra 2e, OpenStax, 2021, §6.7, Example 3 "Finding the Age of a Bone": 20 % of the
// carbon-14 remains, half-life 5 730 years. Exactly t = 5 730·log₂ 5 = 13 304.65…; the book prints ≈ 13 301
// because it rounds k to −0.000121 first, so the exact value is the oracle and the book's within 0.03 %.
func TestExponentialDecayElapsed_boneExample(t *testing.T) {
	got, err := ExponentialDecayElapsed(1, 0.2, 5730)
	if err != nil || !relClose(got, 5730*math.Log2(5)) {
		t.Fatalf("ExponentialDecayElapsed(1, 0.2, 5730) = %v, %v", got, err)
	}
	if math.Abs(got-13301)/13301 > 3e-4 {
		t.Errorf("%v differs from the book's 13 301 by more than its rounding of k explains", got)
	}
}

// Abramson §6.7, Example 2: carbon-14, k = ln(0.5)/5 730.
func TestDecayConstantFromHalfLife_carbon14(t *testing.T) {
	got, err := DecayConstantFromHalfLife(5730)
	if err != nil || !relClose(got, -math.Log(0.5)/5730) {
		t.Errorf("DecayConstantFromHalfLife(5730) = %v, %v", got, err)
	}
}

// Hand computations: 100 with half-life 50 → 50 after 50, 25 after 100, 200 at −50; 100 → 10 over 50 gives
// half-life 50·ln 2/ln 10; the four solves are mutually inverse.
func TestExponentialDecay_cases(t *testing.T) {
	cases := []struct{ initial, elapsed, half, rem float64 }{
		{100, 50, 50, 50}, {100, 100, 50, 25}, {100, -50, 50, 200}, {100, 0, 7, 100}, {0, 10, 3, 0},
	}
	for _, c := range cases {
		if got, err := ExponentialDecayRemaining(c.initial, c.elapsed, c.half); err != nil || !relClose(got, c.rem) {
			t.Errorf("Remaining(%v, %v, %v) = %v, %v; want %v", c.initial, c.elapsed, c.half, got, err, c.rem)
		}
		if got, err := ExponentialDecayInitial(c.rem, c.elapsed, c.half); err != nil || !relClose(got, c.initial) {
			t.Errorf("Initial(%v, %v, %v) = %v, %v; want %v", c.rem, c.elapsed, c.half, got, err, c.initial)
		}
		if c.initial > 0 {
			if got, err := ExponentialDecayElapsed(c.initial, c.rem, c.half); err != nil || math.Abs(got-c.elapsed) > 1e-12*math.Max(1, c.half) {
				t.Errorf("Elapsed(%v, %v, %v) = %v, %v; want %v", c.initial, c.rem, c.half, got, err, c.elapsed)
			}
		}
	}
	if got, err := ExponentialDecayHalfLife(100, 10, 50); err != nil || !relClose(got, 50*math.Ln2/math.Ln10) {
		t.Errorf("HalfLife(100, 10, 50) = %v, %v", got, err)
	}
	if got, err := ExponentialDecayHalfLife(100, 25, 100); err != nil || !relClose(got, 50) {
		t.Errorf("HalfLife(100, 25, 100) = %v, %v", got, err)
	}
	// Ratio beyond the float64 range: ln(MaxFloat64) − ln(1e-300) still finite.
	if got, err := ExponentialDecayHalfLife(math.MaxFloat64, 1e-300, 1); err != nil || !relClose(got, math.Ln2/(math.Log(math.MaxFloat64)-math.Log(1e-300))) {
		t.Errorf("HalfLife huge ratio = %v, %v", got, err)
	}
	if got, err := MeanLifetimeFromHalfLife(math.Ln2); err != nil || !relClose(got, 1) {
		t.Errorf("MeanLifetimeFromHalfLife(ln 2) = %v, %v", got, err)
	}
	if got, err := MeanLifetimeFromHalfLife(50); err != nil || !relClose(got, 50/math.Ln2) {
		t.Errorf("MeanLifetimeFromHalfLife(50) = %v, %v", got, err)
	}
}

func TestExponentialDecay_errors(t *testing.T) {
	nan := math.NaN()
	e := func(_ float64, err error) error { return err }
	tests := []struct {
		name      string
		err, want error
	}{
		{"Remaining NaN", e(ExponentialDecayRemaining(nan, 1, 1)), ErrNotFinite},
		{"Remaining halfLife 0", e(ExponentialDecayRemaining(1, 1, 0)), ErrOutOfRange},
		{"Remaining overflow", e(ExponentialDecayRemaining(1e300, -1e4, 1)), ErrOverflow},
		{"Initial Inf", e(ExponentialDecayInitial(1, math.Inf(1), 1)), ErrNotFinite},
		{"Initial halfLife negative", e(ExponentialDecayInitial(1, 1, -1)), ErrOutOfRange},
		{"Initial overflow", e(ExponentialDecayInitial(1e300, 1e4, 1)), ErrOverflow},
		{"HalfLife growth", e(ExponentialDecayHalfLife(10, 20, 1)), ErrOutOfRange},
		{"HalfLife equal", e(ExponentialDecayHalfLife(10, 10, 1)), ErrOutOfRange},
		{"HalfLife zero remaining", e(ExponentialDecayHalfLife(10, 0, 1)), ErrOutOfRange},
		{"HalfLife zero elapsed", e(ExponentialDecayHalfLife(10, 5, 0)), ErrOutOfRange},
		{"HalfLife NaN", e(ExponentialDecayHalfLife(10, 5, nan)), ErrNotFinite},
		{"HalfLife overflow", e(ExponentialDecayHalfLife(math.Nextafter(1, 2), 1, math.MaxFloat64)), ErrOverflow},
		{"Elapsed zero initial", e(ExponentialDecayElapsed(0, 1, 1)), ErrOutOfRange},
		{"Elapsed NaN", e(ExponentialDecayElapsed(1, nan, 1)), ErrNotFinite},
		{"Elapsed overflow", e(ExponentialDecayElapsed(math.MaxFloat64, 5e-324, math.MaxFloat64)), ErrOverflow},
		{"DecayConstant 0", e(DecayConstantFromHalfLife(0)), ErrOutOfRange},
		{"DecayConstant NaN", e(DecayConstantFromHalfLife(nan)), ErrNotFinite},
		{"DecayConstant overflow", e(DecayConstantFromHalfLife(5e-324)), ErrOverflow},
		{"MeanLifetime negative", e(MeanLifetimeFromHalfLife(-1)), ErrOutOfRange},
		{"MeanLifetime Inf", e(MeanLifetimeFromHalfLife(math.Inf(1))), ErrNotFinite},
		{"MeanLifetime overflow", e(MeanLifetimeFromHalfLife(math.MaxFloat64)), ErrOverflow},
	}
	for _, tc := range tests {
		if !errors.Is(tc.err, tc.want) {
			t.Errorf("%s: err = %v, want %v", tc.name, tc.err, tc.want)
		}
	}
}

func TestExponentialDecay_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() {
		_, _ = ExponentialDecayRemaining(100, 50, 50)
		_, _ = ExponentialDecayElapsed(1, 0.2, 5730)
	}); a != 0 {
		t.Errorf("allocates %v", a)
	}
}

func BenchmarkExponentialDecayRemaining(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = ExponentialDecayRemaining(100, 37, 5730)
	}
}

func BenchmarkExponentialDecayInitial(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = ExponentialDecayInitial(10, 37, 5730)
	}
}

func BenchmarkExponentialDecayHalfLife(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = ExponentialDecayHalfLife(100, 10, 50)
	}
}

func BenchmarkExponentialDecayElapsed(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = ExponentialDecayElapsed(1, 0.2, 5730)
	}
}

func BenchmarkDecayConstantFromHalfLife(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = DecayConstantFromHalfLife(5730)
	}
}

func BenchmarkMeanLifetimeFromHalfLife(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = MeanLifetimeFromHalfLife(5730)
	}
}
