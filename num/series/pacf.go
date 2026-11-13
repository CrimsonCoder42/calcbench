package series

import (
	"fmt"
)

// PACF returns the sample partial autocorrelations φ_00, φ_11, φ_22, …, φ_mm (m = maxLag) of values in
// dst[0:maxLag+1], dst grown when its capacity is short; dst[0] = 1 by convention. φ_kk is the last
// coefficient of the order-k autoregression fitted to the sample autocorrelations r_1…r_k (NUM-SERIES-013),
// obtained by the Durbin–Levinson recursion:
//
//	φ_kk = (r_k − Σ_{j=1}^{k−1} φ_{k−1,j}·r_{k−j}) / (1 − Σ_{j=1}^{k−1} φ_{k−1,j}·r_j)
//	φ_kj = φ_{k−1,j} − φ_kk·φ_{k−1,k−j},  j = 1…k−1
//
// so no k×k Yule–Walker system is solved. The denominator is the prediction-error variance ratio
// v_{k−1}/c_0 = Π_{j<k}(1 − φ_jj²); it is positive for a non-constant series because the divisor-N sample
// autocorrelations are positive definite. values is read, never modified.
// Errors: ErrTooFew when len(values) < 2; ErrNotFinite when a value is NaN or ±Inf; ErrOutOfRange when
// maxLag < 0 or maxLag ≥ len(values); ErrZeroVariance when every value is equal, or when the prediction
// error variance reaches zero at some lag (the series is exactly predictable from its past by rounding).
//
// ID:      NUM-SERIES-014
// Time:    O(n·m + m²)   n = len(values), m = maxLag; the autocorrelations, then the recursion
// Space:   O(m): the result and one scratch slice of 2m + 1 values (one allocation)
// Source:  Durbin, "The Fitting of Time-Series Models", Revue de l'Institut International de Statistique, 1960, 28(3):233–244, doi:10.2307/1401322, §2 (the recursion for the autoregressive coefficients)
// Source:  Box, Jenkins, Reinsel & Ljung, Time Series Analysis: Forecasting and Control, 5th ed., Wiley, 2015, §3.2.6 "Estimation of the partial autocorrelation function" and Appendix A3.2 "Recursive method for calculating estimates of autoregressive parameters"
// Tags:    time-series, partial-autocorrelation, pacf, durbin-levinson, yule-walker, autoregressive, box-jenkins
// Units:   values in any one unit; maxLag in samples; results dimensionless in [−1, 1]
// Example: 1, 2, 3, 4, maxLag 2 → r_1 = 0.25, r_2 = −0.3 → φ_11 = 0.25, φ_22 = (−0.3 − 0.0625)/(1 − 0.0625) = −0.38666… (hand-checked)
func PACF(dst, values []float64, maxLag int) ([]float64, error) {
	if maxLag < 0 || maxLag >= len(values) {
		if len(values) < 2 {
			return dst, fmt.Errorf("%w: len(values) = %d, want ≥ 2", ErrTooFew, len(values))
		}
		return dst, fmt.Errorf("%w: maxLag = %d, want 0 ≤ maxLag < %d", ErrOutOfRange, maxLag, len(values))
	}
	scratch := make([]float64, 2*maxLag+1)
	r, err := ACF(scratch[:0], values, maxLag) // r[0..m], fits in scratch, no allocation
	if err != nil {
		return dst, err
	}
	prev := scratch[maxLag+1:] // φ_{k−1, 1…k−1} in prev[0…k−2]
	dst = resizeFloats(dst, maxLag+1)
	dst[0] = 1
	for k := 1; k <= maxLag; k++ {
		num, den := r[k], 1.0
		for j := 1; j < k; j++ {
			num -= prev[j-1] * r[k-j]
			den -= prev[j-1] * r[j]
		}
		if !(den > 0) {
			return dst, fmt.Errorf("%w: prediction error variance is zero at lag %d", ErrZeroVariance, k)
		}
		phiKK := num / den
		// φ_kj = φ_{k−1,j} − φ_kk·φ_{k−1,k−j}, in place, pairing j with k − j.
		for lo, hi := 0, k-2; lo <= hi; lo, hi = lo+1, hi-1 {
			a, b := prev[lo], prev[hi]
			prev[lo] = a - phiKK*b
			if lo != hi {
				prev[hi] = b - phiKK*a
			}
		}
		prev[k-1] = phiKK
		dst[k] = phiKK
	}
	return dst, nil
}
