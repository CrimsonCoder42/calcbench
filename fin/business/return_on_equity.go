package business

// ReturnOnEquity returns net income divided by average stockholders' equity, as a decimal. It is not
// rounded. A net loss or negative equity is allowed; the sign of the result follows the inputs.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrZeroDenominator when averageEquityCents is 0.
//
// ID:      FIN-BIZ-031
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Franklin, Graybeal and Cooper, Principles of Accounting, Vol. 1: Financial Accounting, OpenStax (Rice University), 2019, Appendix A "Financial Statement Analysis", https://openstax.org/books/principles-financial-accounting/pages/a-financial-statement-analysis (read 2026-11-12)
// Tags:    ratio, profitability, return, roe
// Units:   netIncomeCents and averageEquityCents in cents; the result is a decimal per period
// Example: Banyan Goods $35,000 ÷ $95,000 = 0.368, printed 0.37 (Principles of Accounting Vol. 1, Appendix A)
func ReturnOnEquity(netIncomeCents, averageEquityCents int64) (float64, error) {
	if err := checkCents("netIncomeCents", netIncomeCents); err != nil {
		return 0, err
	}
	if err := checkCents("averageEquityCents", averageEquityCents); err != nil {
		return 0, err
	}
	return ratioCents(netIncomeCents, "averageEquityCents", averageEquityCents)
}
