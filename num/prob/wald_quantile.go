package prob

// WaldQuantile returns the x > 0 with P(X ≤ x) = u for the inverse Gaussian (Wald) distribution, by
// bisection on the CDF (u ≤ ½) or the survival function (u > ½) of WaldCDF (NUM-PROB-139), starting from
// the mean. The survival function is a difference of two positive terms whose ratio tends to 1 far in the
// upper tail, so there the result's relative accuracy is about (x + μ)/(2μ) times worse than in the body.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when u is outside [0, 1];
// ErrOutOfRange when mean ≤ 0, shape ≤ 0, or u is 0 or 1; ErrOverflow when the quantile exceeds MaxFloat64.
//
// ID:      NUM-PROB-140
// Time:    O(log(x/ε))   bisection steps, each one WaldCDF evaluation
// Space:   O(1), no allocation
// Source:  Chhikara, Folks, The Inverse Gaussian Distribution, Marcel Dekker, 1989, ch. 2 (the CDF), inverted by bisection (Burden, Faires, Numerical Analysis, 9th ed., 2011, §2.1)
// Tags:    probability, distribution, wald, inverse-gaussian, quantile, inverse-cdf, continuous
// Units:   u: a probability; mean, shape and result: one unit
// Example: μ = 1, λ = 1, u = WaldCDF(1, 1, 1) ≈ 0.6681 → 1
func WaldQuantile(u, mean, shape float64) (float64, error) {
	if err := checkQuantileOpen(u); err != nil {
		return 0, err
	}
	if err := checkWald(mean, shape); err != nil {
		return 0, err
	}
	return quantileByBisection(u, 0, mean, func(x float64) (float64, float64, error) {
		if x <= 0 {
			return 0, 1, nil
		}
		c, s := waldTails(x, mean, shape)
		return c, s, nil
	})
}
