package anthropometry

// RelativeFatMass returns the relative fat mass (RFM), an estimate of whole-body fat percentage:
// 64 − 20 × (height / waist circumference) + 12 × sex, with sex = 0 for men and 1 for women. Height and waist
// enter only as their ratio, so both are taken in cm. Not rounded.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-ANTHRO-011
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Woolcott OO, Bergman RN. Relative fat mass (RFM) as a new estimator of whole-body fat percentage ─ a cross-sectional study in American adult individuals. Sci Rep. 2018;8(1):10980. doi:10.1038/s41598-018-29362-1. PMID 30030479. Abstract: "64 - (20 × height/waist circumference) + (12 × sex), named as the relative fat mass (RFM); sex = 0 for men and 1 for women"
// Tags:    anthropometry, body-fat, rfm, relative-fat-mass, waist
// Units:   heightCm cm; waistCm cm (the same unit as height); result percent body fat
// Range:   heightCm > 0; waistCm > 0 (population: NHANES 1999–2004 adults, n = 12,581, validated in NHANES 2005–2006; input limits not stated in the abstract)
// Example: man, height 175 cm, waist 88 cm → 64 − 20 × 1.98864 = 24.2273; woman, 165 cm, 80 cm → 34.75 (HAND-DERIVED from the abstract)
func RelativeFatMass(heightCm, waistCm float64, female bool) (float64, error) {
	if err := checkFinite("heightCm", heightCm); err != nil {
		return 0, err
	}
	if err := checkFinite("waistCm", waistCm); err != nil {
		return 0, err
	}
	if err := checkPositive("heightCm", heightCm); err != nil {
		return 0, err
	}
	if err := checkPositive("waistCm", waistCm); err != nil {
		return 0, err
	}
	sex := 0.0
	if female {
		sex = 1
	}
	r := 64 - 20*(heightCm/waistCm) + 12*sex
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}
