package intmath

import (
	"errors"
	"fmt"
	"math"
	"math/bits"
)

// Errors returned by this package. They are wrapped with the offending parameter; test with errors.Is.
var (
	// ErrOverflow: the result does not fit in its type.
	ErrOverflow = errors.New("intmath: result does not fit the integer type")
	// ErrNegative: a parameter that must be ≥ 0 is negative.
	ErrNegative = errors.New("intmath: parameter is negative")
	// ErrDivideByZero: a divisor or denominator is zero.
	ErrDivideByZero = errors.New("intmath: division by zero")
	// ErrOutOfRange: a parameter is outside the range the function accepts.
	ErrOutOfRange = errors.New("intmath: parameter out of range")
	// ErrTooLong: a decimal expansion has more digits than the caller allowed.
	ErrTooLong = errors.New("intmath: expansion longer than the digit limit")
	// ErrEmpty: a digit string or slice is empty.
	ErrEmpty = errors.New("intmath: no digits")
	// ErrInvalidDigit: a digit is outside the radix.
	ErrInvalidDigit = errors.New("intmath: invalid digit")
)

// PowerBySquaring returns base raised to the power exp, by binary exponentiation: about log2(exp)
// squarings and multiplications instead of exp−1 multiplications. 0^0 is 1 (Knuth's convention).
// Errors: ErrNegative when exp < 0; ErrOverflow when the result is outside the int64 range (checked on
// every product, so nothing wraps). (−2)^63 = math.MinInt64 is in range and is returned.
//
// ID:      ALG-INTMTH-001
// Time:    O(log e)   e = exp; at most 2⌊log2 e⌋+1 multiplications
// Space:   O(1), no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §4.6.3 "Evaluation of Powers", Algorithm A (right-to-left binary method for exponentiation)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §31.6 "Powers of an element" (repeated squaring)
// Tags:    integer, power, exponentiation, binary-exponentiation, repeated-squaring, overflow-checked
func PowerBySquaring(base, exp int64) (int64, error) {
	if exp < 0 {
		return 0, fmt.Errorf("%w: exp = %d", ErrNegative, exp)
	}
	negative := base < 0 && exp&1 == 1
	limit := uint64(math.MaxInt64)
	if negative {
		limit = 1 << 63
	}
	m := uint64(base)
	if base < 0 {
		m = -m
	}
	result := uint64(1)
	for e := exp; e > 0; e >>= 1 {
		if e&1 == 1 {
			hi, lo := bits.Mul64(result, m)
			if hi != 0 || lo > limit {
				return 0, fmt.Errorf("%w: %d^%d", ErrOverflow, base, exp)
			}
			result = lo
		}
		if e > 1 {
			// m² is a factor of the result whenever a set bit of e remains, so its overflow is the result's.
			hi, lo := bits.Mul64(m, m)
			if hi != 0 || lo > limit {
				return 0, fmt.Errorf("%w: %d^%d", ErrOverflow, base, exp)
			}
			m = lo
		}
	}
	if negative {
		return -int64(result), nil // result = 2^63 converts to MinInt64 and its negation is itself
	}
	return int64(result), nil
}
