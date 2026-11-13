package prob

// DiscreteUniformCDF returns P(X ≤ k) = (k − a + 1)/(b − a + 1) for X uniform on the integers a…b: 0 below
// a and 1 from b on.
// Errors: as DiscreteUniformPMF.
//
// ID:      NUM-PROB-039
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (discrete uniform distribution on 1…N, shifted here to a…b)
// Tags:    probability, distribution, discrete-uniform, cdf, discrete, dice
// Units:   k, a, b: integer outcomes; result: a probability in [0, 1]
// Example: a fair die (1, 6), k = 4 → 4/6 ≈ 0.6666666666666666
func DiscreteUniformCDF(k, a, b int) (float64, error) {
	_, m, err := discreteUniformCount(a, b)
	if err != nil {
		return 0, err
	}
	switch {
	case k < a:
		return 0, nil
	case k >= b:
		return 1, nil
	}
	return float64(uint64(k)-uint64(a)+1) / m, nil // a ≤ k < b, so k − a + 1 ≤ b − a fits in uint64
}
