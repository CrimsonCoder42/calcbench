package hematology

// CRPUpperReferenceLimitWenerMgDL returns the age- and sex-adjusted upper reference limit of serum C-reactive
// protein, in mg/dL, by Wener 2000: "CRP = age/50 for males, and CRP = age/50 + 0.6 for females", proposed for
// ages 25–70 years. The same abstract gives separate race-specific coefficients for non-Hispanic Black adults
// (age/30 for males, age/50 + 1.0 for females); they are not applied here.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-HEMAT-017
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Wener MH, Daum PR, McQuillan GM. The influence of age, sex, and race on the upper reference limit of serum C-reactive protein concentration. J Rheumatol. 2000;27(10):2351-2359. PMID 11036829, abstract (Results)
// Tags:    hematology, inflammation, c-reactive protein, crp, reference limit
// Units:   ageYears years; result CRP upper reference limit in mg/dL
// Range:   25 ≤ ageYears ≤ 70 (the ages the formulas were proposed for)
// Example: man 50 → 1.0 mg/dL; woman 50 → 1.6 mg/dL; man 25 → 0.5; woman 70 → 2.0 (abstract formulas)
func CRPUpperReferenceLimitWenerMgDL(ageYears float64, female bool) (float64, error) {
	if err := checkFloat("ageYears", ageYears, 25, 70); err != nil {
		return 0, err
	}
	if female {
		return ageYears/50 + 0.6, nil
	}
	return ageYears / 50, nil
}
