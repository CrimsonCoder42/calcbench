package invest

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
	"github.com/CrimsonCoder42/calcbench/num/stats"
)

// checkConfidence returns ErrDomain unless confidence is finite and in (0, 1).
func checkConfidence(confidence float64) error {
	if !reduce.AllFinite(confidence) || confidence <= 0 || confidence >= 1 {
		return fmt.Errorf("%w: confidence = %v, must be in (0, 1)", ErrDomain, confidence)
	}
	return nil
}

// HistoricalVaR returns the historical-simulation value at risk at the given confidence, as a positive
// loss fraction: −Q(1 − confidence), with Q the sample quantile of the observed returns by Hyndman–Fan
// definition method (NUM-STATS-015 Percentile; HyndmanFan1, the inverse of the empirical distribution,
// matches ExpectedShortfall's tail). The confidence level (for example 0.95 or 0.99) and the quantile
// method are the caller's. A negative result means even the tail quantile is a gain. returns is read,
// never modified.
// Errors: ErrTooFew when returns is empty; ErrNotFinite for NaN or ±Inf; ErrDomain when confidence is not
// in (0, 1) or method is not HyndmanFan1 … HyndmanFan9.
//
// ID:      FIN-INVEST-049
// Time:    O(n)   n = len(returns); one linear-time selection
// Space:   O(n) extra (one allocation, the selection's copy)
// Source:  Jorion, Value at Risk: The New Benchmark for Managing Financial Risk, 3rd ed., McGraw-Hill, 2007, ch. 5 "Computing VAR", §5.2 (VAR from the empirical distribution: the loss at the (1 − c) quantile)
// Tags:    invest, risk, value-at-risk, var, historical-simulation, quantile, tail-risk
// Units:   returns decimal per period; confidence a probability in (0, 1); result decimal loss per period (0.03 = a 3 % loss)
// Example: hand-checked: returns −0.05, −0.03, −0.01, 0, 0.01, 0.02, 0.02, 0.03, 0.04, 0.05, confidence 0.8, HyndmanFan1 → −X₍₂₎ = 0.03; HyndmanFan7 → 0.014
func HistoricalVaR(returns []float64, confidence float64, method stats.QuantileMethod) (float64, error) {
	if err := checkSeries("returns", returns, 1); err != nil {
		return 0, err
	}
	if err := checkConfidence(confidence); err != nil {
		return 0, err
	}
	if method < stats.HyndmanFan1 || method > stats.HyndmanFan9 {
		return 0, fmt.Errorf("%w: method = %d, must be HyndmanFan1 … HyndmanFan9", ErrDomain, method)
	}
	q, err := stats.Percentile(returns, 1-confidence, method)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrDomain, err)
	}
	return -q, nil
}
