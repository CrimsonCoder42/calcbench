package interest

import "math"

// growingAnnuityPresentWorth returns PV per unit of first payment for a growing ordinary annuity:
// (1 − qⁿ)/(r − g) with q = (1+g)/(1+r), and n/(1+r) when r = g. With u = (r − g)/(1+g) and
// ℓ = ln(1+u) it is n · (e^(−nℓ) − 1)/(−nℓ) · (ℓ/u) / (1+g), with no cancellation when r is near g.
// It assumes checked rates and n ≥ 0.
func growingAnnuityPresentWorth(r, g float64, n int) float64 {
	u := (r - g) / (1 + g)
	nf := float64(n)
	return nf * expm1OverX(-nf*math.Log1p(u)) * log1pOverX(u) / (1 + g)
}

// checkGrowingAnnuity validates the shared inputs of the growing-annuity functions.
func checkGrowingAnnuity(centsName string, cents int64, ratePerPeriod, growthPerPeriod float64, periods, minPeriods int) error {
	if err := checkCents(centsName, cents); err != nil {
		return err
	}
	if err := checkRate("ratePerPeriod", ratePerPeriod); err != nil {
		return err
	}
	if err := checkRate("growthPerPeriod", growthPerPeriod); err != nil {
		return err
	}
	return checkPeriods("periods", periods, minPeriods)
}

// PresentValueGrowingAnnuity returns the value now of periods end-of-period payments, the first
// firstPaymentCents at the end of period 1 and each later one growthPerPeriod larger than the one
// before, discounted at ratePerPeriod: C·(1 − ((1+g)/(1+r))ⁿ)/(r − g), rounded half-even to the cent.
// r = g gives C·n/(1+r); g > r is valid. Evaluated without cancellation when r is near g.
// The result is the float64 value rounded half-even; it equals the correctly rounded exact value
// except when that value lies within about n·ε·|x| of a half cent, and is within 1 cent of it for
// |result| ≤ 2⁴³ cents.
// Errors: ErrOutOfRange if |firstPaymentCents| > 2⁵³; ErrRate if ratePerPeriod or growthPerPeriod is
// NaN, ±Inf or ≤ −1; ErrPeriods if periods < 0; ErrOverflow if the result exceeds ±2⁵³ cents.
//
// ID:      FIN-INT-022
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  A. Damodaran, Applied Corporate Finance, 3rd ed., Wiley, 2010, Appendix 3 "Time Value of Money", "Growing Annuities": PV = A(1+g)[1 − (1+g)^n/(1+r)^n]/(r − g), and nA when r = g, https://pages.stern.nyu.edu/~adamodar/pdfiles/acf3E/book/app3.pdf
// Tags:    interest, time-value, annuity, growing-annuity, present-value, graduated-payment
// Units:   firstPaymentCents and the result in cents; ratePerPeriod and growthPerPeriod decimal per period; periods count
// Example: first payment $1,545,000 (5,000 oz × $300 × 1.03), r = 10 %, g = 3 %, 20 years → $16,145,980 (Damodaran App. 3, gold mine)
func PresentValueGrowingAnnuity(firstPaymentCents int64, ratePerPeriod, growthPerPeriod float64, periods int) (int64, error) {
	if err := checkGrowingAnnuity("firstPaymentCents", firstPaymentCents, ratePerPeriod, growthPerPeriod, periods, 0); err != nil {
		return 0, err
	}
	if firstPaymentCents == 0 || periods == 0 {
		return 0, nil
	}
	return roundCents(float64(firstPaymentCents) * growingAnnuityPresentWorth(ratePerPeriod, growthPerPeriod, periods))
}
