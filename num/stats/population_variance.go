package stats

// PopulationVariance returns the variance of values taken as a whole population, σ² = Σ(x − μ)²/N with μ
// their mean: the variance E[(X − μ)²] of the distribution that puts probability 1/N on each value. For a
// sample's estimate of a larger population's variance use SampleVariance (divisor N − 1).
// values is read, never modified.
// A result below the smallest float64 rounds to a subnormal or 0 like any float64 result (the spread
// itself is computed on rescaled deviations and is never lost); for such data use the standard deviation.
// Errors: ErrEmpty when values is empty; ErrNotFinite when a value is NaN or ±Inf; ErrOverflow when the
// variance exceeds the float64 range.
//
// ID:      NUM-STATS-017
// Time:    O(n)   n = len(values); three passes (mean, then the centred sums)
// Space:   O(1) extra, no allocation
// Source:  Hogg, McKean, Craig, Introduction to Mathematical Statistics, 8th ed., Pearson, 2019, §1.9 "Some Special Expectations", Definition 1.9.2 (variance σ² = E[(X − μ)²])
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §2.3 "Moments and Moment Generating Functions" (variance)
// Source:  Chan, Golub, LeVeque, "Algorithms for computing the sample variance: analysis and recommendations", The American Statistician, 1983, 37(3):242–247, doi:10.1080/00031305.1983.10483115 (the corrected two-pass algorithm)
// Tags:    variance, population-variance, dispersion, scale, descriptive, two-pass
// Units:   values in any one unit; result in that unit squared
// Example: NumAcc1 (10000001, 10000003, 10000002) → 2/3 (hand-checked: deviations −1, 1, 0)
func PopulationVariance(values []float64) (float64, error) {
	if err := checkSample(values, 1); err != nil {
		return 0, err
	}
	return varianceOf(values, float64(len(values)))
}
