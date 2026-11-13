package intmath

// MinSquaresToSum returns the least number of positive perfect squares that sum to n (0 for n = 0).
// The answer is never more than 4 (Lagrange's four-square theorem), and:
//   - 1 when n is a square;
//   - 4 exactly when n = 4^a·(8b + 7) (Legendre's three-square theorem: those and only those need 4);
//   - 2 when n is a sum of two squares, found by two pointers over a ≤ b with a² + b² = m, where m is n
//     with its factors of 4 removed (4m = x² + y² forces x and y even, so n and m agree);
//   - 3 otherwise.
//
// ID:      ALG-INTMTH-011
// Time:    O(√n)   the two-square search over m ≤ n, at most about √m steps (≈4.3e9 for n near 2⁶⁴, seconds of CPU); every other test is O(log n)
// Space:   O(1), no allocation
// Source:  Hardy, Wright, An Introduction to the Theory of Numbers, 6th ed., Oxford University Press, 2008, §20.5, Theorem 369 (Lagrange: every positive integer is a sum of four squares) and §20.10 (numbers 8m + 7 are not sums of three squares, so g(2) = 4)
// Source:  Gauss, Disquisitiones Arithmeticae, 1801, art. 291–292 (Legendre–Gauss three-square theorem: every n not of the form 4^a(8m + 7) is a sum of three squares)
// Tags:    integer, squares, sum-of-squares, lagrange, legendre, four-square-theorem, perfect-squares
func MinSquaresToSum(n uint64) int {
	if n == 0 {
		return 0
	}
	if s := ISqrt(n); s*s == n {
		return 1
	}
	m := n
	for m%4 == 0 {
		m /= 4
	}
	if m%8 == 7 {
		return 4
	}
	if m%4 != 3 { // squares are 0 or 1 mod 4, so m ≡ 3 (mod 4) is never a sum of two squares
		a, b := uint64(0), ISqrt(m)
		for a <= b {
			rest := m - b*b // b ≤ ⌊√m⌋ keeps b² ≤ m: no underflow, and a² ≤ 2^64 − 1
			switch sq := a * a; {
			case sq == rest:
				return 2
			case sq < rest:
				a++
			default:
				b--
			}
		}
	}
	return 3
}
