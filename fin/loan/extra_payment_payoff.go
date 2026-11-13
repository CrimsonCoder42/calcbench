package loan

import "fmt"

// addInterest adds in to total and returns ErrOverflow when the total leaves ±2⁵³ cents.
func addInterest(total, in int64) (int64, error) {
	total += in
	return total, checkResultCents("total interest", total)
}

// ExtraPaymentPayoff returns how many periods it takes to repay a fully amortizing loan when
// extraPerPeriodCents is paid on top of the level payment every period, and the interest saved compared
// with AmortizationSchedule(principalCents, ratePerPeriod, periods). Each period accrues the opening
// balance × ratePerPeriod rounded half-even (per-row rounding is the definition of the result); the loan
// is paid off in the first period whose balance plus interest does not exceed the level payment plus the
// extra (that period pays just the balance plus interest), and at the latest in period periods.
// extraPerPeriodCents = 0 gives (the schedule's row count, 0). principalCents = 0 gives (0, 0).
// Errors: ErrOutOfRange if |principalCents| or |extraPerPeriodCents| > 2⁵³; ErrNegative if either is
// negative; ErrRate if ratePerPeriod is NaN, ±Inf or ≤ −1; ErrPeriods if periods is outside 1…10⁷;
// ErrOverflow if the payment plus the extra, an interest amount, or a total of interest exceeds ±2⁵³
// cents.
//
// ID:      FIN-LOAN-005
// Time:    O(n)   n = periods ≤ 10⁷ (one pass for the scheduled loan, one for the prepaid loan)
// Space:   O(1), no allocation
// Source:  Dahlquist & Knight, Principles of Finance, OpenStax (Rice University), 2022, §8.3 "Loan Amortization" (extra payment applied to principal pays the 20-year mortgage off early), https://openstax.org/books/principles-finance/pages/8-3-loan-amortization (read 2026-11-12)
// Tags:    loan, prepayment, extra-payment, payoff, mortgage, amortization, interest-saved
// Units:   principalCents, extraPerPeriodCents and interestSavedCents in cents; ratePerPeriod decimal per payment period; periods and periodsToPayoff counts
// Example: hand-checked: $140,000 at 0.3 %/month over 240 (OpenStax §8.3 mortgage) with $100 extra a month → 204 periods, $9,257.23 saved
func ExtraPaymentPayoff(principalCents int64, ratePerPeriod float64, periods int, extraPerPeriodCents int64) (periodsToPayoff int, interestSavedCents int64, err error) {
	payment, err := loanPayment(principalCents, ratePerPeriod, periods, maxSchedulePeriods)
	if err != nil {
		return 0, 0, err
	}
	if err := checkNonNegativeCents("extraPerPeriodCents", extraPerPeriodCents); err != nil {
		return 0, 0, err
	}
	if principalCents == 0 {
		return 0, 0, nil
	}
	paid := payment + extraPerPeriodCents
	if err := checkResultCents("payment + extra", paid); err != nil {
		return 0, 0, err
	}

	var scheduled int64
	balance := principalCents
	for k := 1; k <= periods; k++ {
		row, payoff, err := amortizeRow(balance, payment, ratePerPeriod, k == periods)
		if err != nil {
			return 0, 0, fmt.Errorf("period %d: %w", k, err)
		}
		if scheduled, err = addInterest(scheduled, row.InterestCents); err != nil {
			return 0, 0, err
		}
		if payoff {
			break
		}
		balance = row.BalanceCents
	}

	var prepaid int64
	balance = principalCents
	for k := 1; k <= periods; k++ {
		in, err := periodInterest(balance, ratePerPeriod)
		if err != nil {
			return 0, 0, fmt.Errorf("period %d: %w", k, err)
		}
		if prepaid, err = addInterest(prepaid, in); err != nil {
			return 0, 0, err
		}
		if k == periods || balance+in <= paid {
			return k, scheduled - prepaid, nil
		}
		balance = balance + in - paid
	}
	return periods, scheduled - prepaid, nil // unreachable: the loop returns at k = periods
}
