package business

// ReceivablesTurnover returns net credit sales divided by average accounts receivable, as a decimal. It is
// not rounded. The source's average is (beginning + ending) ÷ 2; the caller computes it.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when netCreditSalesCents or
// averageReceivablesCents is negative; ErrZeroDenominator when averageReceivablesCents is 0.
//
// ID:      FIN-BIZ-014
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Franklin, Graybeal and Cooper, Principles of Accounting, Vol. 1: Financial Accounting, OpenStax (Rice University), 2019, §9.3 "Determine the Efficiency of Receivables Management Using Financial Ratios", https://openstax.org/books/principles-financial-accounting/pages/9-3-determine-the-efficiency-of-receivables-management-using-financial-ratios (read 2026-11-12)
// Source:  Dahlquist and Knight, Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §6.2 "Operating Efficiency Ratios", https://openstax.org/books/principles-finance/pages/6-2-operating-efficiency-ratios (read 2026-11-12)
// Tags:    ratio, efficiency, activity, receivables-turnover
// Units:   netCreditSalesCents and averageReceivablesCents in cents; the result is a dimensionless count of times per period
// Example: BWW 2017 $400,000 ÷ $80,000 = 5 (Principles of Accounting Vol. 1 §9.3)
func ReceivablesTurnover(netCreditSalesCents, averageReceivablesCents int64) (float64, error) {
	if err := checkNonNegativeCents("netCreditSalesCents", netCreditSalesCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("averageReceivablesCents", averageReceivablesCents); err != nil {
		return 0, err
	}
	return ratioCents(netCreditSalesCents, "averageReceivablesCents", averageReceivablesCents)
}
