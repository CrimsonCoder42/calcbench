package chemistry

import (
	"errors"
	"math"
	"testing"
)

// mhTGUpper and mhColumn130to159 are the test's own copies from Martin et al. 2013, Figure 2: the TG stratum
// upper bounds and the factors of the non-HDL-C 130–159 column, row by row, written out so a test does not
// lean on the table in the code under test.
var (
	mhTGUpper = []float64{49, 56, 61, 66, 71, 75, 79, 83, 87, 92, 96, 100, 105, 110, 115, 120,
		126, 132, 138, 146, 154, 163, 173, 185, 201, 220, 247, 292, 399, 13975}
	mhColumn130to159 = []float64{3.3, 3.7, 4.0, 4.1, 4.3, 4.4, 4.5, 4.6, 4.6, 4.7, 4.8, 4.8, 5.0, 5.0, 5.1,
		5.2, 5.3, 5.3, 5.4, 5.6, 5.7, 5.8, 5.9, 6.0, 6.2, 6.4, 6.6, 7.0, 7.5, 8.8}
)

// Worked example, HAND-DERIVED from Figure 2: TC 200, HDL 50, TG 150 → non-HDL 150 (column 130–159),
// TG 150 in row 147–154, factor 5.7: 150 − 150/5.7 = 123.68421052631578. Tolerance 1e-13 relative.
func TestLDLCholesterolMartinHopkins_example(t *testing.T) {
	if got, err := LDLCholesterolMartinHopkins(200, 50, 150); err != nil || !relNear(got, 150-150/5.7, 1e-13) {
		t.Errorf("TC 200, HDL 50, TG 150: got (%v, %v), want %v", got, err, 150-150/5.7)
	}
}

// Every TG row edge on both sides, column 130–159 (non-HDL 150): TG at a printed upper bound uses that row;
// TG half a unit above uses the next row (ASSUMPTION A-M5), as does the next printed integer.
func TestLDLCholesterolMartinHopkins_rowEdges(t *testing.T) {
	check := func(tg, factor float64) {
		t.Helper()
		want := 150 - tg/factor
		if got, err := LDLCholesterolMartinHopkins(200, 50, tg); err != nil || !relNear(got, want, 1e-13) {
			t.Errorf("TG %v: got (%v, %v), want %v (factor %v)", tg, got, err, want, factor)
		}
	}
	check(7, mhColumn130to159[0])
	for i, upper := range mhTGUpper {
		check(upper, mhColumn130to159[i])
		if i+1 < len(mhTGUpper) {
			check(upper+0.5, mhColumn130to159[i+1])
			check(upper+1, mhColumn130to159[i+1])
		}
	}
}

// Every non-HDL-C column edge on both sides, row TG 97–100 (factors 5.4 5.1 4.8 4.7 4.5 4.3), and the four
// corner cells of the table (3.5, 3.1, 11.9, 6.7).
func TestLDLCholesterolMartinHopkins_columnEdgesAndCorners(t *testing.T) {
	tests := []struct{ nonHDL, tg, factor float64 }{
		{99, 100, 5.4}, {99.5, 100, 5.1}, {129, 100, 5.1}, {129.5, 100, 4.8},
		{159, 100, 4.8}, {159.5, 100, 4.7}, {189, 100, 4.7}, {189.5, 100, 4.5},
		{219, 100, 4.5}, {219.5, 100, 4.3}, {1000, 100, 4.3},
		{50, 7, 3.5}, {300, 7, 3.1}, {50, 13975, 11.9}, {300, 13975, 6.7},
	}
	for _, tc := range tests {
		want := tc.nonHDL - tc.tg/tc.factor
		got, err := LDLCholesterolMartinHopkins(50+tc.nonHDL, 50, tc.tg)
		if err != nil || !relNear(got, want, 1e-12) {
			t.Errorf("non-HDL %v, TG %v: got (%v, %v), want %v (factor %v)", tc.nonHDL, tc.tg, got, err, want, tc.factor)
		}
	}
}

func TestLDLCholesterolMartinHopkins_errors(t *testing.T) {
	call := func(tc, hdl, tg float64) func() error {
		return func() error { _, err := LDLCholesterolMartinHopkins(tc, hdl, tg); return err }
	}
	runErrCases(t, []errCase{
		{"TC NaN", call(nan, 50, 150), ErrNotFinite},
		{"TC +Inf", call(pInf, 50, 150), ErrNotFinite},
		{"TC -Inf", call(nInf, 50, 150), ErrNotFinite},
		{"HDL NaN", call(200, nan, 150), ErrNotFinite},
		{"HDL +Inf", call(200, pInf, 150), ErrNotFinite},
		{"HDL -Inf", call(200, nInf, 150), ErrNotFinite},
		{"TG NaN", call(200, 50, nan), ErrNotFinite},
		{"TG +Inf", call(200, 50, pInf), ErrNotFinite},
		{"TG -Inf", call(200, 50, nInf), ErrNotFinite},
		{"NaN before range", call(0, 50, nan), ErrNotFinite},
		{"TC 0", call(0, 50, 150), ErrOutOfRange},
		{"HDL 0", call(200, 0, 150), ErrOutOfRange},
		{"TG 6.999", call(200, 50, 6.999), ErrOutOfRange},
		{"TG 13975.001", call(200, 50, 13975.001), ErrOutOfRange},
		{"TC equals HDL", call(50, 50, 150), ErrOutOfRange},
		{"TC below HDL", call(40, 50, 150), ErrOutOfRange},
	})
}

// The SI twin with the Figure 2 caption factors (× 0.0259 cholesterol, × 0.0113 TG). HAND-DERIVED: TC 5.18,
// HDL 1.295, TG 1.695 mmol/L are 200, 50, 150 mg/dL → (150 − 150/5.7) × 0.0259 = 3.2034210526315787.
// Agreement with the mg/dL form away from stratum edges; tolerance 1e-12 relative (two conversions each way).
func TestLDLCholesterolMartinHopkinsSI_agrees(t *testing.T) {
	if got, err := LDLCholesterolMartinHopkinsSI(5.18, 1.295, 1.695); err != nil || !relNear(got, 3.2034210526315787, 1e-12) {
		t.Errorf("5.18, 1.295, 1.695 mmol/L: got (%v, %v), want 3.2034210526315787", got, err)
	}
	for _, c := range []struct{ tc, hdl, tg float64 }{{180, 45, 90.3}, {240, 60, 210.5}, {300, 35, 500}, {130, 70, 20}} {
		mg, _ := LDLCholesterolMartinHopkins(c.tc, c.hdl, c.tg)
		got, err := LDLCholesterolMartinHopkinsSI(c.tc*0.0259, c.hdl*0.0259, c.tg*0.0113)
		if err != nil || !relNear(got, mg*0.0259, 1e-12) {
			t.Errorf("%v mg/dL: SI (%v, %v), mg/dL × 0.0259 = %v", c, got, err, mg*0.0259)
		}
	}
}

func TestLDLCholesterolMartinHopkinsSI_range(t *testing.T) {
	// TG 0.08 mmol/L ÷ 0.0113 = 7.08 mg/dL (in); 0.079 → 6.99 (out). 157.9 → 13973.5 (in); 158 → 13982.3 (out).
	for _, tg := range []float64{0.08, 157.9} {
		if _, err := LDLCholesterolMartinHopkinsSI(5, 1, tg); err != nil {
			t.Errorf("TG %v mmol/L: err = %v, want nil", tg, err)
		}
	}
	for _, tg := range []float64{0.079, 158} {
		if _, err := LDLCholesterolMartinHopkinsSI(5, 1, tg); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("TG %v mmol/L: err = %v, want ErrOutOfRange", tg, err)
		}
	}
}

func TestLDLCholesterolMartinHopkinsSI_errors(t *testing.T) {
	call := func(tc, hdl, tg float64) func() error {
		return func() error { _, err := LDLCholesterolMartinHopkinsSI(tc, hdl, tg); return err }
	}
	runErrCases(t, []errCase{
		{"TC NaN", call(nan, 1.3, 1.7), ErrNotFinite},
		{"HDL +Inf", call(5.2, pInf, 1.7), ErrNotFinite},
		{"TG -Inf", call(5.2, 1.3, nInf), ErrNotFinite},
		{"TG NaN", call(5.2, 1.3, nan), ErrNotFinite},
		{"TC 0", call(0, 1.3, 1.7), ErrOutOfRange},
		{"HDL negative", call(5.2, -1.3, 1.7), ErrOutOfRange},
		{"TC equals HDL", call(1.3, 1.3, 1.7), ErrOutOfRange},
		// 1e308 ÷ 0.0259 is beyond float64.
		{"result beyond float64", call(1e308, 1.3, 1.7), ErrDomain},
	})
}

func TestLDLCholesterolMartinHopkins_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = LDLCholesterolMartinHopkins(200, 50, 150) }); n != 0 {
		t.Errorf("LDLCholesterolMartinHopkins allocates %v times, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sink, _ = LDLCholesterolMartinHopkinsSI(5.18, 1.295, 1.695) }); n != 0 {
		t.Errorf("LDLCholesterolMartinHopkinsSI allocates %v times, want 0", n)
	}
}

// The worst case scans every TG stratum and every non-HDL-C stratum.
func BenchmarkLDLCholesterolMartinHopkins(b *testing.B) {
	b.ReportAllocs()
	tc, hdl, tg := 400.0, 50.0, 5000.0
	for b.Loop() {
		sink, _ = LDLCholesterolMartinHopkins(tc, hdl, tg)
	}
}

func BenchmarkLDLCholesterolMartinHopkinsSI(b *testing.B) {
	b.ReportAllocs()
	tc, hdl, tg := 10.36, 1.295, 56.5
	for b.Loop() {
		sink, _ = LDLCholesterolMartinHopkinsSI(tc, hdl, tg)
	}
}

// mhFigure2 is a second, independent transcription of Martin et al. 2013 (JAMA 310:2061), Figure 2, the
// 180-cell table of median TG:VLDL-C factors, read from the figure image (PMC4226221, nihms639194f2.jpg) on
// 2026-11-12 by the lead, separately from the builder's table. Rows: TG strata in mhTGUpper order; columns:
// non-HDL-C <100, 100–129, 130–159, 160–189, 190–219, ≥220 mg/dL.
var mhFigure2 = [30][6]float64{
	{3.5, 3.4, 3.3, 3.3, 3.2, 3.1},
	{4.0, 3.9, 3.7, 3.6, 3.6, 3.4},
	{4.3, 4.1, 4.0, 3.9, 3.8, 3.6},
	{4.5, 4.3, 4.1, 4.0, 3.9, 3.9},
	{4.7, 4.4, 4.3, 4.2, 4.1, 3.9},
	{4.8, 4.6, 4.4, 4.2, 4.2, 4.1},
	{4.9, 4.6, 4.5, 4.3, 4.3, 4.2},
	{5.0, 4.8, 4.6, 4.4, 4.3, 4.2},
	{5.1, 4.8, 4.6, 4.5, 4.4, 4.3},
	{5.2, 4.9, 4.7, 4.6, 4.4, 4.3},
	{5.3, 5.0, 4.8, 4.7, 4.5, 4.4},
	{5.4, 5.1, 4.8, 4.7, 4.5, 4.3},
	{5.5, 5.2, 5.0, 4.7, 4.6, 4.5},
	{5.6, 5.3, 5.0, 4.8, 4.6, 4.5},
	{5.7, 5.4, 5.1, 4.9, 4.7, 4.5},
	{5.8, 5.5, 5.2, 5.0, 4.8, 4.6},
	{6.0, 5.5, 5.3, 5.0, 4.8, 4.6},
	{6.1, 5.7, 5.3, 5.1, 4.9, 4.7},
	{6.2, 5.8, 5.4, 5.2, 5.0, 4.7},
	{6.3, 5.9, 5.6, 5.3, 5.0, 4.8},
	{6.5, 6.0, 5.7, 5.4, 5.1, 4.8},
	{6.7, 6.2, 5.8, 5.4, 5.2, 4.9},
	{6.8, 6.3, 5.9, 5.5, 5.3, 5.0},
	{7.0, 6.5, 6.0, 5.7, 5.4, 5.1},
	{7.3, 6.7, 6.2, 5.8, 5.5, 5.2},
	{7.6, 6.9, 6.4, 6.0, 5.6, 5.3},
	{8.0, 7.2, 6.6, 6.2, 5.9, 5.4},
	{8.5, 7.6, 7.0, 6.5, 6.1, 5.6},
	{9.5, 8.3, 7.5, 7.0, 6.5, 5.9},
	{11.9, 10.0, 8.8, 8.1, 7.5, 6.7},
}

// mhNonHDLMid is one non-HDL-C value inside each column stratum.
var mhNonHDLMid = [6]float64{80, 115, 145, 175, 205, 240}

// Every one of the 180 cells: pick TG at the stratum's upper bound and non-HDL inside the column, then check
// LDL = non-HDL − TG/factor with the factor from this independent copy. Tolerance 1e-12 relative (one division,
// one subtraction).
func TestLDLCholesterolMartinHopkins_all180Cells(t *testing.T) {
	const hdl = 50.0
	for r, tg := range mhTGUpper {
		for c, nonHDL := range mhNonHDLMid {
			want := nonHDL - tg/mhFigure2[r][c]
			got, err := LDLCholesterolMartinHopkins(nonHDL+hdl, hdl, tg)
			if err != nil || math.Abs(got-want) > 1e-12*math.Max(1, math.Abs(want)) {
				t.Errorf("row %d (TG ≤ %v), column %d: got (%v, %v), want %v (factor %v)", r, tg, c, got, err, want, mhFigure2[r][c])
			}
		}
	}
}
