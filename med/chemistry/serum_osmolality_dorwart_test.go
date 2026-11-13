package chemistry

import "testing"

// Worked examples, HAND-DERIVED from Dorwart and Chalmers 1975 (abstract: 1.86 Na + Glu/18 + BUN/2.8 + 9):
// Na 140, Glu 90, BUN 14 → 260.4 + 5 + 5 + 9 = 279.4.
// Na 130, Glu 540, BUN 56 → 241.8 + 30 + 20 + 9 = 300.8.
// Tolerance 1e-13 relative (four correctly rounded operations on values near 300).
func TestSerumOsmolalityDorwart_example(t *testing.T) {
	tests := []struct{ na, glu, bun, want float64 }{
		{140, 90, 14, 279.4},
		{130, 540, 56, 300.8},
	}
	for _, tc := range tests {
		if got, err := SerumOsmolalityDorwart(tc.na, tc.glu, tc.bun); err != nil || !relNear(got, tc.want, 1e-13) {
			t.Errorf("SerumOsmolalityDorwart(%v, %v, %v) = (%v, %v), want %v", tc.na, tc.glu, tc.bun, got, err, tc.want)
		}
	}
}

func TestSerumOsmolalityDorwart_errors(t *testing.T) {
	call := func(na, glu, bun float64) func() error {
		return func() error { _, err := SerumOsmolalityDorwart(na, glu, bun); return err }
	}
	runErrCases(t, []errCase{
		{"Na NaN", call(nan, 90, 14), ErrNotFinite},
		{"Na +Inf", call(pInf, 90, 14), ErrNotFinite},
		{"Na -Inf", call(nInf, 90, 14), ErrNotFinite},
		{"glucose NaN", call(140, nan, 14), ErrNotFinite},
		{"glucose +Inf", call(140, pInf, 14), ErrNotFinite},
		{"glucose -Inf", call(140, nInf, 14), ErrNotFinite},
		{"BUN NaN", call(140, 90, nan), ErrNotFinite},
		{"BUN +Inf", call(140, 90, pInf), ErrNotFinite},
		{"BUN -Inf", call(140, 90, nInf), ErrNotFinite},
		{"Na 0", call(0, 90, 14), ErrOutOfRange},
		{"glucose 0", call(140, 0, 14), ErrOutOfRange},
		{"BUN 0", call(140, 90, 0), ErrOutOfRange},
		{"BUN negative", call(140, 90, -1), ErrOutOfRange},
		// 1.86 × 1e308 is beyond float64.
		{"result beyond float64", call(1e308, 90, 14), ErrDomain},
	})
}

// The SI twin: glucose 90 mg/dL = 5 mmol/L and BUN 14 mg/dL = urea 5 mmol/L by the paper's own divisors
// (90/18, 14/2.8), so both forms give 279.4 (spec test). Agreement at glucose/18 and BUN/2.8 on more
// inputs; tolerance 1e-13 relative.
func TestSerumOsmolalityDorwartSI_agrees(t *testing.T) {
	if got, err := SerumOsmolalityDorwartSI(140, 5, 5); err != nil || !relNear(got, 279.4, 1e-13) {
		t.Errorf("Na 140, Glu 5, urea 5: got (%v, %v), want 279.4", got, err)
	}
	for _, glu := range []float64{70, 126, 600} {
		for _, bun := range []float64{7, 20, 90} {
			mg, _ := SerumOsmolalityDorwart(138, glu, bun)
			got, err := SerumOsmolalityDorwartSI(138, glu/18, bun/2.8)
			if err != nil || !relNear(got, mg, 1e-13) {
				t.Errorf("Glu %v BUN %v: SI (%v, %v), mg/dL %v", glu, bun, got, err, mg)
			}
		}
	}
}

func TestSerumOsmolalityDorwartSI_errors(t *testing.T) {
	call := func(na, glu, urea float64) func() error {
		return func() error { _, err := SerumOsmolalityDorwartSI(na, glu, urea); return err }
	}
	runErrCases(t, []errCase{
		{"Na NaN", call(nan, 5, 5), ErrNotFinite},
		{"glucose +Inf", call(140, pInf, 5), ErrNotFinite},
		{"urea -Inf", call(140, 5, nInf), ErrNotFinite},
		{"urea NaN", call(140, 5, nan), ErrNotFinite},
		{"Na negative", call(-140, 5, 5), ErrOutOfRange},
		{"glucose 0", call(140, 0, 5), ErrOutOfRange},
		{"urea 0", call(140, 5, 0), ErrOutOfRange},
		{"result beyond float64", call(1e308, 5, 5), ErrDomain},
	})
}

func TestSerumOsmolalityDorwart_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = SerumOsmolalityDorwart(140, 90, 14) }); n != 0 {
		t.Errorf("SerumOsmolalityDorwart allocates %v times, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sink, _ = SerumOsmolalityDorwartSI(140, 5, 5) }); n != 0 {
		t.Errorf("SerumOsmolalityDorwartSI allocates %v times, want 0", n)
	}
}

func BenchmarkSerumOsmolalityDorwart(b *testing.B) {
	b.ReportAllocs()
	na, glu, bun := 140.0, 90.0, 14.0
	for b.Loop() {
		sink, _ = SerumOsmolalityDorwart(na, glu, bun)
	}
}

func BenchmarkSerumOsmolalityDorwartSI(b *testing.B) {
	b.ReportAllocs()
	na, glu, urea := 140.0, 5.0, 5.0
	for b.Loop() {
		sink, _ = SerumOsmolalityDorwartSI(na, glu, urea)
	}
}
