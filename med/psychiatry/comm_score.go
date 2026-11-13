package psychiatry

import (
	"errors"
	"fmt"
)

// ErrOutOfRange is returned when an item rating or a total is outside the range the source allows.
var ErrOutOfRange = errors.New("psychiatry: input outside the range the source validated")

// COMMScore returns the Current Opioid Misuse Measure total, 0–68: the sum of its 17 items, each rated on the
// past 30 days "from '0 = Never' to '4 = Very Often.'". The source reports "A cutoff score of 9 or higher";
// applying it is the caller's.
// Errors: ErrOutOfRange when any item is outside 0–4, naming the item index.
//
// ID:      MED-PSYCH-001
// Time:    O(1)   17 items, fixed
// Space:   O(1), no allocation
// Source:  Butler SF, Budman SH, Fernandez KC, et al. Development and validation of the Current Opioid Misuse Measure. Pain. 2007;130(1-2):144-156. doi:10.1016/j.pain.2007.01.014, PMID 17493754, abstract ("17 items"), Discussion (0–4 scale) and Results (cut-off)
// Tags:    psychiatry, substance use, opioid, comm, questionnaire, score
// Units:   each item points 0–4; result points 0–68
// Range:   every item 0–4
// Example: nine items rated 1, the rest 0 → 9; all 4 → 68 (the sum)
func COMMScore(items [17]int) (int, error) {
	total := 0
	for i, v := range items {
		if v < 0 || v > 4 {
			return 0, fmt.Errorf("%w: items[%d] = %d, want 0–4", ErrOutOfRange, i, v)
		}
		total += v
	}
	return total, nil
}
