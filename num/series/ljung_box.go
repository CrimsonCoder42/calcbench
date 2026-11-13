package series

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
	"github.com/CrimsonCoder42/calcbench/num/prob"
)

// LjungBox returns the Ljung–Box portmanteau statistic of a series (usually a model's residuals) over lags
// 1…maxLag, Q = n(n + 2)·Σ_{k=1}^{h} r_k²/(n − k) with n = len(values), h = maxLag and r_k the sample
// autocorrelation (NUM-SERIES-012, centred on the sample mean, divisor n), its degrees of freedom
// df = h − fittedParams, and the p-value P(χ²_df > Q) (NUM-PROB-064 ChiSquareSurvival). Under the null
// hypothesis that the series is white noise Q is approximately χ²_df. fittedParams is the number of ARMA
// parameters p + q estimated to produce the residuals (0 for a raw series). values is read, never modified.
// Errors: ErrTooFew when len(values) < 2; ErrNotFinite when a value is NaN or ±Inf; ErrOutOfRange when
// maxLag < 1, maxLag ≥ len(values), fittedParams < 0, fittedParams ≥ maxLag, or df > 10⁶ (the p-value's range); ErrZeroVariance when every
// value is equal; ErrOverflow when Q exceeds the float64 range (n beyond about 10¹⁵⁰).
//
// ID:      NUM-SERIES-015
// Time:    O(n·h)   n = len(values), h = maxLag
// Space:   O(1), no allocation
// Source:  Ljung & Box, "On a Measure of Lack of Fit in Time Series Models", Biometrika, 1978, 65(2):297–303, doi:10.1093/biomet/65.2.297, §2 (the modified statistic Q̃(r̂) = n(n + 2)Σ(n − k)⁻¹r̂_k², referred to χ² with m − p − q degrees of freedom)
// Source:  Box, Jenkins, Reinsel & Ljung, Time Series Analysis: Forecasting and Control, 5th ed., Wiley, 2015, §8.2.2 "Portmanteau lack-of-fit test"
// Tags:    time-series, ljung-box, portmanteau, white-noise, autocorrelation, residual-diagnostics, hypothesis-test, p-value
// Units:   values in any one unit; maxLag and fittedParams counts; q dimensionless; pValue a probability in [0, 1]
// Example: 1, 2, 3, 4, maxLag 1, fittedParams 0 → r_1 = 0.25, Q = 4·6·0.0625/3 = 0.5, df = 1 (hand-checked)
func LjungBox(values []float64, maxLag, fittedParams int) (q float64, df int, pValue float64, err error) {
	f, mean, err := acfPrepare(values)
	if err != nil {
		return 0, 0, 0, err
	}
	n := len(values)
	if maxLag < 1 || maxLag >= n {
		return 0, 0, 0, fmt.Errorf("%w: maxLag = %d, want 1 ≤ maxLag < %d", ErrOutOfRange, maxLag, n)
	}
	if fittedParams < 0 || fittedParams >= maxLag {
		return 0, 0, 0, fmt.Errorf("%w: fittedParams = %d, want 0 ≤ fittedParams < maxLag = %d", ErrOutOfRange, fittedParams, maxLag)
	}
	c0 := autocovarianceScaled(values, f, mean, 0)
	if c0 == 0 {
		return 0, 0, 0, fmt.Errorf("%w: every value is equal", ErrZeroVariance)
	}
	var sum, comp float64
	for k := 1; k <= maxLag; k++ {
		r := autocovarianceScaled(values, f, mean, k) / c0
		sum, comp = reduce.NeumaierAdd(sum, comp, r*r/float64(n-k))
	}
	nf := float64(n)
	q = nf * (nf + 2) * (sum + comp)
	if !reduce.AllFinite(q) {
		return 0, 0, 0, fmt.Errorf("%w: Q", ErrOverflow)
	}
	df = maxLag - fittedParams
	pValue, err = prob.ChiSquareSurvival(q, float64(df))
	if err != nil {
		return 0, 0, 0, fmt.Errorf("%w: chi-square p-value: %v", ErrOutOfRange, err)
	}
	return q, df, pValue, nil
}
