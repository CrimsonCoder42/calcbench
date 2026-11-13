package prob

import "fmt"

// discreteUniformCount validates a ≤ b and returns the span d = b − a, exact in uint64 (the true difference
// is in [0, 2⁶⁴ − 1] for any ints a ≤ b), and m = b − a + 1, the number of support points, as a float64.
// m is up to 2⁶⁴ (a = MinInt, b = MaxInt), which uint64 and int cannot hold but float64 can (rounded to
// 53 bits above 2⁵³), so no int range limits the support.
func discreteUniformCount(a, b int) (d uint64, m float64, err error) {
	if a > b {
		return 0, 0, fmt.Errorf("%w: a = %d > b = %d", ErrOutOfRange, a, b)
	}
	d = uint64(b) - uint64(a)
	return d, float64(d) + 1, nil
}

// DiscreteUniformPMF returns P(X = k) = 1/(b − a + 1) for an integer k in a…b, each equally likely; 0
// elsewhere. One fair n-sided die is a = 1, b = n.
// Every int range is accepted, a = MinInt, b = MaxInt included (m = 2⁶⁴ support points, P = 2⁻⁶⁴); m is
// held in float64, so above 2⁵³ points it is rounded to 53 bits.
// Errors: ErrOutOfRange when a > b.
//
// ID:      NUM-PROB-038
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (discrete uniform distribution on 1…N, shifted here to a…b)
// Tags:    probability, distribution, discrete-uniform, pmf, discrete, dice
// Units:   k, a, b: integer outcomes; result: a probability in [0, 1]
// Example: a fair die (1, 6), k = 4 → 1/6
func DiscreteUniformPMF(k, a, b int) (float64, error) {
	_, m, err := discreteUniformCount(a, b)
	if err != nil {
		return 0, err
	}
	if k < a || k > b {
		return 0, nil
	}
	return 1 / m, nil
}
