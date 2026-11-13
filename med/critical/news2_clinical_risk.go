package critical

import "fmt"

// NEWS2Risk is a NEWS2 clinical-risk band; its value is RCP Chart 2's own label, verbatim.
type NEWS2Risk string

// RCP 2017 Chart 2 "NEWS thresholds and triggers", column "Clinical risk".
const (
	NEWS2RiskLow       NEWS2Risk = "Low"        // aggregate score 0–4
	NEWS2RiskLowMedium NEWS2Risk = "Low–medium" // red score: score of 3 in any individual parameter
	NEWS2RiskMedium    NEWS2Risk = "Medium"     // aggregate score 5–6
	NEWS2RiskHigh      NEWS2Risk = "High"       // aggregate score 7 or more
)

// NEWS2ClinicalRisk returns the NEWS2 clinical-risk band of RCP Chart 2 for an aggregate score and the red-score
// flag (a score of 3 in any individual parameter, NEWS2Scores.RedScore): 7 or more "High"; 5–6 "Medium";
// otherwise "Low–medium" when redScore is true and "Low" when it is not. The chart lists the red score beside
// aggregates 0–4; an aggregate of 5 or more takes its own band whatever the flag.
// Errors: ErrOutOfRange when aggregate is outside 0–20.
//
// ID:      MED-CRIT-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Royal College of Physicians. National Early Warning Score (NEWS) 2. London: RCP, 2017. Chart 2 "NEWS thresholds and triggers", column "Clinical risk"
// Tags:    critical, news2, early warning, category
// Units:   aggregate points 0–20; result the chart's label
// Range:   0 ≤ aggregate ≤ 20
// Example: 4, false → "Low"; 3, true → "Low–medium"; 5 → "Medium"; 7 → "High" (Chart 2)
func NEWS2ClinicalRisk(aggregate int, redScore bool) (NEWS2Risk, error) {
	if aggregate < 0 || aggregate > 20 {
		return "", fmt.Errorf("%w: aggregate = %d, want 0–20", ErrOutOfRange, aggregate)
	}
	switch {
	case aggregate >= 7:
		return NEWS2RiskHigh, nil
	case aggregate >= 5:
		return NEWS2RiskMedium, nil
	case redScore:
		return NEWS2RiskLowMedium, nil
	}
	return NEWS2RiskLow, nil
}
