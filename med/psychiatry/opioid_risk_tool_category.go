package psychiatry

import "fmt"

// ORTCategory is an Opioid Risk Tool band; its value is the source's words, verbatim.
type ORTCategory string

// Webster & Webster 2005 abstract: "scores of 0-3 (low risk), 4-7 (moderate risk), or >or= 8 (high risk)".
const (
	ORTLowRisk      ORTCategory = "low risk"      // 0–3
	ORTModerateRisk ORTCategory = "moderate risk" // 4–7
	ORTHighRisk     ORTCategory = "high risk"     // ≥ 8
)

// OpioidRiskToolCategory returns the Opioid Risk Tool band for a total score: 0–3 "low risk"; 4–7 "moderate
// risk"; ≥ 8 "high risk". The item weights that make the total were not read in the primary, so the total is
// the caller's; only the published bands are here.
// Errors: ErrOutOfRange when total is negative.
//
// ID:      MED-PSYCH-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Webster LR, Webster RM. Predicting aberrant behaviors in opioid-treated patients: preliminary validation of the Opioid Risk Tool. Pain Med. 2005;6(6):432-442. doi:10.1111/j.1526-4637.2005.00072.x, PMID 16336480, abstract
// Tags:    psychiatry, substance use, opioid, ort, category
// Units:   total points; result the source's words
// Range:   total ≥ 0
// Example: 3 → "low risk"; 4 → "moderate risk"; 7 → "moderate risk"; 8 → "high risk" (abstract)
func OpioidRiskToolCategory(total int) (ORTCategory, error) {
	if total < 0 {
		return "", fmt.Errorf("%w: total = %d, want ≥ 0", ErrOutOfRange, total)
	}
	switch {
	case total <= 3:
		return ORTLowRisk, nil
	case total <= 7:
		return ORTModerateRisk, nil
	}
	return ORTHighRisk, nil
}
