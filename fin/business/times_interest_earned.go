package business

// TimesInterestEarned returns earnings before interest and taxes divided by interest expense (interest
// coverage), as a decimal. It is not rounded. A negative EBIT is allowed and gives a negative ratio.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when interestExpenseCents is negative;
// ErrZeroDenominator when interestExpenseCents is 0.
//
// ID:      FIN-BIZ-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist and Knight, Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §6.4 "Solvency Ratios", https://openstax.org/books/principles-finance/pages/6-4-solvency-ratios (read 2026-11-12)
// Source:  Franklin, Graybeal and Cooper, Principles of Accounting, Vol. 1: Financial Accounting, OpenStax (Rice University), 2019, Appendix A "Financial Statement Analysis", https://openstax.org/books/principles-financial-accounting/pages/a-financial-statement-analysis (read 2026-11-12)
// Tags:    ratio, solvency, coverage, times-interest-earned, interest-coverage
// Units:   ebitCents and interestExpenseCents in cents; the result is a dimensionless count of times
// Example: Clear Lake Sporting Goods $43,000 ÷ $2,000 = 21.5 (Principles of Finance §6.4)
func TimesInterestEarned(ebitCents, interestExpenseCents int64) (float64, error) {
	if err := checkCents("ebitCents", ebitCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("interestExpenseCents", interestExpenseCents); err != nil {
		return 0, err
	}
	return ratioCents(ebitCents, "interestExpenseCents", interestExpenseCents)
}
