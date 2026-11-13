package intmath

// IsPalindromeNumber reports whether the decimal representation of x reads the same forwards and
// backwards (121, 7, 0, 1221). Negative numbers are not palindromes (the sign has no mirror). Only the
// low half of the digits is reversed and compared with the high half, so nothing can overflow.
//
// ID:      ALG-INTMTH-005
// Time:    O(d)   d = number of decimal digits of x (≤ 20); d/2 divisions
// Space:   O(1), no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §4.4 "Radix Conversion", Method (1a) (digits by repeated division by 10) and Method (1b) (Horner's rule to rebuild the reversed half)
// Tags:    integer, digits, palindrome, decimal
func IsPalindromeNumber[T integer](x T) bool {
	if x < 0 {
		return false
	}
	u := uint64(x)
	if u != 0 && u%10 == 0 {
		return false // a leading zero cannot mirror the trailing zero
	}
	var rev uint64
	for u > rev {
		rev = rev*10 + u%10
		u /= 10
	}
	return u == rev || u == rev/10 // rev/10 drops the middle digit of an odd-length number
}
