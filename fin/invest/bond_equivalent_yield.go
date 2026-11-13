package invest

import "fmt"

// BondEquivalentYield returns the coupon-equivalent yield (Treasury "investment rate") of a discount
// bill with not more than half a year to maturity, i = ((F − P)/P)·(y/r) as a decimal, where r is
// daysToMaturity and y is daysInYear, the number of days in the year that starts on the issue date
// (365 or 366). Bills with more than half a year to run use Treasury's quadratic, not built here.
// Errors: ErrOutOfRange if an amount is above 2⁵³ cents; ErrNegative if faceCents or priceCents < 0;
// ErrZeroDenominator if priceCents is 0; ErrDomain if faceCents is 0, daysInYear is not 365 or 366,
// or 2·daysToMaturity > daysInYear; ErrPeriods if daysToMaturity < 1.
//
// ID:      FIN-INVEST-009
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  US Treasury (TreasuryDirect), "Price, Yield and Rate Calculations for a Treasury Bill", 2004, p. 2 "Calculate Coupon Equivalent Yield For bills of not more than one half-year to maturity", https://www.treasurydirect.gov/instit/annceresult/press/preanre/2004/ofcalc6decbill.pdf
// Tags:    invest, treasury-bill, t-bill, investment-rate, coupon-equivalent-yield, bond-equivalent-yield, money-market
// Units:   faceCents and priceCents in cents; daysToMaturity and daysInYear days; result decimal per year
// Example: Treasury 2004 p. 2: price 99.937778 per 100, 28 days, 366-day year → i = 0.008138368141143 (0.814 %)
func BondEquivalentYield(faceCents, priceCents int64, daysToMaturity, daysInYear int) (float64, error) {
	if err := checkBill(faceCents, priceCents, daysToMaturity, ErrDomain, ErrZeroDenominator); err != nil {
		return 0, err
	}
	if daysInYear != 365 && daysInYear != 366 {
		return 0, fmt.Errorf("%w: daysInYear = %d, want 365 or 366", ErrDomain, daysInYear)
	}
	if 2*daysToMaturity > daysInYear {
		return 0, fmt.Errorf("%w: daysToMaturity = %d is more than half of daysInYear = %d", ErrDomain, daysToMaturity, daysInYear)
	}
	discount := float64(faceCents - priceCents) // exact: both within 0 … 2⁵³
	return discount / float64(priceCents) * (float64(daysInYear) / float64(daysToMaturity)), nil
}
