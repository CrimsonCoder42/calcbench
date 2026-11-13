package endocrine

import "testing"

// ADA 2026 Tables 2.1 and 2.2: FPG 100–125 mg/dL Prediabetes, ≥ 126 Diabetes; both sides of each edge.
func TestFastingGlucoseCategoryADAMgDL(t *testing.T) {
	checkGlycemic(t, "FastingGlucoseCategoryADAMgDL", FastingGlucoseCategoryADAMgDL, []glycemicCase{
		{40, GlycemicNoCategoryADA}, {99, GlycemicNoCategoryADA}, {99.9, GlycemicNoCategoryADA},
		{100, GlycemicPrediabetesADA}, {125, GlycemicPrediabetesADA}, {125.5, GlycemicPrediabetesADA},
		{125.9, GlycemicPrediabetesADA}, {126, GlycemicDiabetesADA}, {600, GlycemicDiabetesADA},
	})
	checkGlycemicErrors(t, "FastingGlucoseCategoryADAMgDL", FastingGlucoseCategoryADAMgDL)
}

// ADA 2026 Tables 2.1 and 2.2, mmol/L: 5.6–6.9 Prediabetes, ≥ 7.0 Diabetes.
func TestFastingGlucoseCategoryADAMmolL(t *testing.T) {
	checkGlycemic(t, "FastingGlucoseCategoryADAMmolL", FastingGlucoseCategoryADAMmolL, []glycemicCase{
		{3, GlycemicNoCategoryADA}, {5.5, GlycemicNoCategoryADA}, {5.59, GlycemicNoCategoryADA},
		{5.6, GlycemicPrediabetesADA}, {6.9, GlycemicPrediabetesADA}, {6.95, GlycemicPrediabetesADA},
		{6.99, GlycemicPrediabetesADA}, {7.0, GlycemicDiabetesADA}, {30, GlycemicDiabetesADA},
	})
	checkGlycemicErrors(t, "FastingGlucoseCategoryADAMmolL", FastingGlucoseCategoryADAMmolL)
}

// The source prints each mg/dL edge with its mmol/L twin (100 ↔ 5.6, 125 ↔ 6.9, 126 ↔ 7.0); the
// categories agree at those pairs.
func TestFastingGlucoseCategoryADA_twins(t *testing.T) {
	for _, p := range [][2]float64{{99, 5.5}, {100, 5.6}, {125, 6.9}, {126, 7.0}} {
		a, _ := FastingGlucoseCategoryADAMgDL(p[0])
		b, _ := FastingGlucoseCategoryADAMmolL(p[1])
		if a != b {
			t.Errorf("%v mg/dL → %q, %v mmol/L → %q", p[0], a, p[1], b)
		}
	}
}

// parity: ledger bs_fasting_adult, source case 1 (bands 70–99 "Normal", 100–125 "Pre-Diabetic", 126–299
// "Diabetic Range", ≥ 300 "Critically High"); the primary's labels are asserted. The source's < 70
// "Hypoglycemia" band is HypoglycemiaLevelADAMgDL (see its parity test).
func TestFastingGlucoseCategoryADAMgDL_parityAdult(t *testing.T) {
	checkGlycemic(t, "FastingGlucoseCategoryADAMgDL", FastingGlucoseCategoryADAMgDL, []glycemicCase{
		{70, GlycemicNoCategoryADA},   // source "Normal"
		{99.9, GlycemicNoCategoryADA}, // source "Normal"
		{100, GlycemicPrediabetesADA}, // source "Pre-Diabetic"
		{125.9, GlycemicPrediabetesADA},
		{126, GlycemicDiabetesADA}, // source "Diabetic Range"
		{299.9, GlycemicDiabetesADA},
		{300, GlycemicDiabetesADA}, // source "Critically High": no primary band; the primary's "Diabetes"
		{1000, GlycemicDiabetesADA},
	})
}

// parity: ledger bs_fasting_child, source case 1 (bands 70–99 "Normal", 100–125 "Elevated", ≥ 126 "High").
// ADA 2026 section 2: Tables 2.1 and 2.2 "apply to children, adolescents, and adults", so the same
// function and the primary's labels are asserted.
func TestFastingGlucoseCategoryADAMgDL_parityChild(t *testing.T) {
	checkGlycemic(t, "FastingGlucoseCategoryADAMgDL", FastingGlucoseCategoryADAMgDL, []glycemicCase{
		{99.9, GlycemicNoCategoryADA},   // source "Normal"
		{100, GlycemicPrediabetesADA},   // source "Elevated"
		{125.9, GlycemicPrediabetesADA}, // source "Elevated"
		{126, GlycemicDiabetesADA},      // source "High"
	})
}

// parity: ledger ClassifyGlucose, source cases 1–4 (relation "fasting": 99 "normal", 100 and 125
// "prediabetes_range", 126 "diabetes_range"); the primary's labels are asserted.
func TestFastingGlucoseCategoryADAMgDL_parityClassifier(t *testing.T) {
	checkGlycemic(t, "FastingGlucoseCategoryADAMgDL", FastingGlucoseCategoryADAMgDL, []glycemicCase{
		{99, GlycemicNoCategoryADA}, {100, GlycemicPrediabetesADA}, {125, GlycemicPrediabetesADA},
		{126, GlycemicDiabetesADA},
	})
}

func TestFastingGlucoseCategoryADA_noAlloc(t *testing.T) {
	checkGlycemicNoAlloc(t, "FastingGlucoseCategoryADAMgDL", FastingGlucoseCategoryADAMgDL, 110)
	checkGlycemicNoAlloc(t, "FastingGlucoseCategoryADAMmolL", FastingGlucoseCategoryADAMmolL, 6.1)
}

func BenchmarkFastingGlucoseCategoryADAMgDL(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = FastingGlucoseCategoryADAMgDL(110)
	}
}

func BenchmarkFastingGlucoseCategoryADAMmolL(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = FastingGlucoseCategoryADAMmolL(6.1)
	}
}
