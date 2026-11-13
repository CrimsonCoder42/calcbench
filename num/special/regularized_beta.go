package special

import (
	"fmt"
	"math"
)

// RegularizedBeta returns the regularized incomplete beta function I_x(a, b) = B_x(a, b)/B(a, b).
// It evaluates the continued fraction DLMF 8.17.22 by the modified Lentz method, on x when
// x < (a+1)/(a+b+2) and otherwise on 1 − x through the symmetry I_x(a,b) = 1 − I_(1−x)(b,a)
// (DLMF 8.17.4), where the fraction converges fast. I_0 = 0 and I_1 = 1 exactly.
//
// ID:      NUM-SPEC-005
// Time:    O(√max(a,b) + log(1/ε))   a, b = shapes, ε = 2⁻⁵² (continued-fraction terms)
// Space:   O(1), no allocation
// Source:  NIST DLMF, Release 1.2.1, §8.17 (eq. 8.17.2 definition, 8.17.4 symmetry, 8.17.22 continued fraction)
// Source:  Press et al., Numerical Recipes, 3rd ed., Cambridge, 2007, §6.4 (eqs. 6.4.5–6.4.6, betacf)
// Tags:    special-function, incomplete-beta, regularized, beta, student-t, binomial, f-distribution
// Units:   dimensionless
// Range:   0 ≤ x ≤ 1; a, b > 0 and finite. Errors: ErrNotFinite, ErrOutOfRange, ErrNoConvergence.
// Example: I_x(a, 1) = x^a (DLMF 8.17.5 family), so I_0.5(3, 1) = 0.125; I_0.5(a, a) = 0.5
func RegularizedBeta(x, a, b float64) (float64, error) {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return 0, fmt.Errorf("%w: x = %v", ErrNotFinite, x)
	}
	if err := checkBetaArgs(a, b); err != nil {
		return 0, err
	}
	if x < 0 || x > 1 {
		return 0, fmt.Errorf("%w: x = %v, want in [0, 1]", ErrOutOfRange, x)
	}
	if x == 0 {
		return 0, nil
	}
	if x == 1 {
		return 1, nil
	}
	if x < (a+1)/(a+b+2) {
		v, err := betaContinuedFraction(x, a, b)
		return v, err
	}
	v, err := betaContinuedFraction(1-x, b, a)
	if err != nil {
		return 0, err
	}
	return 1 - v, nil
}

// betaContinuedFraction returns x^a (1−x)^b / (a B(a,b)) times the DLMF 8.17.22 continued fraction.
func betaContinuedFraction(x, a, b float64) (float64, error) {
	front := math.Exp(a*math.Log(x) + b*math.Log1p(-x) - logBeta(a, b))
	qab := a + b
	qap := a + 1
	qam := a - 1
	c := 1.0
	d := 1 - qab*x/qap
	if math.Abs(d) < tiny {
		d = tiny
	}
	d = 1 / d
	h := d
	limit := iterationCap(math.Max(a, b))
	for m := 1; m <= limit; m++ {
		fm := float64(m)
		m2 := 2 * fm
		// even step d_(2m) of DLMF 8.17.22
		aa := fm * (b - fm) * x / ((qam + m2) * (a + m2))
		d = 1 + aa*d
		if math.Abs(d) < tiny {
			d = tiny
		}
		c = 1 + aa/c
		if math.Abs(c) < tiny {
			c = tiny
		}
		d = 1 / d
		h *= d * c
		// odd step d_(2m+1)
		aa = -(a + fm) * (qab + fm) * x / ((a + m2) * (qap + m2))
		d = 1 + aa*d
		if math.Abs(d) < tiny {
			d = tiny
		}
		c = 1 + aa/c
		if math.Abs(c) < tiny {
			c = tiny
		}
		d = 1 / d
		del := d * c
		h *= del
		if math.Abs(del-1) <= epsilon {
			return math.Min(1, front*h/a), nil
		}
	}
	return 0, fmt.Errorf("%w: continued fraction for I_x(a, b), x = %v, a = %v, b = %v", ErrNoConvergence, x, a, b)
}
