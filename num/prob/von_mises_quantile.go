package prob

import "math"

// VonMisesQuantile returns the u-quantile of the von Mises distribution on [μ − π, μ + π]: the x with
// VonMisesCDF(x, μ, κ) = u. The CDF is continuous and increasing in θ = x − μ, so it bisects θ on the
// bracket [−π, π] (F = 0 and 1 there) until the midpoint equals an end point (at most 64 halvings, about
// 2⁻⁶²·π) or the series hits u exactly; θ = 0 is the first midpoint, so u = ½ returns μ. κ = 0 (the uniform
// circle) is the closed form θ = 2πu − π. The CDF's error is absolute (a few ulps of 1), so the quantile's
// error is that divided by the density, which is largest where the density is smallest.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when u is outside [0, 1];
// ErrOutOfRange when κ < 0 or κ > 10⁶.
//
// ID:      NUM-PROB-148
// Time:    O(√κ + 1) per CDF evaluation, at most 64 evaluations
// Space:   O(1), no allocation
// Source:  Hill, "Algorithm 518: Incomplete Bessel Function I0. The Von Mises Distribution", ACM Transactions on Mathematical Software, 1977, 3(3):279–284, doi:10.1145/355744.355753 (the CDF series, inverted here)
// Source:  Burden, Faires, Numerical Analysis, 9th ed., Brooks/Cole, 2011, §2.1 (bisection)
// Tags:    probability, distribution, von-mises, circular, quantile, inverse-cdf, continuous, angle, directional
// Units:   u: a probability; location and result: radians; kappa: dimensionless
// Example: any κ, u = 0.5 → μ; κ = 0, u = 0.75 → μ + π/2
func VonMisesQuantile(u, location, kappa float64) (float64, error) {
	if err := checkProbability("u", u); err != nil {
		return 0, err
	}
	if err := checkVonMises(location, location, kappa); err != nil {
		return 0, err
	}
	switch {
	case u == 0:
		return location - math.Pi, nil
	case u == 1:
		return location + math.Pi, nil
	case kappa == 0:
		return location + (2*math.Pi*u - math.Pi), nil
	}
	lo, hi := -math.Pi, math.Pi // F(lo) = 0 < u < 1 = F(hi)
	for range 64 {
		mid := lo + (hi-lo)/2
		if mid == lo || mid == hi {
			break
		}
		f := vonMisesCDFAngle(mid, kappa)
		if f == u {
			return location + mid, nil
		}
		if f < u {
			lo = mid
		} else {
			hi = mid
		}
	}
	return location + (lo + (hi-lo)/2), nil
}
