package business

// LoanToDepositRatio returns net loans and leases divided by total deposits, as a decimal (the UBPR prints
// it as a percent). It is not rounded.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when totalLoansCents or
// totalDepositsCents is negative; ErrZeroDenominator when totalDepositsCents is 0.
//
// ID:      FIN-BIZ-050
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Federal Financial Institutions Examination Council, UBPR User's Guide, "Summary Ratios — Page 1", 2026-08-20, item 32 "Net Loans & Leases to Deposits" (UBPRE600), https://www.cdr.ffiec.gov/CDRDownload/CDR/UserGuide/v186/FFIEC%20UBPR%20User%20Guide%20Summary%20Ratios--Page%201_2026-08-20.PDF (read 2026-11-12)
// Tags:    ratio, bank, liquidity, loan-to-deposit, ubpr
// Units:   totalLoansCents and totalDepositsCents in cents; the result is a dimensionless decimal
// Example: $800,000 ÷ $1,000,000 = 0.8 (hand-checked; the UBPR guide gives no numbers)
func LoanToDepositRatio(totalLoansCents, totalDepositsCents int64) (float64, error) {
	if err := checkNonNegativeCents("totalLoansCents", totalLoansCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("totalDepositsCents", totalDepositsCents); err != nil {
		return 0, err
	}
	return ratioCents(totalLoansCents, "totalDepositsCents", totalDepositsCents)
}
