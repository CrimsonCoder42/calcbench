package prob

// PoissonCDF returns P(X ≤ k) = Σ_{i=0..k} e^(−λ)·λⁱ/i! for a Poisson variable with mean λ; 0 for k < 0.
// It sums the tail on the far side of k from λ (the lower tail when k < λ, else the upper tail, returned
// as 1 − tail) from one saddle-point term outward with the ratio recurrence P(i−1) = P(i)·i/λ or
// P(i+1) = P(i)·λ/(i+1), compensated, and stops when the rest is below 2⁻⁵³ of the sum.
// Errors: as PoissonPMF; ErrOutOfRange also when the sum would need more than 2²⁸ terms (λ above roughly
// 10¹⁴ with k near λ), for which the incomplete-gamma form is not yet available.
//
// ID:      NUM-PROB-042
// Time:    O(min(k + 1, c·(√λ + 1)))   c ≈ 10 in double precision (terms beyond c·√λ from λ are below 2⁻⁵³)
// Space:   O(1), no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.19 "Poisson Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda366j.htm (read 2026-11-12)
// Source:  Loader, "Fast and Accurate Computation of Binomial Probabilities", technical report, 2000 (rev. 2002), https://www.r-project.org/doc/reports/CLoader-dbinom-2002.pdf (read 2026-11-12), §2 eqs. (3)–(7) and Appendix A (stirlerr, bd0)
// Tags:    probability, distribution, poisson, cdf, discrete, counts, compensated-summation
// Units:   k: a count of events; lambda: the mean count per interval; result: a probability in [0, 1]
// Example: λ = 2, k = 2 → 5e⁻² ≈ 0.6766764161830635
func PoissonCDF(k int, lambda float64) (float64, error) {
	if err := checkPoisson(lambda); err != nil {
		return 0, err
	}
	if k < 0 {
		return 0, nil
	}
	if lambda == 0 {
		return 1, nil
	}
	var s compensatedSum
	if float64(k) < lambda {
		t := poissonTerm(float64(k), lambda)
		s.add(t)
		for i, steps := k, 0; i > 0; i, steps = i-1, steps+1 {
			if steps == maxTailTerms {
				return 0, errTailTooLong("Poisson tail")
			}
			r := float64(i) / lambda
			t *= r
			s.add(t)
			if tailDone(t, r, &s) {
				break
			}
		}
		return min(1, s.total()), nil
	}
	// k ≥ λ: the upper tail from k + 1, in float64 so that k = MaxInt cannot overflow the index. Above
	// 2⁵³ the index can no longer step by 1; the term cap is reached long before that matters.
	x := float64(k) + 1
	t := poissonTerm(x, lambda)
	s.add(t)
	for steps := 0; t != 0; x, steps = x+1, steps+1 {
		if steps == maxTailTerms || x >= 1<<53 {
			return 0, errTailTooLong("Poisson tail")
		}
		r := lambda / (x + 1)
		t *= r
		s.add(t)
		if tailDone(t, r, &s) {
			break
		}
	}
	return 1 - min(1, s.total()), nil
}
