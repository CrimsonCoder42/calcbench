package prob

import "fmt"

// NormalIntervalProbability returns P(lo ≤ X ≤ hi) for X normal with mean μ and standard deviation σ:
// Φ(z_hi) − Φ(z_lo) when z_hi ≤ 0, Q(z_lo) − Q(z_hi) when z_lo ≥ 0, and 1 − Φ(z_lo) − Q(z_hi) across 0,
// so that every subtraction is between the small tail values it needs. For an interval in z units pass
// mean 0 and sd 1.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when sd ≤ 0 or lo > hi.
//
// ID:      NUM-PROB-050
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.1 "Normal Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3661.htm (read 2026-11-12)
// Source:  NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §7.2 eqs. 7.2.1–7.2.2 (erf, erfc) and §7.1 (Φ(z) = ½erfc(−z/√2), Q(z) = ½erfc(z/√2))
// Tags:    probability, distribution, normal, gaussian, interval, continuous, z-score
// Units:   lo, hi, mean, sd: one unit; result: a probability in [0, 1]
// Example: lo = −1.96, hi = 1.96, mean 0, sd 1 → 0.950004209703559; z from 8 to 9 → 6.22e-16 (Φ − Φ gives 0)
func NormalIntervalProbability(lo, hi, mean, sd float64) (float64, error) {
	if err := checkNormal(mean, sd); err != nil {
		return 0, err
	}
	if err := checkFinite("lo", lo); err != nil {
		return 0, err
	}
	if err := checkFinite("hi", hi); err != nil {
		return 0, err
	}
	if lo > hi {
		return 0, fmt.Errorf("%w: lo = %v > hi = %v", ErrOutOfRange, lo, hi)
	}
	zlo, zhi := standardize(lo, mean, sd), standardize(hi, mean, sd)
	var p float64
	switch {
	case zhi <= 0:
		p = standardNormalCDF(zhi) - standardNormalCDF(zlo)
	case zlo >= 0:
		p = standardNormalSurvival(zlo) - standardNormalSurvival(zhi)
	default:
		p = 1 - standardNormalCDF(zlo) - standardNormalSurvival(zhi)
	}
	return min(1, max(0, p)), nil
}
