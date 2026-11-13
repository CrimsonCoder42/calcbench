package endocrine

// RandomGlucoseCriterionADAMgDL reports whether a random plasma glucose meets the fourth ADA Table 2.1
// criterion: "In an individual with classic symptoms of hyperglycemia or hyperglycemic crisis, a random
// plasma glucose ≥200 mg/dL (≥11.1 mmol/L)." It is true only when classicSymptoms is true and the glucose
// is ≥ 200 mg/dL. The source defines no prediabetes range for a random glucose.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-ENDO-019
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  American Diabetes Association Professional Practice Committee. 2. Diagnosis and Classification of Diabetes: Standards of Care in Diabetes—2026. Diabetes Care. 2026;49(Suppl 1):S27-S49. doi:10.2337/dc26-s002. PMID 41358893. PMC12690183, Table 2.1, fourth criterion
// Tags:    endocrine, diabetes, glucose, random, ada, threshold
// Units:   randomGlucoseMgDL plasma glucose mg/dL; classicSymptoms the caller's finding (classic symptoms of hyperglycemia or hyperglycemic crisis); result true when the criterion is met
// Range:   randomGlucoseMgDL > 0; nonpregnant individuals (Table 2.1 caption); children, adolescents and adults (section text)
// Example: 199, true → false; 200, true → true; 300, false → false (Table 2.1)
func RandomGlucoseCriterionADAMgDL(randomGlucoseMgDL float64, classicSymptoms bool) (bool, error) {
	if err := checkFinitePositive("randomGlucoseMgDL", randomGlucoseMgDL); err != nil {
		return false, err
	}
	// ADA 2026 Table 2.1: ≥ 200 mg/dL.
	return classicSymptoms && randomGlucoseMgDL >= 200, nil
}

// RandomGlucoseCriterionADAMmolL reports whether a random plasma glucose in mmol/L meets the fourth ADA
// Table 2.1 criterion, from the source's own mmol/L number (≥ 11.1 mmol/L; no conversion is applied), with
// classic symptoms of hyperglycemia or hyperglycemic crisis.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-ENDO-020
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  American Diabetes Association Professional Practice Committee. 2. Diagnosis and Classification of Diabetes: Standards of Care in Diabetes—2026. Diabetes Care. 2026;49(Suppl 1):S27-S49. doi:10.2337/dc26-s002. PMID 41358893. PMC12690183, Table 2.1, fourth criterion "(≥11.1 mmol/L)"
// Tags:    endocrine, diabetes, glucose, random, ada, threshold, si
// Units:   randomGlucoseMmolL plasma glucose mmol/L; classicSymptoms the caller's finding; result true when the criterion is met
// Range:   randomGlucoseMmolL > 0; nonpregnant individuals; children, adolescents and adults
// Example: 11.0, true → false; 11.1, true → true; 15, false → false (Table 2.1)
func RandomGlucoseCriterionADAMmolL(randomGlucoseMmolL float64, classicSymptoms bool) (bool, error) {
	if err := checkFinitePositive("randomGlucoseMmolL", randomGlucoseMmolL); err != nil {
		return false, err
	}
	// ADA 2026 Table 2.1: ≥ 11.1 mmol/L.
	return classicSymptoms && randomGlucoseMmolL >= 11.1, nil
}
