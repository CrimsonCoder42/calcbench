package intmath

// MultiplyDecimalStrings returns the product of two non-negative integers written in decimal ("123" ×
// "45" = "5535"), of any length, by schoolbook long multiplication: each digit of b times all of a,
// added into the partial product with one carry chain per row. Inputs may have leading zeros; the
// result has none (zero is "0").
// Errors: ErrEmpty when a or b is ""; ErrInvalidDigit when a character is not '0'…'9'.
//
// ID:      ALG-INTMTH-015
// Time:    O(m·n)   m = len(a), n = len(b)
// Space:   O(m + n) for the result; one buffer and the string
// Source:  Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §4.3.1 "The Classical Algorithms", Algorithm M (multiplication of nonnegative integers), radix b = 10
// Tags:    integer, decimal, multiplication, big-integer, digit-string, schoolbook, long-multiplication
func MultiplyDecimalStrings(a, b string) (string, error) {
	if err := checkDigits("a", a, 10); err != nil {
		return "", err
	}
	if err := checkDigits("b", b, 10); err != nil {
		return "", err
	}
	m, n := len(a), len(b)
	w := make([]byte, m+n) // w[k] holds the digit of weight 10^(m+n-1-k): most significant first, like the inputs
	for j := n - 1; j >= 0; j-- {
		v := uint(b[j] - '0')
		if v == 0 {
			continue // Algorithm M step M2: a zero multiplier digit adds nothing
		}
		var carry uint
		for i := m - 1; i >= 0; i-- {
			t := uint(a[i]-'0')*v + uint(w[i+j+1]) + carry // ≤ 9·9 + 9 + 9 = 99
			w[i+j+1] = byte(t % 10)
			carry = t / 10
		}
		w[j] = byte(carry) // the position above row j is still zero, as in step M6
	}
	for k := range w {
		w[k] += '0'
	}
	return string(trimLeadingZeros(w)), nil
}
