package interest

import (
	"fmt"
	"math"
)

// solveAnnuityRate returns the rate i with pv = pmt·(P/A,i,n) + fv·(P/F,i,n) for checked inputs
// (pv > 0, pmt ≥ 0, fv ≥ 0, not both 0, n ≥ 1), by the package's bracketed rtsafe solver. Each
// evaluation is O(1) whatever n is: the factors use their closed forms, never a loop over periods.
// A term whose coefficient is 0 is skipped, so 0·Inf never makes a NaN.
func solveAnnuityRate(pv, pmt, fv float64, n int) (float64, error) {
	nf := float64(n)
	return solveRate(func(rate float64) (float64, float64) {
		f, df := -pv, 0.0
		if pmt != 0 {
			f += pmt * uniformSeriesPresentWorth(rate, n)
			// d(P/A)/di = (n(1+i)^−(n+1) − P/A)/i; near i = 0 its limit −n(n+1)/2 avoids cancellation.
			// The derivative only steers Newton; rtsafe bisects whenever it is poor or not finite.
			if math.Abs(nf*rate) < 1e-6 {
				df += pmt * (-nf * (nf + 1) / 2)
			} else {
				df += pmt * (nf*math.Exp(-(nf+1)*math.Log1p(rate)) - uniformSeriesPresentWorth(rate, n)) / rate
			}
		}
		if fv != 0 {
			// d(P/F)/di = −n(1+i)^−(n+1).
			f += fv * math.Exp(-nf*math.Log1p(rate))
			df -= fv * nf * math.Exp(-(nf+1)*math.Log1p(rate))
		}
		return f, df
	})
}

// AnnuityRate returns the rate per period i at which presentValueCents equals the present value of
// paymentCents paid at the end of each of periods periods plus futureValueCents paid at the end of
// the last: PV = PMT·(1 − (1+i)⁻ⁿ)/i + FV·(1+i)⁻ⁿ. This is the yield of a level-coupon bond bought at
// PV (coupon PMT, face FV), the rate of a loan (FV = 0) and the growth rate of a lump sum (PMT = 0).
// With non-negative inputs the right side falls strictly as i rises, so the rate is unique; it is
// found by the scan-and-rtsafe solver described at InternalRateOfReturn. Each step costs O(1) for
// any periods, and nothing is allocated. The result is unrounded.
// Errors: ErrOutOfRange if any amount is outside ±2⁵³ cents; ErrNegative if any amount is negative;
// ErrDomain if presentValueCents = 0 or paymentCents and futureValueCents are both 0; ErrPeriods if
// periods < 1; ErrNoSolution if the rate lies outside (−1 + 1e-12, 1e6].
//
// ID:      FIN-INT-042
// Time:    O(k)   k ≤ 72 scan steps + 200 refinement steps, independent of periods
// Space:   O(1), no allocation
// Source:  R. A. DeFusco, D. W. McLeavey, J. E. Pinto and D. E. Runkle, Quantitative Investment Analysis, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ch. 1 "The Time Value of Money", §6.1 Eq. 11 and §5 Eq. 8 (the equation), §7.1 "Solving for Interest Rates and Growth Rates" (solving it for r)
// Source:  J. Dahlquist and R. Knight, Principles of Finance, OpenStax, 2022, §7.3 "Methods for Solving Time Value of Money Problems" (solving for I/Y), https://openstax.org/books/principles-finance/pages/7-3-methods-for-solving-time-value-of-money-problems; W. H. Press et al., Numerical Recipes, 3rd ed., 2007, §9.4 (solver)
// Tags:    interest, time-value, rate, solve-for-rate, yield, yield-to-maturity, annuity, root-finding
// Units:   presentValueCents, paymentCents and futureValueCents in cents; periods count; result decimal per period
// Example: PV $1,000, PMT 0, FV $1,125.51, N 4 → 3 % (Principles of Finance §7.3); PV $100,000, PMT $733.76, FV 0, N 360 → 0.08/12 (QIA ch. 1 Example 20, reversed)
func AnnuityRate(presentValueCents, paymentCents, futureValueCents int64, periods int) (float64, error) {
	for _, v := range [...]struct {
		name  string
		cents int64
	}{{"presentValueCents", presentValueCents}, {"paymentCents", paymentCents}, {"futureValueCents", futureValueCents}} {
		if err := checkCents(v.name, v.cents); err != nil {
			return 0, err
		}
		if v.cents < 0 {
			return 0, fmt.Errorf("%w: %s = %d", ErrNegative, v.name, v.cents)
		}
	}
	if presentValueCents == 0 {
		return 0, fmt.Errorf("%w: presentValueCents = 0", ErrDomain)
	}
	if paymentCents == 0 && futureValueCents == 0 {
		return 0, fmt.Errorf("%w: paymentCents and futureValueCents are both 0", ErrDomain)
	}
	if err := checkPeriods("periods", periods, 1); err != nil {
		return 0, err
	}
	return solveAnnuityRate(float64(presentValueCents), float64(paymentCents), float64(futureValueCents), periods)
}
