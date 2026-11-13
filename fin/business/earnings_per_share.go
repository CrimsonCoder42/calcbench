package business

import "fmt"

// checkResultCents returns ErrOverflow when a money result is outside ±2⁵³.
func checkResultCents(name string, v int64) error {
	if v > maxExactCents || v < -maxExactCents {
		return fmt.Errorf("%w: %s = %d, want |x| ≤ 2^53", ErrOverflow, name, v)
	}
	return nil
}

// checkShares returns ErrOutOfRange when shares > 2⁵³, ErrNegative when shares < 0 and ErrZeroDenominator
// when shares == 0.
func checkShares(name string, shares int64) error {
	if shares > maxExactCents {
		return fmt.Errorf("%w: %s = %d, want ≤ 2^53", ErrOutOfRange, name, shares)
	}
	if shares < 0 {
		return fmt.Errorf("%w: %s = %d", ErrNegative, name, shares)
	}
	if shares == 0 {
		return fmt.Errorf("%w: %s = 0", ErrZeroDenominator, name)
	}
	return nil
}

// divRoundHalfEven returns num/den rounded half-even, exactly, in integer arithmetic. den must be > 0.
func divRoundHalfEven(num, den int64) int64 {
	q, r := num/den, num%den // r has the sign of num, |r| < den
	if r == 0 {
		return q
	}
	step := int64(1)
	if r < 0 {
		r, step = -r, -1
	}
	twice := 2 * r // < 2·den ≤ 2⁶³ − 1 for den ≤ 2⁵³
	if twice > den || (twice == den && q%2 != 0) {
		q += step
	}
	return q
}

// EarningsPerShare returns (net income − preferred dividends) ÷ weighted average common shares outstanding,
// in cents per share, rounded half-even. The division is exact integer arithmetic (no float64), so the
// rounding is exact. A loss gives a negative EPS.
// Errors: ErrOutOfRange when a money input is outside ±2⁵³ or weightedAverageShares > 2⁵³; ErrNegative when
// preferredDividendsCents or weightedAverageShares is negative; ErrZeroDenominator when
// weightedAverageShares is 0; ErrOverflow when the result is outside ±2⁵³ cents.
//
// ID:      FIN-BIZ-021
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist and Knight, Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §6.5 "Market Value Ratios", https://openstax.org/books/principles-finance/pages/6-5-market-value-ratios (read 2026-11-12)
// Source:  Franklin, Graybeal and Cooper, Principles of Accounting, Vol. 1: Financial Accounting, OpenStax (Rice University), 2019, §14.5 "Discuss the Applicability of Earnings per Share as a Method to Measure Performance", https://openstax.org/books/principles-financial-accounting/pages/14-5-discuss-the-applicability-of-earnings-per-share-as-a-method-to-measure-performance (read 2026-11-12)
// Tags:    per-share, market-value, earnings-per-share, eps
// Units:   netIncomeCents and preferredDividendsCents in cents; weightedAverageShares a share count; the result in cents per share
// Example: Clear Lake Sporting Goods ($35,000 − $5,000) ÷ 8,000 = $3.75 (Principles of Finance §6.5); Sanaron ($50,000 − $10,000) ÷ 5,000 = $8.00 (Principles of Accounting Vol. 1 §14.5)
func EarningsPerShare(netIncomeCents, preferredDividendsCents, weightedAverageShares int64) (int64, error) {
	if err := checkCents("netIncomeCents", netIncomeCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("preferredDividendsCents", preferredDividendsCents); err != nil {
		return 0, err
	}
	if err := checkShares("weightedAverageShares", weightedAverageShares); err != nil {
		return 0, err
	}
	eps := divRoundHalfEven(netIncomeCents-preferredDividendsCents, weightedAverageShares)
	if err := checkResultCents("earnings per share", eps); err != nil {
		return 0, err
	}
	return eps, nil
}
