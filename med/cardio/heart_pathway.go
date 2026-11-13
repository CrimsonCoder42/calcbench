package cardio

import "fmt"

// HEARTPathwayLowRisk reports whether the HEART Pathway's low-risk condition is met: a HEART score of 0–3
// and negative troponin results at 0 and 3 hours after ED presentation.
// Errors: ErrOutOfRange when heartScore is outside 0–10.
//
// ID:      MED-CARDIO-058
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Mahler SA, et al. Circ Cardiovasc Qual Outcomes. 2015;8(2):195-203. doi:10.1161/CIRCOUTCOMES.114.001384. PMID 25737484. PMC4413911, Methods: "low-risk HEART scores (HEART score of 0–3) and negative troponin results", troponin "at 0 and 3 hours after ED presentation"
// Tags:    cardio, chest-pain, acs, heart-pathway, decision-rule
// Units:   heartScore points 0–10; result true when low risk
// Range:   0 ≤ heartScore ≤ 10
// Example: score 3, both troponins negative → true; score 4 → false
func HEARTPathwayLowRisk(heartScore int, troponin0hNegative, troponin3hNegative bool) (bool, error) {
	if heartScore < 0 || heartScore > 10 {
		return false, fmt.Errorf("%w: heartScore = %d, want 0–10", ErrOutOfRange, heartScore)
	}
	return heartScore <= 3 && troponin0hNegative && troponin3hNegative, nil
}
