package stats

// PopulationStandardDeviation returns the standard deviation of values taken as a whole population,
// σ = √(Σ(x − μ)²/N), the square root of NUM-STATS-017. It is finite whenever σ is, even when σ² would
// overflow or underflow. values is read, never modified.
// Errors: ErrEmpty when values is empty; ErrNotFinite when a value is NaN or ±Inf; ErrOverflow when σ
// exceeds the float64 range.
//
// ID:      NUM-STATS-022
// Time:    O(n)   n = len(values); three passes (mean, then the centred sums)
// Space:   O(1) extra, no allocation
// Source:  Hogg, McKean, Craig, Introduction to Mathematical Statistics, 8th ed., Pearson, 2019, §1.9 "Some Special Expectations", Definition 1.9.2 and p. 69 (standard deviation σ, the positive square root of the variance)
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §2.3 "Moments and Moment Generating Functions" (variance and standard deviation)
// Source:  Chan, Golub, LeVeque, "Algorithms for computing the sample variance: analysis and recommendations", The American Statistician, 1983, 37(3):242–247, doi:10.1080/00031305.1983.10483115 (the corrected two-pass algorithm)
// Tags:    standard-deviation, population-standard-deviation, dispersion, scale, descriptive
// Units:   values in any one unit; result in that unit
// Example: 2, 4, 4, 4, 5, 5, 7, 9 → 2 (hand-checked: mean 5, squared deviations sum to 32, 32/8 = 4)
func PopulationStandardDeviation(values []float64) (float64, error) {
	if err := checkSample(values, 1); err != nil {
		return 0, err
	}
	return standardDeviationOf(values, float64(len(values)))
}
