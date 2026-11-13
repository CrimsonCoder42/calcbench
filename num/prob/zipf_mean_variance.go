package prob

import "fmt"

// ZipfMeanVariance returns the mean ζ(a − 1)/ζ(a) (a > 2) and variance ζ(a − 2)/ζ(a) − mean² (a > 3) of the
// Zipf distribution; the variance is computed as (ζ(a − 2)ζ(a) − ζ(a − 1)²)/ζ(a)².
// Errors: ErrNotFinite when a is NaN or ±Inf; ErrOutOfRange when a ≤ 3 (the variance is infinite; for
// 2 < a ≤ 3 the mean alone is finite, use ZipfPMF sums if needed).
//
// ID:      NUM-PROB-146
// Time:    O(1)   three ζ evaluations
// Space:   O(1), no allocation
// Source:  Johnson, Kemp, Kotz, Univariate Discrete Distributions, 3rd ed., Wiley, 2005, ch. 11, §11.2.20 (moments μ'_r = ζ(a − r)/ζ(a))
// Tags:    probability, distribution, zipf, zeta, power-law, mean, variance, moments, discrete
// Units:   a: dimensionless; mean: a rank; variance: a rank squared
// Example: a = 4 → mean ζ(3)/ζ(4) ≈ 1.1106265353261484, variance ζ(2)/ζ(4) − mean² ≈ 0.2863264536645025
func ZipfMeanVariance(a float64) (mean, variance float64, err error) {
	if err := checkZipf(a); err != nil {
		return 0, 0, err
	}
	if a <= 3 {
		return 0, 0, fmt.Errorf("%w: a = %v, variance infinite for a <= 3", ErrOutOfRange, a)
	}
	z0, err := zipfZeta(a)
	if err != nil {
		return 0, 0, err
	}
	z1, err := zipfZeta(a - 1)
	if err != nil {
		return 0, 0, err
	}
	z2, err := zipfZeta(a - 2)
	if err != nil {
		return 0, 0, err
	}
	mean = z1 / z0
	return mean, max(0, (z2*z0-z1*z1)/(z0*z0)), nil
}
