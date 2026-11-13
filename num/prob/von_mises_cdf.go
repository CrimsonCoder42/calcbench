package prob

import "math"

// VonMisesCDF returns P(μ − π ≤ X ≤ x) for the von Mises distribution on [μ − π, μ + π], by the Fourier
// series of Hill's Algorithm 518: with θ = x − μ,
// F = (θ + π)/(2π) + (1/π) Σₚ (I_p(κ)/I₀(κ)) sin(pθ)/p, the Bessel ratios from the backward recurrence
// r_p = 1/(2p/κ + r_(p+1)) nested into the sum (Horner form). Hill switches to a normal approximation above
// κ = 10.5; here the series is used for every κ, truncated at P = 30 + ⌈10√κ⌉ terms, where
// I_P/I₀ ≈ e^(−P²/2κ) < e^(−50). The error is absolute, a few ulps of 1, as in Hill's algorithm; a tiny
// lower-tail probability has no relative accuracy.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when κ < 0, κ > 10⁶, or |x − μ| > π.
//
// ID:      NUM-PROB-143
// Time:    O(√κ + 1)
// Space:   O(1), no allocation
// Source:  Hill, "Algorithm 518: Incomplete Bessel Function I0. The Von Mises Distribution", ACM Transactions on Mathematical Software, 1977, 3(3):279–284, doi:10.1145/355744.355753 (the series branch)
// Tags:    probability, distribution, von-mises, circular, cdf, continuous, angle, directional
// Units:   x, location: radians; kappa: dimensionless; result: a probability in [0, 1]
// Example: any κ, x = μ → 0.5 (symmetry)
func VonMisesCDF(x, location, kappa float64) (float64, error) {
	if err := checkVonMises(x, location, kappa); err != nil {
		return 0, err
	}
	return vonMisesCDFAngle(x-location, kappa), nil
}

// vonMisesCDFAngle is the series of VonMisesCDF at θ = x − μ ∈ [−π, π] for a validated κ.
func vonMisesCDFAngle(theta, kappa float64) float64 {
	base := (theta + math.Pi) / (2 * math.Pi)
	if kappa == 0 {
		return base
	}
	p := 30 + int(math.Ceil(10*math.Sqrt(kappa)))
	r, v := 0.0, 0.0
	for n := p; n >= 1; n-- {
		fn := float64(n)
		r = 1 / (2*fn/kappa + r)
		v = r * (math.Sin(fn*theta)/fn + v)
	}
	return min(1, max(0, base+v/math.Pi))
}
