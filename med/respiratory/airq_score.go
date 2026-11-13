package respiratory

import "fmt"

// AIRQControl is the asthma-control level the AIRQ development study assigns to a total score.
type AIRQControl string

// The AIRQ control levels, worded from the source's abstract: "Cut points of 0-1, 2-4, and 5-10 best
// represented well-, not well-, and very poorly controlled asthma."
const (
	AIRQWellControlled       AIRQControl = "well-controlled"        // total 0–1
	AIRQNotWellControlled    AIRQControl = "not well-controlled"    // total 2–4
	AIRQVeryPoorlyControlled AIRQControl = "very poorly controlled" // total 5–10
)

// AIRQScore returns the Asthma Impairment and Risk Questionnaire total (0–10): the count of "yes" answers
// among its 10 equally weighted yes/no questions. The question text is the instrument's and is not
// reproduced here; the caller passes the answers in the instrument's order.
// Every input is valid, so there is no error result.
//
// ID:      MED-RESP-004
// Time:    O(1) (10 items)
// Space:   O(1), no allocation
// Source:  Murphy KR, Chipps B, Beuther DA, et al. Development of the Asthma Impairment and Risk Questionnaire (AIRQ): A Composite Control Measure. J Allergy Clin Immunol Pract. 2020;8(7):2263-2274.e5. doi:10.1016/j.jaip.2020.02.042. PMID 32387166, abstract: "The final AIRQ comprises 10 equally weighted yes/no impairment and risk questions."
// Tags:    respiratory, asthma, airq, questionnaire, control, score
// Units:   answersYes one bool per question, true = "yes"; result points 0–10
// Range:   ages ≥ 12 (the study enrolled people aged 12 years and over with physician-diagnosed asthma); not enforced, no age input
// Example: 3 "yes" answers → 3 (hand-derived: one point per "yes")
func AIRQScore(answersYes [10]bool) int {
	n := 0
	for _, yes := range answersYes {
		if yes {
			n++
		}
	}
	return n
}

// AIRQControlLevel returns the AIRQ control level for a total score: 0–1 well-controlled, 2–4 not well-controlled,
// 5–10 very poorly controlled.
// Errors: ErrOutOfRange for a total outside 0–10.
//
// ID:      MED-RESP-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Murphy KR, Chipps B, Beuther DA, et al. Development of the Asthma Impairment and Risk Questionnaire (AIRQ): A Composite Control Measure. J Allergy Clin Immunol Pract. 2020;8(7):2263-2274.e5. doi:10.1016/j.jaip.2020.02.042. PMID 32387166, abstract: "Cut points of 0-1, 2-4, and 5-10"
// Tags:    respiratory, asthma, airq, control, category
// Units:   total points 0–10; result the source's control level
// Range:   total 0–10
// Example: 4 → "not well-controlled"; 5 → "very poorly controlled" (the abstract's cut points)
func AIRQControlLevel(total int) (AIRQControl, error) {
	switch {
	case total < 0 || total > 10:
		return "", fmt.Errorf("%w: total = %d, want 0 ≤ total ≤ 10", ErrOutOfRange, total)
	case total <= 1:
		return AIRQWellControlled, nil
	case total <= 4:
		return AIRQNotWellControlled, nil
	default:
		return AIRQVeryPoorlyControlled, nil
	}
}
