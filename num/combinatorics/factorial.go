package combinatorics

import (
	"errors"
	"fmt"
)

// Errors returned by this package. They are wrapped with the offending parameter; test with errors.Is.
var (
	// ErrNegative: a count is negative.
	ErrNegative = errors.New("combinatorics: negative count")
	// ErrOverflow: the exact result does not fit the result type (uint64 or float64), or an input sum
	// does not fit an int.
	ErrOverflow = errors.New("combinatorics: result overflows")
	// ErrOutOfRange: the arguments are outside the domain on which the quantity is defined.
	ErrOutOfRange = errors.New("combinatorics: argument out of range")
	// ErrEmpty: a list of counts has no elements.
	ErrEmpty = errors.New("combinatorics: empty input")
)

// factorials holds n! for n = 0…20; 21! = 51 090 942 171 709 440 000 exceeds 2⁶⁴ − 1.
// Values: OEIS A000142, checked in the tests against math/big.
var factorials = [...]uint64{
	1,
	1,
	2,
	6,
	24,
	120,
	720,
	5040,
	40320,
	362880,
	3628800,
	39916800,
	479001600,
	6227020800,
	87178291200,
	1307674368000,
	20922789888000,
	355687428096000,
	6402373705728000,
	121645100408832000,
	2432902008176640000,
}

// Factorial returns n! = 1·2·…·n exactly, with 0! = 1. 20! = 2 432 902 008 176 640 000 is the largest
// that fits a uint64; past it use new(big.Int).MulRange(1, n), or LogFactorial for the log.
// Errors: ErrNegative when n < 0; ErrOverflow when n > 20.
//
// ID:      NUM-COMBI-001
// Time:    O(1)   (a 21-entry table)
// Space:   O(1), no allocation
// Source:  NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §5.4 eq. 5.4.1 (n! = Γ(n+1))
// Source:  OEIS A000142, https://oeis.org/A000142 (read 2026-11-12)
// Tags:    combinatorics, factorial, counting
// Units:   n: a count of objects; result: a count of orderings
// Example: Factorial(10) = 3 628 800 (OEIS A000142)
func Factorial(n int) (uint64, error) {
	if n < 0 {
		return 0, fmt.Errorf("%w: n", ErrNegative)
	}
	if n >= len(factorials) {
		return 0, fmt.Errorf("%w: n", ErrOverflow)
	}
	return factorials[n], nil
}
