package nutrition

import (
	"math"
	"testing"
)

// gramsNear compares all six fields at 1e-12 relative (each value is energy × an integer percent / 400 or 900).
func gramsNear(got, want MacronutrientGrams) bool {
	return relNear(got.ProteinLowG, want.ProteinLowG, 1e-12) && relNear(got.ProteinHighG, want.ProteinHighG, 1e-12) &&
		relNear(got.CarbohydrateLowG, want.CarbohydrateLowG, 1e-12) && relNear(got.CarbohydrateHighG, want.CarbohydrateHighG, 1e-12) &&
		relNear(got.FatLowG, want.FatLowG, 1e-12) && relNear(got.FatHighG, want.FatHighG, 1e-12)
}

// NASEM 2024 Table 3-1 percents with the 21 CFR 101.9(c)(1)(i)(B) factors 4, 4, 9; HAND-DERIVED: adult 2000
// kcal → protein 2000 × 10/400 = 50 to 2000 × 35/400 = 175; carbohydrate 225–325; fat 400/9–700/9.
func TestAMDRGrams_example(t *testing.T) {
	want := MacronutrientGrams{50, 175, 225, 325, 400.0 / 9, 700.0 / 9}
	if got, err := AMDRGrams(2000, 30); err != nil || !gramsNear(got, want) {
		t.Errorf("got (%+v, %v), want %+v", got, err, want)
	}
}

// parity: ledger amdr_macro_grams, source cases 1–6 (kcal, age → protein, carbohydrate, fat low–high).
func TestAMDRGrams_parity(t *testing.T) {
	tests := []struct {
		name      string
		kcal, age float64
		want      MacronutrientGrams
	}{
		{"source case 1", 2000, 30, MacronutrientGrams{50, 175, 225, 325, 400.0 / 9, 700.0 / 9}},
		{"source case 2", 2400, 19, MacronutrientGrams{60, 210, 270, 390, 480.0 / 9, 840.0 / 9}},
		{"source case 3", 2800, 18.5, MacronutrientGrams{70, 210, 315, 455, 700.0 / 9, 980.0 / 9}},
		{"source case 4", 1800, 10, MacronutrientGrams{45, 135, 202.5, 292.5, 50, 70}},
		{"source case 5", 1400, 4, MacronutrientGrams{35, 105, 157.5, 227.5, 350.0 / 9, 490.0 / 9}},
		{"source case 6", 1000, 3.9, MacronutrientGrams{12.5, 50, 112.5, 162.5, 300.0 / 9, 400.0 / 9}},
	}
	for _, tc := range tests {
		if got, err := AMDRGrams(tc.kcal, tc.age); err != nil || !gramsNear(got, tc.want) {
			t.Errorf("%s: got (%+v, %v), want %+v", tc.name, got, err, tc.want)
		}
	}
}

// parity: ledger amdr_macro_grams, source refusal case (18 months). The source refused ages under 2 (it read
// DGA Table A1-2, which starts at 2); the primary's Table 3-1 row "Children 1–3 years" covers 18 months, so
// the primary's value is asserted: 900 kcal → protein 5–20 % → 11.25–45 g.
func TestAMDRGrams_parityEighteenMonths(t *testing.T) {
	want := MacronutrientGrams{11.25, 45, 101.25, 146.25, 30, 40}
	if got, err := AMDRGrams(900, 1.5); err != nil || !gramsNear(got, want) {
		t.Errorf("got (%+v, %v), want %+v", got, err, want)
	}
}

// Group edges: 1, 3.999 | 4, 18.999 | 19.
func TestAMDRGrams_groupEdges(t *testing.T) {
	tests := []struct {
		age                  float64
		proteinLo, proteinHi float64 // percent
	}{
		{1, 5, 20}, {3.999, 5, 20}, {4, 10, 30}, {18.999, 10, 30}, {19, 10, 35},
	}
	for _, tc := range tests {
		got, err := AMDRGrams(400, tc.age)
		if err != nil || !relNear(got.ProteinLowG, tc.proteinLo, 1e-12) || !relNear(got.ProteinHighG, tc.proteinHi, 1e-12) {
			t.Errorf("age %v: got (%+v, %v), want protein %v–%v g", tc.age, got, err, tc.proteinLo, tc.proteinHi)
		}
	}
}

func TestAMDRGrams_errors(t *testing.T) {
	f := func(kcal, age float64) func() error {
		return func() error { _, err := AMDRGrams(kcal, age); return err }
	}
	runErrCases(t, []errCase{
		{"kcal NaN", f(nan, 30), ErrNotFinite},
		{"kcal +Inf", f(pInf, 30), ErrNotFinite},
		{"age -Inf", f(2000, nInf), ErrNotFinite},
		{"age NaN", f(2000, nan), ErrNotFinite},
		{"kcal 0", f(0, 30), ErrOutOfRange},
		{"kcal negative", f(-1, 30), ErrOutOfRange},
		{"age 0.999", f(2000, 0.999), ErrOutOfRange},
		{"overflow", f(math.MaxFloat64, 30), ErrDomain},
	})
}

func TestAMDRGrams_noAllocation(t *testing.T) {
	var g MacronutrientGrams
	if n := testing.AllocsPerRun(100, func() { g, _ = AMDRGrams(2000, 30) }); n != 0 {
		t.Errorf("AMDRGrams allocates %v times, want 0", n)
	}
	sink = g.FatLowG
}

func BenchmarkAMDRGrams(b *testing.B) {
	b.ReportAllocs()
	kcal, age := 2000.0, 30.0
	var g MacronutrientGrams
	for b.Loop() {
		g, _ = AMDRGrams(kcal, age)
	}
	sink = g.FatLowG
}
