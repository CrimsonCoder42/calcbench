package scales

import (
	"errors"
	"fmt"
)

var (
	// ErrOutOfRange is returned when an item or a total is outside the range the source allows.
	ErrOutOfRange = errors.New("scales: input outside the range the source validated")
	// ErrNotTestable is returned when a Glasgow Coma Scale component was rated Non testable (NT): the
	// components are then the result, and no total exists.
	ErrNotTestable = errors.New("scales: a component is Non testable (NT), so no total exists")
)

// EpworthSleepinessScore returns the Epworth Sleepiness Scale total, 0–24: the sum of the 8 item ratings,
// each rated 0 to 3 on the response key of Johns 1991, Table 1. Licence: the ESS is © MW Johns 1990–1997 and
// "A license is required to use the ESS or the ESS-CHAD, whether or not a license fee is payable"
// (epworthsleepinessscale.com/licenses, read 2026-11-12). Only the scoring arithmetic is here; the item text is
// never embedded and the licensed caller supplies it (D-053). items[i] is the rating of item i+1 in the
// source's order.
// Errors: ErrOutOfRange when any item is outside 0–3, naming the item index.
//
// ID:      MED-SCALE-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johns MW. A new method for measuring daytime sleepiness: the Epworth sleepiness scale. Sleep. 1991;14(6):540-545. doi:10.1093/sleep/14.6.540. PMID 1798888. Table 1
// Tags:    scales, epworth, ess, sleep, sleepiness, score
// Units:   items ratings, points 0–3 each; result points 0–24
// Range:   each item 0–3 (population: 180 adults: 150 with sleep disorders, aged 18–78, and 30 controls)
// Example: all 8 items rated 1 → 8; all rated 3 → 24; all 0 → 0 (HAND-DERIVED; Johns 1991 controls scored 2–10)
func EpworthSleepinessScore(items [8]int) (int, error) {
	total := 0
	for i, v := range items {
		if v < 0 || v > 3 {
			return 0, fmt.Errorf("%w: items[%d] = %d, want 0–3", ErrOutOfRange, i, v)
		}
		total += v
	}
	return total, nil
}
