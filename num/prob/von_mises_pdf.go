package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/num/special"
)

// vonMisesMaxKappa bounds κ: the CDF series needs about 10√κ terms, and above 10⁶ the distribution is
// normal with σ = 1/√κ to far below float64 resolution of the angle.
const vonMisesMaxKappa = 1e6

// checkVonMises validates the location μ, concentration κ ∈ [0, 10⁶] and x ∈ [μ − π, μ + π].
func checkVonMises(x, location, kappa float64) error {
	if err := checkFinite("location", location); err != nil {
		return err
	}
	if err := checkFinite("kappa", kappa); err != nil {
		return err
	}
	if err := checkFinite("x", x); err != nil {
		return err
	}
	if kappa < 0 || kappa > vonMisesMaxKappa {
		return fmt.Errorf("%w: kappa = %v, want 0..1e6", ErrOutOfRange, kappa)
	}
	if math.Abs(x-location) > math.Pi {
		return fmt.Errorf("%w: x = %v, want within π of location %v", ErrOutOfRange, x, location)
	}
	return nil
}

// scaledBesselI0 returns I₀(κ)·e^(−κ) for κ ≥ 0: special.BesselI0 below 700, above it the asymptotic series
// I₀(κ)e^(−κ) ≈ (2πκ)^(−½) Σ ((2j − 1)!!)²/(j!(8κ)ʲ) (Abramowitz & Stegun 9.7.1; at κ = 700 the terms
// fall below 1e-17 relative after six).
func scaledBesselI0(kappa float64) (float64, error) {
	if kappa < 700 {
		i0, err := special.BesselI0(kappa)
		if err != nil {
			return 0, specialErr("BesselI0", err)
		}
		return i0 * math.Exp(-kappa), nil
	}
	sum, term := 1.0, 1.0
	for j := 1; j <= 8; j++ {
		c := float64(2*j - 1)
		term *= c * c / (float64(j) * 8 * kappa)
		sum += term
	}
	return sum / math.Sqrt(2*math.Pi*kappa), nil
}

// VonMisesPDF returns the von Mises (circular normal) density f(x) = e^(κ cos(x − μ)) / (2π I₀(κ)) on
// x ∈ [μ − π, μ + π], location μ, concentration κ ≥ 0 (κ = 0 is uniform, 1/(2π)). It is computed as
// e^(κ(cos(x − μ) − 1)) / (2π I₀(κ)e^(−κ)), so large κ does not overflow.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when κ < 0, κ > 10⁶, or |x − μ| > π.
//
// ID:      NUM-PROB-142
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Mardia, Jupp, Directional Statistics, Wiley, 2000, §3.5.4 (the von Mises density), with I₀ from special.BesselI0 (NUM-SPEC) and A&S 9.7.1 for κ ≥ 700
// Tags:    probability, distribution, von-mises, circular, pdf, continuous, angle, directional
// Units:   x, location: radians; kappa: dimensionless; result: density per radian
// Example: μ = 0, κ = 0, any x → 1/(2π) ≈ 0.15915494309189535
func VonMisesPDF(x, location, kappa float64) (float64, error) {
	if err := checkVonMises(x, location, kappa); err != nil {
		return 0, err
	}
	s, err := scaledBesselI0(kappa)
	if err != nil {
		return 0, err
	}
	return math.Exp(kappa*(math.Cos(x-location)-1)) / (2 * math.Pi * s), nil
}
