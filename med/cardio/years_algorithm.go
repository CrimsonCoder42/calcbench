package cardio

import "fmt"

// YEARSPulmonaryEmbolismExcluded reports whether the YEARS algorithm considers pulmonary embolism excluded:
// no YEARS item and D-dimer < 1000 ng/mL, or one or more YEARS items and D-dimer < 500 ng/mL. The three items
// are clinical signs of deep vein thrombosis, haemoptysis, and pulmonary embolism as the most likely diagnosis.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-CARDIO-062
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  van der Hulle T, Cheung WY, Kooij S, et al. Simplified diagnostic management of suspected pulmonary embolism (the YEARS study): a prospective, multicentre, cohort study. Lancet. 2017;390(10091):289-297. doi:10.1016/S0140-6736(17)30885-1. PMID 28549662. Abstract: "without YEARS items and D-dimer less than 1000 ng/mL", or "one or more YEARS items and D-dimer less than 500 ng/mL", "pulmonary embolism was considered excluded"
// Tags:    cardio, pulmonary-embolism, d-dimer, years, decision-rule
// Units:   yearsItems count 0–3; dDimerNgML ng/mL (the assay units the study reports); result true when excluded
// Range:   yearsItems 0–3; dDimerNgML 0–100000 (persons with clinically suspected PE)
// Example: 0 items, D-dimer 900 → true; 1 item, D-dimer 600 → false
func YEARSPulmonaryEmbolismExcluded(yearsItems int, dDimerNgML float64) (bool, error) {
	if err := checkFinite("dDimerNgML", dDimerNgML); err != nil {
		return false, err
	}
	if yearsItems < 0 || yearsItems > 3 {
		return false, fmt.Errorf("%w: yearsItems = %d, want 0–3", ErrOutOfRange, yearsItems)
	}
	if err := checkRange("dDimerNgML", dDimerNgML, 0, 100000); err != nil {
		return false, err
	}
	if yearsItems == 0 {
		return dDimerNgML < 1000, nil
	}
	return dDimerNgML < 500, nil
}
