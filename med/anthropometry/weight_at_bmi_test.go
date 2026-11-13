package anthropometry

import (
	"math"
	"testing"
)

// Worked examples, HAND-DERIVED from the CDC BMI definition solved for weight:
// BMI 25, 1.778 m → 25 × 1.778² = 25 × 3.161284 = 79.0321 kg; BMI 18.5, 1.778 m → 58.483754 kg; BMI 20, 2 m → 80 kg.
// Tolerance 1e-15 relative: two multiplies, each correctly rounded.
func TestWeightAtBMI_example(t *testing.T) {
	tests := []struct{ bmi, heightM, want float64 }{
		{25, 1.778, 79.0321},
		{18.5, 1.778, 58.483754},
		{20, 2, 80},
	}
	for _, tc := range tests {
		if got, err := WeightAtBMI(tc.bmi, tc.heightM); err != nil || !relNear(got, tc.want, 1e-15) {
			t.Errorf("WeightAtBMI(%v, %v) = (%v, %v), want %v", tc.bmi, tc.heightM, got, err, tc.want)
		}
	}
}

// Worked examples, HAND-DERIVED from the CDC imperial form solved for weight:
// BMI 25, 70 in → 25 × 4900 / 703 = 174.25320056899005 lb; BMI 18.5, 70 in → 128.94736842105263 lb;
// BMI 703, 1 in → 1 lb exactly. Tolerance 1e-15 relative.
func TestWeightAtBMIImperial_example(t *testing.T) {
	tests := []struct{ bmi, heightIn, want float64 }{
		{25, 70, 174.25320056899005},
		{18.5, 70, 128.94736842105263},
		{703, 1, 1},
	}
	for _, tc := range tests {
		if got, err := WeightAtBMIImperial(tc.bmi, tc.heightIn); err != nil || !relNear(got, tc.want, 1e-15) {
			t.Errorf("WeightAtBMIImperial(%v, %v) = (%v, %v), want %v", tc.bmi, tc.heightIn, got, err, tc.want)
		}
	}
}

// Round trip through the forward functions: BMI(WeightAtBMI(b, h), h) = b. Tolerance 1e-14 relative
// (four roundings).
func TestWeightAtBMI_roundTrip(t *testing.T) {
	for _, b := range []float64{15, 18.5, 25, 30, 40} {
		for _, h := range []float64{1.5, 1.7, 1.9} {
			w, err := WeightAtBMI(b, h)
			got, err2 := BodyMassIndex(w, h)
			if err != nil || err2 != nil || !relNear(got, b, 1e-14) {
				t.Errorf("metric round trip b=%v h=%v: %v (%v, %v)", b, h, got, err, err2)
			}
		}
		for _, h := range []float64{60, 66, 75} {
			w, err := WeightAtBMIImperial(b, h)
			got, err2 := BodyMassIndexImperial(w, h)
			if err != nil || err2 != nil || !relNear(got, b, 1e-14) {
				t.Errorf("imperial round trip b=%v h=%v: %v (%v, %v)", b, h, got, err, err2)
			}
		}
	}
}

func TestWeightAtBMI_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"metric NaN bmi", func() error { _, err := WeightAtBMI(nan, 1.8); return err }, ErrNotFinite},
		{"metric +Inf height", func() error { _, err := WeightAtBMI(25, pInf); return err }, ErrNotFinite},
		{"metric zero bmi", func() error { _, err := WeightAtBMI(0, 1.8); return err }, ErrOutOfRange},
		{"metric negative height", func() error { _, err := WeightAtBMI(25, -1.8); return err }, ErrOutOfRange},
		{"metric overflow", func() error { _, err := WeightAtBMI(math.MaxFloat64, 2); return err }, ErrDomain},
		{"metric underflow", func() error { _, err := WeightAtBMI(1e-300, 1e-300); return err }, ErrDomain},
		{"imperial NaN height", func() error { _, err := WeightAtBMIImperial(25, nan); return err }, ErrNotFinite},
		{"imperial -Inf bmi", func() error { _, err := WeightAtBMIImperial(nInf, 70); return err }, ErrNotFinite},
		{"imperial zero height", func() error { _, err := WeightAtBMIImperial(25, 0); return err }, ErrOutOfRange},
		{"imperial negative bmi", func() error { _, err := WeightAtBMIImperial(-25, 70); return err }, ErrOutOfRange},
		{"imperial overflow", func() error { _, err := WeightAtBMIImperial(math.MaxFloat64, 70); return err }, ErrDomain},
		{"imperial underflow", func() error { _, err := WeightAtBMIImperial(1e-300, 1e-300); return err }, ErrDomain},
	})
}

func TestWeightAtBMI_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		sink, _ = WeightAtBMI(25, 1.778)
		sink, _ = WeightAtBMIImperial(25, 70)
	}); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkWeightAtBMI(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sink, _ = WeightAtBMI(25, 1.778)
	}
}

func BenchmarkWeightAtBMIImperial(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sink, _ = WeightAtBMIImperial(25, 70)
	}
}
