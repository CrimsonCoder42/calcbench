package series

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// ErrZeroVariance: every value is equal, so a statistic that divides by the sample variance is undefined.
var ErrZeroVariance = errors.New("series: zero variance")

// acfScaleHigh and acfScaleLow bound max|x| outside which the values are rescaled by an exact power of two
// before the products are formed, so no product or sum can overflow or lose its bits to underflow. They are
// numerical guards, not model constants: r_k is invariant to the scale of the series.
const (
	acfScaleHigh = 0x1p500
	acfScaleLow  = 0x1p-500
)

// acfPrepare checks values and returns the power-of-two factor f (exact scaling) and the mean of f·x.
func acfPrepare(values []float64) (f, mean float64, err error) {
	n := len(values)
	if n < 2 {
		return 0, 0, fmt.Errorf("%w: len(values) = %d, want ≥ 2", ErrTooFew, n)
	}
	maxAbs := 0.0
	for i, v := range values {
		if !reduce.AllFinite(v) {
			return 0, 0, fmt.Errorf("%w: values[%d] = %v", ErrNotFinite, i, v)
		}
		maxAbs = math.Max(maxAbs, math.Abs(v))
	}
	f = 1
	if maxAbs >= acfScaleHigh || (maxAbs > 0 && maxAbs < acfScaleLow) {
		_, e := math.Frexp(maxAbs)
		// A subnormal maxAbs gives −e up to 1074, and 2^1074 overflows to +Inf; cap the exponent at 1023
		// (2^1023 is the largest finite power of two). Then f·maxAbs ∈ [0.5, 1) for normal maxAbs and
		// ∈ [2⁻⁵¹, 0.5) for subnormal maxAbs. Either way f·v is a normal (or zero) value, and multiplying a
		// finite v by a power of two whose product is normal is exact, so the scaling is exact for every
		// finite input, subnormals included.
		f = math.Ldexp(1, min(-e, 1023))
	}
	var sum, comp float64
	for _, v := range values {
		sum, comp = reduce.NeumaierAdd(sum, comp, f*v)
	}
	return f, (sum + comp) / float64(n), nil
}

// autocovarianceScaled returns n·c_k of the scaled series: Σ_{t=0}^{n−1−k} (f·x_t − mean)(f·x_{t+k} − mean),
// compensated.
func autocovarianceScaled(values []float64, f, mean float64, k int) float64 {
	var sum, comp float64
	for t := 0; t+k < len(values); t++ {
		sum, comp = reduce.NeumaierAdd(sum, comp, (f*values[t]-mean)*(f*values[t+k]-mean))
	}
	return sum + comp
}

// Autocorrelation returns the sample autocorrelation of values at one lag, r_k = c_k/c_0, where
// c_k = (1/N)·Σ_{t=1}^{N−k} (z_t − z̄)(z_{t+k} − z̄) and z̄ is the sample mean. The divisor is N at every
// lag (not N − k), the estimator Box and Jenkins recommend: it keeps the sequence r_0, r_1, … positive
// semidefinite. r_0 = 1. Values are in time order and are read, never modified. Sums are compensated;
// values beyond 2⁵⁰⁰ or below 2⁻⁵⁰⁰ in magnitude are rescaled by an exact power of two first.
// Errors: ErrTooFew when len(values) < 2; ErrNotFinite when a value is NaN or ±Inf; ErrOutOfRange when
// lag < 0 or lag ≥ len(values); ErrZeroVariance when every value is equal.
//
// ID:      NUM-SERIES-012
// Time:    O(n)   n = len(values); four passes (check, mean, c_0, c_k)
// Space:   O(1), no allocation
// Source:  Box, Jenkins, Reinsel & Ljung, Time Series Analysis: Forecasting and Control, 5th ed., Wiley, 2015, §2.1.5 "Estimation of autocovariance and autocorrelation functions" (c_k with divisor N, r_k = c_k/c_0)
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.12 "Autocorrelation", https://www.itl.nist.gov/div898/handbook/eda/section3/eda35c.htm (C_h = (1/N)Σ(Y_t − Ȳ)(Y_{t+h} − Ȳ), R_h = C_h/C_0)
// Tags:    time-series, autocorrelation, acf, serial-correlation, lag, box-jenkins
// Units:   values in any one unit; lag in samples; result dimensionless in [−1, 1]
// Example: 1, 2, 3, 4, lag 1 → z̄ = 2.5, c_0 = 5/4, c_1 = 1.25/4 → r_1 = 0.25 (hand-checked)
func Autocorrelation(values []float64, lag int) (float64, error) {
	f, mean, err := acfPrepare(values)
	if err != nil {
		return 0, err
	}
	if lag < 0 || lag >= len(values) {
		return 0, fmt.Errorf("%w: lag = %d, want 0 ≤ lag < %d", ErrOutOfRange, lag, len(values))
	}
	c0 := autocovarianceScaled(values, f, mean, 0)
	if c0 == 0 {
		return 0, fmt.Errorf("%w: every value is equal", ErrZeroVariance)
	}
	return autocovarianceScaled(values, f, mean, lag) / c0, nil
}

// ACF returns the sample autocorrelations r_0, r_1, …, r_maxLag of values (NUM-SERIES-012 at every lag)
// in dst[0:maxLag+1], dst grown when its capacity is short; dst[0] = 1. It reuses dst's storage, so a
// caller that passes a long enough dst allocates nothing. The lags are direct sums, O(n·m): an FFT route
// would be O(n log n) for every lag, but num/signal is not importable from num/series (D-009) and m ≪ n in
// practice.
// Errors: ErrTooFew when len(values) < 2; ErrNotFinite when a value is NaN or ±Inf; ErrOutOfRange when
// maxLag < 0 or maxLag ≥ len(values); ErrZeroVariance when every value is equal.
//
// ID:      NUM-SERIES-013
// Time:    O(n·m)   n = len(values), m = maxLag + 1; direct sums
// Space:   O(m) for the result, no other allocation
// Source:  Box, Jenkins, Reinsel & Ljung, Time Series Analysis: Forecasting and Control, 5th ed., Wiley, 2015, §2.1.5 (r_k = c_k/c_0, k = 0, 1, …, K)
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.12 "Autocorrelation", https://www.itl.nist.gov/div898/handbook/eda/section3/eda35c.htm
// Tags:    time-series, autocorrelation, acf, correlogram, serial-correlation, box-jenkins
// Units:   values in any one unit; maxLag in samples; results dimensionless in [−1, 1]
// Example: 1, 2, 3, 4, maxLag 3 → 1, 0.25, −0.3, −0.45 (hand-checked)
func ACF(dst, values []float64, maxLag int) ([]float64, error) {
	f, mean, err := acfPrepare(values)
	if err != nil {
		return dst, err
	}
	if maxLag < 0 || maxLag >= len(values) {
		return dst, fmt.Errorf("%w: maxLag = %d, want 0 ≤ maxLag < %d", ErrOutOfRange, maxLag, len(values))
	}
	c0 := autocovarianceScaled(values, f, mean, 0)
	if c0 == 0 {
		return dst, fmt.Errorf("%w: every value is equal", ErrZeroVariance)
	}
	dst = resizeFloats(dst, maxLag+1)
	dst[0] = 1
	for k := 1; k <= maxLag; k++ {
		dst[k] = autocovarianceScaled(values, f, mean, k) / c0
	}
	return dst, nil
}

// resizeFloats returns dst with length n, reusing its storage when the capacity allows.
func resizeFloats(dst []float64, n int) []float64 {
	if cap(dst) >= n {
		return dst[:n]
	}
	return make([]float64, n)
}
