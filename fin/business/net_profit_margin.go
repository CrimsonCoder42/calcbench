package business

// NetProfitMargin returns net income divided by net sales (profit margin), as a decimal. It is not rounded.
// A net loss is allowed and gives a negative margin.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when netSalesCents is negative;
// ErrZeroDenominator when netSalesCents is 0.
//
// ID:      FIN-BIZ-020
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist and Knight, Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §6.6 "Profitability Ratios and the DuPont Method", https://openstax.org/books/principles-finance/pages/6-6-profitability-ratios-and-the-dupont-method (read 2026-11-12)
// Source:  Franklin, Graybeal and Cooper, Principles of Accounting, Vol. 1: Financial Accounting, OpenStax (Rice University), 2019, Appendix A "Financial Statement Analysis", https://openstax.org/books/principles-financial-accounting/pages/a-financial-statement-analysis (read 2026-11-12)
// Tags:    ratio, profitability, margin, net-margin, dupont
// Units:   netIncomeCents and netSalesCents in cents; the result is a dimensionless decimal
// Example: Clear Lake Sporting Goods $35,000 ÷ $120,000 = 0.2917, printed 0.29 (Principles of Finance §6.6)
func NetProfitMargin(netIncomeCents, netSalesCents int64) (float64, error) {
	if err := checkCents("netIncomeCents", netIncomeCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("netSalesCents", netSalesCents); err != nil {
		return 0, err
	}
	return ratioCents(netIncomeCents, "netSalesCents", netSalesCents)
}
