package renal

import "fmt"

// STONEGroup is a STONE-score probability group, labelled as Moore 2014 labels it.
type STONEGroup string

const (
	STONELowProbability      STONEGroup = "low probability"      // score 0–5
	STONEModerateProbability STONEGroup = "moderate probability" // score 6–9
	STONEHighProbability     STONEGroup = "high probability"     // score 10–13
)

// STONEProbabilityGroup returns the STONE-score group of Moore 2014: "low probability (score 0-5)", "moderate
// probability (score 6-9)", "high probability (score 10-13)".
// Errors: ErrOutOfRange when score is outside 0–13.
//
// ID:      MED-RENAL-014
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Moore CL, et al. BMJ. 2014;348:g2191. doi:10.1136/bmj.g2191. PMID 24671981, PMC3966515, Results: "low probability (score 0-5) group ... moderate probability (score 6-9) group, and ... high probability (score 10-13) group"
// Tags:    renal, urology, ureteral stone, stone score, category
// Units:   score points 0–13; result label
// Range:   0 ≤ score ≤ 13
// Example: 5 → "low probability"; 6 → "moderate probability"; 10 → "high probability" (from the Results)
func STONEProbabilityGroup(score int) (STONEGroup, error) {
	switch {
	case score < 0 || score > 13:
		return "", fmt.Errorf("%w: score = %d, want 0–13", ErrOutOfRange, score)
	case score <= 5:
		return STONELowProbability, nil
	case score <= 9:
		return STONEModerateProbability, nil
	default:
		return STONEHighProbability, nil
	}
}
