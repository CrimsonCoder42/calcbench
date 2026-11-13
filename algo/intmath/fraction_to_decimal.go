package intmath

import (
	"fmt"
	"math/bits"
	"slices"
	"strconv"
)

// FractionToDecimal returns numerator/denominator written in decimal, with the repeating block (the
// repetend) in parentheses: 1/2 → "0.5", 2/1 → "2", 1/3 → "0.(3)", −1/6 → "-0.1(6)", 22/7 →
// "3.(142857)". Long division produces one fractional digit per step from the remainder r ← 10·r mod
// denominator; the expansion repeats from the first remainder that occurs twice (there are fewer than
// |denominator| possible remainders), and stops when a remainder is 0. The fraction need not be in
// lowest terms; the result is the same. There is no "-0": a zero numerator gives "0".
// maxFractionDigits caps the digits after the point (non-repeating plus one period), because a period
// can be as long as |denominator| − 1 digits.
// Errors: ErrDivideByZero when denominator = 0; ErrNegative when maxFractionDigits < 0; ErrTooLong when
// the expansion needs more than maxFractionDigits fractional digits.
//
// ID:      ALG-INTMTH-013
// Time:    O(L)   L = digits after the point (≤ maxFractionDigits); expected O(1) per digit (one hash-map lookup)
// Space:   O(L) extra (the remainder table), O(L + 20) for the result
// Source:  Hardy, Wright, An Introduction to the Theory of Numbers, 6th ed., Oxford University Press, 2008, §9.2 "Terminating and recurring decimals" (a rational number has a terminating or recurring decimal; the period follows the order of 10 modulo the denominator)
// Source:  Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §4.4 "Radix Conversion", Method (2a) (fraction digits by repeated multiplication by the radix)
// Tags:    integer, fraction, decimal, repeating-decimal, recurring-decimal, long-division, cycle-detection
func FractionToDecimal(numerator, denominator int64, maxFractionDigits int) (string, error) {
	if denominator == 0 {
		return "", fmt.Errorf("%w: denominator = 0", ErrDivideByZero)
	}
	if maxFractionDigits < 0 {
		return "", fmt.Errorf("%w: maxFractionDigits = %d", ErrNegative, maxFractionDigits)
	}
	num, den := magnitude(numerator), magnitude(denominator)
	buf := make([]byte, 0, 24)
	if numerator != 0 && (numerator < 0) != (denominator < 0) {
		buf = append(buf, '-')
	}
	buf = strconv.AppendUint(buf, num/den, 10)
	rem := num % den
	if rem == 0 {
		return string(buf), nil
	}
	buf = append(buf, '.')
	seen := make(map[uint64]int) // remainder → index in buf of the digit it produces
	for digits := 0; rem != 0; digits++ {
		if at, ok := seen[rem]; ok {
			buf = slices.Insert(buf, at, '(')
			buf = append(buf, ')')
			break
		}
		if digits == maxFractionDigits {
			return "", fmt.Errorf("%w: %d/%d needs more than %d fractional digits", ErrTooLong, numerator, denominator, maxFractionDigits)
		}
		seen[rem] = len(buf)
		hi, lo := bits.Mul64(rem, 10) // rem < den ≤ 2^63, so 10·rem may exceed 64 bits but hi < den
		var d uint64
		d, rem = bits.Div64(hi, lo, den)
		buf = append(buf, byte('0'+d))
	}
	return string(buf), nil
}
