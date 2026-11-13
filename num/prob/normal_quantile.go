package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Coefficients of Wichura's PPND16, Applied Statistics 37(3), 1988, as printed in the algorithm's
// published listing (Algorithm AS 241). Region 1: |u − ½| ≤ 0.425, r = 0.180625 − (u − ½)².
var (
	ppndA = [8]float64{
		3.3871328727963666080e0, 1.3314166789178437745e+2, 1.9715909503065514427e+3, 1.3731693765509461125e+4,
		4.5921953931549871457e+4, 6.7265770927008700853e+4, 3.3430575583588128105e+4, 2.5090809287301226727e+3,
	}
	ppndB = [8]float64{
		1, 4.2313330701600911252e+1, 6.8718700749205790830e+2, 5.3941960214247511077e+3,
		2.1213794301586595867e+4, 3.9307895800092710610e+4, 2.8729085735721942674e+4, 5.2264952788528545610e+3,
	}
	// Region 2: r = √(−ln min(u, 1 − u)) ≤ 5, polynomials in r − 1.6.
	ppndC = [8]float64{
		1.42343711074968357734e0, 4.63033784615654529590e0, 5.76949722146069140550e0, 3.64784832476320460504e0,
		1.27045825245236838258e0, 2.41780725177450611770e-1, 2.27238449892691845833e-2, 7.74545014278341407640e-4,
	}
	ppndD = [8]float64{
		1, 2.05319162663775882187e0, 1.67638483018380384940e0, 6.89767334985100004550e-1,
		1.48103976427480074590e-1, 1.51986665636164571966e-2, 5.47593808499534494600e-4, 1.05075007164441684324e-9,
	}
	// Region 3: r > 5, polynomials in r − 5.
	ppndE = [8]float64{
		6.65790464350110377720e0, 5.46378491116411436990e0, 1.78482653991729133580e0, 2.96560571828504891230e-1,
		2.65321895265761230930e-2, 1.24266094738807843860e-3, 2.71155556874348757815e-5, 2.01033439929228813265e-7,
	}
	ppndF = [8]float64{
		1, 5.99832206555887937690e-1, 1.36929880922735805310e-1, 1.48753612908506148525e-2,
		7.86869131145613259100e-4, 1.84631831751005468180e-5, 1.42151175831644588870e-7, 2.04426310338993978564e-15,
	}
)

// horner7 evaluates c[0] + c[1]·x + … + c[7]·x⁷.
func horner7(c *[8]float64, x float64) float64 {
	s := c[7]
	for i := 6; i >= 0; i-- {
		s = s*x + c[i]
	}
	return s
}

// standardNormalQuantile returns Φ⁻¹(u) for 0 < u < 1 by Wichura's PPND16: a ratio of degree-7
// polynomials in one of three regions, with no iteration.
func standardNormalQuantile(u float64) float64 {
	q := u - 0.5
	if math.Abs(q) <= 0.425 {
		r := 0.180625 - q*q
		return q * horner7(&ppndA, r) / horner7(&ppndB, r)
	}
	r := u
	if q > 0 {
		r = 1 - u
	}
	r = math.Sqrt(-math.Log(r))
	var z float64
	if r <= 5 {
		r -= 1.6
		z = horner7(&ppndC, r) / horner7(&ppndD, r)
	} else {
		r -= 5
		z = horner7(&ppndE, r) / horner7(&ppndF, r)
	}
	if q < 0 {
		return -z
	}
	return z
}

// checkOpenUnit validates u finite and strictly inside (0, 1).
func checkOpenUnit(u float64) error {
	if err := checkFinite("u", u); err != nil {
		return err
	}
	if u <= 0 || u >= 1 {
		return fmt.Errorf("%w: u = %v, want 0 < u < 1", ErrOutOfRange, u)
	}
	return nil
}

// NormalQuantile returns x = μ + σ·Φ⁻¹(u), the value below which a normal variable with mean μ and
// standard deviation σ falls with probability u. Φ⁻¹ is Wichura's PPND16 (Algorithm AS 241), accurate to
// about 1 part in 10¹⁶, with no iteration. u = 0 and u = 1 (the infinite quantiles) are errors, not ±Inf.
// For u above ½ the accuracy is limited by 1 − u, which the caller's u already rounds.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when u is outside (0, 1) or sd ≤ 0;
// ErrOverflow when μ + σ·Φ⁻¹(u) exceeds the float64 range.
//
// ID:      NUM-PROB-061
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Wichura, "Algorithm AS 241: The Percentage Points of the Normal Distribution", Journal of the Royal Statistical Society, Series C (Applied Statistics), 1988, 37(3):477–484, doi:10.2307/2347330 (PPND16)
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.7.1 "Cumulative Distribution Function of the Standard Normal Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3671.htm (read 2026-11-12)
// Tags:    probability, distribution, normal, gaussian, quantile, inverse-cdf, ppf, critical-value, z-score
// Units:   u: a probability in (0, 1); mean, sd and result: one unit
// Example: NIST §1.3.6.7.1 critical values: u = 0.975 → 1.960; 0.995 → 2.576; 0.999 → 3.090
func NormalQuantile(u, mean, sd float64) (float64, error) {
	if err := checkOpenUnit(u); err != nil {
		return 0, err
	}
	if err := checkNormal(mean, sd); err != nil {
		return 0, err
	}
	x := mean + sd*standardNormalQuantile(u)
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: quantile for mean = %v, sd = %v", ErrOverflow, mean, sd)
	}
	return x, nil
}
