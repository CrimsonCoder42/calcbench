package invest

import (
	"fmt"
	"math"
)

// checkBond validates the inputs shared by BondPrice, MacaulayDuration, ModifiedDuration and
// Convexity and returns the periodic yield i = yieldPerYear/periodsPerYear and the coupon per period
// c = faceCents·couponRatePerYear/periodsPerYear (unrounded cents).
func checkBond(faceCents int64, couponRatePerYear, yieldPerYear float64, periods, periodsPerYear int) (i, c float64, err error) {
	if err := checkNonNegativeCents("faceCents", faceCents); err != nil {
		return 0, 0, err
	}
	if err := checkFinite("couponRatePerYear", couponRatePerYear); err != nil {
		return 0, 0, fmt.Errorf("%w: %w", ErrRate, err)
	}
	if couponRatePerYear < 0 {
		return 0, 0, fmt.Errorf("%w: couponRatePerYear = %v", ErrNegative, couponRatePerYear)
	}
	if periods < 1 {
		return 0, 0, fmt.Errorf("%w: periods = %d, want ≥ 1", ErrPeriods, periods)
	}
	if periodsPerYear < 1 {
		return 0, 0, fmt.Errorf("%w: periodsPerYear = %d, want ≥ 1", ErrPeriods, periodsPerYear)
	}
	if err := checkFinite("yieldPerYear", yieldPerYear); err != nil {
		return 0, 0, fmt.Errorf("%w: %w", ErrRate, err)
	}
	m := float64(periodsPerYear)
	i = yieldPerYear / m
	if err := checkRate("yieldPerYear/periodsPerYear", i); err != nil {
		return 0, 0, err
	}
	return i, float64(faceCents) * couponRatePerYear / m, nil
}

// BondPrice returns the price of a level-coupon bond on a coupon date (no accrued interest), the
// present value of periods coupons and the face at the periodic yield, rounded half-even to the cent:
// P = c·[1 − (1+i)⁻ⁿ]/i + F·(1+i)⁻ⁿ, with c = F·couponRatePerYear/periodsPerYear (not rounded) and
// i = yieldPerYear/periodsPerYear. i = 0 gives exactly c·n + F. faceCents is read only.
// The result is the float64 value rounded half-even; it equals the correctly rounded exact value except
// when that value lies within about n·ε·|x| of a half cent, and is within 1 cent of it for
// |result| ≤ 2⁴³ cents. (n = periods)
// Errors: ErrOutOfRange if |faceCents| > 2⁵³; ErrNegative if faceCents or couponRatePerYear < 0;
// ErrRate if a rate is NaN or ±Inf or the periodic yield is ≤ −1; ErrPeriods if periods < 1 or
// periodsPerYear < 1; ErrOverflow if the price exceeds 2⁵³ cents.
//
// ID:      FIN-INVEST-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist, Knight et al., Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §10.2 "Bond Valuation", https://openstax.org/books/principles-finance/pages/10-2-bond-valuation
// Source:  Macaulay, Some Theoretical Problems Suggested by the Movements of Interest Rates, Bond Yields and Stock Prices in the United States since 1856, NBER, 1938, ch. 2, p. 48 n. 22, https://www.nber.org/system/files/chapters/c6342/c6342.pdf
// Tags:    invest, bond, price, valuation, present-value, coupon, yield
// Units:   faceCents and the result in cents; couponRatePerYear and yieldPerYear decimal per year (nominal, compounded periodsPerYear times); periods count
// Example: OpenStax §10.2: $1,000 3M bond, 2.25 % annual coupon, 15 years, YTM 1.24 % → $1,137.47; 4 % semiannual, 30 periods, 5 % → $895.35
func BondPrice(faceCents int64, couponRatePerYear, yieldPerYear float64, periods, periodsPerYear int) (int64, error) {
	i, c, err := checkBond(faceCents, couponRatePerYear, yieldPerYear, periods, periodsPerYear)
	if err != nil {
		return 0, err
	}
	n := float64(periods)
	var annuity, discount float64
	if i == 0 {
		annuity, discount = n, 1
	} else {
		x := -n * math.Log1p(i)
		discount = math.Exp(x)       // (1+i)⁻ⁿ
		annuity = -math.Expm1(x) / i // [1 − (1+i)⁻ⁿ]/i without cancellation (Goldberg 1991 §1.4)
	}
	return roundCents(c*annuity + float64(faceCents)*discount)
}
