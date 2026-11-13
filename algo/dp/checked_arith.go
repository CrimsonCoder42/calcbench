package dp

import (
	"errors"
	"math"
	"math/bits"
)

// Sentinel errors of the package. Each is wrapped with the offending parameter (fmt.Errorf("%w: …")),
// so callers test with errors.Is.
var (
	// ErrNegative: a value that must be zero or more is negative.
	ErrNegative = errors.New("dp: value must not be negative")
	// ErrNotPositive: a value that must be one or more is zero or negative.
	ErrNotPositive = errors.New("dp: value must be positive")
	// ErrLengthMismatch: two slices that pair up element by element differ in length.
	ErrLengthMismatch = errors.New("dp: slices differ in length")
	// ErrOverflow: a result, or a partial sum the recurrence must hold, does not fit its integer type.
	ErrOverflow = errors.New("dp: result overflows its integer type")
	// ErrTooLarge: the table the recurrence needs would exceed the package's 128 MiB cap (maxTableBytes).
	ErrTooLarge = errors.New("dp: table would exceed the 128 MiB cap")
	// ErrNotRectangular: the rows of a grid differ in length.
	ErrNotRectangular = errors.New("dp: grid rows differ in length")
	// ErrNotTriangle: row i of a triangle does not have i+1 entries.
	ErrNotTriangle = errors.New("dp: row i of the triangle must have i+1 entries")
	// ErrEmpty: the input has no element where the function needs at least one.
	ErrEmpty = errors.New("dp: input is empty")
	// ErrDuplicate: a set given as a slice repeats a value.
	ErrDuplicate = errors.New("dp: duplicate value")
	// ErrUnbounded: the optimum is unbounded (an item of zero weight and positive value).
	ErrUnbounded = errors.New("dp: the optimum is unbounded")
)

// maxTableBytes is the largest table (beyond the input and an O(input) result) one call allocates:
// 128 MiB, the same policy as algo/sorting's CountingSort cap (lead decision L13). A size whose table
// would be larger returns ErrTooLarge before any allocation; it also keeps make() from panicking on a
// length a slice cannot hold (lead decision L26).
const maxTableBytes = 1 << 27

// tableFits reports whether rows×cols cells of cellBytes bytes each fit under maxTableBytes.
// Negative or zero dimensions fit (nothing is allocated for them).
func tableFits(rows, cols, cellBytes int) bool {
	if rows <= 0 || cols <= 0 {
		return true
	}
	hi, cells := bits.Mul64(uint64(rows), uint64(cols))
	if hi != 0 {
		return false
	}
	hi, size := bits.Mul64(cells, uint64(cellBytes))
	return hi == 0 && size <= maxTableBytes
}

// addInt64 returns a+b and whether it fits in int64.
func addInt64(a, b int64) (int64, bool) {
	s := a + b
	// Overflow iff a and b have the same sign and s has the other.
	return s, (a >= 0) != (b >= 0) || (s >= 0) == (a >= 0)
}

// addUint64 returns a+b and whether it fits in uint64.
func addUint64(a, b uint64) (uint64, bool) {
	s, carry := bits.Add64(a, b, 0)
	return s, carry == 0
}

// mulPositiveInt64 returns a·b for a, b > 0 and whether it fits in int64.
func mulPositiveInt64(a, b int64) (int64, bool) {
	hi, lo := bits.Mul64(uint64(a), uint64(b))
	return int64(lo), hi == 0 && lo <= math.MaxInt64
}
