// Package automotive holds everyday vehicle calculators: fuel economy, fuel consumption per 100 distance
// units, and the fuel cost of a trip.
//
// ID short form: AUTO (IDs are LIF-AUTO-NNN).
//
// Units: a distance and a fuel volume are taken in one consistent pair of units the caller names (miles and
// US gallons, or kilometres and litres); the result is in that pair. No unit is converted here. Money is int64
// cents, rounded once, half to even. NaN or ±Inf in returns ErrNotFinite; an input outside its range returns
// ErrOutOfRange; a result beyond float64 or int64 returns ErrDomain.
//
// Functions:
//
//	LIF-AUTO-001  FuelEconomy           distance per volume of fuel (mpg, km/L)
//	LIF-AUTO-002  FuelConsumptionPer100 fuel per 100 distance units (L/100 km, gal/100 mi)
//	LIF-AUTO-003  FuelCostCents         fuel cost of a trip in cents, from distance, fuel economy and price per volume
package automotive
