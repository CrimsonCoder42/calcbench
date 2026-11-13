package endocrine

import (
	"math"
	"testing"
)

// HAND-DERIVED from METS-IR = ln(2G + TG) × BMI / ln(HDL), mg/dL (Bello-Chavolla 2018 abstract; units per
// Zhang 2022 and Wu 2022). Tolerance 1e-12 relative.
func TestMETSIR_example(t *testing.T) {
	tests := []struct {
		name              string
		glu, tg, bmi, hdl float64
		want              float64
	}{
		{"worked", 100, 150, 28, 50, math.Log(350) * 28 / math.Log(50)},
		{"e-based", (math.E*math.E - 1) / 2, 1, 10, math.E, 20},
		{"HDL just above 1", 90, 120, 25, 1.5, math.Log(300) * 25 / math.Log(1.5)},
	}
	for _, tc := range tests {
		if got, err := METSIR(tc.glu, tc.tg, tc.bmi, tc.hdl); err != nil || !relNear(got, tc.want, 1e-12) {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// parity: ledger mets_ir. The app computes Math.log(2g + tg) × bmi / Math.log(hdl), mg/dL, and returns nothing
// when ln(hdl) ≤ 0. Case 1 is its default (100, 150, 28, 50) → 41.927…
func TestMETSIR_parity(t *testing.T) {
	cases := []struct{ g, tg, bmi, hdl float64 }{{100, 150, 28, 50}, {300, 1000, 60, 10}, {50, 30, 15, 120}}
	for i, c := range cases {
		want := math.Log(2*c.g+c.tg) * c.bmi / math.Log(c.hdl) // parity: ledger mets_ir, source case i+1
		if got, err := METSIR(c.g, c.tg, c.bmi, c.hdl); err != nil || !relNear(got, want, 1e-12) {
			t.Errorf("case %d: got (%v, %v), want %v", i+1, got, err, want)
		}
	}
}

func TestMETSIR_errors(t *testing.T) {
	call := func(g, tg, bmi, hdl float64) func() error {
		return func() error { _, err := METSIR(g, tg, bmi, hdl); return err }
	}
	runErrCases(t, []errCase{
		{"glucose NaN", call(nan, 150, 28, 50), ErrNotFinite},
		{"TG +Inf", call(100, pInf, 28, 50), ErrNotFinite},
		{"BMI -Inf", call(100, 150, nInf, 50), ErrNotFinite},
		{"HDL NaN", call(100, 150, 28, nan), ErrNotFinite},
		{"NaN before range", call(0, 150, 28, nan), ErrNotFinite},
		{"glucose 0", call(0, 150, 28, 50), ErrOutOfRange},
		{"TG negative", call(100, -1, 28, 50), ErrOutOfRange},
		{"BMI 0", call(100, 150, 0, 50), ErrOutOfRange},
		{"HDL 1", call(100, 150, 28, 1), ErrOutOfRange},
		{"HDL 0.5", call(100, 150, 28, 0.5), ErrOutOfRange},
		{"overflow", call(math.MaxFloat64, math.MaxFloat64, 28, 50), ErrDomain},
	})
}

func TestMETSIR_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = METSIR(100, 150, 28, 50) }); n != 0 {
		t.Errorf("METSIR allocates %v times, want 0", n)
	}
}

func BenchmarkMETSIR(b *testing.B) {
	b.ReportAllocs()
	g, tg, bmi, hdl := 100.0, 150.0, 28.0, 50.0
	for b.Loop() {
		sink, _ = METSIR(g, tg, bmi, hdl)
	}
}
