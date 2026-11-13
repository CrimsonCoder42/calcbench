package stats

import "math"

// StandardErrorOfMean returns the standard error of the mean of values, s/√n, with s the sample standard
// deviation (NUM-STATS-019). It is at most half the range (Popoviciu's bound on the variance), so it is
// finite for finite values, even when s is not. values is read, never modified.
// Errors: ErrEmpty when values is empty; ErrTooFew when it has one value; ErrNotFinite when a value is
// NaN or ±Inf.
//
// ID:      NUM-STATS-023
// Time:    O(n)   n = len(values); three passes (mean, then the centred sums)
// Space:   O(1) extra, no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.2 "Confidence Limits for the Mean", https://www.itl.nist.gov/div898/handbook/eda/section3/eda352.htm (read 2026-11-12)
// Source:  NIST, Dataplot Reference Manual, "STANDARD DEVIATION OF THE MEAN", https://www.itl.nist.gov/div898/software/dataplot/refman2/ch2/sd_mean.pdf (read 2026-11-12)
// Tags:    standard-error, sem, mean, precision, descriptive
// Units:   values in any one unit; result in that unit
// Example: NIST StRD Univariate Michelso (n = 100) → 0.0790105478190518/√100 = 0.00790105478190518 (from the certified s)
func StandardErrorOfMean(values []float64) (float64, error) {
	if err := checkSample(values, 2); err != nil {
		return 0, err
	}
	n := float64(len(values))
	sd, e := scaledStandardDeviation(values, n-1)
	r := math.Ldexp(sd/math.Sqrt(n), e)
	return min(r, math.MaxFloat64), nil // the bound holds exactly; this only undoes rounding
}
