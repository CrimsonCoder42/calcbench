package prob

import "fmt"

// tailEpsilon is the relative size below which the rest of a tail sum is dropped: 2⁻⁵³, the unit roundoff.
const tailEpsilon = 0x1p-53

// maxTailTerms caps a tail sum at 2²⁸ terms (about a second), so every CDF returns in bounded time. A sum
// that needs more (a standard deviation above roughly 10⁷) is ErrOutOfRange: the incomplete-beta and
// incomplete-gamma forms that would handle it in O(1) have no home yet (catalog Q6).
const maxTailTerms = 1 << 28

// errTailTooLong is the error for a tail sum that reached maxTailTerms.
func errTailTooLong(what string) error {
	return fmt.Errorf("%w: %s needs more than 2^28 terms (standard deviation above roughly 1e7)", ErrOutOfRange, what)
}

// tailDone reports whether a tail sum may stop after adding term, whose ratio to the previous term was
// ratio. For a log-concave PMF summed away from its mode the ratios only fall, so the rest of the tail is
// at most term·ratio/(1 − ratio), and the sum stops once that is below tailEpsilon·sum.
func tailDone(term, ratio float64, s *compensatedSum) bool {
	if term == 0 {
		return true
	}
	return ratio < 1 && term*ratio <= tailEpsilon*(1-ratio)*s.total()
}

// binomialTails returns (P(X ≤ k), P(X > k)) for 0 ≤ k < n and 0 < p < 1. It sums whichever side of k
// lies away from the mean, starting from one saddle-point term b(k) or b(k+1) and walking outward with the
// ratio recurrence b(i−1)/b(i) = i·q / ((n−i+1)·p) or b(i+1)/b(i) = (n−i)·p / ((i+1)·q), compensated.
// The other side is 1 minus that sum, which does not cancel: the summed side is the smaller one, up to
// one term.
// It returns errTailTooLong after maxTailTerms terms.
func binomialTails(k, n int, p float64) (lower, upper float64, err error) {
	q := 1 - p
	var s compensatedSum
	if float64(k) < float64(n)*p {
		t := binomialTerm(k, n, p)
		s.add(t)
		for i, steps := k, 0; i > 0; i, steps = i-1, steps+1 {
			if steps == maxTailTerms {
				return 0, 0, errTailTooLong("binomial tail")
			}
			r := float64(i) * q / (float64(n-i+1) * p)
			t *= r
			s.add(t)
			if tailDone(t, r, &s) {
				break
			}
		}
		lower = min(1, s.total())
		return lower, 1 - lower, nil
	}
	t := binomialTerm(k+1, n, p)
	s.add(t)
	for i, steps := k+1, 0; i < n; i, steps = i+1, steps+1 {
		if steps == maxTailTerms {
			return 0, 0, errTailTooLong("binomial tail")
		}
		r := float64(n-i) * p / (float64(i+1) * q)
		t *= r
		s.add(t)
		if tailDone(t, r, &s) {
			break
		}
	}
	upper = min(1, s.total())
	return 1 - upper, upper, nil
}

// BinomialCDF returns P(X ≤ k) = Σ_{i=0..k} C(n, i)·pⁱ·(1 − p)ⁿ⁻ⁱ for X binomial with n trials and success
// probability p: 0 for k < 0 and 1 for k ≥ n. It sums the tail on the far side of k from the mean (the
// lower tail when k < np, else the upper tail, returned as 1 − tail), starting from one saddle-point term
// and using the ratio recurrence, with compensation, and stops when the rest is below 2⁻⁵³ of the sum.
// No term is computed in naive order from i = 0, where the leading terms underflow.
// Errors: as BinomialPMF; ErrOutOfRange also when the sum would need more than 2²⁸ terms (σ above roughly
// 10⁷), for which the incomplete-beta form is not yet available.
//
// ID:      NUM-PROB-027
// Time:    O(min(k + 1, n − k, c·(σ + 1)))   σ = √(np(1 − p)); c ≈ 10 in double precision (the terms beyond c·σ from the mean are below 2⁻⁵³ by Bernstein's inequality)
// Space:   O(1), no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.18 "Binomial Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda366i.htm (read 2026-11-12)
// Source:  Loader, "Fast and Accurate Computation of Binomial Probabilities", technical report, 2000 (rev. 2002), https://www.r-project.org/doc/reports/CLoader-dbinom-2002.pdf (read 2026-11-12), §2 eqs. (3)–(7) and Appendix A (stirlerr, bd0)
// Tags:    probability, distribution, binomial, cdf, discrete, compensated-summation
// Units:   k, n: counts of successes and trials; p and result: probabilities in [0, 1]
// Example: n = 10, p = 0.5, k = 5 → 638/1024 = 0.623046875
func BinomialCDF(k, n int, p float64) (float64, error) {
	if err := checkBinomial(n, p); err != nil {
		return 0, err
	}
	switch {
	case k < 0:
		return 0, nil
	case k >= n:
		return 1, nil
	case p == 0:
		return 1, nil
	case p == 1:
		return 0, nil
	}
	lower, _, err := binomialTails(k, n, p)
	return lower, err
}
