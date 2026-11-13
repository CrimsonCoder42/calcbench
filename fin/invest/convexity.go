package invest

// Convexity returns the (modified) convexity of a level-coupon bond on a coupon date, in years²:
// (1/P)·d²P/dy² with the price P written as a function of the nominal annual yield y,
// Σ k(k+1)·cₖvᵏ / (P·m²·(1+i)²), m = periodsPerYear, i = y/m, v = 1/(1+i).
// Errors: ErrOutOfRange if |faceCents| > 2⁵³; ErrNegative if faceCents or couponRatePerYear < 0;
// ErrRate if a rate is NaN or ±Inf or the periodic yield is ≤ −1; ErrPeriods if periods < 1, periods >
// 10 000 000 or periodsPerYear < 1; ErrDomain if the bond has no flows; ErrOverflow if a sum exceeds float64.
//
// ID:      FIN-INVEST-006
// Time:    O(n)   n = periods
// Space:   O(1), no allocation
// Source:  Society of Actuaries, "Notation and terminology used for Exam FM" ("Modified convexity = P″(i)/P(i)"), https://www.soa.org/globalassets/assets/Files/Edu/2019/exam-fm-notation-terminology2.pdf (read 2026-11-12)
// Tags:    invest, bond, convexity, interest-rate-risk, sensitivity
// Units:   faceCents in cents; couponRatePerYear and yieldPerYear decimal per year (nominal); periods count; result years²
// Example: hand-checked: 6 % semiannual 10-year bond at par → 68.774822 years²
func Convexity(faceCents int64, couponRatePerYear, yieldPerYear float64, periods, periodsPerYear int) (float64, error) {
	i, p, _, s2, err := bondDurationParts(faceCents, couponRatePerYear, yieldPerYear, periods, periodsPerYear)
	if err != nil {
		return 0, err
	}
	m := float64(periodsPerYear)
	return finiteResult(s2 / p / (m * m) / ((1 + i) * (1 + i)))
}
