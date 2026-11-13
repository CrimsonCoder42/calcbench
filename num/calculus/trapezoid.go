package calculus

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

var (
	// ErrNotFinite reports a NaN or ±Inf input.
	ErrNotFinite = errors.New("calculus: value is NaN or infinite")
	// ErrTooFew reports fewer samples than the method needs.
	ErrTooFew = errors.New("calculus: too few samples")
	// ErrLengthMismatch reports x and y of different lengths.
	ErrLengthMismatch = errors.New("calculus: inputs have different lengths")
	// ErrOverflow reports a result beyond the float64 range.
	ErrOverflow = errors.New("calculus: result exceeds the float64 range")
	// ErrNotAscending reports abscissae that are not strictly ascending.
	ErrNotAscending = errors.New("calculus: x not strictly ascending")
	// ErrOutOfRange reports a query outside the table.
	ErrOutOfRange = errors.New("calculus: out of range")
)

// checkFinite returns ErrNotFinite, naming the parameter and index, for the first NaN or ±Inf.
func checkFinite(name string, values []float64) error {
	for i, v := range values {
		if !reduce.AllFinite(v) {
			return fmt.Errorf("%w: %s[%d] = %v", ErrNotFinite, name, i, v)
		}
	}
	return nil
}

// neumaierSum is a running sum with Neumaier's correction (Neumaier, ZAMM 1974, 54(1):39–51).
type neumaierSum struct{ sum, comp float64 }

func (s *neumaierSum) add(v float64) {
	s.sum, s.comp = reduce.NeumaierAdd(s.sum, s.comp, v)
}

func (s *neumaierSum) total() float64 { return s.sum + s.comp }

// Trapezoid returns the integral of the samples y over the abscissae x by the composite trapezoidal rule,
// Σ (x[i+1] − x[i])·(y[i] + y[i+1])/2. x need not be evenly spaced; descending x gives the negated
// integral, as ∫ from x[0] to x[n−1] does. Each panel is formed as Δx·(y[i]/2 + y[i+1]/2), so a single
// panel overflows only when its own value does; the sum is compensated. x and y are read, never modified.
// Errors: ErrLengthMismatch when len(y) ≠ len(x); ErrTooFew when there are fewer than 2 samples;
// ErrNotFinite when a value is NaN or ±Inf; ErrOverflow when a step, a panel or the sum is beyond float64.
//
// ID:      NUM-CALC-001
// Time:    O(n)   n = len(x)
// Space:   O(1) extra, no allocation
// Source:  Abramowitz, Stegun, Handbook of Mathematical Functions, 10th printing, 1972, eq. 25.4.2 (trapezoidal rule)
// Source:  Press, Teukolsky, Vetterling, Flannery, Numerical Recipes, 3rd ed., Cambridge, 2007, §4.1, eq. 4.1.11 (extended trapezoidal rule)
// Tags:    calculus, integration, quadrature, trapezoid, trapezoidal-rule, sampled-data
// Units:   x in any unit, y in any unit; result in unit(y)·unit(x)
// Example: y = x² at x = 0, 1, 2, 3 → 9.5 (exact integral 9; error −(b−a)h²f”/12 = −(−0.5)) (hand-checked)
func Trapezoid(y, x []float64) (float64, error) {
	if len(y) != len(x) {
		return 0, fmt.Errorf("%w: len(y) = %d, len(x) = %d", ErrLengthMismatch, len(y), len(x))
	}
	if len(x) < 2 {
		return 0, fmt.Errorf("%w: %d samples, want ≥ 2", ErrTooFew, len(x))
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	if err := checkFinite("y", y); err != nil {
		return 0, err
	}
	var s neumaierSum
	for i := 1; i < len(x); i++ {
		h := x[i] - x[i-1]
		p := h * (y[i-1]/2 + y[i]/2)
		if !reduce.AllFinite(p) {
			return 0, fmt.Errorf("%w: panel %d", ErrOverflow, i)
		}
		s.add(p)
	}
	t := s.total()
	if !reduce.AllFinite(t) {
		return 0, fmt.Errorf("%w: sum of panels", ErrOverflow)
	}
	return t, nil
}

// TrapezoidUniform returns the integral of the samples y, a constant step dx apart, by the composite
// trapezoidal rule dx·(y[0]/2 + y[1] + … + y[n−2] + y[n−1]/2). A negative dx gives the negated integral.
// The sum is compensated. y is read, never modified.
// Errors: ErrTooFew when there are fewer than 2 samples; ErrNotFinite when dx or a value is NaN or ±Inf;
// ErrOverflow when the result is beyond float64.
//
// ID:      NUM-CALC-002
// Time:    O(n)   n = len(y)
// Space:   O(1) extra, no allocation
// Source:  Abramowitz, Stegun, Handbook of Mathematical Functions, 10th printing, 1972, eq. 25.4.2 (trapezoidal rule)
// Source:  Press, Teukolsky, Vetterling, Flannery, Numerical Recipes, 3rd ed., Cambridge, 2007, §4.1, eq. 4.1.11 (extended trapezoidal rule)
// Tags:    calculus, integration, quadrature, trapezoid, trapezoidal-rule, sampled-data
// Units:   dx in any unit, y in any unit; result in unit(y)·unit(dx)
// Example: y = 1, 2, 3 with dx = 0.5 → 2 (hand-checked)
func TrapezoidUniform(y []float64, dx float64) (float64, error) {
	if len(y) < 2 {
		return 0, fmt.Errorf("%w: %d samples, want ≥ 2", ErrTooFew, len(y))
	}
	if !reduce.AllFinite(dx) {
		return 0, fmt.Errorf("%w: dx = %v", ErrNotFinite, dx)
	}
	if err := checkFinite("y", y); err != nil {
		return 0, err
	}
	r := uniformSum(y, 1) * dx
	if !reduce.AllFinite(r) {
		// The sum or its product with dx overflowed: redo on values scaled by 2^-64 (exact).
		r = math.Ldexp(uniformSum(y, 0x1p-64)*dx, 64)
		if !reduce.AllFinite(r) {
			return 0, fmt.Errorf("%w: integral", ErrOverflow)
		}
	}
	return r, nil
}

// uniformSum returns scale·(y[0]/2 + y[1] + … + y[n−2] + y[n−1]/2), compensated; scale is a power of two.
func uniformSum(y []float64, scale float64) float64 {
	n := len(y)
	var s neumaierSum
	s.add(y[0] / 2 * scale)
	for _, v := range y[1 : n-1] {
		s.add(v * scale)
	}
	s.add(y[n-1] / 2 * scale)
	return s.total()
}
