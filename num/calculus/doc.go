// Package calculus holds numerical calculus on sampled data: integration by the trapezoidal rule,
// derivatives of tabulated values by finite differences, and interpolation in a table.
//
// ID short form: CALC (IDs are NUM-CALC-NNN).
//
// Inputs are float64 slices of samples y at abscissae x, each in the caller's unit; a result says its unit
// (an integral is in unit(y)·unit(x), a derivative in unit(y)/unit(x)). No function modifies its input.
// NaN or ±Inf in an input is an error (ErrNotFinite), never a NaN result; a finite input whose result is
// beyond float64 is ErrOverflow. Sums are compensated (Neumaier).
//
// Functions:
//
//	NUM-CALC-001  Trapezoid           ∫y dx over samples at given x, composite trapezoidal rule
//	NUM-CALC-002  TrapezoidUniform    ∫y dx over samples a constant step apart, composite trapezoidal rule
//	NUM-CALC-003  Gradient            dy/dx at every sample, second-order finite differences on any grid
//	NUM-CALC-004  Interp              value at x by piecewise-linear interpolation in a table
package calculus
