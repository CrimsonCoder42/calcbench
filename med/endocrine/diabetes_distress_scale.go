package endocrine

import "fmt"

// DiabetesDistressScore17 returns the Diabetes Distress Scale (DDS17) score: the mean of the 17 item ratings,
// each rated 1 ("no problem") to 6 ("a serious problem"), so the result is 1–6. Item wording is not embedded
// (the instrument is the authors'); the caller passes the 17 ratings in the instrument's order. The subscale
// means are not computed here (their item assignment was not read in an open source).
// Errors: ErrOutOfRange when an item is outside 1–6.
//
// ID:      MED-ENDO-009
// Time:    O(1)   17 items
// Space:   O(1), no allocation
// Source:  Fisher L, Hessler DM, Polonsky WH, Mullan J. When is diabetes distress clinically meaningful? Establishing cut points for the Diabetes Distress Scale. Diabetes Care. 2012;35(2):259-264. doi:10.2337/dc11-1572. PMID 22228744. PMC3263871, Measures: "a Likert scale to score each item from 1 (no problem) to 6 (a serious problem)"; "Mean-item scores are then calculated (DDS17)"
// Source:  Polonsky WH, Fisher L, Earles J, et al. Assessing psychosocial distress in diabetes: development of the diabetes distress scale. Diabetes Care. 2005;28(3):626-631. doi:10.2337/diacare.28.3.626. PMID 15735199 (the instrument; full text not read)
// Tags:    endocrine, diabetes, distress, dds17, questionnaire, score
// Units:   items points 1–6 each; result mean points 1–6
// Range:   each item 1–6
// Example: all items 2 → 34 / 17 = 2.0; all items 3 → 3.0 (HAND-DERIVED)
func DiabetesDistressScore17(items [17]int) (float64, error) {
	sum := 0
	for i, v := range items {
		if v < 1 || v > 6 {
			return 0, fmt.Errorf("%w: items[%d] = %d, want 1–6", ErrOutOfRange, i, v)
		}
		sum += v
	}
	return float64(sum) / 17, nil
}

// DiabetesDistressCategory is a DDS17 category with the label of Fisher 2012, verbatim.
type DiabetesDistressCategory string

const (
	DiabetesDistressLittleOrNone DiabetesDistressCategory = "little or no DD" // DDS < 2.0
	DiabetesDistressModerate     DiabetesDistressCategory = "moderate DD"     // DDS 2.0–2.9
	DiabetesDistressHigh         DiabetesDistressCategory = "high DD"         // DDS ≥ 3.0
)

// DiabetesDistressCategoryFisher returns the Fisher 2012 category for a DDS17 mean-item score: "little or no
// DD" below 2.0, "moderate DD" from 2.0 to below 3.0, "high DD" at 3.0 and above. The source prints the middle
// band as "2.0–2.9"; a mean between 2.9 and 3.0 (for example 50/17 ≈ 2.94) is placed in it, since the source's
// next band starts at 3.0.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-ENDO-010
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Fisher L, Hessler DM, Polonsky WH, Mullan J. Diabetes Care. 2012;35(2):259-264. doi:10.2337/dc11-1572. PMID 22228744. PMC3263871, Conclusions: "'little or no DD' (DDS <2.0), 'moderate DD' (DDS = 2.0–2.9), and 'high DD' (DDS ≥3.0)"
// Tags:    endocrine, diabetes, distress, dds17, category
// Units:   meanScore points 1–6; result category label
// Range:   1 ≤ meanScore ≤ 6
// Example: 1.94 → "little or no DD"; 2.0 → "moderate DD"; 3.0 → "high DD"
func DiabetesDistressCategoryFisher(meanScore float64) (DiabetesDistressCategory, error) {
	if err := checkFinite("meanScore", meanScore); err != nil {
		return "", err
	}
	if !(meanScore >= 1 && meanScore <= 6) {
		return "", fmt.Errorf("%w: meanScore = %v, want 1–6", ErrOutOfRange, meanScore)
	}
	switch {
	case meanScore >= 3.0:
		return DiabetesDistressHigh, nil
	case meanScore >= 2.0:
		return DiabetesDistressModerate, nil
	default:
		return DiabetesDistressLittleOrNone, nil
	}
}
