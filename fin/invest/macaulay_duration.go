package invest

import "fmt"

// maxPeriods caps the period count of the O(n) bond sums (MacaulayDuration, ModifiedDuration,
// Convexity) so one call cannot run unbounded; 10 000 000 periods is far beyond any bond.
const maxPeriods = 10_000_000

// bondMoments returns, for the bond's flows cₖ (c each period, plus faceCents at k = n) and
// v = 1/(1+i): p = Σ cₖvᵏ (the price), s1 = Σ k·cₖvᵏ and s2 = Σ k(k+1)·cₖvᵏ. It evaluates the
// polynomial p(v) and its first two derivatives by one Horner pass (s1 = v·p′, s2 = v²·p″ + 2v·p′).
// The error grows at most linearly in n. Inputs are already checked.
func bondMoments(faceCents int64, i, c float64, periods int) (p, s1, s2 float64) {
	v := 1 / (1 + i)
	// q(v) = Σ_{k=1..n} cₖ v^(k−1); p = v·q. Horner from k = n down to 1.
	var q, dq, ddq float64
	for k := periods; k >= 1; k-- {
		ck := c
		if k == periods {
			ck += float64(faceCents)
		}
		ddq = ddq*v + 2*dq
		dq = dq*v + q
		q = q*v + ck
	}
	// p = v·q, p′ = q + v·q′, p″ = 2q′ + v·q″.
	p = v * q
	dp := q + v*dq
	ddp := 2*dq + v*ddq
	s1 = v * dp
	s2 = v*v*ddp + 2*v*dp
	return p, s1, s2
}

// bondDurationParts checks the bond, computes its moments and returns them with the periodic yield.
// A zero price (face and coupon both 0) is ErrDomain; a non-finite moment is ErrOverflow.
func bondDurationParts(faceCents int64, couponRatePerYear, yieldPerYear float64, periods, periodsPerYear int) (i, p, s1, s2 float64, err error) {
	i, c, err := checkBond(faceCents, couponRatePerYear, yieldPerYear, periods, periodsPerYear)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	if periods > maxPeriods {
		return 0, 0, 0, 0, fmt.Errorf("%w: periods = %d, want ≤ %d", ErrPeriods, periods, maxPeriods)
	}
	p, s1, s2 = bondMoments(faceCents, i, c, periods)
	if p == 0 {
		return 0, 0, 0, 0, fmt.Errorf("%w: the bond has no cash flows (faceCents and coupon are 0)", ErrDomain)
	}
	for _, x := range [...]float64{p, s1, s2} {
		if _, err := finiteResult(x); err != nil {
			return 0, 0, 0, 0, err
		}
	}
	return i, p, s1, s2, nil
}

// MacaulayDuration returns the Macaulay duration of a level-coupon bond on a coupon date, in years:
// the present-value-weighted average time of its flows, D = Σ k·cₖvᵏ / Σ cₖvᵏ periods, divided by
// periodsPerYear. cₖ is the coupon faceCents·couponRatePerYear/periodsPerYear (plus the face at
// k = periods), v = 1/(1 + yieldPerYear/periodsPerYear). The closed form (Macaulay p. 49) is O(1) but
// divides by R − 1 and cancels as the yield goes to 0; this sums the flows once instead.
// Errors: ErrOutOfRange if |faceCents| > 2⁵³; ErrNegative if faceCents or couponRatePerYear < 0;
// ErrRate if a rate is NaN or ±Inf or the periodic yield is ≤ −1; ErrPeriods if periods < 1, periods >
// 10 000 000 or periodsPerYear < 1; ErrDomain if the bond has no flows; ErrOverflow if a sum exceeds float64.
//
// ID:      FIN-INVEST-004
// Time:    O(n)   n = periods
// Space:   O(1), no allocation
// Source:  Macaulay, Some Theoretical Problems Suggested by the Movements of Interest Rates, Bond Yields and Stock Prices in the United States since 1856, NBER, 1938, ch. 2, pp. 48–49 (D), p. 51 (table "Duration in years"), https://www.nber.org/system/files/chapters/c6342/c6342.pdf
// Source:  Society of Actuaries, "Notation and terminology used for Exam FM" (Macaulay duration), https://www.soa.org/globalassets/assets/Files/Edu/2019/exam-fm-notation-terminology2.pdf
// Tags:    invest, bond, duration, macaulay, interest-rate-risk
// Units:   faceCents in cents; couponRatePerYear and yieldPerYear decimal per year (nominal); periods count; result years
// Example: Macaulay 1938 p. 51, semiannual bonds at par: 4 % 1 year .990, 4 % 10 years 8.339, 5 % 10 years 7.989, 6 % 10 years 7.662
func MacaulayDuration(faceCents int64, couponRatePerYear, yieldPerYear float64, periods, periodsPerYear int) (float64, error) {
	_, p, s1, _, err := bondDurationParts(faceCents, couponRatePerYear, yieldPerYear, periods, periodsPerYear)
	if err != nil {
		return 0, err
	}
	return finiteResult(s1 / p / float64(periodsPerYear))
}
