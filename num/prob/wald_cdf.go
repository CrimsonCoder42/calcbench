package prob

import "math"

// scaledNormalTail returns Φ(−b)·e^(b²/2) for b ≥ 0: directly while that stays in range (b < 26), and by
// the asymptotic series (1/(b√(2π)))·Σ (−1)ʲ (2j − 1)!!/b²ʲ above (Abramowitz & Stegun 26.2.12; at b = 26
// the first omitted term is below 1e-17 relative).
func scaledNormalTail(b float64) float64 {
	if b < 26 {
		return standardNormalSurvival(b) * math.Exp(b*b/2)
	}
	inv := 1 / (b * b)
	sum, term := 1.0, 1.0
	for j := 1; j <= 9; j++ {
		term *= -float64(2*j-1) * inv
		sum += term
	}
	return sum / (b * math.Sqrt(2*math.Pi))
}

// waldTails returns (P(X ≤ x), P(X > x)) for x > 0: with a = √(λ/x)(x/μ − 1), b = √(λ/x)(x/μ + 1),
// F = Φ(a) + e^(2λ/μ)Φ(−b) and 1 − F = Φ(−a) − e^(2λ/μ)Φ(−b). The second term is formed as
// e^(−λ(x − μ)²/(2μ²x)) · Φ(−b)e^(b²/2) (the exponents combine exactly), so it never overflows.
func waldTails(x, mean, shape float64) (cdf, survival float64) {
	s := math.Sqrt(shape / x)
	a := s * (x/mean - 1)
	b := s * (x/mean + 1)
	r := (x - mean) / mean
	t := math.Exp(-shape*r*r/(2*x)) * scaledNormalTail(b)
	cdf = min(1, standardNormalCDF(a)+t)
	survival = max(0, standardNormalSurvival(a)-t)
	return cdf, survival
}

// WaldCDF returns P(X ≤ x) = Φ(√(λ/x)(x/μ − 1)) + e^(2λ/μ) Φ(−√(λ/x)(x/μ + 1)) for the inverse Gaussian
// (Wald) distribution with mean μ and shape λ; x ≤ 0 gives 0. Φ is NormalCDF's kernel (NUM-PROB-048); the
// second term is evaluated without forming e^(2λ/μ), so large λ/μ does not overflow.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when mean ≤ 0 or shape ≤ 0.
//
// ID:      NUM-PROB-139
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Chhikara, Folks, The Inverse Gaussian Distribution, Marcel Dekker, 1989, ch. 2 (the CDF in terms of Φ)
// Source:  Abramowitz, Stegun, Handbook of Mathematical Functions, 1964, eq. 26.2.12 (the normal tail series for large arguments)
// Tags:    probability, distribution, wald, inverse-gaussian, cdf, continuous, first-passage-time
// Units:   x, mean, shape: one unit; result: a probability in [0, 1]
// Example: μ = 1, λ = 1, x = 1 → Φ(0) + e²Φ(−2) ≈ 0.6681020012231706
func WaldCDF(x, mean, shape float64) (float64, error) {
	if err := checkWald(mean, shape); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	if x <= 0 {
		return 0, nil
	}
	c, _ := waldTails(x, mean, shape)
	return c, nil
}
