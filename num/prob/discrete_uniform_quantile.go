package prob

import "math"

// DiscreteUniformQuantile returns the smallest integer k in a…b with P(X ≤ k) ≥ u for X uniform on the
// integers a…b: k = a + ⌈u·m⌉ − 1 with m = b − a + 1 (a at u = 0), corrected by single steps until
// DiscreteUniformCDF(k−1) < u ≤ DiscreteUniformCDF(k) holds for the CDF as this package computes it.
// Every int range is accepted, a = MinInt, b = MaxInt included. Above 2⁵³ support points the float64 CDF
// cannot tell neighbouring k apart, so the result is the least k at the CDF's own resolution (on
// MinInt…MaxInt, u = ½ gives −513, whose CDF already rounds to ½; the exact-rational answer is −1).
// Errors: ErrNotFinite when u is NaN or ±Inf; ErrProbability when u is outside [0, 1]; ErrOutOfRange
// when a > b.
//
// ID:      NUM-PROB-058
// Time:    O(1 + m·2⁻⁵²)   m = b − a + 1 (the float64 estimate is off by at most a few units in the last place of m, each a correction step)
// Space:   O(1), no allocation
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (discrete uniform distribution on 1…N, shifted here to a…b), inverted
// Tags:    probability, distribution, discrete-uniform, quantile, inverse-cdf, dice
// Units:   u: a probability; a, b and result: integer outcomes
// Example: a fair die (1, 6), u = 0.5 → 3
func DiscreteUniformQuantile(u float64, a, b int) (int, error) {
	if err := checkProbability("u", u); err != nil {
		return 0, err
	}
	d, mf, err := discreteUniformCount(a, b)
	if err != nil {
		return 0, err
	}
	cdf := func(k int) float64 { // DiscreteUniformCDF for a ≤ k < b
		return float64(uint64(k)-uint64(a)+1) / mf
	}
	// offset = ⌈u·m⌉ − 1 in [0, d]; computed in float64 and clamped to d in uint64 (float64(d) may round
	// up to 2⁶⁴, which uint64 cannot hold), then added in uint64 so that a + offset cannot overflow.
	off := max(0, math.Ceil(u*mf)-1)
	var o uint64
	if off >= float64(d) {
		o = d
	} else {
		o = uint64(off)
	}
	k := int(uint64(a) + o)
	for k > a && cdf(k-1) >= u {
		k--
	}
	for k < b && cdf(k) < u {
		k++
	}
	return k, nil
}
