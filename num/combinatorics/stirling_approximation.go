package combinatorics

import (
	"fmt"
	"math"
)

// stirlingApproximationMax is the largest n whose approximation of n! is a finite float64 (171! > MaxFloat64).
const stirlingApproximationMax = 170

// StirlingApproximation returns Stirling's approximation of n!, √(2πn)·(n/e)ⁿ. It is always low: by
// Robbins's bound n! = √(2πn)·(n/e)ⁿ·e^(rₙ) with 1/(12n+1) < rₙ < 1/(12n), so n!/approximation lies in
// (e^(1/(12n+1)), e^(1/(12n))): 1.0844 at n = 1, 1.008365 at n = 10, 1.000834 at n = 100. It exists to
// teach and to bound; for n! itself use Factorial (exact) or LogFactorial.
// Computed as exp(½ln(2πn) + n·ln n − n), so nⁿ never overflows on the way.
// Errors: ErrOutOfRange when n < 1; ErrOverflow when n > 170 (the value exceeds the float64 range).
//
// ID:      NUM-COMBI-016
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Robbins, "A Remark on Stirling's Formula", American Mathematical Monthly, 1955, 62(1):26–29, doi:10.2307/2308012 (the bound 1/(12n+1) < rₙ < 1/(12n))
// Source:  NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §5.6 eq. 5.6.1, §5.11 eq. 5.11.1
// Tags:    combinatorics, factorial, stirling, approximation, education
// Units:   n: a count of objects; result: approximate count of orderings (dimensionless)
// Example: StirlingApproximation(10) = 3 598 695.619, against 10! = 3 628 800 (ratio 1.008365)
func StirlingApproximation(n int) (float64, error) {
	if n < 1 {
		return 0, fmt.Errorf("%w: n must be ≥ 1", ErrOutOfRange)
	}
	if n > stirlingApproximationMax {
		return 0, fmt.Errorf("%w: n", ErrOverflow)
	}
	x := float64(n)
	return math.Exp(0.5*math.Log(2*math.Pi*x) + x*math.Log(x) - x), nil
}
