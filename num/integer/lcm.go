package integer

import (
	"errors"
	"fmt"
	"math/bits"
)

// ErrOverflow: the exact result does not fit the result type.
var ErrOverflow = errors.New("integer: result overflows")

// LCM returns the least common multiple of a and b, computed as (a / GCD(a, b)) · b so that the only
// product formed is the result itself. LCM(0, b) = LCM(a, 0) = 0.
// Errors: ErrOverflow when the least common multiple exceeds 2⁶⁴ − 1.
//
// ID:      NUM-INTEG-005
// Time:    O(log a + log b)
// Space:   O(1) (no allocation)
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1998, §4.5.2 (gcd(u, v)·lcm(u, v) = u·v)
// Tags:    lcm, least-common-multiple, gcd, number-theory, overflow-checked, no-alloc
func LCM(a, b uint64) (uint64, error) {
	if a == 0 || b == 0 {
		return 0, nil
	}
	hi, lo := bits.Mul64(a/GCD(a, b), b)
	if hi != 0 {
		return 0, fmt.Errorf("%w: lcm(%d, %d) > 2^64 - 1", ErrOverflow, a, b)
	}
	return lo, nil
}
