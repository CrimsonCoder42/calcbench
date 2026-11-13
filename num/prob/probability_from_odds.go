package prob

import "fmt"

// ProbabilityFromOdds returns p = o / (1 + o), the probability of an event whose odds in favour are o
// ("o to 1"). Above 2⁵³ the result rounds to 1, which is the correctly rounded value.
// Errors: ErrNotFinite when odds is NaN or ±Inf; ErrOutOfRange when odds < 0.
//
// ID:      NUM-PROB-013
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Bland, Altman, "Statistics Notes: The odds ratio", BMJ, 2000, 320(7247):1468, doi:10.1136/bmj.320.7247.1468 (odds = P(event)/P(no event))
// Tags:    probability, odds, conversion, betting
// Units:   odds: a ratio ≥ 0; result: a probability in [0, 1]
// Example: 3 → 0.75
func ProbabilityFromOdds(odds float64) (float64, error) {
	if err := checkFinite("odds", odds); err != nil {
		return 0, err
	}
	if odds < 0 {
		return 0, fmt.Errorf("%w: odds = %v, want >= 0", ErrOutOfRange, odds)
	}
	// 1 + o cannot overflow: at o = MaxFloat64 it rounds back to o.
	return odds / (1 + odds), nil
}
