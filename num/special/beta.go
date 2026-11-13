package special

import (
	"fmt"
	"math"
)

// Beta returns the beta function B(a, b) = Γ(a)Γ(b)/Γ(a+b) for a, b > 0. When both are ≥ 1 and
// a+b < 171 the gamma values are formed directly (no overflow: Γ(a)/Γ(a+b) ≤ 1 there); otherwise
// B = exp(ln B) from LogBeta.
//
// ID:      NUM-SPEC-003
// Time:    O(1)   three gamma evaluations
// Space:   O(1), no allocation
// Source:  NIST DLMF, Release 1.2.1, §5.12 (eq. 5.12.1 B(a,b) = Γ(a)Γ(b)/Γ(a+b))
// Source:  Press et al., Numerical Recipes, 3rd ed., Cambridge, 2007, §6.1 (eq. 6.1.8)
// Tags:    special-function, beta, gamma
// Units:   dimensionless
// Range:   a, b > 0 and finite. Errors: ErrNotFinite, ErrOutOfRange, ErrOverflow (B > MaxFloat64, for a or b near 0).
// Example: B(2, 3) = 1!·2!/4! = 1/12; B(½, ½) = π (DLMF 5.12.1 with 5.4.6)
func Beta(a, b float64) (float64, error) {
	if err := checkBetaArgs(a, b); err != nil {
		return 0, err
	}
	if a >= 1 && b >= 1 && a+b < 171 {
		if a < b {
			a, b = b, a
		}
		return math.Gamma(a) / math.Gamma(a+b) * math.Gamma(b), nil
	}
	lb := logBeta(a, b)
	if lb > math.Log(math.MaxFloat64) {
		return 0, fmt.Errorf("%w: B(%v, %v)", ErrOverflow, a, b)
	}
	return math.Exp(lb), nil
}

// LogBeta returns ln B(a, b) = ln Γ(a) + ln Γ(b) − ln Γ(a+b) for a, b > 0. The three log-gamma values
// cancel when a and b are both large, so the absolute error is about ε·ln Γ(a+b).
//
// ID:      NUM-SPEC-004
// Time:    O(1)   three log-gamma evaluations
// Space:   O(1), no allocation
// Source:  NIST DLMF, Release 1.2.1, §5.12 (eq. 5.12.1, taken in logarithms)
// Source:  Press et al., Numerical Recipes, 3rd ed., Cambridge, 2007, §6.1
// Tags:    special-function, beta, log-beta, gamma, log-gamma
// Units:   dimensionless
// Range:   a, b > 0 and finite. Errors: ErrNotFinite, ErrOutOfRange.
// Example: ln B(2, 3) = −ln 12 = −2.4849066497880004
func LogBeta(a, b float64) (float64, error) {
	if err := checkBetaArgs(a, b); err != nil {
		return 0, err
	}
	return logBeta(a, b), nil
}

func logBeta(a, b float64) float64 {
	la, _ := math.Lgamma(a)
	lb, _ := math.Lgamma(b)
	lab, _ := math.Lgamma(a + b)
	return la + lb - lab
}

func checkBetaArgs(a, b float64) error {
	if math.IsNaN(a) || math.IsInf(a, 0) {
		return fmt.Errorf("%w: a = %v", ErrNotFinite, a)
	}
	if math.IsNaN(b) || math.IsInf(b, 0) {
		return fmt.Errorf("%w: b = %v", ErrNotFinite, b)
	}
	if a <= 0 {
		return fmt.Errorf("%w: a = %v, want > 0", ErrOutOfRange, a)
	}
	if b <= 0 {
		return fmt.Errorf("%w: b = %v, want > 0", ErrOutOfRange, b)
	}
	return nil
}
