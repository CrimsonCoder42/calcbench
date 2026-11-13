package bits

import (
	"errors"
	"fmt"
	mathbits "math/bits"
)

// Errors returned by this package. They are wrapped with the offending parameter; test with errors.Is.
var (
	// ErrOverflow: the result does not fit in a uint64.
	ErrOverflow = errors.New("bits: result does not fit in 64 bits")
	// ErrWidth: a bit width is outside 1…64.
	ErrWidth = errors.New("bits: width must be in 1..64")
	// ErrOutOfRange: a value has set bits above the stated width.
	ErrOutOfRange = errors.New("bits: value has bits above the width")
	// ErrNegative: a count is negative.
	ErrNegative = errors.New("bits: count is negative")
	// ErrTooLarge: a requested table is larger than one allocation may be (rule L41).
	ErrTooLarge = errors.New("bits: table exceeds the allocation limit")
	// ErrRange: lo is greater than hi.
	ErrRange = errors.New("bits: lo is greater than hi")
)

// NextPowerOfTwo returns the least power of two that is ≥ x: 1 for x = 0 or 1, x itself when x is a
// power of two, and 2^⌈log2 x⌉ otherwise.
// Errors: ErrOverflow when x > 2^63, whose next power of two (2^64) is not a uint64.
//
// ID:      ALG-BITS-003
// Time:    O(1)   one leading-zero count
// Space:   O(1), no allocation
// Source:  Warren, Hacker's Delight, 2nd ed., Addison-Wesley, 2013, §3-2 "Rounding Up/Down to the Next Power of 2" (clp2: 2^(w − nlz(x − 1)))
// Tags:    bits, power-of-two, round-up, capacity
func NextPowerOfTwo(x uint64) (uint64, error) {
	if x <= 1 {
		return 1, nil
	}
	if x > 1<<63 {
		return 0, fmt.Errorf("%w: x = %d is above 2^63", ErrOverflow, x)
	}
	return 1 << (64 - mathbits.LeadingZeros64(x-1)), nil
}
