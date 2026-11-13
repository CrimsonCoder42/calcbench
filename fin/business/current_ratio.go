package business

import (
	"errors"
	"fmt"
)

// Sentinel errors of package business. Every error is wrapped with the offending parameter name; test with
// errors.Is.
var (
	// ErrNotFinite: a float64 input is NaN or ±Inf.
	ErrNotFinite = errors.New("business: input is not finite")
	// ErrOutOfRange: a money or count input is outside ±2⁵³.
	ErrOutOfRange = errors.New("business: input out of range")
	// ErrOverflow: a money result is outside ±2⁵³ cents.
	ErrOverflow = errors.New("business: result overflows")
	// ErrRate: a rate is outside its documented range.
	ErrRate = errors.New("business: rate out of range")
	// ErrPeriods: a period count is outside its documented range.
	ErrPeriods = errors.New("business: period count out of range")
	// ErrZeroDenominator: the denominator of a ratio is zero.
	ErrZeroDenominator = errors.New("business: zero denominator")
	// ErrNegative: an input that cannot be negative is negative.
	ErrNegative = errors.New("business: negative input")
	// ErrDomain: the inputs are finite but the formula has no meaning for them.
	ErrDomain = errors.New("business: input outside the formula's domain")
)

// maxExactCents is 2⁵³: every int64 with |x| ≤ maxExactCents converts to float64 exactly.
const maxExactCents = 1 << 53

// checkCents returns ErrOutOfRange when |v| > 2⁵³.
func checkCents(name string, v int64) error {
	if v > maxExactCents || v < -maxExactCents {
		return fmt.Errorf("%w: %s = %d, want |x| ≤ 2^53", ErrOutOfRange, name, v)
	}
	return nil
}

// checkNonNegativeCents returns ErrOutOfRange when |v| > 2⁵³ and ErrNegative when v < 0.
func checkNonNegativeCents(name string, v int64) error {
	if err := checkCents(name, v); err != nil {
		return err
	}
	if v < 0 {
		return fmt.Errorf("%w: %s = %d", ErrNegative, name, v)
	}
	return nil
}

// ratioCents returns num/den for two range-checked int64 values: both convert to float64 exactly, so the
// result has one rounding and is always finite (|den| ≥ 1).
func ratioCents(num int64, denName string, den int64) (float64, error) {
	if den == 0 {
		return 0, fmt.Errorf("%w: %s = 0", ErrZeroDenominator, denName)
	}
	return float64(num) / float64(den), nil
}

// CurrentRatio returns current assets divided by current liabilities, as a decimal (2 means 2:1). It is
// not rounded.
// Errors: ErrOutOfRange when an input is outside ±2⁵³; ErrNegative when an input is negative;
// ErrZeroDenominator when currentLiabilitiesCents is 0.
//
// ID:      FIN-BIZ-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Dahlquist and Knight, Principles of Finance, 1st ed., OpenStax (Rice University), 2022, §6.3 "Liquidity Ratios", https://openstax.org/books/principles-finance/pages/6-3-liquidity-ratios (read 2026-11-12)
// Source:  Franklin, Graybeal and Cooper, Principles of Accounting, Vol. 1: Financial Accounting, OpenStax (Rice University), 2019, Appendix A "Financial Statement Analysis", https://openstax.org/books/principles-financial-accounting/pages/a-financial-statement-analysis (read 2026-11-12)
// Tags:    ratio, liquidity, current-ratio, balance-sheet
// Units:   currentAssetsCents and currentLiabilitiesCents in cents; the result is a dimensionless decimal
// Example: Clear Lake Sporting Goods $200,000 ÷ $100,000 = 2 (Principles of Finance §6.3)
func CurrentRatio(currentAssetsCents, currentLiabilitiesCents int64) (float64, error) {
	if err := checkNonNegativeCents("currentAssetsCents", currentAssetsCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("currentLiabilitiesCents", currentLiabilitiesCents); err != nil {
		return 0, err
	}
	return ratioCents(currentAssetsCents, "currentLiabilitiesCents", currentLiabilitiesCents)
}
