package anthropometry

import "testing"

// Category edges from the CDC Adult BMI Categories table (read 2026-11-12): less than 18.5; 18.5 to less than
// 25; 25 to less than 30; 30 to less than 35; 35 to less than 40; 40 or greater. Each edge on both sides.
// The labels are written out here, verbatim from the CDC page, so a typo in a constant fails the test.
func TestBMICategoryCDCAdult_edges(t *testing.T) {
	tests := []struct {
		bmi  float64
		want string
	}{
		{5e-324, "Underweight"},
		{18.4999, "Underweight"},
		{18.5, "Healthy Weight"},
		{24.9999, "Healthy Weight"},
		{25, "Overweight"},
		{29.9999, "Overweight"},
		{30, "Class 1 Obesity"},
		{34.9999, "Class 1 Obesity"},
		{35, "Class 2 Obesity"},
		{39.9999, "Class 2 Obesity"},
		{40, "Class 3 Obesity (Severe Obesity)"},
		{1e308, "Class 3 Obesity (Severe Obesity)"},
	}
	for _, tc := range tests {
		if got, err := BMICategoryCDCAdult(tc.bmi); err != nil || string(got) != tc.want {
			t.Errorf("BMICategoryCDCAdult(%v) = (%q, %v), want %q", tc.bmi, got, err, tc.want)
		}
	}
}

func TestBMICategoryCDCAdult_errors(t *testing.T) {
	call := func(bmi float64) func() error {
		return func() error { _, err := BMICategoryCDCAdult(bmi); return err }
	}
	runErrCases(t, []errCase{
		{"NaN", call(nan), ErrNotFinite},
		{"+Inf", call(pInf), ErrNotFinite},
		{"-Inf", call(nInf), ErrNotFinite},
		{"0", call(0), ErrOutOfRange},
		{"negative", call(-18.5), ErrOutOfRange},
	})
}

func TestBMICategoryCDCAdult_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkCategory, _ = BMICategoryCDCAdult(27) }); n != 0 {
		t.Errorf("BMICategoryCDCAdult allocates %v times, want 0", n)
	}
}

func BenchmarkBMICategoryCDCAdult(b *testing.B) {
	b.ReportAllocs()
	bmi := 27.0
	for b.Loop() {
		sinkCategory, _ = BMICategoryCDCAdult(bmi)
	}
}
