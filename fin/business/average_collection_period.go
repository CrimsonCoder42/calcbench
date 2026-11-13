package business

import "fmt"

// AverageCollectionPeriod returns the number of days, on average, that receivables take to collect:
// daysInPeriod × average receivables ÷ net credit sales, which equals daysInPeriod ÷ receivables turnover
// (FIN-BIZ-014). The sources use a 365-day year; the caller passes the period length. At most two roundings
// (the product, then the quotient); the result is not rounded to whole days.
// Errors: ErrOutOfRange when a money input is outside ±2⁵³; ErrNegative when a money input is negative;
// ErrPeriods when daysInPeriod < 1; ErrZeroDenominator when netCreditSalesCents is 0.
//
// ID:      FIN-BIZ-015
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Franklin, Graybeal and Cooper, Principles of Accounting, Vol. 1: Financial Accounting, OpenStax (Rice University), 2019, §9.3 "Determine the Efficiency of Receivables Management Using Financial Ratios" (number of days' sales in receivables = 365 ÷ accounts receivable turnover), https://openstax.org/books/principles-financial-accounting/pages/9-3-determine-the-efficiency-of-receivables-management-using-financial-ratios (read 2026-11-12)
// Source:  Dahlquist and Knight, Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §19.1 "What Is Working Capital?" (receivables collection period = 365 ÷ accounts receivable turnover), https://openstax.org/books/principles-finance/pages/19-1-what-is-working-capital (read 2026-11-12)
// Tags:    ratio, efficiency, activity, receivables, collection-period, days-sales-outstanding
// Units:   averageReceivablesCents and netCreditSalesCents in cents; daysInPeriod in days; the result in days
// Example: BWW 2017: average receivables $80,000, credit sales $400,000, 365 days → 73 days (Principles of Accounting Vol. 1 §9.3)
func AverageCollectionPeriod(averageReceivablesCents, netCreditSalesCents int64, daysInPeriod int) (float64, error) {
	if err := checkNonNegativeCents("averageReceivablesCents", averageReceivablesCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("netCreditSalesCents", netCreditSalesCents); err != nil {
		return 0, err
	}
	if daysInPeriod < 1 {
		return 0, fmt.Errorf("%w: daysInPeriod = %d, want ≥ 1", ErrPeriods, daysInPeriod)
	}
	if netCreditSalesCents == 0 {
		return 0, fmt.Errorf("%w: netCreditSalesCents = 0", ErrZeroDenominator)
	}
	return float64(daysInPeriod) * float64(averageReceivablesCents) / float64(netCreditSalesCents), nil
}
