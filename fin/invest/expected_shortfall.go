package invest

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
	"github.com/CrimsonCoder42/calcbench/algo/selection"
)

// ExpectedShortfall returns the expected shortfall (conditional VaR) of the empirical distribution of
// returns at the given confidence, as a positive loss fraction: with α = 1 − confidence, w = n·α and
// m = ⌊w⌋, ES = −( X₍₁₎ + … + X₍ₘ₎ + (w − m)·X₍ₘ₊₁₎ )/w, X₍ᵢ₎ the order statistics. This is the
// source's definition applied to the empirical distribution, including its correction for the atom at the
// quantile, so it is continuous in α and coherent; when n·α is an integer it is the mean of the n·α worst
// returns. The confidence level is the caller's. returns is read, never modified.
// Errors: ErrTooFew when returns is empty; ErrNotFinite for NaN or ±Inf; ErrDomain when confidence is not
// in (0, 1); ErrOverflow when the tail sum exceeds float64.
//
// ID:      FIN-INVEST-051
// Time:    O(n)   n = len(returns); one linear-time selection and one pass
// Space:   O(n) extra (one allocation, the selection's copy)
// Source:  Acerbi, Tasche, "On the coherence of expected shortfall", Journal of Banking & Finance, 2002, 26(7):1487–1503, doi:10.1016/S0378-4266(02)00283-2, Definition 2.6 (ES with the quantile-atom correction) and Proposition 4.1 (its empirical estimator)
// Tags:    invest, risk, expected-shortfall, cvar, conditional-var, tail-risk, coherent-risk-measure
// Units:   returns decimal per period; confidence a probability in (0, 1); result decimal loss per period
// Example: hand-checked: returns −0.05, −0.03, −0.01, 0, 0.01, 0.02, 0.02, 0.03, 0.04, 0.05: confidence 0.8 → (0.05 + 0.03)/2 = 0.04; 0.85 → (0.05 + 0.5 × 0.03)/1.5 = 0.04333…
func ExpectedShortfall(returns []float64, confidence float64) (float64, error) {
	if err := checkSeries("returns", returns, 1); err != nil {
		return 0, err
	}
	if err := checkConfidence(confidence); err != nil {
		return 0, err
	}
	n := len(returns)
	w := float64(n) * (1 - confidence)
	m := int(math.Floor(w))
	if m >= n { // only by rounding in n·α; w < n in exact arithmetic
		m = n - 1
	}
	t, err := selection.KthSmallest(returns, m+1)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrDomain, err)
	}
	// Σ of the m smallest = Σ{x < t} + (m − #{x < t})·t, since t is the (m+1)-th smallest.
	var s, c float64
	below := 0
	for _, x := range returns {
		if x < t {
			s, c = reduce.NeumaierAdd(s, c, x)
			below++
		}
	}
	s, c = reduce.NeumaierAdd(s, c, float64(m-below)*t)
	s, c = reduce.NeumaierAdd(s, c, (w-float64(m))*t)
	return finiteResult(-(s + c) / w)
}
