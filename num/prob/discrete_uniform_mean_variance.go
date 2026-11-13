package prob

// DiscreteUniformMeanVariance returns the mean (a + b)/2 and the variance ((b − a + 1)² − 1)/12 of X
// uniform on the integers a…b. The mean is a + ⌊(b − a)/2⌋ in integers (exact, cannot overflow) plus ½
// when b − a is odd, then converted to float64; the variance is (m − 1)(m + 1)/12 in float64 with
// m = b − a + 1. No int overflows and every int range is accepted.
// Errors: as DiscreteUniformPMF.
//
// ID:      NUM-PROB-040
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (discrete uniform distribution on 1…N, shifted here to a…b)
// Tags:    probability, distribution, discrete-uniform, mean, variance, expected-value, dice
// Units:   a, b: integer outcomes; mean: that unit; variance: that unit squared
// Example: a fair die (1, 6) → (3.5, 35/12 ≈ 2.9166666666666665)
func DiscreteUniformMeanVariance(a, b int) (mean, variance float64, err error) {
	d, m, err := discreteUniformCount(a, b)
	if err != nil {
		return 0, 0, err
	}
	// mean = a + (b − a)/2: the integer part a + ⌊d/2⌋ lies in a…b (fits int), plus ½ when d is odd.
	// (float64(a)/2 + float64(b)/2 would round MaxInt up and give 0, not −½, on MinInt…MaxInt.)
	mean = float64(int(uint64(a)+d/2)) + 0.5*float64(d&1)
	return mean, (m - 1) * (m + 1) / 12, nil
}
