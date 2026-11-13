package hepatic

import (
	"math"
	"testing"
)

// meld3Of is the test's own copy of the MELD 3.0 equation of Kim et al. 2021, with the source's bounds
// (bilirubin, INR, creatinine ≥ 1; creatinine ≤ 3.0; sodium 125–137; albumin 1.5–3.5), not rounded.
func meld3Of(female bool, bili, na, inr, cr, alb float64) float64 {
	bili, inr = math.Max(bili, 1), math.Max(inr, 1)
	cr = math.Min(math.Max(cr, 1), 3)
	na = math.Min(math.Max(na, 125), 137)
	alb = math.Min(math.Max(alb, 1.5), 3.5)
	s := 4.56*math.Log(bili) + 0.82*(137-na) - 0.24*(137-na)*math.Log(bili) + 9.09*math.Log(inr) +
		11.14*math.Log(cr) + 1.85*(3.5-alb) - 1.83*(3.5-alb)*math.Log(cr) + 6
	if female {
		s += 1.33
	}
	return s
}

// Worked examples, HAND-DERIVED from the MELD 3.0 equation (Kim et al. 2021), rounded to the nearest integer.
// Unrounded values (independent computation): male baseline 6.0; female baseline 7.33; male bili 2, Na 130,
// INR 1.5, Cr 2, albumin 2.5 → 25.725141963728950; female same → 27.055141963728950; male bili 10, Na 125,
// INR 2.5, Cr 3, albumin 1.5 → 39.955045628203536; female bili 5, Na 133, INR 1.8, Cr 1.4, albumin 3.0 →
// 26.112385888174927. None is within 0.1 of a .5 tie.
func TestMELD3_example(t *testing.T) {
	tests := []struct {
		name                   string
		female                 bool
		bili, na, inr, cr, alb float64
		want                   int
	}{
		{"male, every value at its bound", false, 1, 137, 1, 1, 3.5, 6},
		{"female, every value at its bound", true, 1, 137, 1, 1, 3.5, 7},
		{"male mid values", false, 2, 130, 1.5, 2, 2.5, 26},
		{"female mid values", true, 2, 130, 1.5, 2, 2.5, 27},
		{"male at the far bounds", false, 10, 125, 2.5, 3, 1.5, 40},
		{"female mixed", true, 5, 133, 1.8, 1.4, 3.0, 26},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := MELD3(tc.female, tc.bili, tc.na, tc.inr, tc.cr, tc.alb)
			if err != nil || got != tc.want {
				t.Errorf("got (%v, %v), want %v", got, err, tc.want)
			}
			if r := int(math.Round(meld3Of(tc.female, tc.bili, tc.na, tc.inr, tc.cr, tc.alb))); r != tc.want {
				t.Errorf("test formula gives %v, table says %v", r, tc.want)
			}
		})
	}
}

// The bounds are clamps: a value past a bound gives the same score as the bound itself; a value just inside
// a bound gives the formula's own value. Each pair is checked against the test's formula.
func TestMELD3_bounds(t *testing.T) {
	type in struct{ bili, na, inr, cr, alb float64 }
	base := in{4, 131, 1.6, 2.2, 2.8}
	tests := []struct {
		name string
		v    in
	}{
		{"bilirubin 0.999 → 1", in{0.999, base.na, base.inr, base.cr, base.alb}},
		{"bilirubin 1", in{1, base.na, base.inr, base.cr, base.alb}},
		{"INR 0.999 → 1", in{base.bili, base.na, 0.999, base.cr, base.alb}},
		{"INR 1.001", in{base.bili, base.na, 1.001, base.cr, base.alb}},
		{"creatinine 0.999 → 1", in{base.bili, base.na, base.inr, 0.999, base.alb}},
		{"creatinine 2.999", in{base.bili, base.na, base.inr, 2.999, base.alb}},
		{"creatinine 3.001 → 3", in{base.bili, base.na, base.inr, 3.001, base.alb}},
		{"creatinine 8 → 3", in{base.bili, base.na, base.inr, 8, base.alb}},
		{"sodium 124.999 → 125", in{base.bili, 124.999, base.inr, base.cr, base.alb}},
		{"sodium 125.001", in{base.bili, 125.001, base.inr, base.cr, base.alb}},
		{"sodium 136.999", in{base.bili, 136.999, base.inr, base.cr, base.alb}},
		{"sodium 137.001 → 137", in{base.bili, 137.001, base.inr, base.cr, base.alb}},
		{"albumin 1.499 → 1.5", in{base.bili, base.na, base.inr, base.cr, 1.499}},
		{"albumin 1.501", in{base.bili, base.na, base.inr, base.cr, 1.501}},
		{"albumin 3.499", in{base.bili, base.na, base.inr, base.cr, 3.499}},
		{"albumin 3.501 → 3.5", in{base.bili, base.na, base.inr, base.cr, 3.501}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			want := int(math.Round(meld3Of(false, tc.v.bili, tc.v.na, tc.v.inr, tc.v.cr, tc.v.alb)))
			got, err := MELD3(false, tc.v.bili, tc.v.na, tc.v.inr, tc.v.cr, tc.v.alb)
			if err != nil || got != want {
				t.Errorf("got (%v, %v), want %v", got, err, want)
			}
		})
	}
	// Clamped: creatinine 8 scores the same as 3.0, sodium 100 as 125, albumin 0.5 as 1.5, albumin 6 as 3.5.
	same := func(a, b int, what string) {
		t.Helper()
		if a != b {
			t.Errorf("%s: %v ≠ %v", what, a, b)
		}
	}
	at := func(v in) int { s, _ := MELD3(true, v.bili, v.na, v.inr, v.cr, v.alb); return s }
	same(at(in{4, 131, 1.6, 8, 2.8}), at(in{4, 131, 1.6, 3, 2.8}), "creatinine 8 vs 3")
	same(at(in{4, 100, 1.6, 2.2, 2.8}), at(in{4, 125, 1.6, 2.2, 2.8}), "sodium 100 vs 125")
	same(at(in{4, 131, 1.6, 2.2, 0.5}), at(in{4, 131, 1.6, 2.2, 1.5}), "albumin 0.5 vs 1.5")
	same(at(in{4, 131, 1.6, 2.2, 6}), at(in{4, 131, 1.6, 2.2, 3.5}), "albumin 6 vs 3.5")
}

func TestMELD3_errors(t *testing.T) {
	call := func(bili, na, inr, cr, alb float64) func() error {
		return func() error { _, err := MELD3(false, bili, na, inr, cr, alb); return err }
	}
	runErrCases(t, []errCase{
		{"bilirubin NaN", call(nan, 137, 1, 1, 3.5), ErrNotFinite},
		{"bilirubin +Inf", call(pInf, 137, 1, 1, 3.5), ErrNotFinite},
		{"sodium NaN", call(1, nan, 1, 1, 3.5), ErrNotFinite},
		{"sodium -Inf", call(1, nInf, 1, 1, 3.5), ErrNotFinite},
		{"INR +Inf", call(1, 137, pInf, 1, 3.5), ErrNotFinite},
		{"INR NaN", call(1, 137, nan, 1, 3.5), ErrNotFinite},
		{"creatinine -Inf", call(1, 137, 1, nInf, 3.5), ErrNotFinite},
		{"creatinine NaN", call(1, 137, 1, nan, 3.5), ErrNotFinite},
		{"albumin NaN", call(1, 137, 1, 1, nan), ErrNotFinite},
		{"albumin +Inf", call(1, 137, 1, 1, pInf), ErrNotFinite},
		{"NaN before range", call(0, 137, 1, 1, nan), ErrNotFinite},
		{"bilirubin 0", call(0, 137, 1, 1, 3.5), ErrOutOfRange},
		{"sodium 0", call(1, 0, 1, 1, 3.5), ErrOutOfRange},
		{"INR negative", call(1, 137, -1, 1, 3.5), ErrOutOfRange},
		{"creatinine 0", call(1, 137, 1, 0, 3.5), ErrOutOfRange},
		{"albumin 0", call(1, 137, 1, 1, 0), ErrOutOfRange},
	})
	// Largest finite inputs: logs keep the score finite and well inside int.
	if _, err := MELD3(true, math.MaxFloat64, math.MaxFloat64, math.MaxFloat64, math.MaxFloat64, math.MaxFloat64); err != nil {
		t.Errorf("MaxFloat64 inputs: err = %v, want nil", err)
	}
}

func TestMELD3_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = MELD3(true, 2, 130, 1.5, 2, 2.5) }); n != 0 {
		t.Errorf("MELD3 allocates %v times, want 0", n)
	}
}

func BenchmarkMELD3(b *testing.B) {
	b.ReportAllocs()
	bili, na, inr, cr, alb := 2.0, 130.0, 1.5, 2.0, 2.5
	for b.Loop() {
		sinkInt, _ = MELD3(true, bili, na, inr, cr, alb)
	}
}
