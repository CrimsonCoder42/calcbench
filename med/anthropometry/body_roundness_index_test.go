package anthropometry

import (
	"math"
	"testing"
)

// HAND-DERIVED from the Zhang 2024 form of Thomas 2013's BRI (values computed in float64 from
// 364.2 − 365.5·√(1 − (w/(πh))²)). Tolerance 1e-9 relative: the subtraction 364.2 − 365.5·ε cancels about
// two digits for typical BRI values.
func TestBodyRoundnessIndex_example(t *testing.T) {
	tests := []struct {
		name string
		w, h float64
		want float64
	}{
		{"worked", 90, 175, 3.630668966930216},
		{"rounder", 100, 160, 6.0060070019427485},
		{"60.96 cm / 177.8 cm", 60.96, 177.8, 0.8831471739942458},
	}
	for _, tc := range tests {
		if got, err := BodyRoundnessIndex(tc.w, tc.h); err != nil || !relNear(got, tc.want, 1e-9) {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
	// A circle-like waist just under π × height: ε → 0, BRI → 364.2.
	if got, err := BodyRoundnessIndex(math.Pi*100*(1-1e-12), 100); err != nil || math.Abs(got-364.2) > 1e-3 {
		t.Errorf("edge just inside: got (%v, %v), want ≈ 364.2", got, err)
	}
}

// The ratio has no unit: inches and cm agree (tolerance 1e-12 relative, one extra multiplication by 2.54).
func TestBodyRoundnessIndexImperial_agreesWithMetric(t *testing.T) {
	for _, c := range [][2]float64{{24, 70}, {36, 67}, {44, 64}, {70, 48}} {
		imp, err1 := BodyRoundnessIndexImperial(c[0], c[1])
		met, err2 := BodyRoundnessIndex(c[0]*2.54, c[1]*2.54)
		if err1 != nil || err2 != nil || !relNear(imp, met, 1e-12) {
			t.Errorf("%v: imperial (%v, %v), metric (%v, %v)", c, imp, err1, met, err2)
		}
	}
}

// parity: ledger body_roundness_index. The source test cases (waist, height in inches → BRI printed to 2
// decimals); every case is reproduced and compared after rounding to 2 decimals, as the source prints it.
func TestBodyRoundnessIndexImperial_parity(t *testing.T) {
	cases := []struct {
		w, h float64
		bri  float64
	}{
		{24, 70, 0.88},  // parity: ledger body_roundness_index, source case 1
		{30, 68, 2.32},  // source case 2
		{36, 67, 4.09},  // source case 3
		{38, 67, 4.71},  // source case 4
		{40, 66, 5.57},  // source case 5
		{44, 64, 7.56},  // source case 6
		{70, 48, 40.47}, // source case 7
	}
	for i, c := range cases {
		got, err := BodyRoundnessIndexImperial(c.w, c.h)
		if err != nil || math.Round(got*100)/100 != c.bri {
			t.Errorf("case %d: got (%v, %v), want %.2f", i+1, got, err, c.bri)
		}
	}
}

func TestBodyRoundnessIndex_errors(t *testing.T) {
	call := func(w, h float64) func() error {
		return func() error { _, err := BodyRoundnessIndex(w, h); return err }
	}
	callI := func(w, h float64) func() error {
		return func() error { _, err := BodyRoundnessIndexImperial(w, h); return err }
	}
	runErrCases(t, []errCase{
		{"waist NaN", call(nan, 170), ErrNotFinite},
		{"height +Inf", call(90, pInf), ErrNotFinite},
		{"waist -Inf", call(nInf, 170), ErrNotFinite},
		{"NaN before range", call(0, nan), ErrNotFinite},
		{"waist 0", call(0, 170), ErrOutOfRange},
		{"height negative", call(90, -1), ErrOutOfRange},
		{"waist = π × height", call(math.Pi*100, 100), ErrDomain},
		{"waist > π × height", call(400, 100), ErrDomain},
		{"imperial waist NaN", callI(nan, 70), ErrNotFinite},
		{"imperial height 0", callI(36, 0), ErrOutOfRange},
		{"imperial waist > π × height", callI(300, 70), ErrDomain},
	})
	if got, err := BodyRoundnessIndex(1, math.MaxFloat64); err != nil || !relNear(got, 364.2-365.5, 1e-12) {
		t.Errorf("huge height: got (%v, %v), want −1.3", got, err)
	}
}

func TestBodyRoundnessIndex_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = BodyRoundnessIndex(90, 175) }); n != 0 {
		t.Errorf("BodyRoundnessIndex allocates %v times, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sink, _ = BodyRoundnessIndexImperial(36, 67) }); n != 0 {
		t.Errorf("BodyRoundnessIndexImperial allocates %v times, want 0", n)
	}
}

func BenchmarkBodyRoundnessIndex(b *testing.B) {
	b.ReportAllocs()
	w, h := 90.0, 175.0
	for b.Loop() {
		sink, _ = BodyRoundnessIndex(w, h)
	}
}

func BenchmarkBodyRoundnessIndexImperial(b *testing.B) {
	b.ReportAllocs()
	w, h := 36.0, 67.0
	for b.Loop() {
		sink, _ = BodyRoundnessIndexImperial(w, h)
	}
}
