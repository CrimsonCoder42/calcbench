package integer

import (
	"errors"
	"fmt"
	"math"
	"math/bits"
)

// ErrZeroDenominator: a fraction has denominator 0, or a fraction is divided by 0/d.
var ErrZeroDenominator = errors.New("integer: zero denominator")

// absInt64 returns |x| as a uint64; x = MinInt64 is rejected by every caller before this is reached.
func absInt64(x int64) uint64 {
	if x < 0 {
		return uint64(-x)
	}
	return uint64(x)
}

// fromMagnitude returns ±m as an int64, or false when m > MaxInt64 (MinInt64 is never produced, so that
// every result can be negated).
func fromMagnitude(m uint64, negative bool) (int64, bool) {
	if m > math.MaxInt64 {
		return 0, false
	}
	if negative {
		return -int64(m), true
	}
	return int64(m), true
}

// mulInt64 returns x·y, or false when the exact product is outside ±(2⁶³ − 1).
func mulInt64(x, y int64) (int64, bool) {
	hi, lo := bits.Mul64(absInt64(x), absInt64(y))
	if hi != 0 {
		return 0, false
	}
	return fromMagnitude(lo, (x < 0) != (y < 0))
}

// addInt64 returns x + y, or false when the exact sum is outside ±(2⁶³ − 1).
func addInt64(x, y int64) (int64, bool) {
	s := x + y
	if (x > 0 && y > 0 && s < 0) || (x < 0 && y < 0 && s >= 0) || s == math.MinInt64 {
		return 0, false
	}
	return s, true
}

// checkFraction rejects a zero denominator and the value MinInt64, whose magnitude has no int64.
func checkFraction(name string, num, den int64) error {
	if den == 0 {
		return fmt.Errorf("%w: %s = %d/0", ErrZeroDenominator, name, num)
	}
	if num == math.MinInt64 || den == math.MinInt64 {
		return fmt.Errorf("%w: %s = %d/%d, |MinInt64| does not fit int64", ErrOverflow, name, num, den)
	}
	return nil
}

// reduce returns num/den in lowest terms with a positive denominator; den ≠ 0 and neither is MinInt64.
func reduce(num, den int64) (int64, int64) {
	g := int64(GCD(absInt64(num), absInt64(den)))
	num, den = num/g, den/g
	if den < 0 {
		num, den = -num, -den
	}
	return num, den
}

// FractionReduce returns num/den in lowest terms with a positive denominator: gcd(|num|, den) = 1, and 0
// is returned as 0/1.
// Errors: ErrZeroDenominator when den = 0; ErrOverflow when num or den is math.MinInt64.
//
// ID:      NUM-INTEG-007
// Time:    O(log |num| + log |den|) (one binary gcd)
// Space:   O(1) (no allocation)
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1998, §4.5.1 "Fractions" (u/u′ kept in lowest terms with u′ > 0)
// Tags:    fraction, rational, lowest-terms, simplify, gcd, no-alloc
// Units:   dimensionless integers; the result is (numerator, denominator)
// Example: 21/98 → 3/14; −6/−4 → 3/2
func FractionReduce(num, den int64) (int64, int64, error) {
	if err := checkFraction("fraction", num, den); err != nil {
		return 0, 0, err
	}
	n, d := reduce(num, den)
	return n, d, nil
}

// FractionAdd returns num1/den1 + num2/den2 in lowest terms with a positive denominator. It reduces both
// inputs, then forms the sum by Knuth's method, which divides by gcd(den1, den2) before multiplying so
// intermediate values stay as small as the result allows.
// Errors: ErrZeroDenominator when a denominator is 0; ErrOverflow when an input is math.MinInt64 or the
// reduced sum does not fit int64.
//
// ID:      NUM-INTEG-008
// Time:    O(log M)   M = the largest magnitude among the inputs (a constant number of binary gcds)
// Space:   O(1) (no allocation)
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1998, §4.5.1, eq. (3) (sum with d₁ = gcd(u′, v′), t = u(v′/d₁) + v(u′/d₁), d₂ = gcd(t, d₁))
// Tags:    fraction, rational, add, sum, lowest-terms, overflow-checked, no-alloc
// Units:   dimensionless integers; the result is (numerator, denominator)
// Example: 2/7 + 3/8 = 37/56
func FractionAdd(num1, den1, num2, den2 int64) (int64, int64, error) {
	if err := checkFraction("first", num1, den1); err != nil {
		return 0, 0, err
	}
	if err := checkFraction("second", num2, den2); err != nil {
		return 0, 0, err
	}
	u, up := reduce(num1, den1)
	v, vp := reduce(num2, den2)
	d1 := int64(GCD(uint64(up), uint64(vp)))
	a, ok1 := mulInt64(u, vp/d1)
	b, ok2 := mulInt64(v, up/d1)
	t, ok3 := addInt64(a, b)
	if !ok1 || !ok2 || !ok3 {
		return 0, 0, fmt.Errorf("%w: %d/%d + %d/%d", ErrOverflow, num1, den1, num2, den2)
	}
	if t == 0 {
		return 0, 1, nil
	}
	d2 := int64(GCD(absInt64(t), uint64(d1)))
	den, ok := mulInt64(up/d1, vp/d2)
	if !ok {
		return 0, 0, fmt.Errorf("%w: %d/%d + %d/%d", ErrOverflow, num1, den1, num2, den2)
	}
	return t / d2, den, nil
}

// FractionSub returns num1/den1 − num2/den2 in lowest terms with a positive denominator (FractionAdd of
// the negated second fraction).
// Errors: as FractionAdd.
//
// ID:      NUM-INTEG-009
// Time:    O(log M)   M = the largest magnitude among the inputs
// Space:   O(1) (no allocation)
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1998, §4.5.1, eq. (3) (subtraction is addition of −v/v′)
// Tags:    fraction, rational, subtract, difference, lowest-terms, overflow-checked, no-alloc
// Units:   dimensionless integers; the result is (numerator, denominator)
// Example: 2/7 − 3/8 = −5/56
func FractionSub(num1, den1, num2, den2 int64) (int64, int64, error) {
	if err := checkFraction("second", num2, den2); err != nil {
		return 0, 0, err
	}
	return FractionAdd(num1, den1, -num2, den2)
}

// FractionMul returns (num1/den1)·(num2/den2) in lowest terms with a positive denominator, cancelling
// gcd(num1, den2) and gcd(num2, den1) before multiplying.
// Errors: ErrZeroDenominator when a denominator is 0; ErrOverflow when an input is math.MinInt64 or the
// reduced product does not fit int64.
//
// ID:      NUM-INTEG-010
// Time:    O(log M)   M = the largest magnitude among the inputs
// Space:   O(1) (no allocation)
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1998, §4.5.1, eq. (2) (product (u/d₁)(v/d₂) / ((u′/d₂)(v′/d₁)), d₁ = gcd(u, v′), d₂ = gcd(u′, v))
// Tags:    fraction, rational, multiply, product, lowest-terms, overflow-checked, no-alloc
// Units:   dimensionless integers; the result is (numerator, denominator)
// Example: 2/7 × 3/8 = 3/28
func FractionMul(num1, den1, num2, den2 int64) (int64, int64, error) {
	if err := checkFraction("first", num1, den1); err != nil {
		return 0, 0, err
	}
	if err := checkFraction("second", num2, den2); err != nil {
		return 0, 0, err
	}
	u, up := reduce(num1, den1)
	v, vp := reduce(num2, den2)
	if u == 0 || v == 0 {
		return 0, 1, nil
	}
	d1 := int64(GCD(absInt64(u), uint64(vp)))
	d2 := int64(GCD(uint64(up), absInt64(v)))
	num, ok1 := mulInt64(u/d1, v/d2)
	den, ok2 := mulInt64(up/d2, vp/d1)
	if !ok1 || !ok2 {
		return 0, 0, fmt.Errorf("%w: %d/%d × %d/%d", ErrOverflow, num1, den1, num2, den2)
	}
	return num, den, nil
}

// FractionDiv returns (num1/den1) ÷ (num2/den2) in lowest terms with a positive denominator: the product
// with the reciprocal den2/num2.
// Errors: ErrZeroDenominator when a denominator is 0 or the divisor num2/den2 is 0; ErrOverflow as
// FractionMul.
//
// ID:      NUM-INTEG-011
// Time:    O(log M)   M = the largest magnitude among the inputs
// Space:   O(1) (no allocation)
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1998, §4.5.1, eq. (2) (division is multiplication by v′/v, v ≠ 0)
// Tags:    fraction, rational, divide, quotient, reciprocal, lowest-terms, overflow-checked, no-alloc
// Units:   dimensionless integers; the result is (numerator, denominator)
// Example: 2/7 ÷ 3/8 = 16/21
func FractionDiv(num1, den1, num2, den2 int64) (int64, int64, error) {
	if err := checkFraction("second", num2, den2); err != nil {
		return 0, 0, err
	}
	if num2 == 0 {
		return 0, 0, fmt.Errorf("%w: divisor %d/%d is zero", ErrZeroDenominator, num2, den2)
	}
	return FractionMul(num1, den1, den2, num2)
}
