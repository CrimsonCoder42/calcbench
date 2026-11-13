package endocrine

import "testing"

const (
	bgOther = AUSDRISKBackgroundOther
	bgSEPI  = AUSDRISKBackgroundSouthernEuropeanOrPacificIslander
	bgAATSI = AUSDRISKBackgroundAsianOrAboriginalTorresStraitIslander
)

// HAND-DERIVED from Chen 2010 Box 3. Maximum: 8 + 3 + 2 + 3 + 6 + 2 + 2 + 2 + 7 = 35 (the source's maximum).
func TestAUSDRISKScore_example(t *testing.T) {
	tests := []struct {
		name string
		f    AUSDRISKFindings
		want int
	}{
		{"worked", AUSDRISKFindings{AgeYears: 47, Male: true, PhysicallyInactive: true, WaistCm: 105}, 13},
		{"minimum", AUSDRISKFindings{AgeYears: 25, WaistCm: 70}, 0},
		{"maximum", AUSDRISKFindings{70, true, bgAATSI, true, true, true, true, true, 120}, 35},
	}
	for _, tc := range tests {
		if got, err := AUSDRISKScore(tc.f); err != nil || got != tc.want {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// Both sides of every age edge and every waist edge (Box 3).
func TestAUSDRISKScore_edges(t *testing.T) {
	age := func(a float64) AUSDRISKFindings { return AUSDRISKFindings{AgeYears: a, WaistCm: 70} }
	waist := func(male bool, bg AUSDRISKBackground, cm float64) AUSDRISKFindings {
		return AUSDRISKFindings{AgeYears: 30, Male: male, Background: bg, WaistCm: cm}
	}
	tests := []struct {
		name string
		f    AUSDRISKFindings
		want int
	}{
		{"age 34.99", age(34.99), 0}, {"age 35", age(35), 2},
		{"age 44.99", age(44.99), 2}, {"age 45", age(45), 4},
		{"age 54.99", age(54.99), 4}, {"age 55", age(55), 6},
		{"age 64.99", age(64.99), 6}, {"age 65", age(65), 8},
		{"other man 101.9", waist(true, bgOther, 101.9), 3}, {"other man 102", waist(true, bgOther, 102), 7},
		{"other man 109.9", waist(true, bgOther, 109.9), 7}, {"other man 110", waist(true, bgOther, 110), 10},
		{"other woman 87.9", waist(false, bgOther, 87.9), 0}, {"other woman 88", waist(false, bgOther, 88), 4},
		{"other woman 99.9", waist(false, bgOther, 99.9), 4}, {"other woman 100", waist(false, bgOther, 100), 7},
		{"SE/PI uses other cut-offs", waist(false, bgSEPI, 87.9), 2},
		{"asian man 89.9", waist(true, bgAATSI, 89.9), 5}, {"asian man 90", waist(true, bgAATSI, 90), 9},
		{"asian man 99.9", waist(true, bgAATSI, 99.9), 9}, {"asian man 100", waist(true, bgAATSI, 100), 12},
		{"asian woman 79.9", waist(false, bgAATSI, 79.9), 2}, {"asian woman 80", waist(false, bgAATSI, 80), 6},
		{"asian woman 89.9", waist(false, bgAATSI, 89.9), 6}, {"asian woman 90", waist(false, bgAATSI, 90), 9},
	}
	for _, tc := range tests {
		if got, err := AUSDRISKScore(tc.f); err != nil || got != tc.want {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// parity: ledger ausdrisk, source cases 1–10 (the app's picker: age category 0–4 → ages 30/40/50/60/70; ethnicity
// 0 Southern European, 1 Asian, 2 Aboriginal/Torres Strait Islander, 3 Pacific Islander, 4 Other; waist inches ×
// 2.54 → cm). Every source value agrees with Box 3.
func TestAUSDRISKScore_parity(t *testing.T) {
	ages := []float64{30, 40, 50, 60, 70}
	bgs := []AUSDRISKBackground{bgSEPI, bgAATSI, bgAATSI, bgSEPI, bgOther}
	cases := []struct {
		in   [9]float64 // age_cat, sex(1 male), ethnicity, parent, glucose, antihtn, smoker, inactive, waist_in
		want int
	}{
		{[9]float64{0, 0, 4, 0, 0, 0, 0, 0, 30}, 0},    // source case 1
		{[9]float64{2, 1, 4, 0, 0, 0, 0, 0, 38}, 7},    // source case 2
		{[9]float64{2, 1, 4, 1, 0, 0, 0, 0, 30}, 10},   // source case 3
		{[9]float64{4, 0, 3, 0, 0, 0, 0, 0, 34.5}, 10}, // source case 4
		{[9]float64{4, 0, 3, 0, 0, 0, 0, 0, 35}, 14},   // source case 5
		{[9]float64{2, 1, 1, 0, 0, 0, 0, 1, 36}, 15},   // source case 6
		{[9]float64{3, 0, 2, 1, 0, 0, 0, 0, 32}, 15},   // source case 7
		{[9]float64{3, 1, 4, 1, 0, 1, 0, 0, 41}, 18},   // source case 8
		{[9]float64{0, 1, 4, 0, 1, 0, 1, 1, 43.5}, 20}, // source case 9
		{[9]float64{4, 1, 1, 1, 1, 1, 1, 1, 40}, 35},   // source case 10
	}
	for i, tc := range cases {
		v := tc.in
		f := AUSDRISKFindings{
			AgeYears: ages[int(v[0])], Male: v[1] == 1, Background: bgs[int(v[2])],
			ParentalDiabetes: v[3] == 1, HighBloodGlucoseHistory: v[4] == 1, AntihypertensiveMedication: v[5] == 1,
			CurrentSmoker: v[6] == 1, PhysicallyInactive: v[7] == 1, WaistCm: v[8] * 2.54,
		}
		if got, err := AUSDRISKScore(f); err != nil || got != tc.want {
			t.Errorf("source case %d: got (%v, %v), want %v", i+1, got, err, tc.want)
		}
	}
}

func TestAUSDRISKScore_errors(t *testing.T) {
	call := func(f AUSDRISKFindings) func() error {
		return func() error { _, err := AUSDRISKScore(f); return err }
	}
	runErrCases(t, []errCase{
		{"age NaN", call(AUSDRISKFindings{AgeYears: nan, WaistCm: 80}), ErrNotFinite},
		{"waist +Inf", call(AUSDRISKFindings{AgeYears: 40, WaistCm: pInf}), ErrNotFinite},
		{"waist -Inf", call(AUSDRISKFindings{AgeYears: 40, WaistCm: nInf}), ErrNotFinite},
		{"age 24.99", call(AUSDRISKFindings{AgeYears: 24.99, WaistCm: 80}), ErrOutOfRange},
		{"waist 0", call(AUSDRISKFindings{AgeYears: 40, WaistCm: 0}), ErrOutOfRange},
		{"background unknown", call(AUSDRISKFindings{AgeYears: 40, WaistCm: 80, Background: 3}), ErrOutOfRange},
		{"background negative", call(AUSDRISKFindings{AgeYears: 40, WaistCm: 80, Background: -1}), ErrOutOfRange},
	})
}

func TestAUSDRISKScore_noAllocation(t *testing.T) {
	f := AUSDRISKFindings{AgeYears: 47, Male: true, WaistCm: 105}
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = AUSDRISKScore(f) }); n != 0 {
		t.Errorf("AUSDRISKScore allocates %v times, want 0", n)
	}
}

func BenchmarkAUSDRISKScore(b *testing.B) {
	b.ReportAllocs()
	f := AUSDRISKFindings{47, true, bgAATSI, true, false, true, false, true, 95}
	for b.Loop() {
		sinkInt, _ = AUSDRISKScore(f)
	}
}
