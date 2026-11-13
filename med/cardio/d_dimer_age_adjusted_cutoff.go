package cardio

// DDimerAgeAdjustedCutoffUgL returns Douma's age-adjusted D-dimer cut-off, in µg/L: age × 10 for persons aged
// over 50, the conventional 500 µg/L otherwise. (ADJUST-PE, Righini 2014, applies age × 10 from 50 years; at
// exactly 50 both give 500.)
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when ageYears is outside 18–120.
//
// ID:      MED-CARDIO-042
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Douma RA, le Gal G, Söhne M, et al. Potential of an age adjusted D-dimer cut-off value to improve the exclusion of pulmonary embolism in older patients. BMJ. 2010;340:c1475. doi:10.1136/bmj.c1475. PMID 20354012. PMC2847688: "(patient's age×10) μg/l in patients aged >50"; conventional "500 μg/l"
// Source:  Schouten HJ, Geersing GJ, Koek HL, et al. Diagnostic accuracy of conventional or age adjusted D-dimer cut-off values in older patients with suspected venous thromboembolism: systematic review and meta-analysis. BMJ. 2013;346:f2492. doi:10.1136/bmj.f2492. PMID 23645857. PMC3643284: "for a patient aged 78 years, the D-dimer concentration would be considered normal below 780 μg/L"
// Tags:    cardio, vte, pulmonary-embolism, d-dimer, age-adjusted
// Units:   ageYears years; result µg/L (the assay's own units; the papers validated VIDAS and Tinaquant)
// Range:   ageYears 18–120 (adults with suspected PE in the three cohorts)
// Example: age 78 → 780 µg/L (Schouten 2013); age 40 → 500
func DDimerAgeAdjustedCutoffUgL(ageYears float64) (float64, error) {
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if err := checkRange("ageYears", ageYears, 18, 120); err != nil {
		return 0, err
	}
	if ageYears > 50 {
		return ageYears * 10, nil
	}
	return 500, nil
}
