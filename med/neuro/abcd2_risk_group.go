package neuro

import "fmt"

// ABCD2RiskGroup is an ABCD² stratum; its value is the source's words, verbatim.
type ABCD2RiskGroup string

// Johnston 2007 abstract: "high risk (score 6-7 ...), ... moderate risk (score 4-5 ...), and ... low risk (score 0-3 ...)".
const (
	ABCD2LowRisk      ABCD2RiskGroup = "low risk"      // 0–3
	ABCD2ModerateRisk ABCD2RiskGroup = "moderate risk" // 4–5
	ABCD2HighRisk     ABCD2RiskGroup = "high risk"     // 6–7
)

// ABCD2RiskGroupOf returns the ABCD² stratum for a total: 0–3 "low risk"; 4–5 "moderate risk"; 6–7 "high risk".
// Errors: ErrOutOfRange when total is outside 0–7.
//
// ID:      MED-NEURO-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnston SC, Rothwell PM, Nguyen-Huynh MN, et al. Validation and refinement of scores to predict very early stroke risk after transient ischaemic attack. Lancet. 2007;369(9558):283-292. doi:10.1016/S0140-6736(07)60150-0, PMID 17258668, abstract (strata and their score ranges)
// Tags:    neuro, stroke, tia, abcd2, category
// Units:   total points 0–7; result the source's words
// Range:   total 0–7
// Example: 3 → "low risk"; 4 → "moderate risk"; 6 → "high risk" (abstract)
func ABCD2RiskGroupOf(total int) (ABCD2RiskGroup, error) {
	if total < 0 || total > 7 {
		return "", fmt.Errorf("%w: total = %d, want 0–7", ErrOutOfRange, total)
	}
	switch {
	case total <= 3:
		return ABCD2LowRisk, nil
	case total <= 5:
		return ABCD2ModerateRisk, nil
	}
	return ABCD2HighRisk, nil
}
