package renal

import (
	"errors"
	"math"
	"testing"
)

// Worked examples, HAND-DERIVED from Inker et al. 2021, Table 2, eGFRcr-cys(AS):
// female 50 y, Scr 0.7, Scys 0.8: both ratios are 1, so every power term is 1:
// 135 × 0.9961^50 × 0.963 = 106.9318450861178.
// male 60 y, Scr 1.2, Scys 1.0: 135 × (1.2/0.9)^−0.544 × (1.0/0.8)^−0.778 × 0.9961^60 = 76.76221740595214.
// The remaining rows write the four branches out. Tolerance 1e-12 relative (a few correctly rounded operations).
func TestEGFRCKDEPI2021CystatinC_example(t *testing.T) {
	tests := []struct {
		name          string
		scr, cys, age float64
		female        bool
		want          float64
	}{
		{"female 50 y Scr 0.7 Scys 0.8", 0.7, 0.8, 50, true, 106.9318450861178},
		{"male 60 y Scr 1.2 Scys 1.0", 1.2, 1.0, 60, false, 76.76221740595214},
		{"male both below kappa", 0.6, 0.6, 40, false,
			135 * math.Pow(0.6/0.9, -0.144) * math.Pow(0.6/0.8, -0.323) * math.Pow(0.9961, 40)},
		{"female both above kappa", 1.5, 1.6, 70, true,
			135 * math.Pow(1.5/0.7, -0.544) * math.Pow(1.6/0.8, -0.778) * math.Pow(0.9961, 70) * 0.963},
		{"female Scr below, Scys above", 0.5, 1.2, 30, true,
			135 * math.Pow(0.5/0.7, -0.219) * math.Pow(1.2/0.8, -0.778) * math.Pow(0.9961, 30) * 0.963},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := EGFRCKDEPI2021CystatinC(tc.scr, tc.cys, tc.age, tc.female)
			if err != nil || !relNear(got, tc.want, 1e-12) {
				t.Errorf("got (%v, %v), want %v", got, err, tc.want)
			}
		})
	}
}

func TestEGFRCKDEPI2021CystatinC_range(t *testing.T) {
	if _, err := EGFRCKDEPI2021CystatinC(1, 1, 18, false); err != nil {
		t.Errorf("age 18: err = %v, want nil", err)
	}
	if _, err := EGFRCKDEPI2021CystatinC(1, 1, 17.999, false); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("age 17.999: err = %v, want ErrOutOfRange", err)
	}
	// The smallest positive inputs still give a finite result (each power term is at most about e^240).
	if got, err := EGFRCKDEPI2021CystatinC(5e-324, 5e-324, 18, true); err != nil || math.IsInf(got, 0) {
		t.Errorf("smallest inputs: got (%v, %v), want finite", got, err)
	}
}

func TestEGFRCKDEPI2021CystatinC_errors(t *testing.T) {
	call := func(scr, cys, age float64) func() error {
		return func() error { _, err := EGFRCKDEPI2021CystatinC(scr, cys, age, false); return err }
	}
	runErrCases(t, []errCase{
		{"Scr NaN", call(nan, 1, 50), ErrNotFinite},
		{"Scr +Inf", call(pInf, 1, 50), ErrNotFinite},
		{"Scr -Inf", call(nInf, 1, 50), ErrNotFinite},
		{"Scys NaN", call(1, nan, 50), ErrNotFinite},
		{"Scys +Inf", call(1, pInf, 50), ErrNotFinite},
		{"Scys -Inf", call(1, nInf, 50), ErrNotFinite},
		{"age NaN", call(1, 1, nan), ErrNotFinite},
		{"age +Inf", call(1, 1, pInf), ErrNotFinite},
		{"age -Inf", call(1, 1, nInf), ErrNotFinite},
		{"Scr 0", call(0, 1, 50), ErrOutOfRange},
		{"Scys 0", call(1, 0, 50), ErrOutOfRange},
		{"Scys negative", call(1, -1, 50), ErrOutOfRange},
		{"age 17", call(1, 1, 17), ErrOutOfRange},
	})
}

// The SI twin agrees with the mg/dL form at Scr × 88.4 µmol/L; tolerance 1e-12 relative.
func TestEGFRCKDEPI2021CystatinCSI_agrees(t *testing.T) {
	for _, scr := range []float64{0.4, 0.7, 0.9, 1.3, 4} {
		for _, female := range []bool{false, true} {
			want, _ := EGFRCKDEPI2021CystatinC(scr, 1.1, 45, female)
			got, err := EGFRCKDEPI2021CystatinCSI(umolOf(scr), 1.1, 45, female)
			if err != nil || !relNear(got, want, 1e-12) {
				t.Errorf("Scr %v female %v: SI (%v, %v), mg/dL %v", scr, female, got, err, want)
			}
		}
	}
	// HAND-DERIVED: 61.88 µmol/L = 0.7 mg/dL, Scys 0.8, female 50 y → 135 × 0.9961^50 × 0.963.
	if got, err := EGFRCKDEPI2021CystatinCSI(61.88, 0.8, 50, true); err != nil || !relNear(got, 106.9318450861178, 1e-12) {
		t.Errorf("61.88 µmol/L: got (%v, %v)", got, err)
	}
}

func TestEGFRCKDEPI2021CystatinCSI_errors(t *testing.T) {
	call := func(scr, cys, age float64) func() error {
		return func() error { _, err := EGFRCKDEPI2021CystatinCSI(scr, cys, age, true); return err }
	}
	runErrCases(t, []errCase{
		{"Scr NaN", call(nan, 1, 50), ErrNotFinite},
		{"Scr +Inf", call(pInf, 1, 50), ErrNotFinite},
		{"Scys -Inf", call(80, nInf, 50), ErrNotFinite},
		{"Scys NaN", call(80, nan, 50), ErrNotFinite},
		{"age NaN", call(80, 1, nan), ErrNotFinite},
		{"age -Inf", call(80, 1, nInf), ErrNotFinite},
		{"Scr 0", call(0, 1, 50), ErrOutOfRange},
		{"Scys 0", call(80, 0, 50), ErrOutOfRange},
		{"age 17.999", call(80, 1, 17.999), ErrOutOfRange},
		{"Scr underflows to 0 mg/dL", call(5e-324, 1, 50), ErrDomain},
	})
}

func TestEGFRCKDEPI2021CystatinC_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = EGFRCKDEPI2021CystatinC(1.2, 1.0, 60, false) }); n != 0 {
		t.Errorf("EGFRCKDEPI2021CystatinC allocates %v times, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sink, _ = EGFRCKDEPI2021CystatinCSI(106, 1.0, 60, false) }); n != 0 {
		t.Errorf("EGFRCKDEPI2021CystatinCSI allocates %v times, want 0", n)
	}
}

func BenchmarkEGFRCKDEPI2021CystatinC(b *testing.B) {
	b.ReportAllocs()
	scr, cys, age := 1.2, 1.0, 60.0
	for b.Loop() {
		sink, _ = EGFRCKDEPI2021CystatinC(scr, cys, age, false)
	}
}

func BenchmarkEGFRCKDEPI2021CystatinCSI(b *testing.B) {
	b.ReportAllocs()
	scr, cys, age := 106.0, 1.0, 60.0
	for b.Loop() {
		sink, _ = EGFRCKDEPI2021CystatinCSI(scr, cys, age, false)
	}
}
