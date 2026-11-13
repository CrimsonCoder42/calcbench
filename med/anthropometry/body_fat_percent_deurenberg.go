package anthropometry

import "fmt"

// checkDeurenbergBMI returns ErrOutOfRange unless 13.9 ≤ bmiKgM2 ≤ 40.9 (the BMI range of the study's 1229
// subjects, Deurenberg et al. 1991).
func checkDeurenbergBMI(bmiKgM2 float64) error {
	if bmiKgM2 < 13.9 || bmiKgM2 > 40.9 {
		return fmt.Errorf("%w: bmiKgM2 = %v, want 13.9 ≤ bmiKgM2 ≤ 40.9", ErrOutOfRange, bmiKgM2)
	}
	return nil
}

// sexIndicator is the source's sex term: 1 for male, 0 for female.
func sexIndicator(male bool) float64 {
	if male {
		return 1
	}
	return 0
}

// BodyFatPercentDeurenbergAdult returns the body-fat percentage predicted from BMI, age and sex by the
// Deurenberg adult formula: 1.20 × BMI + 0.23 × age − 10.8 × sex − 5.4, with sex = 1 for male and 0 for
// female (the source's coding). The source reports R² 0.79 and SEE 4.1% for adults.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-ANTHRO-008
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Deurenberg P, Weststrate JA, Seidell JC. Body mass index as a measure of body fatness: age- and sex-specific prediction formulas. Br J Nutr. 1991;65(2):105-114. doi:10.1079/bjn19910073. PMID 2043597. Abstract (adults, R² 0.79, SEE 4.1%)
// Tags:    anthropometry, body-fat, bmi, deurenberg, adult
// Units:   bmiKgM2 kg/m²; ageYears years (decimal); male the source's sex term; result percent body fat
// Range:   15 < ageYears ≤ 83; 13.9 ≤ bmiKgM2 ≤ 40.9 (the study's 1229 subjects, ages 7–83, BMI 13.9–40.9; the adult formula above age 15)
// Example: male, 40 y, BMI 25 → 30 + 9.2 − 10.8 − 5.4 = 23.0% (HAND-DERIVED)
func BodyFatPercentDeurenbergAdult(bmiKgM2, ageYears float64, male bool) (float64, error) {
	if err := checkFinite("bmiKgM2", bmiKgM2); err != nil {
		return 0, err
	}
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if err := checkDeurenbergBMI(bmiKgM2); err != nil {
		return 0, err
	}
	if !(ageYears > 15) || ageYears > 83 {
		return 0, fmt.Errorf("%w: ageYears = %v, want 15 < ageYears ≤ 83", ErrOutOfRange, ageYears)
	}
	return 1.20*bmiKgM2 + 0.23*ageYears - 10.8*sexIndicator(male) - 5.4, nil
}

// BodyFatPercentDeurenbergChild returns the body-fat percentage predicted from BMI, age and sex by the
// Deurenberg child formula: 1.51 × BMI − 0.70 × age − 3.6 × sex + 1.4, with sex = 1 for male and 0 for
// female (the source's coding). The source reports R² 0.38 and SEE 4.4% for children aged 15 and under.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-ANTHRO-009
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Deurenberg P, Weststrate JA, Seidell JC. Body mass index as a measure of body fatness: age- and sex-specific prediction formulas. Br J Nutr. 1991;65(2):105-114. doi:10.1079/bjn19910073. PMID 2043597. Abstract (children ≤ 15 y, R² 0.38, SEE 4.4%)
// Tags:    anthropometry, body-fat, bmi, deurenberg, pediatric
// Units:   bmiKgM2 kg/m²; ageYears years (decimal); male the source's sex term; result percent body fat
// Range:   7 ≤ ageYears ≤ 15; 13.9 ≤ bmiKgM2 ≤ 40.9 (the study's subjects, ages 7–83, BMI 13.9–40.9; the child formula to age 15)
// Example: female, 10 y, BMI 18 → 27.18 − 7 + 1.4 = 21.58% (HAND-DERIVED)
func BodyFatPercentDeurenbergChild(bmiKgM2, ageYears float64, male bool) (float64, error) {
	if err := checkFinite("bmiKgM2", bmiKgM2); err != nil {
		return 0, err
	}
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if err := checkDeurenbergBMI(bmiKgM2); err != nil {
		return 0, err
	}
	if ageYears < 7 || ageYears > 15 {
		return 0, fmt.Errorf("%w: ageYears = %v, want 7 ≤ ageYears ≤ 15", ErrOutOfRange, ageYears)
	}
	return 1.51*bmiKgM2 - 0.70*ageYears - 3.6*sexIndicator(male) + 1.4, nil
}
