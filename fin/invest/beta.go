package invest

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/num/stats"
)

// Beta returns the asset's beta against the market, βᵢ = Cov(Rᵢ, R_M)/Var(R_M), from paired per-period
// returns: the sample covariance (NUM-STATS-038) over the sample variance of the market (NUM-STATS-021);
// the n − 1 divisors cancel. Equal to the slope of the least-squares regression of Rᵢ on R_M. The inputs
// are read, never modified.
// Errors: ErrLengthMismatch when the lengths differ; ErrTooFew when there are fewer than 2 pairs;
// ErrNotFinite for NaN or ±Inf; ErrZeroDenominator when the market returns are all equal; ErrOverflow
// when an intermediate or the result exceeds float64.
//
// ID:      FIN-INVEST-043
// Time:    O(n)   n = number of periods; two three-pass reductions
// Space:   O(1), no allocation
// Source:  Bodie, Kane, Marcus, Investments, 10th ed., McGraw-Hill, 2014, ch. 9 "The Capital Asset Pricing Model", §9.1 (βᵢ = Cov(rᵢ, r_M)/σ²_M)
// Source:  Sharpe, "Capital Asset Prices: A Theory of Market Equilibrium under Conditions of Risk", Journal of Finance, 1964, 19(3):425–442, doi:10.1111/j.1540-6261.1964.tb02865.x
// Tags:    invest, risk, beta, capm, systematic-risk, covariance, regression
// Units:   returns decimal per period; result dimensionless
// Example: hand-checked: market 0.01, 0.02, −0.01, 0.03 and asset = 2 × market + 0.001 → 2
func Beta(assetReturns, marketReturns []float64) (float64, error) {
	if err := checkPairedSeries("assetReturns", assetReturns, "marketReturns", marketReturns, 2); err != nil {
		return 0, err
	}
	v, err := stats.SampleVariance(marketReturns)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrOverflow, err)
	}
	if v == 0 {
		return 0, fmt.Errorf("%w: marketReturns have zero variance", ErrZeroDenominator)
	}
	c, err := stats.SampleCovariance(assetReturns, marketReturns)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrOverflow, err)
	}
	return finiteResult(c / v)
}
