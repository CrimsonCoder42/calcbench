package intmath

import "fmt"

// trimLeadingZeros returns digits without leading '0' characters, keeping one digit for zero.
func trimLeadingZeros(digits []byte) []byte {
	i := 0
	for i < len(digits)-1 && digits[i] == '0' {
		i++
	}
	return digits[i:]
}

// checkDigits returns ErrEmpty or ErrInvalidDigit (naming the parameter) unless s is a non-empty string
// of ASCII digits below radix.
func checkDigits(name, s string, radix byte) error {
	if s == "" {
		return fmt.Errorf("%w: %s is empty", ErrEmpty, name)
	}
	for i := range len(s) {
		if s[i] < '0' || s[i]-'0' >= radix {
			return fmt.Errorf("%w: %s[%d] = %q in radix %d", ErrInvalidDigit, name, i, s[i], radix)
		}
	}
	return nil
}

// AddBinaryStrings returns the sum of two non-negative integers written in binary ("1011" + "11" =
// "1110"), of any length, by schoolbook addition with a carry from the least significant digit. Inputs
// may have leading zeros; the result has none (zero is "0").
// Errors: ErrEmpty when a or b is ""; ErrInvalidDigit when a character is not '0' or '1'.
//
// ID:      ALG-INTMTH-014
// Time:    O(m + n)   m = len(a), n = len(b)
// Space:   O(max(m, n)) for the result; one buffer and the string
// Source:  Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §4.3.1 "The Classical Algorithms", Algorithm A (addition of nonnegative integers), radix b = 2
// Tags:    integer, binary, addition, big-integer, digit-string, carry
func AddBinaryStrings(a, b string) (string, error) {
	if err := checkDigits("a", a, 2); err != nil {
		return "", err
	}
	if err := checkDigits("b", b, 2); err != nil {
		return "", err
	}
	out := make([]byte, max(len(a), len(b))+1)
	var carry byte
	i, j := len(a)-1, len(b)-1
	for k := len(out) - 1; k >= 0; k-- {
		sum := carry
		if i >= 0 {
			sum += a[i] - '0'
			i--
		}
		if j >= 0 {
			sum += b[j] - '0'
			j--
		}
		out[k] = '0' + sum&1
		carry = sum >> 1
	}
	return string(trimLeadingZeros(out)), nil
}
