package economics

import (
	"errors"
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// The package's sentinel errors. They live in the file of the first function that uses them
// (doc/ASSUMPTIONS.md A-008). Every error a function returns wraps one of these with the name of the
// offending parameter; test with errors.Is.
var (
	// ErrNotFinite: a float64 input is NaN or ±Inf. A non-finite rate also wraps ErrRate.
	ErrNotFinite = errors.New("economics: value is NaN or infinite")
	// ErrOutOfRange: a money input is outside ±2⁵³ cents, or another input is outside its stated range.
	ErrOutOfRange = errors.New("economics: input out of range")
	// ErrOverflow: a float64 result is NaN or ±Inf (beyond the float64 range).
	ErrOverflow = errors.New("economics: result overflows")
	// ErrRate: a rate is NaN, ±Inf, or not greater than −1.
	ErrRate = errors.New("economics: rate must be finite and greater than -1")
	// ErrNegative: an input that must not be negative is negative.
	ErrNegative = errors.New("economics: negative input")
	// ErrZeroDenominator: the input used as the denominator is zero.
	ErrZeroDenominator = errors.New("economics: zero denominator")
	// ErrDomain: the inputs are valid one by one but the formula is undefined for them together.
	ErrDomain = errors.New("economics: inputs outside the formula's domain")
)

// maxExactCents is the largest money magnitude, in cents, that float64 holds exactly (2⁵³).
const maxExactCents = 1 << 53

// checkNonNegativeCents returns ErrOutOfRange when |cents| > 2⁵³ and ErrNegative when cents < 0.
func checkNonNegativeCents(name string, cents int64) error {
	if cents > maxExactCents || cents < -maxExactCents {
		return fmt.Errorf("%w: %s = %d is outside ±2^53 cents", ErrOutOfRange, name, cents)
	}
	if cents < 0 {
		return fmt.Errorf("%w: %s = %d", ErrNegative, name, cents)
	}
	return nil
}

// checkFinite returns ErrNotFinite for NaN or ±Inf.
func checkFinite(name string, v float64) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, v)
	}
	return nil
}

// ConsumerPriceIndex returns the price index of a fixed basket of goods and services, 100 times the
// basket's cost in the current period over its cost in the base period, so the base period is 100. The
// index is not rounded (BLS publishes three decimals; rounding is the caller's). The basket must be the
// same in both periods: this is the textbook fixed-basket ratio, not the BLS modified-Laspeyres
// aggregation of basic indexes.
// Errors: ErrOutOfRange if either cost is outside ±2⁵³ cents; ErrNegative if either cost is negative;
// ErrZeroDenominator if basketCostBaseCents is 0.
//
// ID:      FIN-ECON-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Greenlaw, Shapiro & MacDonald, Principles of Economics, 3rd ed., OpenStax (Rice University), 2022, §22.1 "Tracking Inflation" (index numbers: the base year has an index number equal to 100), https://openstax.org/books/principles-economics-3e/pages/22-1-tracking-inflation (read 2026-11-12)
// Source:  U.S. Bureau of Labor Statistics, Handbook of Methods, "Consumer Price Index: Calculation" (modified Laspeyres index, quantity weights fixed for a year, index base period 1982–1984), https://www.bls.gov/opub/hom/cpi/calculation.htm (read 2026-11-12)
// Tags:    economics, inflation, price-index, cpi, fixed-basket, laspeyres
// Units:   both costs in cents; result in index points (base period = 100)
// Example: OpenStax §22.1, base Period 3 basket $107: Period 4 $117.50 → 109.8, Period 1 $100 → 93.4
func ConsumerPriceIndex(basketCostCurrentCents, basketCostBaseCents int64) (float64, error) {
	if err := checkNonNegativeCents("basketCostCurrentCents", basketCostCurrentCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("basketCostBaseCents", basketCostBaseCents); err != nil {
		return 0, err
	}
	if basketCostBaseCents == 0 {
		return 0, fmt.Errorf("%w: basketCostBaseCents = 0", ErrZeroDenominator)
	}
	return float64(basketCostCurrentCents) / float64(basketCostBaseCents) * 100, nil
}
