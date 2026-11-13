package stats

// SampleStandardDeviation returns the sample standard deviation of values, s = √(Σ(x − x̄)²/(n − 1)),
// the square root of NUM-STATS-021. It is accurate whenever s is in the float64 range, even when s²
// would overflow or underflow. values is read, never modified.
// Errors: ErrEmpty when values is empty; ErrTooFew when it has one value; ErrNotFinite when a value is
// NaN or ±Inf; ErrOverflow when s exceeds the float64 range.
//
// ID:      NUM-STATS-019
// Time:    O(n)   n = len(values); three passes (mean, then the centred sums)
// Space:   O(1) extra, no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.6 "Measures of Scale", https://www.itl.nist.gov/div898/handbook/eda/section3/eda356.htm (read 2026-11-12)
// Source:  Chan, Golub, LeVeque, "Algorithms for computing the sample variance: analysis and recommendations", The American Statistician, 1983, 37(3):242–247, doi:10.1080/00031305.1983.10483115 (the corrected two-pass algorithm)
// Tags:    standard-deviation, sample-standard-deviation, dispersion, scale, descriptive
// Units:   values in any one unit; result in that unit
// Example: NIST StRD Univariate Michelso (n = 100) → 0.0790105478190518, the certified sample standard deviation
func SampleStandardDeviation(values []float64) (float64, error) {
	if err := checkSample(values, 2); err != nil {
		return 0, err
	}
	return standardDeviationOf(values, float64(len(values)-1))
}
