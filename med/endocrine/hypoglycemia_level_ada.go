package endocrine

// HypoglycemiaLevelADA is a glucose-defined level of the ADA "Classification of hypoglycemia" (Table 6.4):
// "Level 1" or "Level 2". Level 3 ("A severe event characterized by altered mental and/or physical status
// requiring assistance for treatment of hypoglycemia, irrespective of glucose level") is not defined by a
// glucose value and is never returned. A glucose at or above the Level 1 upper limit has no level:
// HypoglycemiaNoLevelADA, the empty string.
type HypoglycemiaLevelADA string

// The ADA 2025 Standards, section 6, Table 6.4.
const (
	// HypoglycemiaLevel1ADA: "Glucose <70 mg/dL (<3.9 mmol/L) and ≥54 mg/dL (≥3.0 mmol/L)".
	HypoglycemiaLevel1ADA HypoglycemiaLevelADA = "Level 1"
	// HypoglycemiaLevel2ADA: "Glucose <54 mg/dL (<3.0 mmol/L)".
	HypoglycemiaLevel2ADA HypoglycemiaLevelADA = "Level 2"
	// HypoglycemiaNoLevelADA: glucose ≥ 70 mg/dL (≥ 3.9 mmol/L); the table names no level.
	HypoglycemiaNoLevelADA HypoglycemiaLevelADA = ""
)

// hypoglycemiaLevel checks g (finite, > 0) and returns Level 2 for g < level2Below, Level 1 for
// g < level1Below, else no level.
func hypoglycemiaLevel(name string, g, level2Below, level1Below float64) (HypoglycemiaLevelADA, error) {
	if err := checkFinitePositive(name, g); err != nil {
		return HypoglycemiaNoLevelADA, err
	}
	switch {
	case g < level2Below:
		return HypoglycemiaLevel2ADA, nil
	case g < level1Below:
		return HypoglycemiaLevel1ADA, nil
	}
	return HypoglycemiaNoLevelADA, nil
}

// HypoglycemiaLevelADAMgDL returns the ADA hypoglycemia level of a glucose in mg/dL: "Level 2" for
// < 54 mg/dL, "Level 1" for < 70 and ≥ 54 mg/dL, else no level (""). The table's inequalities are strict
// and closed as printed, so there is no gap between printed numbers.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-ENDO-021
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  American Diabetes Association Professional Practice Committee. 6. Glycemic Goals and Hypoglycemia: Standards of Care in Diabetes—2025. Diabetes Care. 2025;48(Suppl 1):S128-S145. doi:10.2337/dc25-s006. PMID 39651981. PMC11635034, Table 6.4 "Classification of hypoglycemia"
// Tags:    endocrine, diabetes, hypoglycemia, glucose, ada, category
// Units:   glucoseMgDL glucose mg/dL; result a level label
// Range:   glucoseMgDL > 0; people with diabetes (the section's population)
// Example: 53.9 → "Level 2"; 54 → "Level 1"; 69.9 → "Level 1"; 70 → "" (Table 6.4)
func HypoglycemiaLevelADAMgDL(glucoseMgDL float64) (HypoglycemiaLevelADA, error) {
	// ADA 2025 Table 6.4: Level 2 < 54 mg/dL, Level 1 < 70 mg/dL.
	return hypoglycemiaLevel("glucoseMgDL", glucoseMgDL, 54, 70)
}

// HypoglycemiaLevelADAMmolL returns the ADA hypoglycemia level of a glucose in mmol/L, from the source's
// own mmol/L numbers (no conversion is applied): "Level 2" for < 3.0 mmol/L, "Level 1" for < 3.9 and
// ≥ 3.0 mmol/L, else no level ("").
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-ENDO-022
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  American Diabetes Association Professional Practice Committee. 6. Glycemic Goals and Hypoglycemia: Standards of Care in Diabetes—2025. Diabetes Care. 2025;48(Suppl 1):S128-S145. doi:10.2337/dc25-s006. PMID 39651981. PMC11635034, Table 6.4 "Classification of hypoglycemia" (mmol/L values)
// Tags:    endocrine, diabetes, hypoglycemia, glucose, ada, category, si
// Units:   glucoseMmolL glucose mmol/L; result a level label
// Range:   glucoseMmolL > 0; people with diabetes (the section's population)
// Example: 2.9 → "Level 2"; 3.0 → "Level 1"; 3.8 → "Level 1"; 3.9 → "" (Table 6.4)
func HypoglycemiaLevelADAMmolL(glucoseMmolL float64) (HypoglycemiaLevelADA, error) {
	// ADA 2025 Table 6.4: Level 2 < 3.0 mmol/L, Level 1 < 3.9 mmol/L.
	return hypoglycemiaLevel("glucoseMmolL", glucoseMmolL, 3.0, 3.9)
}
