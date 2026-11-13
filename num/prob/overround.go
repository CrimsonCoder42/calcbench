package prob

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// checkDecimalOdds validates a market of decimal odds: at least two outcomes, each finite and ≥ 1.
func checkDecimalOdds(decimalOdds []float64) error {
	if len(decimalOdds) < 2 {
		return fmt.Errorf("%w: len(decimalOdds) = %d, want >= 2", ErrOutOfRange, len(decimalOdds))
	}
	for i, d := range decimalOdds {
		if !reduce.AllFinite(d) {
			return fmt.Errorf("%w: decimalOdds[%d] = %v", ErrNotFinite, i, d)
		}
		if d < 1 {
			return fmt.Errorf("%w: decimalOdds[%d] = %v, want >= 1", ErrOutOfRange, i, d)
		}
	}
	return nil
}

// Overround returns the overround of one market, Σᵢ 1/dᵢ − 1 (the booksum minus 1), as a fraction:
// 0.05 is 5 %. It is negative for an arbitrage market. The caller passes every outcome of the market;
// completeness cannot be checked. The sum, with the −1, is compensated (Neumaier).
// Errors: ErrOutOfRange when there are fewer than 2 outcomes or a decimal odd is below 1; ErrNotFinite
// when one is NaN or ±Inf.
//
// ID:      NUM-PROB-016
// Time:    O(n)   n = len(decimalOdds)
// Space:   O(1), no allocation
// Source:  Wheatcroft, "Forecasting football matches by predicting match statistics", Journal of Sports Analytics, 2021, 7(2):77–97, doi:10.3233/JSA-200462; author preprint arXiv:2001.09097v1 (read 2026-11-12), §2.1 "Betting odds", eq. (1) (overround π = Σᵢ 1/Oᵢ − 1)
// Source:  Egidi, Pauli, Torelli, "Combining historical data and bookmakers' odds in modelling football scores", Statistical Modelling, 2018, 18(5–6):436–459, doi:10.1177/1471082X18798414; author preprint arXiv:1802.08848 (read 2026-11-12), §2 "Transforming the betting odds into probabilities", eq. (2.1) (booksum β = Σᵢ of the inverse odds)
// Tags:    probability, odds, overround, booksum, margin, vig, betting
// Units:   decimalOdds: units returned per unit staked, each ≥ 1; result: a dimensionless fraction
// Example: [1.9, 1.9] → 2/1.9 − 1 ≈ 0.0526
func Overround(decimalOdds []float64) (float64, error) {
	if err := checkDecimalOdds(decimalOdds); err != nil {
		return 0, err
	}
	var s compensatedSum
	for _, d := range decimalOdds {
		s.add(1 / d)
	}
	s.add(-1)
	return s.total(), nil
}
