package renal

import (
	"errors"
	"testing"
)

// Worked examples, HAND-DERIVED from Schwartz et al. 2009 (abstract, "0.413*(height/serum creatinine)"):
// 100 cm, Scr 0.413 → 100; 120 cm, Scr 0.6 → 0.413 × 120 / 0.6 = 82.6; 150 cm, Scr 1.0 → 61.95.
// Tolerance 1e-14 relative (two correctly rounded operations).
func TestEGFRSchwartzBedside_example(t *testing.T) {
	tests := []struct{ height, scr, age, want float64 }{
		{100, 0.413, 5, 100},
		{120, 0.6, 8, 82.6},
		{150, 1.0, 14, 61.95},
	}
	for _, tc := range tests {
		if got, err := EGFRSchwartzBedside(tc.height, tc.scr, tc.age); err != nil || !relNear(got, tc.want, 1e-14) {
			t.Errorf("EGFRSchwartzBedside(%v, %v, %v) = (%v, %v), want %v", tc.height, tc.scr, tc.age, got, err, tc.want)
		}
	}
}

func TestEGFRSchwartzBedside_range(t *testing.T) {
	for _, age := range []float64{1, 16} {
		if _, err := EGFRSchwartzBedside(100, 0.5, age); err != nil {
			t.Errorf("age %v: err = %v, want nil", age, err)
		}
	}
	for _, age := range []float64{0.999, 16.001} {
		if _, err := EGFRSchwartzBedside(100, 0.5, age); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("age %v: err = %v, want ErrOutOfRange", age, err)
		}
	}
}

func TestEGFRSchwartzBedside_errors(t *testing.T) {
	call := func(height, scr, age float64) func() error {
		return func() error { _, err := EGFRSchwartzBedside(height, scr, age); return err }
	}
	runErrCases(t, []errCase{
		{"height NaN", call(nan, 0.5, 5), ErrNotFinite},
		{"height +Inf", call(pInf, 0.5, 5), ErrNotFinite},
		{"height -Inf", call(nInf, 0.5, 5), ErrNotFinite},
		{"Scr NaN", call(100, nan, 5), ErrNotFinite},
		{"Scr +Inf", call(100, pInf, 5), ErrNotFinite},
		{"Scr -Inf", call(100, nInf, 5), ErrNotFinite},
		{"age NaN", call(100, 0.5, nan), ErrNotFinite},
		{"age +Inf", call(100, 0.5, pInf), ErrNotFinite},
		{"age -Inf", call(100, 0.5, nInf), ErrNotFinite},
		{"height 0", call(0, 0.5, 5), ErrOutOfRange},
		{"height negative", call(-1, 0.5, 5), ErrOutOfRange},
		{"Scr 0", call(100, 0, 5), ErrOutOfRange},
		{"age 0", call(100, 0.5, 0), ErrOutOfRange},
		{"age 17", call(100, 0.5, 17), ErrOutOfRange},
		{"result beyond float64", call(1e308, 1e-10, 5), ErrDomain},
	})
}

// The SI twin: 88.4 µmol/L = 1.0 mg/dL, 100 cm → 41.3 (HAND-DERIVED), and agreement at Scr × 88.4;
// tolerance 1e-14 relative.
func TestEGFRSchwartzBedsideSI_agrees(t *testing.T) {
	if got, err := EGFRSchwartzBedsideSI(100, 88.4, 10); err != nil || !relNear(got, 41.3, 1e-14) {
		t.Errorf("88.4 µmol/L: got (%v, %v), want 41.3", got, err)
	}
	for _, scr := range []float64{0.2, 0.413, 0.9, 3} {
		want, _ := EGFRSchwartzBedside(130, scr, 9)
		got, err := EGFRSchwartzBedsideSI(130, umolOf(scr), 9)
		if err != nil || !relNear(got, want, 1e-14) {
			t.Errorf("Scr %v: SI (%v, %v), mg/dL %v", scr, got, err, want)
		}
	}
}

func TestEGFRSchwartzBedsideSI_errors(t *testing.T) {
	call := func(height, scr, age float64) func() error {
		return func() error { _, err := EGFRSchwartzBedsideSI(height, scr, age); return err }
	}
	runErrCases(t, []errCase{
		{"height NaN", call(nan, 40, 5), ErrNotFinite},
		{"Scr +Inf", call(100, pInf, 5), ErrNotFinite},
		{"Scr NaN", call(100, nan, 5), ErrNotFinite},
		{"age -Inf", call(100, 40, nInf), ErrNotFinite},
		{"height 0", call(0, 40, 5), ErrOutOfRange},
		{"Scr 0", call(100, 0, 5), ErrOutOfRange},
		{"age 0.999", call(100, 40, 0.999), ErrOutOfRange},
		{"age 16.001", call(100, 40, 16.001), ErrOutOfRange},
		{"Scr underflows to 0 mg/dL", call(100, 5e-324, 5), ErrDomain},
	})
}

func TestEGFRSchwartzBedside_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = EGFRSchwartzBedside(120, 0.6, 8) }); n != 0 {
		t.Errorf("EGFRSchwartzBedside allocates %v times, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sink, _ = EGFRSchwartzBedsideSI(120, 53, 8) }); n != 0 {
		t.Errorf("EGFRSchwartzBedsideSI allocates %v times, want 0", n)
	}
}

func BenchmarkEGFRSchwartzBedside(b *testing.B) {
	b.ReportAllocs()
	height, scr, age := 120.0, 0.6, 8.0
	for b.Loop() {
		sink, _ = EGFRSchwartzBedside(height, scr, age)
	}
}

func BenchmarkEGFRSchwartzBedsideSI(b *testing.B) {
	b.ReportAllocs()
	height, scr, age := 120.0, 53.0, 8.0
	for b.Loop() {
		sink, _ = EGFRSchwartzBedsideSI(height, scr, age)
	}
}
