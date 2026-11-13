package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// compensatedSum is a running sum with Neumaier's correction: the true total of the added terms is
// sum + comp up to the error of comp itself, O(n·u²) for n terms and unit roundoff u. The zero value is
// an empty sum. Each step is reduce.NeumaierAdd (ALG-RED-002); the FMA product correction stays here.
// Source: Neumaier, ZAMM, 1974, 54(1):39–51, doi:10.1002/zamm.19740540106; Higham, 2002, §4.3.
type compensatedSum struct {
	sum, comp float64
}

// add adds x and keeps the exact rounding error of that addition in comp.
func (s *compensatedSum) add(x float64) {
	s.sum, s.comp = reduce.NeumaierAdd(s.sum, s.comp, x)
}

// addProduct adds x·y and the exact rounding error of that product (recovered with one FMA).
func (s *compensatedSum) addProduct(x, y float64) {
	p := x * y
	s.add(p)
	s.comp += math.FMA(x, y, -p)
}

// total returns the corrected sum. It is NaN or ±Inf when a partial sum overflowed.
func (s *compensatedSum) total() float64 {
	return s.sum + s.comp
}

// distributionTolerance returns the largest |Σpᵢ − 1| accepted for n probabilities: n·2⁻⁵², the
// first-order error bound n·u (u = 2⁻⁵³) of a recursive sum, doubled to cover probabilities that were
// themselves computed with one rounding each. (Package policy; see the package comment.)
func distributionTolerance(n int) float64 {
	return float64(n) * 0x1p-52
}

// checkDistribution validates every element of probs as a probability and their sum as 1 within
// distributionTolerance. It returns ErrNotFinite, ErrProbability or ErrNotDistribution, naming the index.
func checkDistribution(name string, probs []float64) error {
	var s compensatedSum
	for i, p := range probs {
		if !reduce.AllFinite(p) {
			return fmt.Errorf("%w: %s[%d] = %v", ErrNotFinite, name, i, p)
		}
		if p < 0 || p > 1 {
			return fmt.Errorf("%w: %s[%d] = %v", ErrProbability, name, i, p)
		}
		s.add(p)
	}
	s.add(-1)
	if d := s.total(); math.Abs(d) > distributionTolerance(len(probs)) {
		return fmt.Errorf("%w: %s sums to 1 %+g", ErrNotDistribution, name, d)
	}
	return nil
}

// checkProbabilities validates every element of probs as a probability, naming the index.
func checkProbabilities(name string, probs []float64) error {
	for i, p := range probs {
		if !reduce.AllFinite(p) {
			return fmt.Errorf("%w: %s[%d] = %v", ErrNotFinite, name, i, p)
		}
		if p < 0 || p > 1 {
			return fmt.Errorf("%w: %s[%d] = %v", ErrProbability, name, i, p)
		}
	}
	return nil
}

// checkPairedLengths returns ErrOutOfRange when the two slices differ in length or are empty.
func checkPairedLengths(nameA string, a int, nameB string, b int) error {
	if a != b {
		return fmt.Errorf("%w: len(%s) = %d, len(%s) = %d", ErrOutOfRange, nameA, a, nameB, b)
	}
	if a == 0 {
		return fmt.Errorf("%w: %s is empty", ErrOutOfRange, nameA)
	}
	return nil
}

// TotalProbability returns P(A) = Σᵢ P(A | Bᵢ)·P(Bᵢ), the law of total probability over a partition
// B₁…Bₙ. pB must be a distribution: its sum is 1 within the package tolerance n·2⁻⁵². The products are
// summed with Neumaier compensation, each product's own rounding error recovered by an FMA; a result
// that overshoots 1 by rounding (at most the tolerance) is returned as 1.
// Errors: ErrOutOfRange when the lengths differ or are 0; ErrNotFinite for a NaN or ±Inf element;
// ErrProbability for an element outside [0, 1]; ErrNotDistribution when pB does not sum to 1.
//
// ID:      NUM-PROB-007
// Time:    O(n)   n = len(pB)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §3.3 "Bayes's Formula" (total probability over a partition)
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §1.2 "Basics of Probability Theory"
// Source:  Neumaier, "Rundungsfehleranalyse einiger Verfahren zur Summation endlicher Summen", ZAMM, 1974, 54(1):39–51, doi:10.1002/zamm.19740540106
// Tags:    probability, total-probability, partition, conditional, compensated-summation, rules
// Units:   every input and the result: probabilities in [0, 1]
// Example: two urns chosen with probability 0.5 each, red drawn with 0.3 and 0.6 → 0.45
func TotalProbability(pAGivenB, pB []float64) (float64, error) {
	if err := checkPairedLengths("pAGivenB", len(pAGivenB), "pB", len(pB)); err != nil {
		return 0, err
	}
	if err := checkProbabilities("pAGivenB", pAGivenB); err != nil {
		return 0, err
	}
	if err := checkDistribution("pB", pB); err != nil {
		return 0, err
	}
	var s compensatedSum
	for i, c := range pAGivenB {
		s.addProduct(c, pB[i])
	}
	return min(1, max(0, s.total())), nil
}
