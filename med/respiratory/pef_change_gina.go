package respiratory

import "math"

// PEFChangeConsistentWithAsthmaGINA reports whether a change in peak expiratory flow meets GINA's
// threshold for adults: "in adults with respiratory symptoms typical of asthma, an increase or decrease in
// FEV1 of >12% and >200 mL from baseline, or (if spirometry is not available) a change in PEF of at least
// 20%, is accepted as being consistent with asthma." The PEF part only: |pefChangePercent| ≥ 20, an
// increase (for example after a bronchodilator) or a decrease. pefChangePercent is (later − baseline) /
// baseline × 100, computed by the caller; GINA notes PEF may vary by up to 20 % between meters, so both
// readings come from one meter. Exactly ±20 % is true ("at least").
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for pefChangePercent < −100 or ageYears < 18.
//
// ID:      MED-RESP-053
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Global Initiative for Asthma. Global Strategy for Asthma Management and Prevention, 2023 update. https://ginasthma.org/wp-content/uploads/2023/07/GINA-2023-Full-report-23_07_06-WMS.pdf (read 2026-11-12), section "How much variation in expiratory airflow is consistent with asthma?" (and Box 1-2 footnote †)
// Tags:    respiratory, asthma, peak-flow, pef, reversibility, bronchodilator, gina, criterion
// Units:   pefChangePercent percent of the baseline PEF (signed); ageYears years (decimal); result true when |change| ≥ 20 %
// Range:   pefChangePercent ≥ −100 (a PEF cannot fall below zero); ageYears ≥ 18 (the sentence states adults)
// Example: +20 % → true; +19.9 % → false; −25 % → true; +35 % → true
func PEFChangeConsistentWithAsthmaGINA(pefChangePercent, ageYears float64) (bool, error) {
	if err := checkClosedRange("pefChangePercent", pefChangePercent, -100, math.Inf(1)); err != nil {
		return false, err
	}
	if err := checkClosedRange("ageYears", ageYears, 18, math.Inf(1)); err != nil {
		return false, err
	}
	// GINA 2023: "a change in PEF of at least 20%".
	return math.Abs(pefChangePercent) >= 20, nil
}
