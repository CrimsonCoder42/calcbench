package renal

import (
	"math"
	"testing"
)

// Pierce 2021 abstract (PMC9083470): K for the height/Scr equation runs from 35.7 at age 1 to 50.8 from 18 in
// males and from 33.1 to 41.4 in females. With height/Scr = 1 the eGFR equals K, so the abstract's end points
// are checked to its printed precision (0.05).
func TestEGFRCKiDU25Creatinine_abstractK(t *testing.T) {
	tests := []struct {
		age    float64
		female bool
		want   float64
	}{
		{1, false, 35.7}, {1, true, 33.1}, {18, false, 50.8}, {18, true, 41.4}, {25, false, 50.8}, {25, true, 41.4},
	}
	for _, tc := range tests {
		got, err := EGFRCKiDU25Creatinine(1, 1, tc.age, tc.female)
		if err != nil || math.Abs(got-tc.want) > 0.05 {
			t.Errorf("age %v female %v: (%v, %v), want %v ± 0.05", tc.age, tc.female, got, err, tc.want)
		}
	}
}

// Table 2 (PMC9083470), HAND-DERIVED: each age-dependent K at the interior points and both sides of 12 and 18.
// Tolerance 1e-9 relative (exact arithmetic of the printed constants).
func TestEGFRCKiDU25Creatinine_table2(t *testing.T) {
	tests := []struct {
		age    float64
		female bool
		want   float64
	}{
		{6, false, 39.0 * math.Pow(1.008, -6)},
		{11.99, false, 39.0 * math.Pow(1.008, -0.01)},
		{12, false, 39.0},
		{15, false, 39.0 * math.Pow(1.045, 3)},
		{17.99, false, 39.0 * math.Pow(1.045, 5.99)},
		{6, true, 36.1 * math.Pow(1.008, -6)},
		{12, true, 36.1},
		{15, true, 36.1 * math.Pow(1.023, 3)},
		{17.99, true, 36.1 * math.Pow(1.023, 5.99)},
	}
	for _, tc := range tests {
		got, err := EGFRCKiDU25Creatinine(1.2, 0.6, tc.age, tc.female)
		if err != nil || !relNear(got, tc.want*2, 1e-9) {
			t.Errorf("age %v female %v: (%v, %v), want %v", tc.age, tc.female, got, err, tc.want*2)
		}
	}
}

// Worked example in the doc comment: male aged 1, 0.75 m, Scr 0.3 → 39.0 × 1.008^−11 × 2.5 = 89.3180 (HAND-DERIVED).
func TestEGFRCKiDU25Creatinine_example(t *testing.T) {
	got, err := EGFRCKiDU25Creatinine(0.75, 0.3, 1, false)
	if err != nil || math.Abs(got-89.3180) > 1e-3 {
		t.Errorf("got (%v, %v), want 89.3180", got, err)
	}
}

// The SI twin divides by 88.4 once; both forms agree to 1e-12 relative.
func TestEGFRCKiDU25CreatinineSI_agrees(t *testing.T) {
	for _, age := range []float64{1, 5, 12, 14.5, 18, 25} {
		for _, female := range []bool{false, true} {
			a, err1 := EGFRCKiDU25Creatinine(1.3, 0.7, age, female)
			b, err2 := EGFRCKiDU25CreatinineSI(1.3, umolOf(0.7), age, female)
			if err1 != nil || err2 != nil || !relNear(a, b, 1e-12) {
				t.Errorf("age %v female %v: %v vs %v (%v, %v)", age, female, a, b, err1, err2)
			}
		}
	}
	// Doc example: female 20, 1.60 m, 70.72 µmol/L → 82.8.
	if got, err := EGFRCKiDU25CreatinineSI(1.60, 70.72, 20, true); err != nil || !relNear(got, 82.8, 1e-9) {
		t.Errorf("got (%v, %v), want 82.8", got, err)
	}
}

// Table 2, cystatin C: constants from 18, both sides of 12 (female) and 15 (male) and the age-18 values.
func TestEGFRCKiDU25CystatinC_table2(t *testing.T) {
	tests := []struct {
		age    float64
		female bool
		want   float64
	}{
		{20, false, 77.1}, {20, true, 68.3},
		{10, false, 87.2 * math.Pow(1.011, -5)},
		{15, false, 87.2},
		{16, false, 87.2 * 0.960},
		{17.99, false, 87.2 * math.Pow(0.960, 2.99)},
		{11, true, 79.9 * math.Pow(1.004, -1)},
		{12, true, 79.9},
		{13, true, 79.9 * 0.974},
	}
	for _, tc := range tests {
		got, err := EGFRCKiDU25CystatinC(1, tc.age, tc.female)
		if err != nil || !relNear(got, tc.want, 1e-9) {
			t.Errorf("age %v female %v: (%v, %v), want %v", tc.age, tc.female, got, err, tc.want)
		}
	}
	// K is continuous at 18 to the table's printed precision (0.1).
	k18m := 87.2 * math.Pow(0.960, 3)
	k18f := 79.9 * math.Pow(0.974, 6)
	if math.Abs(k18m-77.1) > 0.1 || math.Abs(k18f-68.3) > 0.1 {
		t.Errorf("continuity at 18: %v, %v", k18m, k18f)
	}
}

// parity: ledger ckid_u25_egfr, source has no _test.go; the app's K functions are the same Table 2 pieces
// (male 15 y, height 1.6 m, Scr 0.8 → 39.0 × 1.045^3 × 2 = 89.0157…).
func TestEGFRCKiDU25Creatinine_parity(t *testing.T) {
	want := 39.0 * math.Pow(1.045, 3) * 1.6 / 0.8
	if got, err := EGFRCKiDU25Creatinine(1.6, 0.8, 15, false); err != nil || !relNear(got, want, 1e-12) {
		t.Errorf("got (%v, %v), want %v", got, err, want)
	}
}

func TestEGFRCKiDU25_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"height NaN", func() error { _, e := EGFRCKiDU25Creatinine(nan, 1, 10, false); return e }, ErrNotFinite},
		{"scr +Inf", func() error { _, e := EGFRCKiDU25Creatinine(1, pInf, 10, false); return e }, ErrNotFinite},
		{"age -Inf", func() error { _, e := EGFRCKiDU25Creatinine(1, 1, nInf, false); return e }, ErrNotFinite},
		{"age 0.99", func() error { _, e := EGFRCKiDU25Creatinine(1, 1, 0.99, false); return e }, ErrOutOfRange},
		{"age 25.01", func() error { _, e := EGFRCKiDU25Creatinine(1, 1, 25.01, false); return e }, ErrOutOfRange},
		{"height 0", func() error { _, e := EGFRCKiDU25Creatinine(0, 1, 10, false); return e }, ErrOutOfRange},
		{"scr 0", func() error { _, e := EGFRCKiDU25Creatinine(1, 0, 10, false); return e }, ErrOutOfRange},
		{"overflow", func() error { _, e := EGFRCKiDU25Creatinine(math.MaxFloat64, 1e-300, 10, false); return e }, ErrDomain},
		{"SI NaN", func() error { _, e := EGFRCKiDU25CreatinineSI(1, nan, 10, false); return e }, ErrNotFinite},
		{"SI scr -1", func() error { _, e := EGFRCKiDU25CreatinineSI(1, -1, 10, false); return e }, ErrOutOfRange},
		{"cys NaN", func() error { _, e := EGFRCKiDU25CystatinC(nan, 10, false); return e }, ErrNotFinite},
		{"cys 0", func() error { _, e := EGFRCKiDU25CystatinC(0, 10, false); return e }, ErrOutOfRange},
		{"cys age 26", func() error { _, e := EGFRCKiDU25CystatinC(1, 26, false); return e }, ErrOutOfRange},
		{"cys overflow", func() error { _, e := EGFRCKiDU25CystatinC(1e-320, 10, false); return e }, ErrDomain},
	})
	// Edges accepted.
	for _, age := range []float64{1, 25} {
		if _, err := EGFRCKiDU25CystatinC(1, age, true); err != nil {
			t.Errorf("age %v: %v", age, err)
		}
	}
}

func TestEGFRCKiDU25_noAlloc(t *testing.T) {
	n := testing.AllocsPerRun(100, func() {
		sink, _ = EGFRCKiDU25Creatinine(1.3, 0.7, 14, true)
		sink, _ = EGFRCKiDU25CreatinineSI(1.3, 62, 14, true)
		sink, _ = EGFRCKiDU25CystatinC(1.1, 14, true)
	})
	if n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkEGFRCKiDU25Creatinine(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sink, _ = EGFRCKiDU25Creatinine(1.3, 0.7, 14, true)
	}
}

func BenchmarkEGFRCKiDU25CreatinineSI(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sink, _ = EGFRCKiDU25CreatinineSI(1.3, 62, 14, true)
	}
}

func BenchmarkEGFRCKiDU25CystatinC(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sink, _ = EGFRCKiDU25CystatinC(1.1, 14, true)
	}
}
