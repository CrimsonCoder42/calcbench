package float

import (
	"errors"
	"fmt"
	"math"
)

// Sentinel errors for package float, wrapped with the offending parameter's name.
var (
	// ErrBadTolerance: a tolerance is negative, NaN or infinite.
	ErrBadTolerance = errors.New("float: tolerance must be finite and ≥ 0")
	// ErrLengthMismatch: two slices compared element by element differ in length.
	ErrLengthMismatch = errors.New("float: slices differ in length")
	// ErrNotFinite: an input is NaN or ±Inf where a finite value is required.
	ErrNotFinite = errors.New("float: value is NaN or infinite")
)

// IsClose reports whether a and b are close: |a − b| ≤ max(relTol·max(|a|, |b|), absTol). The test is
// symmetric in a and b (build decision D3), unlike numpy's isclose, which scales by |b| only. NaN is
// close to nothing (not even NaN); +Inf is close only to +Inf and −Inf only to −Inf.
//
// ID:      NUM-FLT-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §4.2.2 (eqs. 21–24, "essentially equal" with the max scale)
// Source:  PEP 485, "A Function for testing approximate equality" (math.isclose), 2015, https://peps.python.org/pep-0485/
// Tags:    floating-point, comparison, tolerance, isclose, approximate-equality
// Units:   relTol dimensionless; absTol in the units of a and b
// Range:   relTol, absTol finite and ≥ 0. Errors: ErrBadTolerance.
// Example: IsClose(1, 1+1e-10, 1e-9, 0) = true; IsClose(1e-12, 0, 1e-9, 0) = false (PEP 485: use absTol near 0)
func IsClose(a, b, relTol, absTol float64) (bool, error) {
	if err := checkTolerances(relTol, absTol); err != nil {
		return false, err
	}
	return isClose(a, b, relTol, absTol), nil
}

func isClose(a, b, relTol, absTol float64) bool {
	if a == b { // equal values, including equal infinities
		return true
	}
	if math.IsNaN(a) || math.IsNaN(b) || math.IsInf(a, 0) || math.IsInf(b, 0) {
		return false
	}
	diff := math.Abs(a - b) // may be +Inf for huge opposite values: then not close, correctly
	return diff <= math.Max(relTol*math.Max(math.Abs(a), math.Abs(b)), absTol)
}

func checkTolerances(relTol, absTol float64) error {
	if math.IsNaN(relTol) || math.IsInf(relTol, 0) || relTol < 0 {
		return fmt.Errorf("%w: relTol = %v", ErrBadTolerance, relTol)
	}
	if math.IsNaN(absTol) || math.IsInf(absTol, 0) || absTol < 0 {
		return fmt.Errorf("%w: absTol = %v", ErrBadTolerance, absTol)
	}
	return nil
}

// AllClose reports whether a[i] and b[i] are close, by IsClose, for every i. Two empty slices are
// all close.
//
// ID:      NUM-FLT-004
// Time:    O(n)   n = len(a); stops at the first pair not close
// Space:   O(1), no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., 1997, §4.2.2 (eqs. 21–24), applied elementwise
// Source:  PEP 485, 2015, https://peps.python.org/pep-0485/
// Tags:    floating-point, comparison, tolerance, allclose, approximate-equality, slice
// Units:   relTol dimensionless; absTol in the units of the elements
// Range:   len(a) == len(b); tolerances finite and ≥ 0. Errors: ErrBadTolerance, ErrLengthMismatch.
// Example: AllClose([1, 2], [1, 2+1e-12], 1e-9, 0) = true
func AllClose(a, b []float64, relTol, absTol float64) (bool, error) {
	if err := checkTolerances(relTol, absTol); err != nil {
		return false, err
	}
	if len(a) != len(b) {
		return false, fmt.Errorf("%w: len(a) = %d, len(b) = %d", ErrLengthMismatch, len(a), len(b))
	}
	for i := range a {
		if !isClose(a[i], b[i], relTol, absTol) {
			return false, nil
		}
	}
	return true, nil
}
