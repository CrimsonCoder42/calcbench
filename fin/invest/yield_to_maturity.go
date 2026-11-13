package invest

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

// fromInterest re-labels an error returned by package interest with this package's sentinel of the
// same meaning, keeping the original message, so callers test errors.Is against invest.Err… only.
func fromInterest(err error) error {
	if err == nil {
		return nil
	}
	pairs := [...]struct{ from, to error }{
		{interest.ErrNoConvergence, ErrNoConvergence},
		{interest.ErrNoSolution, ErrNoSolution},
		{interest.ErrRate, ErrRate},
		{interest.ErrPeriods, ErrPeriods},
		{interest.ErrOutOfRange, ErrOutOfRange},
		{interest.ErrOverflow, ErrOverflow},
		{interest.ErrDomain, ErrDomain},
		{interest.ErrNegative, ErrNegative},
		{interest.ErrNotFinite, ErrNotFinite},
	}
	for _, p := range pairs {
		if errors.Is(err, p.from) {
			return fmt.Errorf("%w: %v", p.to, err)
		}
	}
	return err
}

// YieldToMaturity returns the bond's yield to maturity as a nominal annual rate, periodsPerYear times
// the periodic rate i at which the present value of the remaining coupons and the face equals
// priceCents (bought on a coupon date): price = c·(P/A,i,n) + F·(P/F,i,n). The coupon per period,
// faceCents·couponRatePerYear/periodsPerYear, is rounded half-even to the cent (a paid coupon is whole
// cents) and i is found by interest.AnnuityRate (FIN-INT-042: closed-form factors, bracketing from 0,
// then Newton guarded by bisection, |Δi| ≤ 1e-12·max(1, |i|)). The right side falls strictly as i
// rises, so the root is unique. Each evaluation is O(1) whatever periods is; nothing is allocated.
// Errors: ErrOutOfRange if |priceCents| or |faceCents| > 2⁵³ or the coupon per period exceeds 2⁵³
// cents; ErrNegative if priceCents, faceCents or couponRatePerYear < 0; ErrDomain if priceCents is 0
// or the bond pays nothing (face and coupon 0); ErrRate if couponRatePerYear is NaN or ±Inf;
// ErrPeriods if periods < 1 or periodsPerYear < 1; ErrNoSolution if the periodic yield lies outside
// (−1 + 1e-12, 1e6]; ErrNoConvergence if the root finder does not converge.
//
// ID:      FIN-INVEST-003
// Time:    O(k)   k ≤ 272 O(1) evaluations (interest.AnnuityRate), independent of periods
// Space:   O(1), no allocation
// Source:  Dahlquist, Knight et al., Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §10.2 "Bond Valuation" (yield to maturity), https://openstax.org/books/principles-finance/pages/10-2-bond-valuation
// Tags:    invest, bond, yield, yield-to-maturity, ytm, irr, root-finding
// Units:   priceCents and faceCents in cents; couponRatePerYear decimal per year; periods count; result decimal nominal per year
// Example: OpenStax §10.2 Coca-Cola: price $952.06, par $1,000, $5.00 per half-year, 14 periods → 0.8651 % per half-year, 1.73 % per year
func YieldToMaturity(priceCents, faceCents int64, couponRatePerYear float64, periods, periodsPerYear int) (float64, error) {
	if err := checkNonNegativeCents("priceCents", priceCents); err != nil {
		return 0, err
	}
	if priceCents == 0 {
		return 0, fmt.Errorf("%w: priceCents = 0", ErrDomain)
	}
	// The yield argument of checkBond is unused here; 0 is always valid.
	_, c, err := checkBond(faceCents, couponRatePerYear, 0, periods, periodsPerYear)
	if err != nil {
		return 0, err
	}
	couponCents := math.RoundToEven(c)
	if couponCents > maxExactCents {
		return 0, fmt.Errorf("%w: coupon per period %v cents is outside 2^53", ErrOutOfRange, couponCents)
	}
	i, err := interest.AnnuityRate(priceCents, int64(couponCents), faceCents, periods)
	if err != nil {
		return 0, fromInterest(err)
	}
	return i * float64(periodsPerYear), nil
}
