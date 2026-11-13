package endocrine

// OGTTTwoHourCategoryADAMgDL returns the ADA category of the 2-h plasma glucose (2-h PG) during a 75-g oral
// glucose tolerance test, in mg/dL: "Diabetes" for 2-h PG ≥ 200 mg/dL (Table 2.1), "Prediabetes" for
// 140–199 mg/dL (Table 2.2, which marks the range "(IGT)"), else no category (""). A value between the
// printed 199 and 200 (199.5) is "Prediabetes". The tables define it for the 75-g OGTT only.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-ENDO-017
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  American Diabetes Association Professional Practice Committee. 2. Diagnosis and Classification of Diabetes: Standards of Care in Diabetes—2026. Diabetes Care. 2026;49(Suppl 1):S27-S49. doi:10.2337/dc26-s002. PMID 41358893. PMC12690183, Table 2.1 "2-h PG ≥200 mg/dL (≥11.1 mmol/L) during OGTT"; Table 2.2 "2-h PG during 75-g OGTT 140 mg/dL (7.8 mmol/L) to 199 mg/dL (11.0 mmol/L) (IGT)"
// Tags:    endocrine, diabetes, prediabetes, glucose, ogtt, igt, ada, category
// Units:   twoHourGlucoseMgDL plasma glucose mg/dL, 2 h after a 75-g glucose load; result a category label
// Range:   twoHourGlucoseMgDL > 0; nonpregnant individuals (the tables' captions); "children, adolescents, and adults" (section text)
// Example: 139 → ""; 140 → "Prediabetes"; 199 → "Prediabetes"; 200 → "Diabetes" (Tables 2.1, 2.2)
func OGTTTwoHourCategoryADAMgDL(twoHourGlucoseMgDL float64) (GlycemicCategoryADA, error) {
	// ADA 2026 Table 2.2 lower limit 140 mg/dL, Table 2.1 threshold 200 mg/dL.
	return glycemicCategory("twoHourGlucoseMgDL", twoHourGlucoseMgDL, 140, 200)
}

// OGTTTwoHourCategoryADAMmolL returns the ADA category of the 2-h plasma glucose during a 75-g OGTT in
// mmol/L, from the source's own mmol/L numbers (no conversion is applied): "Diabetes" for ≥ 11.1 mmol/L
// (Table 2.1), "Prediabetes" for 7.8–11.0 mmol/L (Table 2.2), else no category (""). 11.05 is "Prediabetes".
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-ENDO-018
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  American Diabetes Association Professional Practice Committee. 2. Diagnosis and Classification of Diabetes: Standards of Care in Diabetes—2026. Diabetes Care. 2026;49(Suppl 1):S27-S49. doi:10.2337/dc26-s002. PMID 41358893. PMC12690183, Table 2.1 "2-h PG ≥200 mg/dL (≥11.1 mmol/L) during OGTT"; Table 2.2 "2-h PG during 75-g OGTT 140 mg/dL (7.8 mmol/L) to 199 mg/dL (11.0 mmol/L) (IGT)"
// Tags:    endocrine, diabetes, prediabetes, glucose, ogtt, igt, ada, category, si
// Units:   twoHourGlucoseMmolL plasma glucose mmol/L, 2 h after a 75-g glucose load; result a category label
// Range:   twoHourGlucoseMmolL > 0; nonpregnant individuals; children, adolescents and adults
// Example: 7.7 → ""; 7.8 → "Prediabetes"; 11.0 → "Prediabetes"; 11.1 → "Diabetes" (Tables 2.1, 2.2)
func OGTTTwoHourCategoryADAMmolL(twoHourGlucoseMmolL float64) (GlycemicCategoryADA, error) {
	// ADA 2026 Table 2.2 lower limit 7.8 mmol/L, Table 2.1 threshold 11.1 mmol/L.
	return glycemicCategory("twoHourGlucoseMmolL", twoHourGlucoseMmolL, 7.8, 11.1)
}
