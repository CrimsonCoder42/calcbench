package prob

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// DecimalOddsFromProbability returns the fair decimal (European) odds d = 1/p: the total returned per
// unit staked on an event of probability p, with no bookmaker's margin.
// Errors: ErrNotFinite when p is NaN or ±Inf; ErrProbability when p is outside [0, 1]; ErrOutOfRange
// when p = 0; ErrOverflow when 1/p exceeds the float64 range (p subnormal).
//
// ID:      NUM-PROB-014
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Wheatcroft, "Forecasting football matches by predicting match statistics", Journal of Sports Analytics, 2021, 7(2):77–97, doi:10.3233/JSA-200462; author preprint arXiv:2001.09097v1 (read 2026-11-12), §2.1 "Betting odds" (odds-implied probability rᵢ = 1/Oᵢ, so fair odds are 1/p)
// Source:  Wunderlich, Memmert, "The Betting Odds Rating System: Using soccer forecasts to forecast soccer", PLOS ONE, 2018, 13(6):e0198668, doi:10.1371/journal.pone.0198668 (read 2026-11-12), Method, "Transferring betting odds to probabilities" (with no bookmaker margin the inverse odds are the probability)
// Tags:    probability, odds, decimal-odds, betting, conversion
// Units:   p: a probability in (0, 1]; result: decimal odds ≥ 1 (units returned per unit staked)
// Example: 0.4 → 2.5
func DecimalOddsFromProbability(p float64) (float64, error) {
	if err := checkProbability("p", p); err != nil {
		return 0, err
	}
	if p == 0 {
		return 0, fmt.Errorf("%w: p = 0, odds infinite", ErrOutOfRange)
	}
	d := 1 / p
	if !reduce.AllFinite(d) {
		return 0, fmt.Errorf("%w: 1/p for p = %v", ErrOverflow, p)
	}
	return d, nil
}
