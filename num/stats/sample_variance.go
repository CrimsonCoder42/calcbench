package stats

// SampleVariance returns the sample variance of values, s² = Σ(x − x̄)²/(n − 1), the unbiased estimate
// of a population's variance from a sample. The sum of squares is NUM-STATS-025's corrected two-pass form.
// values is read, never modified.
// A result below the smallest float64 rounds to a subnormal or 0 like any float64 result (the spread
// itself is computed on rescaled deviations and is never lost); for such data use the standard deviation.
// Errors: ErrEmpty when values is empty; ErrTooFew when it has one value; ErrNotFinite when a value is
// NaN or ±Inf; ErrOverflow when the variance exceeds the float64 range.
//
// ID:      NUM-STATS-021
// Time:    O(n)   n = len(values); three passes (mean, then the centred sums)
// Space:   O(1) extra, no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.6 "Measures of Scale", https://www.itl.nist.gov/div898/handbook/eda/section3/eda356.htm (read 2026-11-12)
// Source:  Chan, Golub, LeVeque, "Algorithms for computing the sample variance: analysis and recommendations", The American Statistician, 1983, 37(3):242–247, doi:10.1080/00031305.1983.10483115 (the corrected two-pass algorithm)
// Tags:    variance, sample-variance, dispersion, scale, descriptive, two-pass
// Units:   values in any one unit; result in that unit squared
// Example: NIST StRD Univariate NumAcc4 (n = 1001) → 0.01, the square of the certified s = 0.1
func SampleVariance(values []float64) (float64, error) {
	if err := checkSample(values, 2); err != nil {
		return 0, err
	}
	return varianceOf(values, float64(len(values)-1))
}
