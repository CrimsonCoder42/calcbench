package loan

import "fmt"

// RemainingBalance returns the balance of a fully amortizing loan after paymentsMade of its periods level
// payments, equal to AmortizationSchedule(principalCents, ratePerPeriod, periods)[paymentsMade−1]
// .BalanceCents exactly (principalCents when paymentsMade = 0; 0 when paymentsMade = periods, and 0 for
// every paymentsMade at or after an early payoff row, where the schedule has ended). It walks the same
// rows with the same per-row half-even rounding, which is the definition of the result; the closed form
// would differ from the schedule by cents.
// Errors: ErrOutOfRange if |principalCents| > 2⁵³; ErrNegative if principalCents < 0; ErrRate if
// ratePerPeriod is NaN, ±Inf or ≤ −1; ErrPeriods if periods is outside 1…10⁷ or paymentsMade is outside
// 0…periods; ErrOverflow if a payment or interest amount exceeds ±2⁵³ cents.
//
// ID:      FIN-LOAN-002
// Time:    O(k)   k = paymentsMade ≤ 10⁷
// Space:   O(1), no allocation
// Source:  Dahlquist & Knight, Principles of Finance, OpenStax (Rice University), 2022, §8.3 "Loan Amortization", Table 8.4 (balance column), https://openstax.org/books/principles-finance/pages/8-3-loan-amortization (read 2026-11-12)
// Tags:    loan, amortization, balance, remaining-balance, payoff, mortgage
// Units:   principalCents and the result in cents; ratePerPeriod decimal per payment period; periods and paymentsMade counts
// Example: OpenStax §8.3 Table 8.4: $32,000, 0.5 %/month, 36 payments, after 3 payments → $29,547.27
func RemainingBalance(principalCents int64, ratePerPeriod float64, periods, paymentsMade int) (int64, error) {
	payment, err := loanPayment(principalCents, ratePerPeriod, periods, maxSchedulePeriods)
	if err != nil {
		return 0, err
	}
	if err := checkPeriods("paymentsMade", paymentsMade, 0, periods); err != nil {
		return 0, err
	}
	balance := principalCents
	for k := 1; k <= paymentsMade && balance != 0; k++ {
		row, payoff, err := amortizeRow(balance, payment, ratePerPeriod, k == periods)
		if err != nil {
			return 0, fmt.Errorf("period %d: %w", k, err)
		}
		balance = row.BalanceCents
		if payoff {
			return 0, nil
		}
	}
	return balance, nil
}
