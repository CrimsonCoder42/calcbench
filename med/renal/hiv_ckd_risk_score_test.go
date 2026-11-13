package renal

import "testing"

// Scherzer 2014 Table 2 (PMC4188545): both sides of every printed edge. HAND-DERIVED.
func TestHIVCKDRiskScore_table2(t *testing.T) {
	type in struct {
		age, glu, sbp float64
		htn           bool
		tg            float64
		prot          bool
		cd4           float64
	}
	base := in{30, 100, 120, false, 150, false, 500}
	tests := []struct {
		name string
		mod  func(*in)
		want int
	}{
		{"zero", func(*in) {}, 0},
		{"age 39.9", func(v *in) { v.age = 39.9 }, 0},
		{"age 40", func(v *in) { v.age = 40 }, 2},
		{"age 49.9", func(v *in) { v.age = 49.9 }, 2},
		{"age 50", func(v *in) { v.age = 50 }, 4},
		{"age 59.9", func(v *in) { v.age = 59.9 }, 4},
		{"age 60", func(v *in) { v.age = 60 }, 6},
		{"age 90", func(v *in) { v.age = 90 }, 6},
		{"glucose 140", func(v *in) { v.glu = 140 }, 0},
		{"glucose 140.1", func(v *in) { v.glu = 140.1 }, 2},
		{"sbp 140", func(v *in) { v.sbp = 140 }, 0},
		{"sbp 141", func(v *in) { v.sbp = 141 }, 1},
		{"hypertension", func(v *in) { v.htn = true }, 2},
		{"tg 200", func(v *in) { v.tg = 200 }, 0},
		{"tg 201", func(v *in) { v.tg = 201 }, 1},
		{"proteinuria", func(v *in) { v.prot = true }, 2},
		{"cd4 200", func(v *in) { v.cd4 = 200 }, 0},
		{"cd4 199", func(v *in) { v.cd4 = 199 }, 1},
		{"maximum 15", func(v *in) { *v = in{70, 200, 160, true, 300, true, 100} }, 15},
	}
	for _, tc := range tests {
		v := base
		tc.mod(&v)
		if got, err := HIVCKDRiskScore(v.age, v.glu, v.sbp, v.htn, v.tg, v.prot, v.cd4); err != nil || got != tc.want {
			t.Errorf("%s: (%d, %v), want %d", tc.name, got, err, tc.want)
		}
	}
}

// Doc comment example: 4+2+1+2 = 9 (HAND-DERIVED).
func TestHIVCKDRiskScore_example(t *testing.T) {
	if got, err := HIVCKDRiskScore(55, 150, 145, true, 180, false, 350); err != nil || got != 9 {
		t.Errorf("got (%d, %v), want 9", got, err)
	}
}

// parity: ledger hiv_ckd_prediction, source case 1. The app gives 1 point per age decade (40–49 → 1, …, ≥70 → 4)
// and 1 point each for glucose, hypertension and proteinuria; the paper gives 0/2/4/6 and 2 points each. Age 45
// with glucose > 140, hypertension and proteinuria: app 1+1+1+1 = 4; paper (asserted) 2+2+2+2 = 8.
func TestHIVCKDRiskScore_parity(t *testing.T) {
	if got, err := HIVCKDRiskScore(45, 150, 120, true, 150, true, 500); err != nil || got != 8 {
		t.Errorf("got (%d, %v), want 8 (paper; app shows 4)", got, err)
	}
}

func TestHIVCKDRiskScore_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"age NaN", func() error { _, e := HIVCKDRiskScore(nan, 1, 1, false, 1, false, 1); return e }, ErrNotFinite},
		{"glu +Inf", func() error { _, e := HIVCKDRiskScore(30, pInf, 1, false, 1, false, 1); return e }, ErrNotFinite},
		{"sbp -Inf", func() error { _, e := HIVCKDRiskScore(30, 1, nInf, false, 1, false, 1); return e }, ErrNotFinite},
		{"tg NaN", func() error { _, e := HIVCKDRiskScore(30, 1, 1, false, nan, false, 1); return e }, ErrNotFinite},
		{"cd4 NaN", func() error { _, e := HIVCKDRiskScore(30, 1, 1, false, 1, false, nan); return e }, ErrNotFinite},
		{"age 18.9", func() error { _, e := HIVCKDRiskScore(18.9, 1, 1, false, 1, false, 1); return e }, ErrOutOfRange},
		{"age 90.1", func() error { _, e := HIVCKDRiskScore(90.1, 1, 1, false, 1, false, 1); return e }, ErrOutOfRange},
		{"cd4 -1", func() error { _, e := HIVCKDRiskScore(30, 1, 1, false, 1, false, -1); return e }, ErrOutOfRange},
	})
	if _, err := HIVCKDRiskScore(19, 1, 1, false, 1, false, 1); err != nil {
		t.Errorf("age 19: %v", err)
	}
}

func TestHIVCKDRiskScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkB6Int, _ = HIVCKDRiskScore(55, 150, 145, true, 180, false, 350) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkHIVCKDRiskScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkB6Int, _ = HIVCKDRiskScore(55, 150, 145, true, 180, false, 350)
	}
}
