package anthropometry

import "fmt"

// BMICategory is a CDC adult BMI category; its value is the CDC's own label, verbatim.
type BMICategory string

// The CDC adult BMI categories, labels verbatim from Centers for Disease Control and Prevention, Adult BMI
// Categories, https://www.cdc.gov/bmi/adult-calculator/bmi-categories.html (last reviewed 2024-03-19; read
// 2026-11-12). The page's umbrella label "Obesity" (30 or greater) covers the three classes; the function
// returns the class.
const (
	BMIUnderweight   BMICategory = "Underweight"                      // less than 18.5
	BMIHealthyWeight BMICategory = "Healthy Weight"                   // 18.5 to less than 25 (the CDC's label)
	BMIOverweight    BMICategory = "Overweight"                       // 25 to less than 30
	BMIClass1Obesity BMICategory = "Class 1 Obesity"                  // 30 to less than 35
	BMIClass2Obesity BMICategory = "Class 2 Obesity"                  // 35 to less than 40
	BMIClass3Obesity BMICategory = "Class 3 Obesity (Severe Obesity)" // 40 or greater
)

// BMICategoryCDCAdult returns the CDC adult BMI category for a BMI in kg/m², with the CDC's own label:
// < 18.5 "Underweight"; 18.5 to < 25 "Healthy Weight"; 25 to < 30 "Overweight"; 30 to < 35 "Class 1
// Obesity"; 35 to < 40 "Class 2 Obesity"; ≥ 40 "Class 3 Obesity (Severe Obesity)". These are the CDC's
// ADULT categories; the function takes no age and does not apply them to children.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for bmiKgM2 ≤ 0.
//
// ID:      MED-ANTHRO-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Centers for Disease Control and Prevention. Adult BMI Categories. https://www.cdc.gov/bmi/adult-calculator/bmi-categories.html (last reviewed 2024-03-19; read 2026-11-12)
// Tags:    anthropometry, bmi, category, cdc, adult
// Units:   bmiKgM2 kg/m²; result the CDC's category label
// Range:   bmiKgM2 > 0; the CDC page's adult categories
// Example: 18.4999 → "Underweight"; 18.5 → "Healthy Weight"; 40 → "Class 3 Obesity (Severe Obesity)" (CDC table)
func BMICategoryCDCAdult(bmiKgM2 float64) (BMICategory, error) {
	if err := checkFinite("bmiKgM2", bmiKgM2); err != nil {
		return "", err
	}
	if !(bmiKgM2 > 0) {
		return "", fmt.Errorf("%w: bmiKgM2 = %v, want > 0", ErrOutOfRange, bmiKgM2)
	}
	switch {
	case bmiKgM2 < 18.5:
		return BMIUnderweight, nil
	case bmiKgM2 < 25:
		return BMIHealthyWeight, nil
	case bmiKgM2 < 30:
		return BMIOverweight, nil
	case bmiKgM2 < 35:
		return BMIClass1Obesity, nil
	case bmiKgM2 < 40:
		return BMIClass2Obesity, nil
	}
	return BMIClass3Obesity, nil
}
