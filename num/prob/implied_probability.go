package prob

import "fmt"

// ImpliedProbability returns π = 1/d, the probability implied by decimal odds d. A bookmaker's quoted
// odds include the overround, so the implied probabilities of a market sum to more than 1;
// NormalizedImpliedProbabilities removes it.
// Errors: ErrNotFinite when decimalOdds is NaN or ±Inf; ErrOutOfRange when decimalOdds < 1.
//
// ID:      NUM-PROB-015
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Wheatcroft, "Forecasting football matches by predicting match statistics", Journal of Sports Analytics, 2021, 7(2):77–97, doi:10.3233/JSA-200462; author preprint arXiv:2001.09097v1 (read 2026-11-12), §2.1 "Betting odds" (odds-implied probability rᵢ = 1/Oᵢ)
// Source:  Egidi, Pauli, Torelli, "Combining historical data and bookmakers' odds in modelling football scores", Statistical Modelling, 2018, 18(5–6):436–459, doi:10.1177/1471082X18798414; author preprint arXiv:1802.08848 (read 2026-11-12), §2 "Transforming the betting odds into probabilities" (the inverse odd is the unfair probability of the event)
// Tags:    probability, odds, decimal-odds, implied-probability, betting, conversion
// Units:   decimalOdds: units returned per unit staked, ≥ 1; result: a probability in (0, 1]
// Example: 2.5 → 0.4; 1.9 → 0.5263157894736842
func ImpliedProbability(decimalOdds float64) (float64, error) {
	if err := checkFinite("decimalOdds", decimalOdds); err != nil {
		return 0, err
	}
	if decimalOdds < 1 {
		return 0, fmt.Errorf("%w: decimalOdds = %v, want >= 1", ErrOutOfRange, decimalOdds)
	}
	return 1 / decimalOdds, nil
}
