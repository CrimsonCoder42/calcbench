package chemistry

import "testing"

// Worked examples, HAND-DERIVED from Figge et al. 1998 (abstract: adjusted AG = observed AG + 0.25 ×
// (normal albumin − observed albumin), albumin in g/L; the abstract gives no number for "normal albumin",
// the source's term, so the tests pass their own):
// AG 12, albumin 20, normalAlbuminGL 44 → 12 + 0.25 × 24 = 18.
// AG 8, albumin 44, normalAlbuminGL 44 → 8 (no adjustment).
// AG 10, albumin 50, normal 40 → 10 + 0.25 × (−10) = 7.5.
// AG −2, albumin 15, normal 45 → −2 + 7.5 = 5.5 (a negative observed gap is any finite value, Range).
// Tolerance 1e-14 absolute (exact binary fractions, two roundings at most).
func TestAdjustedAnionGapFigge_example(t *testing.T) {
	tests := []struct{ ag, alb, normal, want float64 }{
		{12, 20, 44, 18},
		{8, 44, 44, 8},
		{10, 50, 40, 7.5},
		{-2, 15, 45, 5.5},
	}
	for _, tc := range tests {
		if got, err := AdjustedAnionGapFigge(tc.ag, tc.alb, tc.normal); err != nil || !near(got, tc.want, 1e-14) {
			t.Errorf("AdjustedAnionGapFigge(%v, %v, %v) = (%v, %v), want %v", tc.ag, tc.alb, tc.normal, got, err, tc.want)
		}
	}
}

func TestAdjustedAnionGapFigge_errors(t *testing.T) {
	call := func(ag, alb, normal float64) func() error {
		return func() error { _, err := AdjustedAnionGapFigge(ag, alb, normal); return err }
	}
	runErrCases(t, []errCase{
		{"AG NaN", call(nan, 20, 44), ErrNotFinite},
		{"AG +Inf", call(pInf, 20, 44), ErrNotFinite},
		{"AG -Inf", call(nInf, 20, 44), ErrNotFinite},
		{"albumin NaN", call(12, nan, 44), ErrNotFinite},
		{"albumin +Inf", call(12, pInf, 44), ErrNotFinite},
		{"albumin -Inf", call(12, nInf, 44), ErrNotFinite},
		{"normal albumin NaN", call(12, 20, nan), ErrNotFinite},
		{"normal albumin +Inf", call(12, 20, pInf), ErrNotFinite},
		{"normal albumin -Inf", call(12, 20, nInf), ErrNotFinite},
		{"albumin 0", call(12, 0, 44), ErrOutOfRange},
		{"albumin negative", call(12, -20, 44), ErrOutOfRange},
		{"normal albumin 0", call(12, 20, 0), ErrOutOfRange},
		// 1.7e308 + 0.25 × 1e308 is beyond float64.
		{"result beyond float64", call(1.7e308, 1, 1e308), ErrDomain},
	})
}

// The g/dL twin: factor 2.5 ("if given in g/dL, the factor is 2.5"). HAND-DERIVED: AG 12, albumin 2.0,
// normal 4.4 → 12 + 2.5 × 2.4 = 18, the same as the g/L example. Agreement with the g/L form at albumin × 10;
// tolerance 1e-13 absolute (values below 100; 4.4 and 2.4 are not exact binary fractions).
func TestAdjustedAnionGapFiggeGDL_agrees(t *testing.T) {
	if got, err := AdjustedAnionGapFiggeGDL(12, 2.0, 4.4); err != nil || !near(got, 18, 1e-13) {
		t.Errorf("AG 12, 2.0, 4.4 g/dL: got (%v, %v), want 18", got, err)
	}
	for _, alb := range []float64{1.2, 2.5, 3.9, 5.1} {
		gl, _ := AdjustedAnionGapFigge(14, alb*10, 40)
		got, err := AdjustedAnionGapFiggeGDL(14, alb, 4.0)
		if err != nil || !near(got, gl, 1e-13) {
			t.Errorf("albumin %v g/dL: (%v, %v), g/L form %v", alb, got, err, gl)
		}
	}
}

func TestAdjustedAnionGapFiggeGDL_errors(t *testing.T) {
	call := func(ag, alb, normal float64) func() error {
		return func() error { _, err := AdjustedAnionGapFiggeGDL(ag, alb, normal); return err }
	}
	runErrCases(t, []errCase{
		{"AG NaN", call(nan, 2, 4.4), ErrNotFinite},
		{"albumin +Inf", call(12, pInf, 4.4), ErrNotFinite},
		{"normal albumin -Inf", call(12, 2, nInf), ErrNotFinite},
		{"normal albumin NaN", call(12, 2, nan), ErrNotFinite},
		{"albumin 0", call(12, 0, 4.4), ErrOutOfRange},
		{"normal albumin negative", call(12, 2, -4.4), ErrOutOfRange},
		// 2.5 × 1e308 is beyond float64.
		{"result beyond float64", call(0, 1, 1e308), ErrDomain},
	})
}

func TestAdjustedAnionGapFigge_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = AdjustedAnionGapFigge(12, 20, 44) }); n != 0 {
		t.Errorf("AdjustedAnionGapFigge allocates %v times, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sink, _ = AdjustedAnionGapFiggeGDL(12, 2, 4.4) }); n != 0 {
		t.Errorf("AdjustedAnionGapFiggeGDL allocates %v times, want 0", n)
	}
}

func BenchmarkAdjustedAnionGapFigge(b *testing.B) {
	b.ReportAllocs()
	ag, alb, normalAlbumin := 12.0, 20.0, 44.0
	for b.Loop() {
		sink, _ = AdjustedAnionGapFigge(ag, alb, normalAlbumin)
	}
}

func BenchmarkAdjustedAnionGapFiggeGDL(b *testing.B) {
	b.ReportAllocs()
	ag, alb, normalAlbumin := 12.0, 2.0, 4.4
	for b.Loop() {
		sink, _ = AdjustedAnionGapFiggeGDL(ag, alb, normalAlbumin)
	}
}
