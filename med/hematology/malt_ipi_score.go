package hematology

// MALTIPIScore returns MALT-IPI, 0–3: the number of its risk factors present: "age ≥70 years", "Ann Arbor
// stage III or IV", and "an elevated lactate dehydrogenase level".
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when ageYears is outside 0–130.
//
// ID:      MED-HEMAT-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Thieblemont C, Cascione L, Conconi A, et al. A MALT lymphoma prognostic index. Blood. 2017;130(12):1409-1417. doi:10.1182/blood-2017-03-771915, PMID 28720586, abstract
// Tags:    hematology, lymphoma, malt, malt-ipi, prognostic score
// Units:   ageYears years; findings present/absent; result factors 0–3
// Range:   0 ≤ ageYears ≤ 130 (physical limits)
// Example: age 70, stage III, LDH normal → 2; age 69.9, stage I, LDH normal → 0
func MALTIPIScore(ageYears float64, stageIIIorIV, ldhElevated bool) (int, error) {
	if err := checkFloat("ageYears", ageYears, 0, 130); err != nil {
		return 0, err
	}
	return pts(ageYears >= 70, 1) + pts(stageIIIorIV, 1) + pts(ldhElevated, 1), nil
}
