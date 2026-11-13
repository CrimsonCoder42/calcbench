package reduce

import "math"

// NeumaierAdd adds x to a running compensated sum and returns the new sum and correction. Start with
// sum = comp = 0, call it once per value, and read the total as newSum + newComp. It is the step Sum
// uses, exposed for streaming sums whose values arrive one at a time. Inputs are not checked: a NaN or Inf
// propagates, and the caller tests the total (AllFinite).
//
// ID:      ALG-RED-002
// Time:    O(1)
// Space:   O(1) (no allocation)
// Source:  Neumaier, "Rundungsfehleranalyse einiger Verfahren zur Summation endlicher Summen", ZAMM, 1974, 54(1):39–51, doi:10.1002/zamm.19740540106
// Source:  Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §4.3
// Tags:    sum, summation, compensated, neumaier, streaming, running-sum, no-alloc
func NeumaierAdd(sum, comp, x float64) (newSum, newComp float64) {
	t := sum + x
	// The larger operand is exact in t; the error of the addition is what the smaller one lost.
	if math.Abs(sum) >= math.Abs(x) {
		comp += (sum - t) + x
	} else {
		comp += (x - t) + sum
	}
	return t, comp
}
