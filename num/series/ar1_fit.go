package series

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
	"github.com/CrimsonCoder42/calcbench/num/stats"
)

// AR1Fit returns the intercept c and coefficient φ of the first-order autoregression x_t = c + φ·x_{t−1} + ε_t
// fitted to a series by conditional least squares: the ordinary least-squares regression of x_2…x_n on
// x_1…x_{n−1} with a constant (NUM-STATS-009 LinearRegression), which is the conditional maximum-likelihood
// estimate for Gaussian ε. It is not the Yule–Walker estimate φ = r_1 (NUM-SERIES-012), which divides by
// the variance of the whole series and shrinks φ toward 0 in short samples. The process mean is
// c/(1 − φ) when φ ≠ 1. values is in time order and is read, never modified.
// Errors: ErrTooFew when len(values) < 3; ErrNotFinite when a value is NaN or ±Inf; ErrZeroVariance when
// x_1…x_{n−1} are all equal; ErrOverflow when c or φ exceeds the float64 range.
//
// ID:      NUM-SERIES-018
// Time:    O(n)   n = len(values)
// Space:   O(1), no allocation
// Source:  Hamilton, Time Series Analysis, Princeton University Press, 1994, §5.2 "The Likelihood Function for a Gaussian AR(1) Process" (the conditional MLE of c and φ is the OLS regression of y_t on a constant and y_{t−1})
// Source:  Box, Jenkins, Reinsel & Ljung, Time Series Analysis: Forecasting and Control, 5th ed., Wiley, 2015, §7.1 (conditional least-squares estimates of the parameters)
// Tags:    time-series, autoregressive, ar1, least-squares, mean-reversion, ornstein-uhlenbeck, fit
// Units:   values in any one unit; intercept in that unit; phi dimensionless (per sampling step)
// Example: 1, 2, 4, 8 → regression of (2, 4, 8) on (1, 2, 4) → intercept 0, phi 2 (hand-checked)
func AR1Fit(values []float64) (intercept, phi float64, err error) {
	n := len(values)
	if n < 3 {
		return 0, 0, fmt.Errorf("%w: len(values) = %d, want ≥ 3", ErrTooFew, n)
	}
	for i, v := range values {
		if !reduce.AllFinite(v) {
			return 0, 0, fmt.Errorf("%w: values[%d] = %v", ErrNotFinite, i, v)
		}
	}
	lagged := values[:n-1]
	constant := true
	for _, v := range lagged[1:] {
		if v != lagged[0] {
			constant = false
			break
		}
	}
	if constant {
		return 0, 0, fmt.Errorf("%w: values[0…n−2] are all equal", ErrZeroVariance)
	}
	intercept, phi, err = stats.LinearRegression(lagged, values[1:])
	if err != nil { // inputs are checked above, so only an overflow (or a variance lost to it) can remain
		return 0, 0, fmt.Errorf("%w: AR(1) fit: %v", ErrOverflow, err)
	}
	return intercept, phi, nil
}

// MeanReversionHalfLife returns the half-life of a stationary AR(1) process with coefficient φ: the number of
// steps h after which the expected deviation from the mean has halved, φ^h = ½, so h = ln 2/(−ln φ). It
// follows from the AR(1) impulse response ∂x_{t+j}/∂ε_t = φ^j. Only 0 < φ < 1 has a half-life: φ ≥ 1 does
// not revert, and φ ≤ 0 oscillates in sign, so neither is accepted.
// Errors: ErrNotFinite when phi is NaN or ±Inf; ErrOutOfRange when phi ≤ 0 or phi ≥ 1.
//
// ID:      NUM-SERIES-019
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Hamilton, Time Series Analysis, Princeton University Press, 1994, §1.1, eq. [1.1.10] (dynamic multiplier ∂y_{t+j}/∂w_t = φ^j of a first-order difference equation); the half-life solves φ^h = ½
// Tags:    time-series, half-life, mean-reversion, ar1, ornstein-uhlenbeck, pairs-trading
// Units:   phi dimensionless (per sampling step); result in sampling steps (the period of the series AR1Fit was fitted to)
// Example: phi 0.5 → 1 step; phi 0.9 → ln 2/−ln 0.9 ≈ 6.578813478960585 steps (hand-checked)
func MeanReversionHalfLife(phi float64) (float64, error) {
	if !reduce.AllFinite(phi) {
		return 0, fmt.Errorf("%w: phi = %v", ErrNotFinite, phi)
	}
	if phi <= 0 || phi >= 1 {
		return 0, fmt.Errorf("%w: phi = %v, want 0 < phi < 1", ErrOutOfRange, phi)
	}
	return math.Ln2 / -math.Log(phi), nil
}
