package business

// BankEfficiencyRatio returns total overhead (noninterest) expense divided by the sum of net interest
// income and noninterest income, as a decimal (the UBPR prints it as a percent: 0.6 here is 60 there).
// The UBPR uses net interest income on a tax-equivalent basis; pass that figure to match it. Either income
// may be negative; the sum is exact in int64 and adds one rounding when its magnitude exceeds 2⁵³.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when noninterestExpenseCents is
// negative; ErrZeroDenominator when the two incomes sum to 0.
//
// ID:      FIN-BIZ-049
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Federal Financial Institutions Examination Council, UBPR User's Guide, "Noninterest Income and Expenses — Page 3", 2026-08-20, item 32 "Efficiency Ratio" (UBPRE088: total overhead expense as a percentage of net interest income (TE) plus noninterest income), https://www.cdr.ffiec.gov/CDRDownload/CDR/UserGuide/v186/FFIEC%20UBPR%20User%20Guide%20Noninterest%20Income%20and%20Expenses--Page%203_2026-08-20.PDF (read 2026-11-12)
// Tags:    ratio, bank, efficiency-ratio, overhead, ubpr
// Units:   every input in cents; the result is a dimensionless decimal
// Example: $60,000 ÷ ($80,000 + $20,000) = 0.6 (hand-checked; the UBPR guide gives no numbers)
func BankEfficiencyRatio(noninterestExpenseCents, netInterestIncomeCents, noninterestIncomeCents int64) (float64, error) {
	if err := checkNonNegativeCents("noninterestExpenseCents", noninterestExpenseCents); err != nil {
		return 0, err
	}
	if err := checkCents("netInterestIncomeCents", netInterestIncomeCents); err != nil {
		return 0, err
	}
	if err := checkCents("noninterestIncomeCents", noninterestIncomeCents); err != nil {
		return 0, err
	}
	return ratioCents(noninterestExpenseCents, "netInterestIncomeCents + noninterestIncomeCents", netInterestIncomeCents+noninterestIncomeCents)
}
