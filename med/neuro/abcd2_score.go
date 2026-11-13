package neuro

import "math"

// ABCD2Score returns the ABCD² score (0–7) after transient ischaemic attack: age ≥ 60 years 1; blood pressure
// ≥ 140/90 mm Hg 1; clinical features: unilateral weakness 2, speech impairment without weakness 1; duration
// ≥ 60 min 2 or 10–59 min 1; diabetes 1.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-NEURO-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnston SC, Rothwell PM, Nguyen-Huynh MN, et al. Validation and refinement of scores to predict very early stroke risk after transient ischaemic attack. Lancet. 2007;369(9558):283-292. doi:10.1016/S0140-6736(07)60150-0, PMID 17258668, abstract ("age >or=60 years [1 point]; blood pressure >or=140/90 mm Hg [1]; clinical features: unilateral weakness [2], speech impairment without weakness [1]; duration >or=60 min [2] or 10-59 min [1]; and diabetes [1]")
// Tags:    neuro, stroke, tia, abcd2, score
// Units:   ageYears years; bpAtLeast140Over90 the caller's reading of "≥ 140/90 mm Hg"; durationMinutes minutes; result points 0–7
// Range:   ageYears ≥ 0; durationMinutes ≥ 0 (the abstract states no limits)
// Example: age 70, BP 150/95, unilateral weakness, 45 min, no diabetes → 1 + 1 + 2 + 1 + 0 = 5 (hand-derived from the abstract)
func ABCD2Score(ageYears float64, bpAtLeast140Over90, unilateralWeakness, speechImpairmentWithoutWeakness bool, durationMinutes float64, diabetes bool) (int, error) {
	if err := checkRange("ageYears", ageYears, 0, math.MaxFloat64); err != nil {
		return 0, err
	}
	if err := checkRange("durationMinutes", durationMinutes, 0, math.MaxFloat64); err != nil {
		return 0, err
	}
	n := 0
	if ageYears >= 60 {
		n++
	}
	if bpAtLeast140Over90 {
		n++
	}
	switch {
	case unilateralWeakness:
		n += 2
	case speechImpairmentWithoutWeakness:
		n++
	}
	switch {
	case durationMinutes >= 60:
		n += 2
	case durationMinutes >= 10:
		n++
	}
	if diabetes {
		n++
	}
	return n, nil
}
