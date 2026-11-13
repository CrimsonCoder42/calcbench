package hematology

import "fmt"

// MALTIPIGroup is a MALT-IPI risk group; its value is the source's word, verbatim.
type MALTIPIGroup string

// Thieblemont 2017 abstract: "low, intermediate, and high risk (corresponding to the presence of 0, 1, or ≥2
// of these factors, respectively)".
const (
	MALTIPILow          MALTIPIGroup = "low"          // 0
	MALTIPIIntermediate MALTIPIGroup = "intermediate" // 1
	MALTIPIHigh         MALTIPIGroup = "high"         // ≥ 2
)

// MALTIPIRiskGroup returns the MALT-IPI risk group for a score: 0 "low"; 1 "intermediate"; ≥ 2 "high".
// Errors: ErrOutOfRange when score is outside 0–3.
//
// ID:      MED-HEMAT-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Thieblemont C, Cascione L, Conconi A, et al. A MALT lymphoma prognostic index. Blood. 2017;130(12):1409-1417. doi:10.1182/blood-2017-03-771915, PMID 28720586, abstract
// Tags:    hematology, lymphoma, malt, malt-ipi, category
// Units:   score factors 0–3; result the source's word
// Range:   0 ≤ score ≤ 3
// Example: 0 → "low"; 1 → "intermediate"; 2 → "high"; 3 → "high" (abstract)
func MALTIPIRiskGroup(score int) (MALTIPIGroup, error) {
	if score < 0 || score > 3 {
		return "", fmt.Errorf("%w: score = %d, want 0–3", ErrOutOfRange, score)
	}
	switch score {
	case 0:
		return MALTIPILow, nil
	case 1:
		return MALTIPIIntermediate, nil
	}
	return MALTIPIHigh, nil
}
