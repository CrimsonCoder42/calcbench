package invest

// ModifiedDuration returns the modified duration of a level-coupon bond on a coupon date, in years:
// −(1/P)·dP/dy with the price P written as a function of the nominal annual yield y, which equals the
// Macaulay duration in years divided by (1 + y/periodsPerYear). It is the first-order relative price
// change per unit change in yieldPerYear.
// Errors: ErrOutOfRange if |faceCents| > 2⁵³; ErrNegative if faceCents or couponRatePerYear < 0;
// ErrRate if a rate is NaN or ±Inf or the periodic yield is ≤ −1; ErrPeriods if periods < 1, periods >
// 10 000 000 or periodsPerYear < 1; ErrDomain if the bond has no flows; ErrOverflow if a sum exceeds float64.
//
// ID:      FIN-INVEST-005
// Time:    O(n)   n = periods
// Space:   O(1), no allocation
// Source:  Society of Actuaries, "Notation and terminology used for Exam FM" ("Modified duration = −P′(i)/P(i)"), https://www.soa.org/globalassets/assets/Files/Edu/2019/exam-fm-notation-terminology2.pdf (read 2026-11-12)
// Tags:    invest, bond, duration, modified-duration, interest-rate-risk, sensitivity
// Units:   faceCents in cents; couponRatePerYear and yieldPerYear decimal per year (nominal); periods count; result years
// Example: hand-checked from Macaulay 1938 p. 51: 6 % semiannual 10-year bond at par, 7.661900 / 1.03 = 7.438737
func ModifiedDuration(faceCents int64, couponRatePerYear, yieldPerYear float64, periods, periodsPerYear int) (float64, error) {
	i, p, s1, _, err := bondDurationParts(faceCents, couponRatePerYear, yieldPerYear, periods, periodsPerYear)
	if err != nil {
		return 0, err
	}
	return finiteResult(s1 / p / float64(periodsPerYear) / (1 + i))
}
