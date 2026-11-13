// Package logs holds logarithms, exponentials and transcendental (and related irrational) constants.
//
// ID short form: LOGS (IDs are NUM-LOGS-NNN).
//
// The standard library already covers most of this subject (math.Log, math.Log1p, math.Exp, math.Expm1,
// math.Ldexp, math.Pow10, math.E, math.Pi); those are not wrapped here. This package holds only what
// the standard library does not name, or what another calcbench package needs under a cited source.
//
// Functions and constants:
//
//	NUM-LOGS-001  GoldenRatio  the golden ratio φ = (1+√5)/2
//	NUM-LOGS-002  LogAddExp    log(eᵃ + eᵇ) without overflow or underflow
//	NUM-LOGS-003  LogAddExp2   log₂(2ᵃ + 2ᵇ) without overflow or underflow
//	NUM-LOGS-004  LogSumExp    log(Σ e^{xᵢ}) of a slice, shifted by the maximum
//	NUM-LOGS-005  ExponentialDecayRemaining  quantity left after time t, A₀·2^(−t/T½)
//	NUM-LOGS-006  ExponentialDecayInitial    starting quantity from what is left, A·2^(t/T½)
//	NUM-LOGS-007  ExponentialDecayHalfLife   half-life from start, end and elapsed time
//	NUM-LOGS-008  ExponentialDecayElapsed    time to decay from A₀ to A, T½·log₂(A₀/A)
//	NUM-LOGS-009  DecayConstantFromHalfLife  decay constant λ = ln 2/T½
//	NUM-LOGS-010  MeanLifetimeFromHalfLife   mean lifetime τ = T½/ln 2
package logs
