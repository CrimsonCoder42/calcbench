package respiratory

import (
	"fmt"
	"math"
)

// PESARDScore returns the PE-SARD bleeding score (0–5 points, in steps of 0.5): syncope +1.5, anaemia
// (haemoglobin < 12 g/dL) +2.5, renal dysfunction (glomerular filtration rate < 60 mL/min) +1.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-RESP-021
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Chopard R, Piazza G, Falvo N, et al. An Original Risk Score to Predict Early Major Bleeding in Acute Pulmonary Embolism: The Syncope, Anemia, Renal Dysfunction (PE-SARD) Bleeding Score. Chest. 2021;160(5):1832-1843. doi:10.1016/j.chest.2021.06.048. PMID 34217683, abstract: "Syncope (+1.5); Anemia, defined as hemoglobin <12 g/dL (+2.5); and Renal Dysfunction, defined as glomerular filtration rate <60 mL/min (+1 point)"
// Tags:    respiratory, pulmonary-embolism, bleeding, pe-sard, score
// Units:   hemoglobinGDL g/dL; gfrMLMin mL/min as the source states it; result points 0–5 (exact halves)
// Range:   hemoglobinGDL > 0; gfrMLMin ≥ 0 (no limits stated in the abstract)
// Example: syncope, Hb 11 g/dL, GFR 75 → 1.5 + 2.5 = 4 (hand-derived from the abstract's points)
func PESARDScore(syncope bool, hemoglobinGDL, gfrMLMin float64) (float64, error) {
	if err := checkPositiveFinite("hemoglobinGDL", hemoglobinGDL); err != nil {
		return 0, err
	}
	if err := checkClosedRange("gfrMLMin", gfrMLMin, 0, math.Inf(1)); err != nil {
		return 0, err
	}
	s := 0.0 // sums of 1.5, 2.5 and 1 are exact in binary floating point
	if syncope {
		s += 1.5
	}
	if hemoglobinGDL < 12 {
		s += 2.5
	}
	if gfrMLMin < 60 {
		s += 1
	}
	return s, nil
}

// PESARDRisk is the PE-SARD bleeding-risk group for a score.
type PESARDRisk string

// The groups, verbatim from the abstract: "low bleeding-risk (score, 0 point), ... intermediate-risk (score,
// 1-2.5 points), and ... high-risk (score >2.5 points)".
const (
	PESARDLow          PESARDRisk = "low bleeding-risk"
	PESARDIntermediate PESARDRisk = "intermediate-risk"
	PESARDHigh         PESARDRisk = "high-risk"
)

// PESARDRiskGroup returns the PE-SARD bleeding-risk group for a score: 0 low, 1–2.5 intermediate, > 2.5 high.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a value that is not a PE-SARD score
// (0, 1, 1.5, 2.5, 3.5, 4 or 5).
//
// ID:      MED-RESP-022
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Chopard R, Piazza G, Falvo N, et al. An Original Risk Score to Predict Early Major Bleeding in Acute Pulmonary Embolism: The Syncope, Anemia, Renal Dysfunction (PE-SARD) Bleeding Score. Chest. 2021;160(5):1832-1843. doi:10.1016/j.chest.2021.06.048. PMID 34217683, abstract (the three groups and their point ranges)
// Tags:    respiratory, pulmonary-embolism, bleeding, pe-sard, category
// Units:   score points (PESARDScore's result); result the source's group
// Range:   score one of 0, 1, 1.5, 2.5, 3.5, 4, 5
// Example: 0 → "low bleeding-risk"; 2.5 → "intermediate-risk"; 3.5 → "high-risk" (the abstract's ranges)
func PESARDRiskGroup(score float64) (PESARDRisk, error) {
	if math.IsNaN(score) || math.IsInf(score, 0) {
		return "", fmt.Errorf("%w: score = %v", ErrNotFinite, score)
	}
	switch score {
	case 0:
		return PESARDLow, nil
	case 1, 1.5, 2.5:
		return PESARDIntermediate, nil
	case 3.5, 4, 5:
		return PESARDHigh, nil
	}
	return "", fmt.Errorf("%w: score = %v, want one of 0, 1, 1.5, 2.5, 3.5, 4, 5", ErrOutOfRange, score)
}
