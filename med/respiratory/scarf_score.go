package respiratory

import (
	"fmt"
	"math"
)

// SCARFScore returns the Sequential Clinical Assessment of Respiratory Function score (0–4) for rib fractures:
// "One point was assigned for: <50% predicted, respiratory rate >20, numeric pain score ≥5, and inadequate cough."
// The "<50% predicted" measurement is the caller's, passed as belowHalfPredicted.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-RESP-025
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Hardin KS, Leasia KN, Haenel J, Moore EE, Burlew CC, Pieracci FM. The Sequential Clinical Assessment of Respiratory Function (SCARF) score: A dynamic pulmonary physiologic score that predicts adverse outcomes in critically ill rib fracture patients. J Trauma Acute Care Surg. 2019;87(6):1260-1268. doi:10.1097/TA.0000000000002480. PMID 31425473, abstract
// Tags:    respiratory, rib-fracture, trauma, scarf, score
// Units:   respiratoryRatePerMin breaths/min; numericPainScore 0–10; result points 0–4
// Range:   respiratoryRatePerMin ≥ 0; numericPainScore 0–10
// Example: below 50 % predicted, RR 24, pain 4, adequate cough → 2 (hand-derived from the abstract's items)
func SCARFScore(belowHalfPredicted bool, respiratoryRatePerMin float64, numericPainScore int, inadequateCough bool) (int, error) {
	if err := checkClosedRange("respiratoryRatePerMin", respiratoryRatePerMin, 0, math.Inf(1)); err != nil {
		return 0, err
	}
	if numericPainScore < 0 || numericPainScore > 10 {
		return 0, fmt.Errorf("%w: numericPainScore = %d, want 0 ≤ numericPainScore ≤ 10", ErrOutOfRange, numericPainScore)
	}
	n := 0
	if belowHalfPredicted {
		n++
	}
	if respiratoryRatePerMin > 20 {
		n++
	}
	if numericPainScore >= 5 {
		n++
	}
	if inadequateCough {
		n++
	}
	return n, nil
}
