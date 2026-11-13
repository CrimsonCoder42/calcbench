package renal

import (
	"errors"
	"math"
	"testing"
)

// Shared by the package's tests.

// sink keeps benchmark results observable.
var (
	sink         float64
	sinkCategory GFRCategory
)

// relNear reports |got − want| ≤ tol·|want| (tol·1 when want is 0).
func relNear(got, want, tol float64) bool {
	return math.Abs(got-want) <= tol*math.Max(math.Abs(want), 1e-300)
}

// umolOf is the test's own mg/dL → µmol/L conversion (KDIGO 2024 p. S126: 88.4), written out so a test
// does not lean on the code under test.
func umolOf(mgDL float64) float64 { return mgDL * 88.4 }

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

// Worked examples, HAND-DERIVED from Inker et al. 2021, Table 2, eGFRcr(AS):
// female 50 y, Scr 0.7: Scr/κ = 1, so both power terms are 1: 142 × 0.9938^50 × 1.012 = 105.29760114914401.
// male 60 y, Scr 1.2: Scr/κ = 1.2/0.9 > 1: 142 × (1.2/0.9)^−1.200 × 0.9938^60 = 69.23112756748166.
// male 40 y, Scr 0.5: Scr/κ = 0.5/0.9 < 1: 142 × (0.5/0.9)^−0.302 × 0.9938^40.
// female 30 y, Scr 0.5: 142 × (0.5/0.7)^−0.241 × 0.9938^30 × 1.012.
// Tolerance 1e-12 relative: the same few correctly rounded pow and multiply operations.
func TestEGFRCKDEPI2021_example(t *testing.T) {
	tests := []struct {
		name     string
		scr, age float64
		female   bool
		want     float64
	}{
		{"female 50 y Scr 0.7", 0.7, 50, true, 105.29760114914401},
		{"male 60 y Scr 1.2", 1.2, 60, false, 69.23112756748166},
		{"male 40 y Scr 0.5", 0.5, 40, false, 142 * math.Pow(0.5/0.9, -0.302) * math.Pow(0.9938, 40)},
		{"female 30 y Scr 0.5", 0.5, 30, true, 142 * math.Pow(0.5/0.7, -0.241) * math.Pow(0.9938, 30) * 1.012},
		{"female 70 y Scr 2.0", 2.0, 70, true, 142 * math.Pow(2.0/0.7, -1.2) * math.Pow(0.9938, 70) * 1.012},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := EGFRCKDEPI2021(tc.scr, tc.age, tc.female)
			if err != nil || !relNear(got, tc.want, 1e-12) {
				t.Errorf("got (%v, %v), want %v", got, err, tc.want)
			}
		})
	}
}

// Range edges: age 18 is in, 17.999 is out; creatinine just above 0 is in.
func TestEGFRCKDEPI2021_range(t *testing.T) {
	if _, err := EGFRCKDEPI2021(1, 18, false); err != nil {
		t.Errorf("age 18: err = %v, want nil", err)
	}
	if _, err := EGFRCKDEPI2021(1, 17.999, false); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("age 17.999: err = %v, want ErrOutOfRange", err)
	}
	if _, err := EGFRCKDEPI2021(1e-9, 40, true); err != nil {
		t.Errorf("Scr 1e-9: err = %v, want nil", err)
	}
	// Large age underflows the 0.9938^age term to 0; that is a finite result, not an error.
	if got, err := EGFRCKDEPI2021(1, 1e6, false); err != nil || got != 0 {
		t.Errorf("age 1e6: got (%v, %v), want (0, nil)", got, err)
	}
}

func TestEGFRCKDEPI2021_errors(t *testing.T) {
	call := func(scr, age float64) func() error {
		return func() error { _, err := EGFRCKDEPI2021(scr, age, false); return err }
	}
	runErrCases(t, []errCase{
		{"Scr NaN", call(nan, 50), ErrNotFinite},
		{"Scr +Inf", call(pInf, 50), ErrNotFinite},
		{"Scr -Inf", call(nInf, 50), ErrNotFinite},
		{"age NaN", call(1, nan), ErrNotFinite},
		{"age +Inf", call(1, pInf), ErrNotFinite},
		{"age -Inf", call(1, nInf), ErrNotFinite},
		{"NaN before range", call(-1, nan), ErrNotFinite},
		{"Scr 0", call(0, 50), ErrOutOfRange},
		{"Scr negative", call(-1, 50), ErrOutOfRange},
		{"age 17", call(1, 17), ErrOutOfRange},
	})
}

// The SI twin agrees with the mg/dL form at Scr × 88.4 µmol/L; tolerance 1e-12 relative (one extra
// division and multiplication by 88.4).
func TestEGFRCKDEPI2021SI_agrees(t *testing.T) {
	for _, scr := range []float64{0.3, 0.7, 0.9, 1.2, 2.5, 8} {
		for _, female := range []bool{false, true} {
			want, _ := EGFRCKDEPI2021(scr, 55, female)
			got, err := EGFRCKDEPI2021SI(umolOf(scr), 55, female)
			if err != nil || !relNear(got, want, 1e-12) {
				t.Errorf("Scr %v female %v: SI (%v, %v), mg/dL %v", scr, female, got, err, want)
			}
		}
	}
	// HAND-DERIVED: 61.88 µmol/L = 0.7 mg/dL, female 50 y → 142 × 0.9938^50 × 1.012.
	if got, err := EGFRCKDEPI2021SI(61.88, 50, true); err != nil || !relNear(got, 105.29760114914401, 1e-12) {
		t.Errorf("61.88 µmol/L: got (%v, %v)", got, err)
	}
}

func TestEGFRCKDEPI2021SI_errors(t *testing.T) {
	call := func(scr, age float64) func() error {
		return func() error { _, err := EGFRCKDEPI2021SI(scr, age, true); return err }
	}
	runErrCases(t, []errCase{
		{"Scr NaN", call(nan, 50), ErrNotFinite},
		{"Scr +Inf", call(pInf, 50), ErrNotFinite},
		{"Scr -Inf", call(nInf, 50), ErrNotFinite},
		{"age NaN", call(80, nan), ErrNotFinite},
		{"age +Inf", call(80, pInf), ErrNotFinite},
		{"age -Inf", call(80, nInf), ErrNotFinite},
		{"Scr 0", call(0, 50), ErrOutOfRange},
		{"Scr negative", call(-5, 50), ErrOutOfRange},
		{"age 17.999", call(80, 17.999), ErrOutOfRange},
		// 5e-324 µmol/L ÷ 88.4 underflows to 0 mg/dL; 0^−0.241 is +Inf.
		{"Scr underflows", call(5e-324, 50), ErrDomain},
	})
	if _, err := EGFRCKDEPI2021SI(80, 18, true); err != nil {
		t.Errorf("age 18: err = %v, want nil", err)
	}
}

func TestEGFRCKDEPI2021_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = EGFRCKDEPI2021(1.2, 60, false) }); n != 0 {
		t.Errorf("EGFRCKDEPI2021 allocates %v times, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sink, _ = EGFRCKDEPI2021SI(106, 60, false) }); n != 0 {
		t.Errorf("EGFRCKDEPI2021SI allocates %v times, want 0", n)
	}
}

func BenchmarkEGFRCKDEPI2021(b *testing.B) {
	b.ReportAllocs()
	scr, age := 1.2, 60.0
	for b.Loop() {
		sink, _ = EGFRCKDEPI2021(scr, age, false)
	}
}

func BenchmarkEGFRCKDEPI2021SI(b *testing.B) {
	b.ReportAllocs()
	scr, age := 106.0, 60.0
	for b.Loop() {
		sink, _ = EGFRCKDEPI2021SI(scr, age, false)
	}
}
