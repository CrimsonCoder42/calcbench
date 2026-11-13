package invest

import "fmt"

// ZeroGrowthStockPrice returns the value of a level dividend paid forever, P = D/r, rounded half-even
// to the cent: a zero-growth common share or a preferred share. The first dividend is one period away.
// The result is the float64 value rounded half-even; it equals the correctly rounded exact value except
// when that value lies within about n·ε·|x| of a half cent, and is within 1 cent of it for
// |result| ≤ 2⁴³ cents. (n = 1)
// Errors: ErrOutOfRange if |dividendCents| > 2⁵³; ErrNegative if dividendCents < 0; ErrRate if
// requiredReturn is NaN, ±Inf or ≤ −1; ErrDomain if requiredReturn ≤ 0; ErrOverflow if the price
// exceeds 2⁵³ cents.
//
// ID:      FIN-INVEST-015
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist, Knight et al., Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §11.2 "Dividend Discount Models (DDMs)" (zero growth, PV = D/r), https://openstax.org/books/principles-finance/pages/11-2-dividend-discount-models-ddms
// Tags:    invest, stock, valuation, dividend-discount-model, zero-growth, preferred-stock, perpetuity
// Units:   dividendCents and the result in cents per share; requiredReturn decimal per period
// Example: OpenStax §11.2: $2.00 a year at 6 % → $33.33
func ZeroGrowthStockPrice(dividendCents int64, requiredReturn float64) (int64, error) {
	if err := checkNonNegativeCents("dividendCents", dividendCents); err != nil {
		return 0, err
	}
	if err := checkRate("requiredReturn", requiredReturn); err != nil {
		return 0, err
	}
	if requiredReturn <= 0 {
		return 0, fmt.Errorf("%w: requiredReturn = %v, want > 0", ErrDomain, requiredReturn)
	}
	return roundCents(float64(dividendCents) / requiredReturn)
}
