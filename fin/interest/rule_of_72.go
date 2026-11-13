package interest

import "fmt"

// RuleOf72 returns the approximate number of periods for an amount to double at ratePerPeriod by the
// Rule of 72: 72 divided by the rate in percent, 72/(100·i). It is a rule of thumb; DoublingTime gives
// the exact ln 2/ln(1+i) (at 7 % the rule gives 10.29, the exact value 10.24).
// Errors: ErrRate if ratePerPeriod is NaN, ±Inf or ≤ −1; ErrDomain if ratePerPeriod ≤ 0.
//
// ID:      FIN-INT-033
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  J. Dahlquist and R. Knight, Principles of Finance, OpenStax, 2022, §7.4 "Applications of TVM in Finance", "The Rule of 72", Eq. 7.37: Years to double = 72/Interest rate, https://openstax.org/books/principles-finance/pages/7-4-applications-of-tvm-in-finance
// Source:  R. A. DeFusco, D. W. McLeavey, J. E. Pinto and D. E. Runkle, Quantitative Investment Analysis, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ch. 1 "The Time Value of Money", §7.2, footnote 9
// Tags:    interest, rule-of-thumb, rule-of-72, doubling-time, approximation
// Units:   ratePerPeriod decimal per period; result periods
// Example: 9 % → 8 years (Principles of Finance §7.4); 7 % → 10.3 (QIA ch. 1, footnote 9)
func RuleOf72(ratePerPeriod float64) (float64, error) {
	if err := checkRate("ratePerPeriod", ratePerPeriod); err != nil {
		return 0, err
	}
	if ratePerPeriod <= 0 {
		return 0, fmt.Errorf("%w: ratePerPeriod = %v, an amount doubles only at a rate > 0", ErrDomain, ratePerPeriod)
	}
	return factorResult(72 / (100 * ratePerPeriod))
}
