package psychiatry

import "fmt"

// GAD7Score returns the GAD-7 anxiety severity score, 0–21: the sum of its seven items, each scored 0, 1, 2 or 3
// for the four response options.
// Errors: ErrOutOfRange when any item is outside 0–3, naming the item index.
//
// ID:      MED-PSYCH-003
// Time:    O(1)   7 items, fixed
// Space:   O(1), no allocation
// Source:  Spitzer RL, Kroenke K, Williams JBW, Löwe B. A brief measure for assessing generalized anxiety disorder: the GAD-7. Arch Intern Med. 2006;166(10):1092-1097. doi:10.1001/archinte.166.10.1092, PMID 16717171, abstract ("A 7-item anxiety scale (GAD-7)")
// Source:  Instructions for Patient Health Questionnaire (PHQ) and GAD-7 Measures (developers' manual, phqscreeners.com, instructions.pdf, read 2026-11-12), p. 3 table ("Seven items, each of which is scored 0 to 3, providing a 0 to 21 severity score") and p. 6 ("assigning scores of 0, 1, 2, and 3")
// Tags:    psychiatry, anxiety, gad-7, questionnaire, score
// Units:   each item points 0–3; result points 0–21
// Range:   every item 0–3
// Example: items 3, 2, 3, 1, 0, 1, 2 → 12 (the sum)
func GAD7Score(items [7]int) (int, error) {
	total := 0
	for i, v := range items {
		if v < 0 || v > 3 {
			return 0, fmt.Errorf("%w: items[%d] = %d, want 0–3", ErrOutOfRange, i, v)
		}
		total += v
	}
	return total, nil
}
