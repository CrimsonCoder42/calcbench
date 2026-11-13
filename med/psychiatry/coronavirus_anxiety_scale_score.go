package psychiatry

import "fmt"

// CoronavirusAnxietyScaleScore returns the Coronavirus Anxiety Scale total, 0–20: the sum of its five items, each
// rated "from 0 (not at all) to 4 (nearly every day over the last 2 weeks)". The cited cut-off, a total of 9 or
// more, is the caller's. Items not embedded.
// Errors: ErrOutOfRange when an item is outside 0–4, naming its index.
//
// ID:      MED-PSYCH-010
// Time:    O(1)   5 items, fixed
// Space:   O(1), no allocation
// Source:  Lee SA. Coronavirus Anxiety Scale: a brief mental health screener for COVID-19 related anxiety. Death Stud. 2020;44(7):393-401. doi:10.1080/07481187.2020.1748481, PMID 32299304 (full text not open)
// Source:  Wu J, Qiu L, Xiong W, et al. BMJ Open. 2023, PMC10680000 ("recorded as 0, 1, 2, 3 and 4 points … cumulative score ranged from 0 to 20 points"; ref 27 Lee 2020)
// Source:  Zhang D, et al. BMC Psychol. 2023;11:224. doi:10.1186/s40359-023-01251-x, PMC10405510 ("a 5-item, single-factor tool"; "The total scores of the scale range from 0 to 20")
// Verified: secondary (Zhang 2023 PMC10405510, ref 17 Lee 2020; Wu 2023 PMC10680000, ref 27 Lee 2020)
// Tags:    psychiatry, anxiety, covid-19, cas, questionnaire, score
// Units:   item points 0–4; result points 0–20
// Range:   each item 0–4
// Example: 2, 2, 2, 2, 1 → 9 (the sum)
func CoronavirusAnxietyScaleScore(items [5]int) (int, error) {
	total := 0
	for i, v := range items {
		if v < 0 || v > 4 {
			return 0, fmt.Errorf("%w: items[%d] = %d, want 0–4", ErrOutOfRange, i, v)
		}
		total += v
	}
	return total, nil
}
