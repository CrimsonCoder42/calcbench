package stream

import (
	"errors"
	"fmt"
	"math"
)

// ErrTickSize: the tick size is not a finite positive number.
var ErrTickSize = errors.New("stream: tick size must be finite and positive")

// twoTo63 is 2⁶³ = −math.MinInt64, the first float64 above the int64 range.
const twoTo63 = float64(1 << 63)

// PriceToTicks returns the whole number of ticks nearest to price, round(price/tickSize), rounding a
// quotient exactly halfway between two integers away from zero (IEEE 754 roundTiesToAway, Go's
// math.Round): 0.375 at tick 0.25 is 1.5 ticks, so 2; −0.375 is −2. The quotient is one rounded
// float64 division, so a decimal price that is not exact in binary (1.15 at tick 0.05 is
// 22.999999999999996) lands on the nearest tick (23), and a price within an ulp of a half-tick may round
// either way; a caller that needs exact decimal ticks keeps prices as int64 ticks throughout.
// Errors: ErrTickSize when tickSize is not finite and positive; ErrNotFinite when price is NaN or ±Inf;
// ErrOverflow when the tick count is outside the int64 range.
//
// ID:      ALG-STREAM-012
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  U.S. Securities and Exchange Commission, Regulation NMS, Rule 612 "Minimum pricing increment", 17 CFR § 242.612 (prices are whole multiples of the minimum increment: price = ticks × tickSize)
// Source:  IEEE Std 754-2019, IEEE Standard for Floating-Point Arithmetic, §4.3.1 (roundTiesToAway)
// Tags:    tick, tick-size, price, minimum-increment, rounding, round-half-away-from-zero, trading, market-microstructure
// Units:   price and tickSize in the same price unit; result in ticks
// Example: PriceToTicks(101.23, 0.05) = 2025 (101.23/0.05 = 2024.6)
func PriceToTicks(price, tickSize float64) (int64, error) {
	if !(tickSize > 0) || math.IsInf(tickSize, 0) {
		return 0, fmt.Errorf("%w: tickSize = %v", ErrTickSize, tickSize)
	}
	if math.IsNaN(price) || math.IsInf(price, 0) {
		return 0, fmt.Errorf("%w: price = %v", ErrNotFinite, price)
	}
	ticks := math.Round(price / tickSize)
	if !(ticks >= -twoTo63 && ticks < twoTo63) {
		return 0, fmt.Errorf("%w: price/tickSize = %v/%v outside int64", ErrOverflow, price, tickSize)
	}
	return int64(ticks), nil
}

// TicksToPrice returns ticks × tickSize, the price of a whole number of ticks, as one rounded float64
// product (a tick count beyond 2⁵³ in magnitude is first rounded to float64).
// Errors: ErrTickSize when tickSize is not finite and positive; ErrOverflow when the product exceeds the
// float64 range.
//
// ID:      ALG-STREAM-013
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  U.S. Securities and Exchange Commission, Regulation NMS, Rule 612 "Minimum pricing increment", 17 CFR § 242.612 (price = ticks × tickSize)
// Tags:    tick, tick-size, price, minimum-increment, trading, market-microstructure
// Units:   ticks in ticks; tickSize in a price unit; result in that price unit
// Example: TicksToPrice(2025, 0.25) = 506.25
func TicksToPrice(ticks int64, tickSize float64) (float64, error) {
	if !(tickSize > 0) || math.IsInf(tickSize, 0) {
		return 0, fmt.Errorf("%w: tickSize = %v", ErrTickSize, tickSize)
	}
	price := float64(ticks) * tickSize
	if math.IsInf(price, 0) {
		return 0, fmt.Errorf("%w: %d ticks × %v", ErrOverflow, ticks, tickSize)
	}
	return price, nil
}

// RoundToTick returns price rounded to the nearest whole tick, TicksToPrice(PriceToTicks(price,
// tickSize), tickSize): halfway cases go away from zero, as PriceToTicks documents. The result is the
// float64 nearest to ticks × tickSize, which for a decimal tick such as 0.05 is not exactly decimal.
// Errors: those of PriceToTicks and TicksToPrice (ErrTickSize, ErrNotFinite, ErrOverflow).
//
// ID:      ALG-STREAM-014
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  U.S. Securities and Exchange Commission, Regulation NMS, Rule 612 "Minimum pricing increment", 17 CFR § 242.612 (price = ticks × tickSize)
// Source:  IEEE Std 754-2019, IEEE Standard for Floating-Point Arithmetic, §4.3.1 (roundTiesToAway)
// Tags:    tick, tick-size, price, round-to-tick, rounding, round-half-away-from-zero, trading, market-microstructure
// Units:   price and tickSize in the same price unit; result in that unit
// Example: RoundToTick(101.23, 0.05) = 101.25; RoundToTick(0.375, 0.25) = 0.5
func RoundToTick(price, tickSize float64) (float64, error) {
	ticks, err := PriceToTicks(price, tickSize)
	if err != nil {
		return 0, err
	}
	return TicksToPrice(ticks, tickSize)
}
