package combinatorics

import (
	"fmt"
	"math"
)

// LogFactorial returns ln n!, the natural log of n factorial, finite for every int (ln((2⁶³−1)!) ≈ 3.9e20).
// For n ≤ 20 it is the log of the exact factorial; above, ln Γ(n+1).
// Errors: ErrNegative when n < 0.
//
// ID:      NUM-COMBI-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §5.4 eq. 5.4.1 (n! = Γ(n+1)), §5.11 eq. 5.11.1
// Tags:    combinatorics, factorial, log-space, gamma
// Units:   n: a count of objects; result: natural log of a count
// Example: LogFactorial(100) = 363.739375555563 (ln of the exact 100!)
func LogFactorial(n int) (float64, error) {
	if n < 0 {
		return 0, fmt.Errorf("%w: n", ErrNegative)
	}
	if n < len(factorials) {
		return math.Log(float64(factorials[n])), nil
	}
	lg, _ := math.Lgamma(float64(n) + 1) // the sign is +1 for every positive argument
	return lg, nil
}
