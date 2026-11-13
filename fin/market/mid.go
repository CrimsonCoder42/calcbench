package market

import (
	"errors"
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Errors returned by this package. They are wrapped with the offending parameter; test with errors.Is.
var (
	// ErrEmpty: a slice input has no elements.
	ErrEmpty = errors.New("market: empty input")
	// ErrTooFew: a slice has fewer elements than the formula needs.
	ErrTooFew = errors.New("market: too few values")
	// ErrLengthMismatch: paired slices differ in length.
	ErrLengthMismatch = errors.New("market: slice lengths differ")
	// ErrNotFinite: a float input is NaN or ±Inf.
	ErrNotFinite = errors.New("market: value is NaN or Inf")
	// ErrNegative: a price, size or volume is negative.
	ErrNegative = errors.New("market: negative value")
	// ErrOutOfRange: a parameter is outside the values the function accepts (tick count above 2⁵³,
	// tick size ≤ 0, a count < 1).
	ErrOutOfRange = errors.New("market: parameter out of range")
	// ErrDomain: the inputs are valid one by one but the formula is undefined for them (a crossed
	// quote, zero total size, a positive Roll covariance).
	ErrDomain = errors.New("market: formula undefined for these inputs")
	// ErrOverflow: every input is finite but the result exceeds the float64 range.
	ErrOverflow = errors.New("market: result overflows float64")
)

// maxTicks is the largest tick count accepted: 2⁵³, so every tick count converts to float64 exactly.
const maxTicks = 1 << 53

// checkTickSize rejects a tick size that is not finite and > 0.
func checkTickSize(tickSize float64) error {
	if !reduce.AllFinite(tickSize) {
		return fmt.Errorf("%w: tickSize = %v", ErrNotFinite, tickSize)
	}
	if tickSize <= 0 {
		return fmt.Errorf("%w: tickSize = %v, want > 0", ErrOutOfRange, tickSize)
	}
	return nil
}

// checkTicks rejects a price in ticks that is negative or above 2⁵³.
func checkTicks(name string, ticks int64) error {
	if ticks < 0 {
		return fmt.Errorf("%w: %s = %d", ErrNegative, name, ticks)
	}
	if ticks > maxTicks {
		return fmt.Errorf("%w: %s = %d, want ≤ 2^53", ErrOutOfRange, name, ticks)
	}
	return nil
}

// checkQuote validates a two-sided quote: both prices in [0, 2⁵³] ticks and askTicks ≥ bidTicks (a
// locked quote is allowed, a crossed one is ErrDomain).
func checkQuote(bidTicks, askTicks int64) error {
	if err := checkTicks("bidTicks", bidTicks); err != nil {
		return err
	}
	if err := checkTicks("askTicks", askTicks); err != nil {
		return err
	}
	if askTicks < bidTicks {
		return fmt.Errorf("%w: askTicks = %d < bidTicks = %d (crossed quote)", ErrDomain, askTicks, bidTicks)
	}
	return nil
}

// checkSize rejects a size or volume that is NaN, ±Inf or negative.
func checkSize(name string, size float64) error {
	if !reduce.AllFinite(size) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, size)
	}
	if size < 0 {
		return fmt.Errorf("%w: %s = %v", ErrNegative, name, size)
	}
	return nil
}

// Mid returns the mid-price M = (bid + ask)/2 of a quote, in price units: the midpoint in ticks
// times tickSize. A locked quote (bid = ask) returns that price.
// Errors: ErrNegative if a price is negative; ErrOutOfRange if a price is above 2⁵³ ticks or tickSize
// ≤ 0; ErrNotFinite if tickSize is NaN or ±Inf; ErrDomain if askTicks < bidTicks (crossed);
// ErrOverflow if the result exceeds float64.
//
// ID:      FIN-MKT-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Stoikov, "The micro-price: a high-frequency estimator of future prices", Quantitative Finance 18(12):1959–1966, 2018, doi:10.1080/14697688.2018.1489139, §2 (mid-price M = (P^b + P^a)/2)
// Tags:    market, microstructure, mid-price, midpoint, quote, order-book
// Units:   bidTicks, askTicks in ticks; tickSize in price units per tick; result in price units
// Example: hand-checked: bid 2419, ask 2421 ticks, tick 0.01 → 24.20
func Mid(bidTicks, askTicks int64, tickSize float64) (float64, error) {
	if err := checkQuote(bidTicks, askTicks); err != nil {
		return 0, err
	}
	if err := checkTickSize(tickSize); err != nil {
		return 0, err
	}
	m := (float64(bidTicks) + float64(askTicks-bidTicks)/2) * tickSize
	if !reduce.AllFinite(m) {
		return 0, fmt.Errorf("%w: mid", ErrOverflow)
	}
	return m, nil
}
