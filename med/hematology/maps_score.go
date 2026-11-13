package hematology

// MAPSScore returns the Mayo Alliance Prognostic System clinical-molecular score for systemic mastocytosis,
// 0–6 (Figure 3, Results): advanced SM (vs indolent or smoldering SM) 2; age > 60 years 1; platelets
// < 150 × 10⁹/L 1; serum alkaline phosphatase above the normal range 1; an adverse mutation (ASXL1, RUNX1 or
// NRAS) 1. The paper's risk tiers are not built here (their labels were not read).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-HEMAT-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Pardanani A, Shah S, Mannelli F, et al. Mayo alliance prognostic system for mastocytosis: clinical and hybrid clinical-molecular models. Blood Adv. 2018;2(21):2964-2972. doi:10.1182/bloodadvances.2018026245, PMID 30413432, Figure 3 and Results
// Tags:    hematology, mastocytosis, maps, prognostic score
// Units:   ageYears years; plateletsE9PerL × 10⁹/L; findings present/absent; result points 0–6
// Range:   0 ≤ ageYears ≤ 130; 0 ≤ plateletsE9PerL ≤ 5000 (physical limits)
// Example: advanced SM, age 61, platelets 149, ALP high, ASXL1 → 6; indolent, age 60, platelets 150 → 0
func MAPSScore(advancedSM bool, ageYears, plateletsE9PerL float64, alpAboveNormal, adverseMutation bool) (int, error) {
	if err := checkFloat("ageYears", ageYears, 0, 130); err != nil {
		return 0, err
	}
	if err := checkFloat("plateletsE9PerL", plateletsE9PerL, 0, 5000); err != nil {
		return 0, err
	}
	return pts(advancedSM, 2) + pts(ageYears > 60, 1) + pts(plateletsE9PerL < 150, 1) +
		pts(alpAboveNormal, 1) + pts(adverseMutation, 1), nil
}
