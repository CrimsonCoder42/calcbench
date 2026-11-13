package hepatic

import "testing"

// Worked examples, HAND-DERIVED from the formula quoted in Kim et al. 2021:
// MELDNa = MELD + 1.32 × (137 − Na) − 0.033 × MELD × (137 − Na), Na bounded to 125–137.
// Inputs 1 (MELD 6.43), Na 137 → 6.43. Na 125 → 6.43 + 15.84 − 0.033 × 6.43 × 12 = 19.72372.
// Na 120 is clamped to 125, Na 140 to 137. Tolerance 1e-13 relative.
func TestMELDNa_example(t *testing.T) {
	m := meldOf(2, 2, 1.5)
	tests := []struct {
		name              string
		bili, inr, cr, na float64
		want              float64
	}{
		{"inputs 1, Na 137", 1, 1, 1, 137, 6.43},
		{"inputs 1, Na 125", 1, 1, 1, 125, 6.43 + 1.32*12 - 0.033*6.43*12},
		{"Na 120 clamped to 125", 1, 1, 1, 120, 6.43 + 1.32*12 - 0.033*6.43*12},
		{"Na 124.999 clamped to 125", 1, 1, 1, 124.999, 6.43 + 1.32*12 - 0.033*6.43*12},
		{"Na 125.5 inside", 1, 1, 1, 125.5, 6.43 + 1.32*11.5 - 0.033*6.43*11.5},
		{"Na 136.5 inside", 1, 1, 1, 136.5, 6.43 + 1.32*0.5 - 0.033*6.43*0.5},
		{"Na 137.001 clamped to 137", 1, 1, 1, 137.001, 6.43},
		{"Na 140 clamped to 137", 1, 1, 1, 140, 6.43},
		{"bili 2, INR 1.5, Cr 2, Na 130", 2, 1.5, 2, 130, m + 1.32*7 - 0.033*m*7},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := MELDNa(tc.bili, tc.inr, tc.cr, tc.na, false)
			if err != nil || !relNear(got, tc.want, 1e-13) {
				t.Errorf("got (%v, %v), want %v", got, err, tc.want)
			}
		})
	}
	// Dialysis passes through to MELD: creatinine set to 4.0.
	m4 := meldOf(4, 1, 1)
	if got, err := MELDNa(1, 1, 1, 130, true); err != nil || !relNear(got, m4+1.32*7-0.033*m4*7, 1e-13) {
		t.Errorf("dialysis: got (%v, %v)", got, err)
	}
}

func TestMELDNa_errors(t *testing.T) {
	call := func(bili, inr, cr, na float64) func() error {
		return func() error { _, err := MELDNa(bili, inr, cr, na, false); return err }
	}
	runErrCases(t, []errCase{
		{"bilirubin NaN", call(nan, 1, 1, 137), ErrNotFinite},
		{"INR +Inf", call(1, pInf, 1, 137), ErrNotFinite},
		{"creatinine -Inf", call(1, 1, nInf, 137), ErrNotFinite},
		{"sodium NaN", call(1, 1, 1, nan), ErrNotFinite},
		{"sodium +Inf", call(1, 1, 1, pInf), ErrNotFinite},
		{"sodium -Inf", call(1, 1, 1, nInf), ErrNotFinite},
		{"NaN before range", call(0, 1, 1, nan), ErrNotFinite},
		{"bilirubin 0", call(0, 1, 1, 137), ErrOutOfRange},
		{"INR negative", call(1, -1, 1, 137), ErrOutOfRange},
		{"creatinine 0", call(1, 1, 0, 137), ErrOutOfRange},
		{"sodium 0", call(1, 1, 1, 0), ErrOutOfRange},
		{"sodium negative", call(1, 1, 1, -137), ErrOutOfRange},
	})
}

func TestMELDNa_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = MELDNa(2, 1.5, 2, 130, false) }); n != 0 {
		t.Errorf("MELDNa allocates %v times, want 0", n)
	}
}

func BenchmarkMELDNa(b *testing.B) {
	b.ReportAllocs()
	bili, inr, cr, na := 2.0, 1.5, 2.0, 130.0
	for b.Loop() {
		sink, _ = MELDNa(bili, inr, cr, na, false)
	}
}
