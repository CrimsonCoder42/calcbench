package endocrine

// GlycemicCategoryADA is a category of the ADA Standards of Care diagnostic tables: "Diabetes" (Table 2.1,
// "Criteria for the diagnosis of diabetes in nonpregnant individuals") or "Prediabetes" (Table 2.2,
// "Criteria defining prediabetes in nonpregnant individuals"). A value below the Table 2.2 range has no
// category in the source: GlycemicNoCategoryADA, the empty string.
type GlycemicCategoryADA string

// The ADA 2026 Standards, section 2, Tables 2.1 and 2.2.
const (
	// GlycemicDiabetesADA: the value meets its Table 2.1 criterion. Table 2.1's footnote: "In the absence
	// of unequivocal hyperglycemia, diagnosis requires two abnormal results from different tests, which may
	// be obtained at the same time (e.g., A1C and FPG), or the same test at two different time points."
	GlycemicDiabetesADA GlycemicCategoryADA = "Diabetes"
	// GlycemicPrediabetesADA: the value is in its Table 2.2 range.
	GlycemicPrediabetesADA GlycemicCategoryADA = "Prediabetes"
	// GlycemicNoCategoryADA: the value is below the Table 2.2 lower limit; the source names no category.
	GlycemicNoCategoryADA GlycemicCategoryADA = ""
)

// glycemicCategory checks v (finite, > 0) and returns Diabetes for v ≥ diabetesFrom, Prediabetes for
// v ≥ prediabetesFrom, else no category. The tables print closed ranges at the source's precision
// (100–125, ≥126); a value between two printed numbers (125.5) takes the lower category, the one whose
// printed lower edge it has passed.
func glycemicCategory(name string, v, prediabetesFrom, diabetesFrom float64) (GlycemicCategoryADA, error) {
	if err := checkFinitePositive(name, v); err != nil {
		return GlycemicNoCategoryADA, err
	}
	switch {
	case v >= diabetesFrom:
		return GlycemicDiabetesADA, nil
	case v >= prediabetesFrom:
		return GlycemicPrediabetesADA, nil
	}
	return GlycemicNoCategoryADA, nil
}

// A1CCategoryADAPercent returns the ADA category of a hemoglobin A1C in percent (NGSP): "Diabetes" for
// A1C ≥ 6.5 % (Table 2.1), "Prediabetes" for 5.7–6.4 % (Table 2.2), else no category (""). A value between
// the printed 6.4 and 6.5 (6.45) is "Prediabetes". Table 2.1 adds: "The test should be performed in a
// laboratory using a method that is NGSP certified and standardized to the DCCT assay."
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-ENDO-013
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  American Diabetes Association Professional Practice Committee. 2. Diagnosis and Classification of Diabetes: Standards of Care in Diabetes—2026. Diabetes Care. 2026;49(Suppl 1):S27-S49. doi:10.2337/dc26-s002. PMID 41358893. PMC12690183, Table 2.1 "A1C ≥6.5% (≥48 mmol/mol)"; Table 2.2 "A1C 5.7–6.4% (39–47 mmol/mol)"
// Tags:    endocrine, diabetes, prediabetes, a1c, hba1c, ada, category
// Units:   a1cPercent percent (NGSP); result a category label
// Range:   a1cPercent > 0; nonpregnant individuals (the tables' captions); "children, adolescents, and adults" (section text)
// Example: 5.6 → ""; 5.7 → "Prediabetes"; 6.4 → "Prediabetes"; 6.5 → "Diabetes" (Tables 2.1, 2.2)
func A1CCategoryADAPercent(a1cPercent float64) (GlycemicCategoryADA, error) {
	// ADA 2026 Table 2.2 lower limit 5.7 %, Table 2.1 threshold 6.5 %.
	return glycemicCategory("a1cPercent", a1cPercent, 5.7, 6.5)
}

// A1CCategoryADAMmolMol returns the ADA category of a hemoglobin A1C in mmol/mol (IFCC), from the source's
// own mmol/mol numbers (no conversion is applied): "Diabetes" for ≥ 48 mmol/mol (Table 2.1),
// "Prediabetes" for 39–47 mmol/mol (Table 2.2), else no category (""). 47.5 is "Prediabetes".
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-ENDO-014
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  American Diabetes Association Professional Practice Committee. 2. Diagnosis and Classification of Diabetes: Standards of Care in Diabetes—2026. Diabetes Care. 2026;49(Suppl 1):S27-S49. doi:10.2337/dc26-s002. PMID 41358893. PMC12690183, Table 2.1 "A1C ≥6.5% (≥48 mmol/mol)"; Table 2.2 "A1C 5.7–6.4% (39–47 mmol/mol)"
// Tags:    endocrine, diabetes, prediabetes, a1c, hba1c, ada, category, si, ifcc
// Units:   a1cMmolMol mmol/mol (IFCC); result a category label
// Range:   a1cMmolMol > 0; nonpregnant individuals; children, adolescents and adults
// Example: 38 → ""; 39 → "Prediabetes"; 47 → "Prediabetes"; 48 → "Diabetes" (Tables 2.1, 2.2)
func A1CCategoryADAMmolMol(a1cMmolMol float64) (GlycemicCategoryADA, error) {
	// ADA 2026 Table 2.2 lower limit 39 mmol/mol, Table 2.1 threshold 48 mmol/mol.
	return glycemicCategory("a1cMmolMol", a1cMmolMol, 39, 48)
}
