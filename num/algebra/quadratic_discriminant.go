package algebra

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// QuadraticDiscriminant returns D = b² − 4ac for a·x² + b·x + c, with a relative error of at most 2u
// (u = 2⁻⁵³) whenever D is a normal float64. It uses Kahan's algorithm for a 2×2 determinant with a
// fused multiply-add: w = 4a·c, e = fma(−4a, c, w) (the exact rounding error of w), f = fma(b, b, −w),
// D = f + e. The naive b*b − 4*a*c can lose every digit when b² ≈ 4ac. Inputs are first scaled by
// exact powers of two so that b² and 4ac cannot overflow or underflow on the way. a = 0 is allowed.
// Errors: ErrNotFinite when a, b or c is NaN or ±Inf; ErrOverflow when D itself is outside the float64
// range. D is 0 with any error.
//
// ID:      NUM-ALGEB-006
// Time:    O(1)
// Space:   O(1)
// Source:  Jeannerod, Louvet, Muller, "Further analysis of Kahan's algorithm for the accurate computation of 2×2 determinants", Mathematics of Computation, 2013, 82(284):2245–2264, doi:10.1090/S0025-5718-2013-02679-8
// Source:  Boldo, "Kahan's algorithm for a correct discriminant computation at last formally proven", IEEE Transactions on Computers, 2009, 58(2):220–225, doi:10.1109/TC.2008.200
// Source:  Goldberg, "What every computer scientist should know about floating-point arithmetic", ACM Computing Surveys, 1991, 23(1):5–48, doi:10.1145/103162.103163, §1.4 "Cancellation", p. 10 (the example)
// Tags:    quadratic, discriminant, fma, kahan, algebra
// Units:   a, b, c and D are dimensionless
// Example: a = 1.22, b = 3.34, c = 2.28 → 0.0292 (exact for those decimals; float64 result within 2 ulps
// of the exact discriminant of the rounded inputs)
func QuadraticDiscriminant(a, b, c float64) (float64, error) {
	if err := checkQuadratic(a, b, c); err != nil {
		return 0, err
	}
	d, exp, _ := scaledDiscriminant(a, b, c)
	disc := math.Ldexp(d, exp)
	if !reduce.AllFinite(disc) {
		return 0, fmt.Errorf("%w: b² − 4ac", ErrOverflow)
	}
	return disc, nil
}

// checkQuadratic returns ErrNotFinite, wrapped with the parameter name, for the first NaN or ±Inf among
// a, b, c.
func checkQuadratic(a, b, c float64) error {
	if !reduce.AllFinite(a) {
		return fmt.Errorf("%w: a = %v", ErrNotFinite, a)
	}
	if !reduce.AllFinite(b) {
		return fmt.Errorf("%w: b = %v", ErrNotFinite, b)
	}
	if !reduce.AllFinite(c) {
		return fmt.Errorf("%w: c = %v", ErrNotFinite, c)
	}
	return nil
}

// scaledDiscriminant returns d, exp and bScaled with b² − 4ac = d·2^exp and b = bScaled·2^(exp/2).
// exp is even. a, b, c must be finite. d and bScaled are at most about 4 in magnitude, so √d, bScaled ± √d
// and their quotients by Frexp fractions cannot overflow.
//
// Each of b and 4ac is split as fraction·2^k with math.Frexp (exact). exp is the larger of 2·k(b) and
// k(a)+k(c)+2, rounded up to even; the dominant term then has magnitude in [1/16, 1] and the other is
// scaled down by an exact power of two, underflowing only where it is below 2⁻¹⁰⁷⁰ relative to the
// dominant one (and so below its rounding error). Kahan's FMA determinant is then exact up to the final
// two roundings.
func scaledDiscriminant(a, b, c float64) (d float64, exp int, bScaled float64) {
	fa, ka := math.Frexp(a)
	fb, kb := math.Frexp(b)
	fc, kc := math.Frexp(c)
	hasB := b != 0
	hasAC := a != 0 && c != 0
	if !hasB && !hasAC {
		return 0, 0, 0
	}
	expB := 2 * kb
	expAC := ka + kc + 2 // 4ac = fa·fc·2^(ka+kc+2)
	switch {
	case !hasAC:
		exp = expB
	case !hasB:
		exp = expAC
	default:
		exp = max(expB, expAC)
	}
	if exp&1 != 0 {
		exp++
	}
	if hasB {
		bScaled = math.Ldexp(fb, kb-exp/2)
	}
	var aScaled float64
	if hasAC {
		aScaled = math.Ldexp(fa, expAC-exp) // aScaled·fc = 4ac·2^−exp
	}
	w := aScaled * fc
	e := math.FMA(-aScaled, fc, w)
	f := math.FMA(bScaled, bScaled, -w)
	return f + e, exp, bScaled
}
