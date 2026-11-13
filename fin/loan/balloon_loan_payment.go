package loan

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

// BalloonLoanPayment returns the level end-of-period payment on a loan of principalCents over periods
// periods at ratePerPeriod that leaves balloonCents owed after the last level payment (the balloon, paid
// with it). The loan equals the present value of the payments plus the present value of the balloon,
// P = A·(P/A,i,n) + B·(P/F,i,n) (the bond-price identity with the payment as coupon and the balloon as
// par), so A = (P − B·(1+i)⁻ⁿ)·(A/P,i,n), rounded half-even once. balloonCents = 0 gives the fully
// amortizing payment; balloonCents = principalCents gives the interest-only payment P·i. When B is close
// to P·(1+i)ⁿ the subtraction cancels and the payment carries the absolute error of that difference
// (well under a cent for amounts below 2⁵³ cents).
// The result is the float64 value rounded half-even; it equals the correctly rounded exact value except
// when that value lies within about n·ε·|x| of a half cent, and is within 1 cent of it for |result| ≤ 2⁴³
// cents (n = the few floating-point operations here, ε = 2⁻⁵³, x = the exact value in cents).
// Errors: ErrOutOfRange if |principalCents| or |balloonCents| > 2⁵³; ErrNegative if either is negative;
// ErrRate if ratePerPeriod is NaN, ±Inf or ≤ −1; ErrPeriods if periods < 1; ErrDomain if balloonCents
// exceeds principalCents·(1+i)ⁿ (the payment would be negative); ErrOverflow if the payment or a factor
// exceeds its range.
//
// ID:      FIN-LOAN-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist & Knight, Principles of Finance, OpenStax (Rice University), 2022, §10.2 "Bond Valuation", eq. 10.9 (price = par × 1/(1+r)ⁿ + coupon × (1 − 1/(1+r)ⁿ)/r), https://openstax.org/books/principles-finance/pages/10-2-bond-valuation (read 2026-11-12)
// Source:  Consumer Financial Protection Bureau, "What is a balloon payment? When is one allowed?", https://www.consumerfinance.gov/ask-cfpb/what-is-a-balloon-payment-when-is-one-allowed-en-104/ (read 2026-11-12)
// Tags:    loan, balloon, payment, mortgage, partially-amortizing, annuity
// Units:   principalCents, balloonCents and the result in cents; ratePerPeriod decimal per payment period; periods count
// Example: hand-checked: $100,000, balloon $50,000, 0.5 %/month, 60 payments → $1,216.64
func BalloonLoanPayment(principalCents, balloonCents int64, ratePerPeriod float64, periods int) (int64, error) {
	if err := checkNonNegativeCents("principalCents", principalCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("balloonCents", balloonCents); err != nil {
		return 0, err
	}
	if err := checkRate("ratePerPeriod", ratePerPeriod); err != nil {
		return 0, err
	}
	if err := checkPeriods("periods", periods, 1, math.MaxInt); err != nil {
		return 0, err
	}
	discount, err := interest.SinglePaymentPresentWorth(ratePerPeriod, periods)
	if err != nil {
		return 0, fromInterest(err)
	}
	recovery, err := interest.CapitalRecovery(ratePerPeriod, periods)
	if err != nil {
		return 0, fromInterest(err)
	}
	payment := math.RoundToEven((float64(principalCents) - float64(balloonCents)*discount) * recovery)
	if payment < 0 {
		return 0, fmt.Errorf("%w: balloonCents = %d exceeds principalCents·(1+i)^n", ErrDomain, balloonCents)
	}
	return roundCents("payment", payment)
}
