package loan

import "fmt"

// DebtToIncomeRatio returns the debt-to-income ratio, all monthly debt payments divided by gross monthly
// income, as an unrounded decimal (0.33 = 33 %). Which debts count is the lender's rule, not this
// function's.
// Errors: ErrOutOfRange if either input is outside ±2⁵³ cents; ErrNegative if either is negative;
// ErrZeroDenominator if grossMonthlyIncomeCents is 0.
//
// ID:      FIN-LOAN-009
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Consumer Financial Protection Bureau, "What is a debt-to-income ratio?" ("all your monthly debt payments divided by your gross monthly income"), https://www.consumerfinance.gov/ask-cfpb/what-is-a-debt-to-income-ratio-en-1791/ (read 2026-11-12)
// Source:  Fannie Mae, Selling Guide, B3-6-02 "Debt-to-Income Ratios" (04/02/2025), https://selling-guide.fanniemae.com/sel/b3-6-02/debt-income-ratios (read 2026-11-12)
// Tags:    loan, mortgage, ratio, dti, debt-to-income, underwriting
// Units:   both inputs in cents per month; result decimal
// Example: CFPB: debts $1,500 + $100 + $400 = $2,000 on gross income $6,000 → 0.333… ("33 percent")
func DebtToIncomeRatio(monthlyDebtPaymentsCents, grossMonthlyIncomeCents int64) (float64, error) {
	if err := checkNonNegativeCents("monthlyDebtPaymentsCents", monthlyDebtPaymentsCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("grossMonthlyIncomeCents", grossMonthlyIncomeCents); err != nil {
		return 0, err
	}
	if grossMonthlyIncomeCents == 0 {
		return 0, fmt.Errorf("%w: grossMonthlyIncomeCents = 0", ErrZeroDenominator)
	}
	return float64(monthlyDebtPaymentsCents) / float64(grossMonthlyIncomeCents), nil
}
