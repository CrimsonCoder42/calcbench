package prob

import "fmt"

// OddsFromProbability returns the odds in favour of an event, p / (1 − p): "o to 1". The odds against
// are 1/o.
// Errors: ErrNotFinite when p is NaN or ±Inf; ErrProbability when p is outside [0, 1]; ErrOutOfRange
// when p = 1 (infinite odds).
//
// ID:      NUM-PROB-012
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Bland, Altman, "Statistics Notes: The odds ratio", BMJ, 2000, 320(7247):1468, doi:10.1136/bmj.320.7247.1468 (odds = P(event)/P(no event))
// Tags:    probability, odds, conversion, betting
// Units:   p: a probability in [0, 1); result: odds in favour, a ratio ≥ 0
// Example: 0.75 → 3 (3 to 1)
func OddsFromProbability(p float64) (float64, error) {
	if err := checkProbability("p", p); err != nil {
		return 0, err
	}
	if p == 1 {
		return 0, fmt.Errorf("%w: p = 1, odds infinite", ErrOutOfRange)
	}
	return p / (1 - p), nil
}
