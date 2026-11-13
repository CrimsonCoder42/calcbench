package renal

import (
	"errors"
	"testing"
)

// Worked examples, HAND-DERIVED from Cockcroft and Gault 1976 (abstract formula):
// male 40 y, 72 kg, Scr 1.0: (140 − 40) × 72 / (72 × 1.0) = 100 mL/min exactly.
// female, same: 100 × 0.85 = 85.
// male 92 y, 60 kg, Scr 2.0: 48 × 60 / 144 = 20.
// male 18 y, 90 kg, Scr 0.9: 122 × 90 / 64.8 = 169.44444444444446 (10980 / 64.8).
// Tolerance 1e-14 relative (three or four correctly rounded operations).
func TestCreatinineClearanceCockcroftGault_example(t *testing.T) {
	tests := []struct {
		name             string
		age, weight, scr float64
		female           bool
		want             float64
	}{
		{"male 40 y 72 kg Scr 1.0", 40, 72, 1.0, false, 100},
		{"female 40 y 72 kg Scr 1.0", 40, 72, 1.0, true, 85},
		{"male 92 y 60 kg Scr 2.0", 92, 60, 2.0, false, 20},
		{"male 18 y 90 kg Scr 0.9", 18, 90, 0.9, false, 10980.0 / 64.8},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := CreatinineClearanceCockcroftGault(tc.age, tc.weight, tc.scr, tc.female)
			if err != nil || !relNear(got, tc.want, 1e-14) {
				t.Errorf("got (%v, %v), want %v", got, err, tc.want)
			}
		})
	}
}

func TestCreatinineClearanceCockcroftGault_range(t *testing.T) {
	for _, age := range []float64{18, 92} {
		if _, err := CreatinineClearanceCockcroftGault(age, 70, 1, false); err != nil {
			t.Errorf("age %v: err = %v, want nil", age, err)
		}
	}
	for _, age := range []float64{17.999, 92.001} {
		if _, err := CreatinineClearanceCockcroftGault(age, 70, 1, false); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("age %v: err = %v, want ErrOutOfRange", age, err)
		}
	}
}

func TestCreatinineClearanceCockcroftGault_errors(t *testing.T) {
	call := func(age, weight, scr float64) func() error {
		return func() error { _, err := CreatinineClearanceCockcroftGault(age, weight, scr, false); return err }
	}
	runErrCases(t, []errCase{
		{"age NaN", call(nan, 70, 1), ErrNotFinite},
		{"age +Inf", call(pInf, 70, 1), ErrNotFinite},
		{"age -Inf", call(nInf, 70, 1), ErrNotFinite},
		{"weight NaN", call(40, nan, 1), ErrNotFinite},
		{"weight +Inf", call(40, pInf, 1), ErrNotFinite},
		{"weight -Inf", call(40, nInf, 1), ErrNotFinite},
		{"Scr NaN", call(40, 70, nan), ErrNotFinite},
		{"Scr +Inf", call(40, 70, pInf), ErrNotFinite},
		{"Scr -Inf", call(40, 70, nInf), ErrNotFinite},
		{"NaN before range", call(10, 70, nan), ErrNotFinite},
		{"weight 0", call(40, 0, 1), ErrOutOfRange},
		{"weight negative", call(40, -70, 1), ErrOutOfRange},
		{"Scr 0", call(40, 70, 0), ErrOutOfRange},
		{"Scr negative", call(40, 70, -1), ErrOutOfRange},
		{"result beyond float64", call(40, 1e308, 1e-300), ErrDomain},
	})
}

// The SI twin: 88.4 µmol/L is 1.0 mg/dL (spec test), and agreement at Scr × 88.4; tolerance 1e-14 relative.
func TestCreatinineClearanceCockcroftGaultSI_agrees(t *testing.T) {
	if got, err := CreatinineClearanceCockcroftGaultSI(40, 72, 88.4, false); err != nil || !relNear(got, 100, 1e-14) {
		t.Errorf("88.4 µmol/L: got (%v, %v), want 100", got, err)
	}
	for _, scr := range []float64{0.5, 1, 1.7, 6} {
		want, _ := CreatinineClearanceCockcroftGault(65, 80, scr, true)
		got, err := CreatinineClearanceCockcroftGaultSI(65, 80, umolOf(scr), true)
		if err != nil || !relNear(got, want, 1e-14) {
			t.Errorf("Scr %v: SI (%v, %v), mg/dL %v", scr, got, err, want)
		}
	}
}

func TestCreatinineClearanceCockcroftGaultSI_errors(t *testing.T) {
	call := func(age, weight, scr float64) func() error {
		return func() error { _, err := CreatinineClearanceCockcroftGaultSI(age, weight, scr, true); return err }
	}
	runErrCases(t, []errCase{
		{"age NaN", call(nan, 70, 88), ErrNotFinite},
		{"weight +Inf", call(40, pInf, 88), ErrNotFinite},
		{"Scr -Inf", call(40, 70, nInf), ErrNotFinite},
		{"Scr NaN", call(40, 70, nan), ErrNotFinite},
		{"age 17.999", call(17.999, 70, 88), ErrOutOfRange},
		{"age 92.001", call(92.001, 70, 88), ErrOutOfRange},
		{"weight 0", call(40, 0, 88), ErrOutOfRange},
		{"Scr 0", call(40, 70, 0), ErrOutOfRange},
		{"Scr underflows to 0 mg/dL", call(40, 70, 5e-324), ErrDomain},
	})
}

func TestCreatinineClearanceCockcroftGault_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = CreatinineClearanceCockcroftGault(40, 72, 1, true) }); n != 0 {
		t.Errorf("CreatinineClearanceCockcroftGault allocates %v times, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sink, _ = CreatinineClearanceCockcroftGaultSI(40, 72, 88.4, true) }); n != 0 {
		t.Errorf("CreatinineClearanceCockcroftGaultSI allocates %v times, want 0", n)
	}
}

func BenchmarkCreatinineClearanceCockcroftGault(b *testing.B) {
	b.ReportAllocs()
	age, weight, scr := 40.0, 72.0, 1.0
	for b.Loop() {
		sink, _ = CreatinineClearanceCockcroftGault(age, weight, scr, true)
	}
}

func BenchmarkCreatinineClearanceCockcroftGaultSI(b *testing.B) {
	b.ReportAllocs()
	age, weight, scr := 40.0, 72.0, 88.4
	for b.Loop() {
		sink, _ = CreatinineClearanceCockcroftGaultSI(age, weight, scr, true)
	}
}
