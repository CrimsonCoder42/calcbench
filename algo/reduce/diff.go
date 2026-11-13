package reduce

import (
	"fmt"
	"math"
)

// Diff returns the order-th forward difference of values: order 1 gives out[i] = values[i+1] - values[i],
// and order k applies that k times, so the result has len(values) - order elements. Order 0 returns a
// copy. The differences are computed in the output buffer, which needs room for len(values) elements:
// dst is reused when cap(dst) >= len(values) (see the package comment; dst may be values itself).
// Errors: ErrOutOfRange when order < 0 or order > len(values); ErrNotFinite (with the index) when a value
// is NaN or ±Inf, checked before anything is written; ErrOverflow when a difference exceeds ±MaxFloat64
// (dst may then have been written). On error the returned slice is nil.
//
// ID:      ALG-RED-007
// Time:    O(n·order)   n = len(values)
// Space:   O(1) extra (no allocation when cap(dst) >= n, else one); O(n) for the result buffer
// Source:  Abramowitz, Stegun, Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 1964, §25.1 (forward differences, 25.1.1)
// Tags:    diff, difference, forward-difference, finite-difference, discrete-derivative, dst
func Diff(dst, values []float64, order int) ([]float64, error) {
	n := len(values)
	if order < 0 || order > n {
		return nil, fmt.Errorf("%w: order = %d, want 0..%d", ErrOutOfRange, order, n)
	}
	if err := checkFinite(values); err != nil {
		return nil, err
	}
	out := resize(dst, n)
	copy(out, values) // a no-op when out and values share their first element.
	m := n
	for range order {
		for i := 0; i+1 < m; i++ {
			out[i] = out[i+1] - out[i]
		}
		m--
	}
	out = out[:m]
	// Inf - finite stays Inf and Inf - Inf is NaN, so an overflow at any order is still visible here.
	for i, x := range out {
		if math.IsNaN(x) || math.IsInf(x, 0) {
			return nil, fmt.Errorf("%w: difference %d", ErrOverflow, i)
		}
	}
	return out, nil
}
