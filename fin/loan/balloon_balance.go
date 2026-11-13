package loan

import "fmt"

// BalloonBalance returns the balance owed after paymentsMade end-of-period payments of paymentCents on a
// loan of principalCents at ratePerPeriod: each period adds the opening balance × ratePerPeriod rounded
// half-even and subtracts the payment, as in AmortizationSchedule but with no last-row adjustment; this
// per-row half-even rounding is the definition of the result. When the payments stop before the loan is
// repaid this is the balloon payment due. When a payment covers the balance plus its interest the loan
// is paid off and the result is 0 from then on (never negative), as in AmortizationSchedule. A payment
// below the interest makes the balance grow (negative amortization).
// Errors: ErrOutOfRange if |principalCents| or |paymentCents| > 2⁵³; ErrNegative if principalCents or
// paymentCents < 0; ErrRate if ratePerPeriod is NaN, ±Inf or ≤ −1; ErrPeriods if paymentsMade is outside
// 0…10⁷; ErrOverflow if the balance or an interest amount leaves ±2⁵³ cents.
//
// ID:      FIN-LOAN-003
// Time:    O(k)   k = paymentsMade ≤ 10⁷
// Space:   O(1), no allocation
// Source:  Consumer Financial Protection Bureau, "What is a balloon payment? When is one allowed?" ("a large, one-time payment at the end of the loan term"), https://www.consumerfinance.gov/ask-cfpb/what-is-a-balloon-payment-when-is-one-allowed-en-104/ (read 2026-11-12)
// Source:  Dahlquist & Knight, Principles of Finance, OpenStax (Rice University), 2022, §8.3 "Loan Amortization", Table 8.4 (the balance recursion), https://openstax.org/books/principles-finance/pages/8-3-loan-amortization (read 2026-11-12)
// Tags:    loan, balloon, balance, mortgage, amortization, partially-amortizing
// Units:   principalCents, paymentCents and the result in cents; ratePerPeriod decimal per payment period; paymentsMade count
// Example: hand-checked: $100,000 at 0.5 %/month, payment $599.55 (30-year level payment), after 60 payments → $93,054.37
func BalloonBalance(principalCents, paymentCents int64, ratePerPeriod float64, paymentsMade int) (int64, error) {
	if err := checkNonNegativeCents("principalCents", principalCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("paymentCents", paymentCents); err != nil {
		return 0, err
	}
	if err := checkRate("ratePerPeriod", ratePerPeriod); err != nil {
		return 0, err
	}
	if err := checkPeriods("paymentsMade", paymentsMade, 0, maxSchedulePeriods); err != nil {
		return 0, err
	}
	balance := principalCents
	for k := 1; k <= paymentsMade && balance != 0; k++ {
		in, err := periodInterest(balance, ratePerPeriod)
		if err != nil {
			return 0, fmt.Errorf("period %d: %w", k, err)
		}
		if balance+in <= paymentCents {
			return 0, nil
		}
		balance = balance + in - paymentCents
		if err := checkResultCents("balance", balance); err != nil {
			return 0, fmt.Errorf("period %d: %w", k, err)
		}
	}
	return balance, nil
}
