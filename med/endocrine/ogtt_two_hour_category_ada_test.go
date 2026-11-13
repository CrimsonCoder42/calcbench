package endocrine

import "testing"

// ADA 2026 Tables 2.1 and 2.2: 2-h PG 140–199 mg/dL Prediabetes, ≥ 200 Diabetes; both sides of each edge.
func TestOGTTTwoHourCategoryADAMgDL(t *testing.T) {
	checkGlycemic(t, "OGTTTwoHourCategoryADAMgDL", OGTTTwoHourCategoryADAMgDL, []glycemicCase{
		{60, GlycemicNoCategoryADA}, {139, GlycemicNoCategoryADA}, {139.9, GlycemicNoCategoryADA},
		{140, GlycemicPrediabetesADA}, {199, GlycemicPrediabetesADA}, {199.5, GlycemicPrediabetesADA},
		{199.9, GlycemicPrediabetesADA}, {200, GlycemicDiabetesADA}, {400, GlycemicDiabetesADA},
	})
	checkGlycemicErrors(t, "OGTTTwoHourCategoryADAMgDL", OGTTTwoHourCategoryADAMgDL)
}

// ADA 2026 Tables 2.1 and 2.2, mmol/L: 7.8–11.0 Prediabetes, ≥ 11.1 Diabetes.
func TestOGTTTwoHourCategoryADAMmolL(t *testing.T) {
	checkGlycemic(t, "OGTTTwoHourCategoryADAMmolL", OGTTTwoHourCategoryADAMmolL, []glycemicCase{
		{4, GlycemicNoCategoryADA}, {7.7, GlycemicNoCategoryADA}, {7.79, GlycemicNoCategoryADA},
		{7.8, GlycemicPrediabetesADA}, {11.0, GlycemicPrediabetesADA}, {11.05, GlycemicPrediabetesADA},
		{11.09, GlycemicPrediabetesADA}, {11.1, GlycemicDiabetesADA}, {25, GlycemicDiabetesADA},
	})
	checkGlycemicErrors(t, "OGTTTwoHourCategoryADAMmolL", OGTTTwoHourCategoryADAMmolL)
}

// The source's paired values (140 ↔ 7.8, 199 ↔ 11.0, 200 ↔ 11.1) give the same category.
func TestOGTTTwoHourCategoryADA_twins(t *testing.T) {
	for _, p := range [][2]float64{{139, 7.7}, {140, 7.8}, {199, 11.0}, {200, 11.1}} {
		a, _ := OGTTTwoHourCategoryADAMgDL(p[0])
		b, _ := OGTTTwoHourCategoryADAMmolL(p[1])
		if a != b {
			t.Errorf("%v mg/dL → %q, %v mmol/L → %q", p[0], a, p[1], b)
		}
	}
}

// parity: ledger bs_ogtt_2h_adult, source case 1 (bands 70–139 "Normal OGTT", 140–199 "Prediabetes —
// Impaired Glucose Tolerance", > 200 "Diabetes (OGTT)" with 200 exclusive). Source case 2 disagrees with
// the primary: the source puts exactly 200 mg/dL in its prediabetes band (min_exclusive), Table 2.1 says
// "2-h PG ≥200 mg/dL"; the primary's "Diabetes" is asserted.
func TestOGTTTwoHourCategoryADAMgDL_parity(t *testing.T) {
	checkGlycemic(t, "OGTTTwoHourCategoryADAMgDL", OGTTTwoHourCategoryADAMgDL, []glycemicCase{
		{139.9, GlycemicNoCategoryADA},  // case 1: source "Normal OGTT"
		{140, GlycemicPrediabetesADA},   // case 1: source "Prediabetes — Impaired Glucose Tolerance"
		{199.9, GlycemicPrediabetesADA}, // case 1
		{200, GlycemicDiabetesADA},      // case 2: source "Prediabetes — Impaired Glucose Tolerance"; primary ≥200
		{200.1, GlycemicDiabetesADA},    // case 1: source "Diabetes (OGTT)"
	})
}

// parity: ledger ClassifyGlucose, source cases 5–8 (relation "after_meal" read as the 2-h OGTT bands:
// 139 "normal", 140 and 199 "prediabetes_range", 200 "diabetes_range"); the primary's labels are asserted.
func TestOGTTTwoHourCategoryADAMgDL_parityClassifier(t *testing.T) {
	checkGlycemic(t, "OGTTTwoHourCategoryADAMgDL", OGTTTwoHourCategoryADAMgDL, []glycemicCase{
		{139, GlycemicNoCategoryADA}, {140, GlycemicPrediabetesADA}, {199, GlycemicPrediabetesADA},
		{200, GlycemicDiabetesADA},
	})
}

func TestOGTTTwoHourCategoryADA_noAlloc(t *testing.T) {
	checkGlycemicNoAlloc(t, "OGTTTwoHourCategoryADAMgDL", OGTTTwoHourCategoryADAMgDL, 150)
	checkGlycemicNoAlloc(t, "OGTTTwoHourCategoryADAMmolL", OGTTTwoHourCategoryADAMmolL, 8.3)
}

func BenchmarkOGTTTwoHourCategoryADAMgDL(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = OGTTTwoHourCategoryADAMgDL(150)
	}
}

func BenchmarkOGTTTwoHourCategoryADAMmolL(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = OGTTTwoHourCategoryADAMmolL(8.3)
	}
}
