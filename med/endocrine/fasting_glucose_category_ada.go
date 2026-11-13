package endocrine

// FastingGlucoseCategoryADAMgDL returns the ADA category of a fasting plasma glucose (FPG) in mg/dL:
// "Diabetes" for FPG ≥ 126 mg/dL (Table 2.1), "Prediabetes" for 100–125 mg/dL (Table 2.2, which marks
// the range "(IFG)"), else no category (""). A value between the printed 125 and 126 (125.5) is
// "Prediabetes". Table 2.1: "Fasting is defined as no caloric intake for at least 8 h."
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-ENDO-015
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  American Diabetes Association Professional Practice Committee. 2. Diagnosis and Classification of Diabetes: Standards of Care in Diabetes—2026. Diabetes Care. 2026;49(Suppl 1):S27-S49. doi:10.2337/dc26-s002. PMID 41358893. PMC12690183, Table 2.1 "FPG ≥126 mg/dL (≥7.0 mmol/L)"; Table 2.2 "FPG 100 mg/dL (5.6 mmol/L) to 125 mg/dL (6.9 mmol/L) (IFG)"
// Tags:    endocrine, diabetes, prediabetes, glucose, fasting, fpg, ifg, ada, category
// Units:   fastingGlucoseMgDL plasma glucose mg/dL; result a category label
// Range:   fastingGlucoseMgDL > 0; nonpregnant individuals (the tables' captions); "children, adolescents, and adults" (section text)
// Example: 99 → ""; 100 → "Prediabetes"; 125 → "Prediabetes"; 126 → "Diabetes" (Tables 2.1, 2.2)
func FastingGlucoseCategoryADAMgDL(fastingGlucoseMgDL float64) (GlycemicCategoryADA, error) {
	// ADA 2026 Table 2.2 lower limit 100 mg/dL, Table 2.1 threshold 126 mg/dL.
	return glycemicCategory("fastingGlucoseMgDL", fastingGlucoseMgDL, 100, 126)
}

// FastingGlucoseCategoryADAMmolL returns the ADA category of a fasting plasma glucose in mmol/L, from the
// source's own mmol/L numbers (no conversion is applied): "Diabetes" for ≥ 7.0 mmol/L (Table 2.1),
// "Prediabetes" for 5.6–6.9 mmol/L (Table 2.2), else no category (""). 6.95 is "Prediabetes".
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-ENDO-016
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  American Diabetes Association Professional Practice Committee. 2. Diagnosis and Classification of Diabetes: Standards of Care in Diabetes—2026. Diabetes Care. 2026;49(Suppl 1):S27-S49. doi:10.2337/dc26-s002. PMID 41358893. PMC12690183, Table 2.1 "FPG ≥126 mg/dL (≥7.0 mmol/L)"; Table 2.2 "FPG 100 mg/dL (5.6 mmol/L) to 125 mg/dL (6.9 mmol/L) (IFG)"
// Tags:    endocrine, diabetes, prediabetes, glucose, fasting, fpg, ifg, ada, category, si
// Units:   fastingGlucoseMmolL plasma glucose mmol/L; result a category label
// Range:   fastingGlucoseMmolL > 0; nonpregnant individuals; children, adolescents and adults
// Example: 5.5 → ""; 5.6 → "Prediabetes"; 6.9 → "Prediabetes"; 7.0 → "Diabetes" (Tables 2.1, 2.2)
func FastingGlucoseCategoryADAMmolL(fastingGlucoseMmolL float64) (GlycemicCategoryADA, error) {
	// ADA 2026 Table 2.2 lower limit 5.6 mmol/L, Table 2.1 threshold 7.0 mmol/L.
	return glycemicCategory("fastingGlucoseMmolL", fastingGlucoseMmolL, 5.6, 7.0)
}
