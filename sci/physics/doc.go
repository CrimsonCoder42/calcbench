// Package physics holds the defining relations of introductory physics: density (ρ = m/V) and average speed
// (v = d/t), each solvable for any one of its three quantities, and average power from work (P = F·d/t).
//
// ID short form: PHYS (IDs are SCI-PHYS-NNN).
//
// Units: SI throughout, named in every parameter (kilograms, cubic meters, meters, seconds, newtons, watts). Conversions from
// other units are the caller's, through unit/<quantity>. NaN or ±Inf returns ErrNotFinite, a negative
// quantity ErrNegative, a zero divisor ErrZero, a result beyond float64 ErrOverflow. Results are not rounded.
//
// Functions:
//
//	SCI-PHYS-001  Density                ρ = m / V
//	SCI-PHYS-002  MassFromDensity        m = ρ · V
//	SCI-PHYS-003  VolumeFromDensity      V = m / ρ
//	SCI-PHYS-004  AverageSpeed           v = d / t
//	SCI-PHYS-005  DistanceFromSpeed      d = v · t
//	SCI-PHYS-006  DurationFromSpeed      t = d / v
//	SCI-PHYS-007  PowerFromWork          P = F · d / t
package physics
