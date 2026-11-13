package loan

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

// AffordableLoanAmount returns the largest loan whose total monthly housing payment fits a housing ratio:
// the principal-and-interest allowance is grossMonthlyIncomeCents × maxHousingRatio less
// monthlyTaxesAndInsuranceCents, and the loan is the present value of that allowance paid for periods
// periods at ratePerPeriod, allowance·(1 − (1+i)⁻ⁿ)/i, rounded half-even once. When taxes and insurance
// use up the allowance the result is 0 (no error). maxHousingRatio is the caller's limit (no lender's
// threshold is built in); a payment period of one month is assumed, since the ratio is monthly.
// The result is the float64 value rounded half-even; it equals the correctly rounded exact value except
// when that value lies within about n·ε·|x| of a half cent, and is within 1 cent of it for |result| ≤ 2⁴³
// cents (n = the few floating-point operations here, ε = 2⁻⁵³, x = the exact value in cents).
// Errors: ErrOutOfRange if a money input is outside ±2⁵³ cents or maxHousingRatio is outside 0…1;
// ErrNegative if a money input is negative; ErrNotFinite if maxHousingRatio is NaN or ±Inf; ErrRate if
// ratePerPeriod is NaN, ±Inf or ≤ −1; ErrPeriods if periods < 1; ErrOverflow if the loan exceeds 2⁵³
// cents.
//
// ID:      FIN-LOAN-010
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  U.S. Department of Housing and Urban Development, Handbook 4000.1 Glossary and Acronyms (last revised 8/12/2026), "PTI - Total Mortgage Payment to Effective Income Ratio", https://www.hud.gov/sites/default/files/Housing/documents/40001-hsgh-Update-18-Glossary.pdf (read 2026-11-12)
// Source:  Consumer Financial Protection Bureau, "On a mortgage, what's the difference between my principal and interest payment and my total monthly payment?", https://www.consumerfinance.gov/ask-cfpb/what-is-a-piti-en-1941/ (read 2026-11-12); Dahlquist & Knight, Principles of Finance, OpenStax, 2022, §8.2 eq. 8.14, https://openstax.org/books/principles-finance/pages/8-2-annuities (read 2026-11-12)
// Tags:    loan, mortgage, affordability, housing-ratio, pti, qualification, present-value
// Units:   grossMonthlyIncomeCents, monthlyTaxesAndInsuranceCents and the result in cents; maxHousingRatio decimal; ratePerPeriod decimal per month; periods count of months
// Example: hand-checked: income $6,000, ratio 0.28, taxes and insurance $400 → allowance $1,280; at 0.5 %/month over 360 → $213,493.27
func AffordableLoanAmount(grossMonthlyIncomeCents int64, maxHousingRatio float64, monthlyTaxesAndInsuranceCents int64, ratePerPeriod float64, periods int) (int64, error) {
	if err := checkNonNegativeCents("grossMonthlyIncomeCents", grossMonthlyIncomeCents); err != nil {
		return 0, err
	}
	if math.IsNaN(maxHousingRatio) || math.IsInf(maxHousingRatio, 0) {
		return 0, fmt.Errorf("%w: maxHousingRatio = %v", ErrNotFinite, maxHousingRatio)
	}
	if maxHousingRatio < 0 || maxHousingRatio > 1 {
		return 0, fmt.Errorf("%w: maxHousingRatio = %v, want 0…1", ErrOutOfRange, maxHousingRatio)
	}
	if err := checkNonNegativeCents("monthlyTaxesAndInsuranceCents", monthlyTaxesAndInsuranceCents); err != nil {
		return 0, err
	}
	if err := checkRate("ratePerPeriod", ratePerPeriod); err != nil {
		return 0, err
	}
	if err := checkPeriods("periods", periods, 1, math.MaxInt); err != nil {
		return 0, err
	}
	allowance := float64(grossMonthlyIncomeCents)*maxHousingRatio - float64(monthlyTaxesAndInsuranceCents)
	if allowance <= 0 {
		return 0, nil
	}
	factor, err := interest.UniformSeriesPresentWorth(ratePerPeriod, periods)
	if err != nil {
		return 0, fromInterest(err)
	}
	return roundCents("loan", allowance*factor)
}
