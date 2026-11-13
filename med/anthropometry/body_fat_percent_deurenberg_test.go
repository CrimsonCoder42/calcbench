package anthropometry

import "testing"

// Worked examples, HAND-DERIVED from the Deurenberg et al. 1991 abstract, adult formula
// BF% = 1.20 × BMI + 0.23 × age − 10.8 × sex − 5.4 (sex 1 male, 0 female):
// male 40 y, BMI 25 → 30 + 9.2 − 10.8 − 5.4 = 23.0; female 40 y, BMI 25 → 30 + 9.2 − 5.4 = 33.8;
// female 83 y, BMI 40.9 → 49.08 + 19.09 − 5.4 = 62.77.
// Tolerance 1e-13 relative: four terms, each rounded once, summed.
func TestBodyFatPercentDeurenbergAdult_example(t *testing.T) {
	tests := []struct {
		bmi, age float64
		male     bool
		want     float64
	}{
		{25, 40, true, 23.0},
		{25, 40, false, 33.8},
		{40.9, 83, false, 62.77},
	}
	for _, tc := range tests {
		if got, err := BodyFatPercentDeurenbergAdult(tc.bmi, tc.age, tc.male); err != nil || !relNear(got, tc.want, 1e-13) {
			t.Errorf("BodyFatPercentDeurenbergAdult(%v, %v, %v) = (%v, %v), want %v", tc.bmi, tc.age, tc.male, got, err, tc.want)
		}
	}
}

// Worked examples, HAND-DERIVED from the same abstract, child formula
// BF% = 1.51 × BMI − 0.70 × age − 3.6 × sex + 1.4:
// female 10 y, BMI 18 → 27.18 − 7 + 1.4 = 21.58; male 10 y, BMI 18 → 21.58 − 3.6 = 17.98;
// male 7 y, BMI 13.9 → 20.989 − 4.9 − 3.6 + 1.4 = 13.889.
func TestBodyFatPercentDeurenbergChild_example(t *testing.T) {
	tests := []struct {
		bmi, age float64
		male     bool
		want     float64
	}{
		{18, 10, false, 21.58},
		{18, 10, true, 17.98},
		{13.9, 7, true, 13.889},
	}
	for _, tc := range tests {
		if got, err := BodyFatPercentDeurenbergChild(tc.bmi, tc.age, tc.male); err != nil || !relNear(got, tc.want, 1e-13) {
			t.Errorf("BodyFatPercentDeurenbergChild(%v, %v, %v) = (%v, %v), want %v", tc.bmi, tc.age, tc.male, got, err, tc.want)
		}
	}
}

// The age edge between the formulas: 15 is the child formula's, 15.0001 the adult's; never both, never neither.
func TestBodyFatPercentDeurenberg_ageHandOver(t *testing.T) {
	if _, err := BodyFatPercentDeurenbergChild(20, 15, true); err != nil {
		t.Errorf("child at 15: err = %v, want nil", err)
	}
	if _, err := BodyFatPercentDeurenbergAdult(20, 15, true); err == nil {
		t.Errorf("adult at 15: err = nil, want ErrOutOfRange")
	}
	if _, err := BodyFatPercentDeurenbergAdult(20, 15.0001, true); err != nil {
		t.Errorf("adult at 15.0001: err = %v, want nil", err)
	}
	if _, err := BodyFatPercentDeurenbergChild(20, 15.0001, true); err == nil {
		t.Errorf("child at 15.0001: err = nil, want ErrOutOfRange")
	}
}

// Range edges, both sides: BMI 13.9 and 40.9 inside, 13.8999 and 40.9001 outside; adult age 83 inside,
// 83.0001 outside; child age 7 inside, 6.9999 outside.
func TestBodyFatPercentDeurenberg_range(t *testing.T) {
	for _, bmi := range []float64{13.9, 40.9} {
		if _, err := BodyFatPercentDeurenbergAdult(bmi, 40, false); err != nil {
			t.Errorf("adult BMI %v: err = %v, want nil", bmi, err)
		}
		if _, err := BodyFatPercentDeurenbergChild(bmi, 10, false); err != nil {
			t.Errorf("child BMI %v: err = %v, want nil", bmi, err)
		}
	}
	if _, err := BodyFatPercentDeurenbergAdult(25, 83, true); err != nil {
		t.Errorf("adult age 83: err = %v, want nil", err)
	}
	if _, err := BodyFatPercentDeurenbergChild(18, 7, true); err != nil {
		t.Errorf("child age 7: err = %v, want nil", err)
	}
}

func TestBodyFatPercentDeurenbergAdult_errors(t *testing.T) {
	call := func(bmi, age float64) func() error {
		return func() error { _, err := BodyFatPercentDeurenbergAdult(bmi, age, true); return err }
	}
	runErrCases(t, []errCase{
		{"bmi NaN", call(nan, 40), ErrNotFinite},
		{"bmi +Inf", call(pInf, 40), ErrNotFinite},
		{"bmi -Inf", call(nInf, 40), ErrNotFinite},
		{"age NaN", call(25, nan), ErrNotFinite},
		{"age +Inf", call(25, pInf), ErrNotFinite},
		{"age -Inf", call(25, nInf), ErrNotFinite},
		{"bmi just below 13.9", call(13.8999, 40), ErrOutOfRange},
		{"bmi just above 40.9", call(40.9001, 40), ErrOutOfRange},
		{"age 15", call(25, 15), ErrOutOfRange},
		{"age just above 83", call(25, 83.0001), ErrOutOfRange},
	})
}

func TestBodyFatPercentDeurenbergChild_errors(t *testing.T) {
	call := func(bmi, age float64) func() error {
		return func() error { _, err := BodyFatPercentDeurenbergChild(bmi, age, false); return err }
	}
	runErrCases(t, []errCase{
		{"bmi NaN", call(nan, 10), ErrNotFinite},
		{"bmi +Inf", call(pInf, 10), ErrNotFinite},
		{"bmi -Inf", call(nInf, 10), ErrNotFinite},
		{"age NaN", call(18, nan), ErrNotFinite},
		{"age +Inf", call(18, pInf), ErrNotFinite},
		{"age -Inf", call(18, nInf), ErrNotFinite},
		{"bmi just below 13.9", call(13.8999, 10), ErrOutOfRange},
		{"bmi just above 40.9", call(40.9001, 10), ErrOutOfRange},
		{"age just below 7", call(18, 6.9999), ErrOutOfRange},
		{"age just above 15", call(18, 15.0001), ErrOutOfRange},
	})
}

func TestBodyFatPercentDeurenberg_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = BodyFatPercentDeurenbergAdult(25, 40, true) }); n != 0 {
		t.Errorf("BodyFatPercentDeurenbergAdult allocates %v times, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sink, _ = BodyFatPercentDeurenbergChild(18, 10, false) }); n != 0 {
		t.Errorf("BodyFatPercentDeurenbergChild allocates %v times, want 0", n)
	}
}

func BenchmarkBodyFatPercentDeurenbergAdult(b *testing.B) {
	b.ReportAllocs()
	bmi, age := 25.0, 40.0
	for b.Loop() {
		sink, _ = BodyFatPercentDeurenbergAdult(bmi, age, true)
	}
}

func BenchmarkBodyFatPercentDeurenbergChild(b *testing.B) {
	b.ReportAllocs()
	bmi, age := 18.0, 10.0
	for b.Loop() {
		sink, _ = BodyFatPercentDeurenbergChild(bmi, age, false)
	}
}
