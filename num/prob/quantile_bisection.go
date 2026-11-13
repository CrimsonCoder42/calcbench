package prob

import (
	"fmt"
	"math"
)

// bisectionMaxSteps caps quantileByBisection (a defect guard). Halving [0, hi] reaches the smallest
// subnormal in under 1,100 steps; geometric halving of a positive bracket needs under 70 more.
const bisectionMaxSteps = 1300

// cdfSurvival returns P(X ≤ x) and P(X > x) for one distribution, each computed directly so neither tail
// loses its relative accuracy.
type cdfSurvival func(x float64) (cdf, survival float64, err error)

// quantileByBisection returns the x ≥ lo with P(X ≤ x) = u for a continuous distribution with
// P(X ≤ lo) = 0, 0 < u < 1. For u ≤ ½ it solves cdf(x) = u, for u > ½ survival(x) = 1 − u (exact by
// Sterbenz), so the upper tail keeps its accuracy. hi starts at start and doubles until it brackets the
// root (ErrOverflow past MaxFloat64); then the bracket is halved, arithmetically while lo ≤ 0 and
// geometrically (√(lo·hi)) once lo > 0, until hi − lo ≤ 4 ulps of hi.
// Source: Burden, Faires, Numerical Analysis, 9th ed., Brooks/Cole, 2011, §2.1 (bisection).
func quantileByBisection(u, lo, start float64, f cdfSurvival) (float64, error) {
	return quantileByBisectionTail(u, 1-u, lo, start, f)
}

// quantileByBisectionTail is quantileByBisection with the upper tail 1 − u supplied by the caller, for a
// caller that holds it more accurately than 1 − u can be formed (a tiny two-sided tail, say).
func quantileByBisectionTail(u, tail, lo, start float64, f cdfSurvival) (float64, error) {
	upper := u > 0.5
	h := func(x float64) (float64, error) {
		c, s, err := f(x)
		if upper {
			return tail - s, err
		}
		return c - u, err
	}
	hi := start
	for {
		v, err := h(hi)
		if err != nil {
			return 0, err
		}
		if v >= 0 {
			break
		}
		lo, hi = hi, 2*hi
		if math.IsInf(hi, 0) {
			return 0, fmt.Errorf("%w: quantile for u = %v", ErrOverflow, u)
		}
	}
	for range bisectionMaxSteps {
		if hi-lo <= 4*0x1p-52*math.Abs(hi) {
			break
		}
		var mid float64
		if lo > 0 {
			mid = math.Sqrt(lo) * math.Sqrt(hi)
		} else {
			mid = lo + (hi-lo)/2
		}
		if mid <= lo || mid >= hi {
			break
		}
		v, err := h(mid)
		if err != nil {
			return 0, err
		}
		if v < 0 {
			lo = mid
		} else {
			hi = mid
		}
	}
	return hi, nil
}
