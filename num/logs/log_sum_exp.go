package logs

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// LogSumExp returns log(Σ e^{xᵢ}) of values without overflow or underflow, by the shifted form
// a + log1p(Σ_{i≠k} e^{xᵢ − a}) with a = xₖ = max values: every exponent is ≤ 0, and the largest term (1)
// is kept out of the sum so log1p keeps the small ones. The sum is compensated (Neumaier).
// −Inf entries stand for log 0 and add nothing; all −Inf gives −Inf. values is read, never modified.
// Errors: ErrEmpty when values is empty; ErrNotFinite when a value is NaN or +Inf.
//
// ID:      NUM-LOGS-004
// Time:    O(n)   n = len(values); two passes (max, then the sum)
// Space:   O(1) extra, no allocation
// Source:  Blanchard, Higham, Higham, "Accurately computing the log-sum-exp and softmax functions", IMA Journal of Numerical Analysis, 2021, 41(4):2311–2330, doi:10.1093/imanum/draa038 (§2, Algorithm 2, eq. 2.3)
// Source:  Neumaier, "Rundungsfehleranalyse einiger Verfahren zur Summation endlicher Summen", ZAMM, 1974, 54(1):39–51, doi:10.1002/zamm.19740540106 (the compensated sum)
// Tags:    logarithm, log-sum-exp, softmax, log-domain, probability, numerically-stable
// Units:   values natural logarithms of positive quantities; result the natural logarithm of their sum
// Example: LogSumExp(log 1, log 2, log 3, log 4) = log 10 = 2.302585092994046; LogSumExp(1000, 1000, 1000) = 1000 + log 3 (hand-checked)
func LogSumExp(values []float64) (float64, error) {
	if len(values) == 0 {
		return 0, fmt.Errorf("%w: values", ErrEmpty)
	}
	k := 0
	for i, v := range values {
		if math.IsNaN(v) || math.IsInf(v, 1) {
			return 0, fmt.Errorf("%w: values[%d] = %v", ErrNotFinite, i, v)
		}
		if v > values[k] {
			k = i
		}
	}
	a := values[k]
	if math.IsInf(a, -1) {
		return a, nil
	}
	var sum, comp float64
	for i, v := range values {
		if i == k {
			continue
		}
		x := math.Exp(v - a) // in [0, 1]
		sum, comp = reduce.NeumaierAdd(sum, comp, x)
	}
	return a + math.Log1p(sum+comp), nil
}
