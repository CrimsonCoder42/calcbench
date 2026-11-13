package automotive

import (
	"errors"
	"fmt"
	"math"
)

var (
	// ErrNotFinite is returned when an input is NaN or ±Inf.
	ErrNotFinite = errors.New("automotive: value is NaN or infinite")
	// ErrOutOfRange is returned when an input is outside its range.
	ErrOutOfRange = errors.New("automotive: input outside its range")
	// ErrDomain is returned when valid inputs give a result beyond float64 or int64.
	ErrDomain = errors.New("automotive: result beyond the representable range")
)

// checkPositiveFinite returns ErrNotFinite or ErrOutOfRange unless v is finite and > 0.
func checkPositiveFinite(name string, v float64) error {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, v)
	}
	if !(v > 0) {
		return fmt.Errorf("%w: %s = %v, want > 0", ErrOutOfRange, name, v)
	}
	return nil
}

// checkResult returns ErrDomain when v is not finite or is 0 from positive inputs (underflow).
func checkResult(v float64) error {
	if math.IsInf(v, 0) || math.IsNaN(v) || v == 0 {
		return fmt.Errorf("%w: result = %v", ErrDomain, v)
	}
	return nil
}

// FuelEconomy returns the fuel economy: the distance travelled per volume of fuel consumed, distance / fuel,
// in the caller's units (miles per US gallon, kilometres per litre).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when distance ≤ 0 or fuelVolume ≤ 0; ErrDomain when the
// result overflows or underflows.
//
// ID:      LIF-AUTO-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  40 CFR 600.002, definition "Fuel economy" (1): "the average number of miles traveled by an automobile or group of automobiles per volume of fuel consumed" (read 2026-11-13 at law.cornell.edu/cfr/text/40/600.002)
// Tags:    automotive, fuel-economy, mpg, fuel
// Units:   distance mi or km; fuelVolume US gal or L; result distance per volume (mpg or km/L)
// Range:   distance > 0; fuelVolume > 0
// Example: 300 mi on 12 gal → 25 mpg (HAND-DERIVED)
func FuelEconomy(distance, fuelVolume float64) (float64, error) {
	if err := checkPositiveFinite("distance", distance); err != nil {
		return 0, err
	}
	if err := checkPositiveFinite("fuelVolume", fuelVolume); err != nil {
		return 0, err
	}
	e := distance / fuelVolume
	if err := checkResult(e); err != nil {
		return 0, err
	}
	return e, nil
}

// FuelConsumptionPer100 returns the fuel consumed per 100 units of distance, 100 × fuelVolume / distance, in the
// caller's units (litres per 100 km, US gallons per 100 mi): the reciprocal of FuelEconomy (LIF-AUTO-001)
// times 100.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when fuelVolume ≤ 0 or distance ≤ 0; ErrDomain when the
// result overflows or underflows.
//
// ID:      LIF-AUTO-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  40 CFR 600.002, definition "Fuel economy" (1) (distance per volume of fuel, read 2026-11-13); this is its reciprocal scaled to 100 distance units
// Tags:    automotive, fuel-consumption, l-per-100km, fuel
// Units:   fuelVolume L or US gal; distance km or mi; result volume per 100 distance units
// Range:   fuelVolume > 0; distance > 0
// Example: 45 L over 600 km → 7.5 L/100 km (HAND-DERIVED)
func FuelConsumptionPer100(fuelVolume, distance float64) (float64, error) {
	if err := checkPositiveFinite("fuelVolume", fuelVolume); err != nil {
		return 0, err
	}
	if err := checkPositiveFinite("distance", distance); err != nil {
		return 0, err
	}
	c := 100 * fuelVolume / distance
	if err := checkResult(c); err != nil {
		return 0, err
	}
	return c, nil
}

// FuelCostCents returns the fuel cost of a trip, in cents: distance / fuelEconomy × priceCentsPerVolume, the
// fuel volume the trip needs times its price, rounded half to even to the whole cent. distance and fuelEconomy
// use one pair of units (mi and mpg, or km and km/L), and the price is per that volume unit.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when distance < 0, fuelEconomy ≤ 0 or
// priceCentsPerVolume < 0; ErrDomain when the cost exceeds int64 cents.
//
// ID:      LIF-AUTO-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  40 CFR 600.002, definition "Fuel economy" (1) (distance per volume of fuel, read 2026-11-13): fuel used = distance / fuel economy; cost = fuel used × unit price
// Tags:    automotive, fuel-cost, trip, money, fuel
// Units:   distance mi or km; fuelEconomy mpg or km/L (same distance unit); priceCentsPerVolume cents per US gal or per L; result cents
// Range:   distance ≥ 0; fuelEconomy > 0; priceCentsPerVolume ≥ 0
// Example: 300 mi at 25 mpg, 350 c/gal → 12 gal × 350 = 4200 c (HAND-DERIVED)
func FuelCostCents(distance, fuelEconomy float64, priceCentsPerVolume int64) (int64, error) {
	if math.IsNaN(distance) || math.IsInf(distance, 0) {
		return 0, fmt.Errorf("%w: distance = %v", ErrNotFinite, distance)
	}
	if distance < 0 {
		return 0, fmt.Errorf("%w: distance = %v, want ≥ 0", ErrOutOfRange, distance)
	}
	if err := checkPositiveFinite("fuelEconomy", fuelEconomy); err != nil {
		return 0, err
	}
	if priceCentsPerVolume < 0 {
		return 0, fmt.Errorf("%w: priceCentsPerVolume = %d, want ≥ 0", ErrOutOfRange, priceCentsPerVolume)
	}
	c := math.RoundToEven(distance / fuelEconomy * float64(priceCentsPerVolume))
	if math.IsInf(c, 0) || c >= math.MaxInt64 {
		return 0, fmt.Errorf("%w: cost = %v cents", ErrDomain, c)
	}
	return int64(c), nil
}
