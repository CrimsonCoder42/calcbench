package prob

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// ExpectedValue returns E[X] = Σᵢ xᵢ·pᵢ for a discrete distribution with outcomes xᵢ and probabilities
// pᵢ. probabilities must be a distribution (sum 1 within the package tolerance n·2⁻⁵²). Each product's
// rounding error is recovered with an FMA and the sum is compensated (Neumaier). A partial sum that
// overflows is redone with the outcomes scaled by 2⁻², exactly, so a finite mean of finite outcomes is
// always returned. Money stays in the caller's integer minor units, converted once.
// Errors: ErrOutOfRange when the lengths differ or are 0; ErrNotFinite for a NaN or ±Inf element;
// ErrProbability, ErrNotDistribution as TotalProbability; ErrOverflow when the mean itself overflows.
//
// ID:      NUM-PROB-018
// Time:    O(n)   n = len(outcomes)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.3 (expected value)
// Source:  Neumaier, "Rundungsfehleranalyse einiger Verfahren zur Summation endlicher Summen", ZAMM, 1974, 54(1):39–51, doi:10.1002/zamm.19740540106
// Tags:    probability, expectation, expected-value, mean, discrete, compensated-summation
// Units:   outcomes and result: the caller's unit; probabilities: in [0, 1]
// Example: a fair die (1…6, each 1/6) → 3.5
func ExpectedValue(outcomes, probabilities []float64) (float64, error) {
	if err := validateDiscrete(outcomes, probabilities); err != nil {
		return 0, err
	}
	return discreteMean(outcomes, probabilities)
}

// validateDiscrete checks the paired lengths, every outcome finite, and probabilities a distribution.
func validateDiscrete(outcomes, probabilities []float64) error {
	if err := checkPairedLengths("outcomes", len(outcomes), "probabilities", len(probabilities)); err != nil {
		return err
	}
	for i, x := range outcomes {
		if !reduce.AllFinite(x) {
			return fmt.Errorf("%w: outcomes[%d] = %v", ErrNotFinite, i, x)
		}
	}
	return checkDistribution("probabilities", probabilities)
}

// discreteMean returns Σ xᵢ·pᵢ for validated inputs, compensated, with one scaled retry on overflow.
func discreteMean(outcomes, probabilities []float64) (float64, error) {
	var s compensatedSum
	for i, x := range outcomes {
		s.addProduct(x, probabilities[i])
	}
	if m := s.total(); reduce.AllFinite(m) {
		return m, nil
	}
	// Partial sums are bounded by max|xᵢ|·Σpᵢ ≤ max|xᵢ|·(1 + n·2⁻⁵²), so a quarter of every outcome
	// cannot overflow; the scaling by 2⁻² is exact for normal numbers.
	s = compensatedSum{}
	for i, x := range outcomes {
		s.addProduct(x*0.25, probabilities[i])
	}
	m := s.total() * 4
	if !reduce.AllFinite(m) {
		return 0, fmt.Errorf("%w: expected value", ErrOverflow)
	}
	return m, nil
}
