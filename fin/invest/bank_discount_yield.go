package invest

import "fmt"

// checkBill validates a discount instrument's face, price and days to maturity: both amounts within
// 0 … 2⁵³ cents, faceCents and priceCents not zero, days ≥ 1. zeroFace and zeroPrice are the errors
// for a zero face and a zero price (the one in the denominator is ErrZeroDenominator).
func checkBill(faceCents, priceCents int64, days int, zeroFace, zeroPrice error) error {
	if err := checkNonNegativeCents("faceCents", faceCents); err != nil {
		return err
	}
	if err := checkNonNegativeCents("priceCents", priceCents); err != nil {
		return err
	}
	if faceCents == 0 {
		return fmt.Errorf("%w: faceCents = 0", zeroFace)
	}
	if priceCents == 0 {
		return fmt.Errorf("%w: priceCents = 0", zeroPrice)
	}
	if days < 1 {
		return fmt.Errorf("%w: daysToMaturity = %d, want ≥ 1", ErrPeriods, days)
	}
	return nil
}

// BankDiscountYield returns the discount rate of a discount instrument such as a Treasury bill,
// d = ((F − P)/F)·(360/r): the discount as a fraction of face value, annualized on a 360-day year,
// as a decimal. r is daysToMaturity. A price above face gives a negative rate.
// Errors: ErrOutOfRange if an amount is above 2⁵³ cents; ErrNegative if faceCents or priceCents < 0;
// ErrZeroDenominator if faceCents is 0; ErrDomain if priceCents is 0; ErrPeriods if daysToMaturity < 1.
//
// ID:      FIN-INVEST-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  US Treasury (TreasuryDirect), "Price, Yield and Rate Calculations for a Treasury Bill", 2004, p. 1 "Convert Price to Discount Rate", https://www.treasurydirect.gov/instit/annceresult/press/preanre/2004/ofcalc6decbill.pdf
// Source:  TreasuryDirect, "Understanding Pricing and Interest Rates" (Price = Face value (1 − (discount rate × time)/360)), https://treasurydirect.gov/marketable-securities/understanding-pricing/ (read 2026-11-12)
// Tags:    invest, treasury-bill, t-bill, discount-rate, bank-discount-yield, money-market
// Units:   faceCents and priceCents in cents; daysToMaturity days; result decimal per 360-day year
// Example: Treasury 2004 p. 1: price 99.937778 per 100, 28 days → d = 0.007999971428571 (0.800 %)
func BankDiscountYield(faceCents, priceCents int64, daysToMaturity int) (float64, error) {
	if err := checkBill(faceCents, priceCents, daysToMaturity, ErrZeroDenominator, ErrDomain); err != nil {
		return 0, err
	}
	discount := float64(faceCents - priceCents) // exact: both within 0 … 2⁵³
	return discount / float64(faceCents) * (360 / float64(daysToMaturity)), nil
}
