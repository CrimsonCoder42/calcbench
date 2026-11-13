package loan

import "fmt"

// LoanToValueRatio returns the loan-to-value ratio, loanCents / propertyValueCents, as an unrounded
// decimal (0.9 = 90 %). The property value is the caller's (for a purchase Fannie Mae uses the lower of
// the sales price and the appraised value); lender delivery rounding (Fannie Mae truncates to two
// decimals and rounds up to a whole percent) is not applied.
// Errors: ErrOutOfRange if |loanCents| or |propertyValueCents| > 2⁵³; ErrNegative if either is negative;
// ErrZeroDenominator if propertyValueCents is 0.
//
// ID:      FIN-LOAN-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Fannie Mae, Selling Guide, B2-1.2-01 "Loan-to-Value (LTV) Ratios" (06/01/2022) ("Divide the original loan amount by the property value."), https://selling-guide.fanniemae.com/sel/b2-1.2-01/loan-value-ltv-ratios (read 2026-11-12)
// Source:  U.S. Department of Housing and Urban Development, Handbook 4000.1 Glossary and Acronyms (last revised 8/12/2026), "Loan-to-Value (LTV)", https://www.hud.gov/sites/default/files/Housing/documents/40001-hsgh-Update-18-Glossary.pdf (read 2026-11-12)
// Tags:    loan, mortgage, ratio, ltv, loan-to-value, underwriting
// Units:   loanCents and propertyValueCents in cents; result decimal
// Example: hand-checked: $360,000 loan on a $400,000 property → 0.9
func LoanToValueRatio(loanCents, propertyValueCents int64) (float64, error) {
	if err := checkNonNegativeCents("loanCents", loanCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("propertyValueCents", propertyValueCents); err != nil {
		return 0, err
	}
	if propertyValueCents == 0 {
		return 0, fmt.Errorf("%w: propertyValueCents = 0", ErrZeroDenominator)
	}
	return float64(loanCents) / float64(propertyValueCents), nil
}
