package hepatic

// BARDScore returns the BARD score, 0–4, for nonalcoholic fatty liver disease: BMI ≥ 28 kg/m² scores 1, an
// AST/ALT ratio (AAR) ≥ 0.8 scores 2, diabetes mellitus scores 1. No band is returned.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-HEPAT-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Harrison SA, Oliver D, Arnold HL, Gogia S, Neuschwander-Tetri BA. Development and validation of a simple NAFLD clinical scoring system for identifying patients without advanced disease. Gut. 2008;57(10):1441-1447. doi:10.1136/gut.2007.146019. PMID 18390575. Abstract: "three variables were combined in a weighted sum (BMI > or = 28 = 1 point, AAR of > or = 0.8 = 2 points, DM = 1 point)"
// Tags:    hepatic, liver, nafld, fibrosis, bard, ast, alt, bmi, score
// Units:   bmiKgM2 kg/m²; astUL U/L; altUL U/L; result points 0–4
// Range:   every numeric input > 0 (population: 827 people with NAFLD, two tertiary centres; input limits not stated in the source)
// Example: BMI 30, AST 40, ALT 45, no diabetes → AAR 0.889 → 1 + 2 = 3 (HAND-DERIVED from the abstract)
func BARDScore(bmiKgM2, astUL, altUL float64, diabetes bool) (int, error) {
	if err := checkFinite("bmiKgM2", bmiKgM2); err != nil {
		return 0, err
	}
	if err := checkFinite("astUL", astUL); err != nil {
		return 0, err
	}
	if err := checkFinite("altUL", altUL); err != nil {
		return 0, err
	}
	if err := checkPositive("bmiKgM2", bmiKgM2); err != nil {
		return 0, err
	}
	if err := checkPositive("astUL", astUL); err != nil {
		return 0, err
	}
	if err := checkPositive("altUL", altUL); err != nil {
		return 0, err
	}
	// The ratio of two positive finite values is never NaN; an overflow to +Inf or underflow to 0 still compares
	// correctly with 0.8.
	aar := astUL / altUL
	points := 0
	if bmiKgM2 >= 28 {
		points++
	}
	if aar >= 0.8 {
		points += 2
	}
	if diabetes {
		points++
	}
	return points, nil
}
