package physics

import (
	"errors"
	"fmt"
	"math"
)

var (
	// ErrNotFinite reports a NaN or infinite input.
	ErrNotFinite = errors.New("physics: value is NaN or infinite")
	// ErrNegative reports a negative quantity.
	ErrNegative = errors.New("physics: value is negative")
	// ErrZero reports a zero divisor.
	ErrZero = errors.New("physics: divisor is zero")
	// ErrOverflow reports a result beyond float64.
	ErrOverflow = errors.New("physics: result overflows float64")
)

// checkNonNegative rejects NaN, ±Inf and negative values, naming the parameter.
func checkNonNegative(name string, v float64) error {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, v)
	}
	if v < 0 {
		return fmt.Errorf("%w: %s = %v", ErrNegative, name, v)
	}
	return nil
}

// quotient returns num/den for checked non-negative inputs, rejecting a zero den and an overflowing result.
func quotient(num float64, numName string, den float64, denName string) (float64, error) {
	if err := checkNonNegative(numName, num); err != nil {
		return 0, err
	}
	if err := checkNonNegative(denName, den); err != nil {
		return 0, err
	}
	if den == 0 {
		return 0, fmt.Errorf("%w: %s = 0", ErrZero, denName)
	}
	q := num / den
	if math.IsInf(q, 0) {
		return 0, fmt.Errorf("%w: %s / %s", ErrOverflow, numName, denName)
	}
	return q, nil
}

// product returns a·b for checked non-negative inputs, rejecting an overflowing result.
func product(a float64, aName string, b float64, bName string) (float64, error) {
	if err := checkNonNegative(aName, a); err != nil {
		return 0, err
	}
	if err := checkNonNegative(bName, b); err != nil {
		return 0, err
	}
	p := a * b
	if math.IsInf(p, 0) {
		return 0, fmt.Errorf("%w: %s · %s", ErrOverflow, aName, bName)
	}
	return p, nil
}

// Density returns the mass density of a body, ρ = m / V, in kg/m³.
// Errors: ErrNotFinite; ErrNegative for a negative mass or volume; ErrZero for volumeCubicMeters = 0; ErrOverflow.
//
// ID:      SCI-PHYS-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ling SJ, Sanny J, Moebs W. University Physics Volume 1. OpenStax, 2016, https://openstax.org/details/books/university-physics-volume-1, §14.1 "Fluids, Density, and Pressure" (ρ = m/V); BIPM, The International System of Units (SI Brochure), 9th ed., 2019, Table 5 (mass density, kilogram per cubic meter)
// Tags:    density, mass, volume, physics, fluids
// Units:   massKg kg; volumeCubicMeters m³; result kg/m³
// Example: 8900 kg in 1 m³ → 8900 kg/m³; 1000 kg in 1 m³ → 1000 kg/m³
func Density(massKg, volumeCubicMeters float64) (float64, error) {
	return quotient(massKg, "massKg", volumeCubicMeters, "volumeCubicMeters")
}

// MassFromDensity returns the mass of a body of the given density and volume, m = ρ · V, in kg.
// Errors: ErrNotFinite; ErrNegative for a negative density or volume; ErrOverflow.
//
// ID:      SCI-PHYS-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ling SJ, Sanny J, Moebs W. University Physics Volume 1. OpenStax, 2016, https://openstax.org/details/books/university-physics-volume-1, §14.1 "Fluids, Density, and Pressure" (ρ = m/V, solved for m)
// Tags:    density, mass, volume, physics, fluids
// Units:   densityKgPerCubicMeter kg/m³; volumeCubicMeters m³; result kg
// Example: 8900 kg/m³ × 0.5 m³ → 4450 kg
func MassFromDensity(densityKgPerCubicMeter, volumeCubicMeters float64) (float64, error) {
	return product(densityKgPerCubicMeter, "densityKgPerCubicMeter", volumeCubicMeters, "volumeCubicMeters")
}

// VolumeFromDensity returns the volume of a body of the given mass and density, V = m / ρ, in m³.
// Errors: ErrNotFinite; ErrNegative for a negative mass or density; ErrZero for density 0; ErrOverflow.
//
// ID:      SCI-PHYS-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ling SJ, Sanny J, Moebs W. University Physics Volume 1. OpenStax, 2016, https://openstax.org/details/books/university-physics-volume-1, §14.1 "Fluids, Density, and Pressure" (ρ = m/V, solved for V)
// Tags:    density, mass, volume, physics, fluids
// Units:   massKg kg; densityKgPerCubicMeter kg/m³; result m³
// Example: 8900 kg at 8900 kg/m³ → 1 m³
func VolumeFromDensity(massKg, densityKgPerCubicMeter float64) (float64, error) {
	return quotient(massKg, "massKg", densityKgPerCubicMeter, "densityKgPerCubicMeter")
}
