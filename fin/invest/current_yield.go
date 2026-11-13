package invest

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Sentinel errors of package invest. Every error a function returns wraps one of these with the name
// of the offending parameter (doc/ASSUMPTIONS.md A-008: they live with the first function that uses them).
var (
	// ErrNotFinite: an input is NaN or ±Inf.
	ErrNotFinite = errors.New("invest: value is NaN or infinite")
	// ErrOutOfRange: a money input is outside ±2⁵³ cents.
	ErrOutOfRange = errors.New("invest: input out of range")
	// ErrOverflow: a result is outside ±2⁵³ cents or beyond float64.
	ErrOverflow = errors.New("invest: result overflows")
	// ErrRate: a rate is NaN, ±Inf or ≤ −1.
	ErrRate = errors.New("invest: rate must be finite and greater than -1")
	// ErrPeriods: a period or day count is out of range.
	ErrPeriods = errors.New("invest: period count out of range")
	// ErrZeroDenominator: the formula would divide by zero.
	ErrZeroDenominator = errors.New("invest: zero denominator")
	// ErrNegative: an input that must not be negative is negative.
	ErrNegative = errors.New("invest: negative input")
	// ErrDomain: inputs outside the formula's domain.
	ErrDomain = errors.New("invest: inputs outside the formula's domain")
	// ErrNoConvergence: the root finder did not converge.
	ErrNoConvergence = errors.New("invest: root finder did not converge")
	// ErrNoSolution: no solution exists for the inputs.
	ErrNoSolution = errors.New("invest: no solution")
)

// maxExactCents is the largest money magnitude, in cents, that float64 holds exactly (2⁵³).
const maxExactCents = 1 << 53

// checkCents returns ErrOutOfRange when |cents| > 2⁵³.
func checkCents(name string, cents int64) error {
	if cents > maxExactCents || cents < -maxExactCents {
		return fmt.Errorf("%w: %s = %d is outside ±2^53 cents", ErrOutOfRange, name, cents)
	}
	return nil
}

// checkNonNegativeCents returns ErrOutOfRange when cents > 2⁵³ and ErrNegative when cents < 0.
func checkNonNegativeCents(name string, cents int64) error {
	if err := checkCents(name, cents); err != nil {
		return err
	}
	if cents < 0 {
		return fmt.Errorf("%w: %s = %d", ErrNegative, name, cents)
	}
	return nil
}

// checkRate returns ErrRate (also wrapping ErrNotFinite for NaN and ±Inf) unless rate is finite and > −1.
func checkRate(name string, rate float64) error {
	if !reduce.AllFinite(rate) {
		return fmt.Errorf("%w: %w: %s = %v", ErrRate, ErrNotFinite, name, rate)
	}
	if rate <= -1 {
		return fmt.Errorf("%w: %s = %v", ErrRate, name, rate)
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

// finiteResult returns ErrOverflow when a float64 result is NaN or ±Inf.
func finiteResult(v float64) (float64, error) {
	if !reduce.AllFinite(v) {
		return 0, fmt.Errorf("%w: result %v is beyond float64", ErrOverflow, v)
	}
	return v, nil
}

// roundCents rounds a float64 cent amount half-even (math.RoundToEven) to int64 cents.
// A NaN, an infinity or a magnitude above 2⁵³ after rounding is ErrOverflow.
func roundCents(v float64) (int64, error) {
	r := math.RoundToEven(v)
	if math.IsNaN(r) || r > maxExactCents || r < -maxExactCents {
		return 0, fmt.Errorf("%w: result %v cents is outside ±2^53", ErrOverflow, v)
	}
	return int64(r), nil
}

// CurrentYield returns the bond's current yield, the annual coupon interest divided by the price paid,
// as a decimal (0.05 = 5 %). Both inputs are in cents; the ratio is not rounded.
// Errors: ErrOutOfRange if an input is outside ±2⁵³ cents; ErrNegative if annualCouponCents or
// priceCents is negative; ErrZeroDenominator if priceCents is 0.
//
// ID:      FIN-INVEST-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Municipal Securities Rulemaking Board (MSRB), "Municipal Bond Basics", https://www.msrb.org/node/13886 ("Current Yield is the ratio of the annual dollar amount of interest paid on a bond to its purchase price"; read 2026-11-12)
// Tags:    invest, bond, yield, current-yield, coupon
// Units:   annualCouponCents and priceCents in cents; result decimal per year
// Example: hand-checked: coupon $50.00, price $952.06 → 5000/95206 = 0.052517698…
func CurrentYield(annualCouponCents, priceCents int64) (float64, error) {
	if err := checkNonNegativeCents("annualCouponCents", annualCouponCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("priceCents", priceCents); err != nil {
		return 0, err
	}
	if priceCents == 0 {
		return 0, fmt.Errorf("%w: priceCents = 0", ErrZeroDenominator)
	}
	return float64(annualCouponCents) / float64(priceCents), nil
}
