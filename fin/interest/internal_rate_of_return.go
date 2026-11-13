package interest

import (
	"fmt"
	"math"
)

// The root finder's fixed policy (doc/finance_2026-11-12/BRIEF.md, numeric policy, IRR/YTM).
const (
	rootRelTolerance = 1e-12 // stop when |Δrate| ≤ 1e-12·max(1, |rate|)
	rootMaxIter      = 200   // Newton/bisection steps after bracketing
	bracketFirstStep = 0.01  // first trial rate on each side of 0
	bracketMaxRate   = 1e6   // highest trial rate
	bracketMinGap    = 1e-12 // lowest trial rate is −1 + bracketMinGap
)

// rateFunc returns f(rate) and df/drate. The solvers below call it but never keep it, so a closure
// passed as a rateFunc stays on the caller's stack.
type rateFunc func(rate float64) (f, df float64)

// solveRate finds a root of f on (−1, 1e6]. Phase 1 scans outward from rate 0 in lockstep, upward by
// doubling (0.01, 0.02, … up to 1e6) and downward (−0.01, −0.02, … −0.64, then halving the distance
// to −1 down to −1 + 1e-12), and stops at the first interval whose ends differ in sign, the upward one
// first when both sides find one at the same step. Phase 2 refines it with rtsafe (Newton steps,
// bisection when a step would leave the bracket or shrinks too slowly). Returns ErrNoSolution when
// no sign change is found and ErrNoConvergence when 200 steps do not meet the tolerance.
func solveRate(f rateFunc) (float64, error) {
	f0, _ := f(0)
	if f0 == 0 {
		return 0, nil
	}
	if math.IsNaN(f0) {
		return 0, fmt.Errorf("%w: f(0) is NaN", ErrNoSolution)
	}
	upPrev, fUpPrev := 0.0, f0
	downPrev, fDownPrev := 0.0, f0
	upRate := bracketFirstStep
	downDist := bracketFirstStep // distance below 0 while < 0.5, then the gap to −1 is halved
	downGap := 1 - downDist
	upDone, downDone := false, false
	for !upDone || !downDone {
		if !upDone {
			if upRate >= bracketMaxRate {
				upRate, upDone = bracketMaxRate, true
			}
			fu, _ := f(upRate)
			switch {
			case fu == 0:
				return upRate, nil
			case math.IsNaN(fu):
				upDone = true
			case (fu < 0) != (fUpPrev < 0):
				return rtsafe(f, upPrev, fUpPrev, upRate)
			default:
				upPrev, fUpPrev = upRate, fu
				upRate *= 2
			}
		}
		if !downDone {
			if downGap <= bracketMinGap {
				downGap, downDone = bracketMinGap, true
			}
			r := -1 + downGap
			fd, _ := f(r)
			switch {
			case fd == 0:
				return r, nil
			case math.IsNaN(fd):
				downDone = true
			case (fd < 0) != (fDownPrev < 0):
				return rtsafe(f, r, fd, downPrev)
			default:
				downPrev, fDownPrev = r, fd
				if downDist < 0.5 {
					downDist *= 2
					downGap = 1 - downDist
				} else {
					downGap /= 2
				}
			}
		}
	}
	return 0, fmt.Errorf("%w: no sign change on (−1, 1e6]", ErrNoSolution)
}

// rtsafe refines a root of f bracketed by lo and hi (f(lo) = fLo and f(hi) differ in sign): Newton
// steps, replaced by bisection whenever the Newton step would leave the bracket or would not halve
// the step before last (Press et al. 2007, §9.4).
func rtsafe(f rateFunc, lo, fLo, hi float64) (float64, error) {
	xl, xh := lo, hi // f(xl) < 0 < f(xh)
	if fLo > 0 {
		xl, xh = hi, lo
	}
	rts := xl + (xh-xl)/2
	dxOld := math.Abs(xh - xl)
	dx := dxOld
	fv, dfv := f(rts)
	for range rootMaxIter {
		if fv == 0 {
			return rts, nil
		}
		a := (rts-xh)*dfv - fv
		b := (rts-xl)*dfv - fv
		outside := math.IsNaN(a) || math.IsNaN(b) || (a > 0 && b > 0) || (a < 0 && b < 0)
		if outside || math.IsInf(fv, 0) || math.IsNaN(dfv) || math.Abs(2*fv) > math.Abs(dxOld*dfv) {
			dxOld = dx
			dx = (xh - xl) / 2
			rts = xl + dx
		} else {
			dxOld = dx
			dx = fv / dfv
			rts -= dx
		}
		if math.Abs(dx) <= rootRelTolerance*math.Max(1, math.Abs(rts)) {
			return rts, nil
		}
		fv, dfv = f(rts)
		if fv < 0 {
			xl = rts
		} else {
			xh = rts
		}
	}
	return 0, fmt.Errorf("%w: %d steps, last rate %v", ErrNoConvergence, rootMaxIter, rts)
}

// InternalRateOfReturn returns the rate per period r at which Σₜ cashFlowsCents[t]/(1+r)ᵗ = 0, with
// cashFlowsCents[0] at t = 0. cashFlowsCents is read only. The rate is found in two phases: a scan
// outward from 0 (upward by doubling from 0.01 to 1e6, downward toward −1 + 1e-12) for the first
// interval where NPV changes sign, then Newton steps guarded by bisection (rtsafe) to
// |Δr| ≤ 1e-12·max(1, |r|), at most 200 steps. A series whose flows change sign more than once can
// have several IRRs (Descartes' rule of signs); this returns the root bracketed nearest 0 (above 0
// before below at the same scan step). Two roots inside one scan interval cancel and are not seen.
// Errors: ErrDomain if there are fewer than 2 flows; ErrOutOfRange if a flow is outside ±2⁵³ cents;
// ErrNoSolution if the flows are all ≥ 0 or all ≤ 0, or no sign change of NPV is found on (−1, 1e6];
// ErrNoConvergence if the tolerance is not met in 200 steps.
//
// ID:      FIN-INT-039
// Time:    O(n·k)   n = len(cashFlowsCents), k ≤ 72 scan steps + 200 refinement steps
// Space:   O(1), no allocation
// Source:  J. Dahlquist and R. Knight, Principles of Finance, OpenStax, 2022, §16.3 "Internal Rate of Return (IRR) Method" (definition; Eq. 16.3; multiple IRRs), https://openstax.org/books/principles-finance/pages/16-3-internal-rate-of-return-irr-method
// Source:  W. H. Press, S. A. Teukolsky, W. T. Vetterling and B. P. Flannery, Numerical Recipes, 3rd ed., Cambridge University Press, 2007, §9.4 "Newton-Raphson Method Using Derivative" (rtsafe), p. 456 (method)
// Tags:    interest, time-value, irr, internal-rate-of-return, root-finding, capital-budgeting, yield
// Units:   cashFlowsCents in cents; result decimal rate per period
// Example: the Table 16.3 flows → 14.09 % (Principles of Finance §16.3, Table 16.7)
func InternalRateOfReturn(cashFlowsCents []int64) (float64, error) {
	if err := checkCashFlows(cashFlowsCents); err != nil {
		return 0, err
	}
	if len(cashFlowsCents) < 2 {
		return 0, fmt.Errorf("%w: cashFlowsCents has %d flow, want ≥ 2", ErrDomain, len(cashFlowsCents))
	}
	hasPos, hasNeg := false, false
	for _, c := range cashFlowsCents {
		hasPos = hasPos || c > 0
		hasNeg = hasNeg || c < 0
	}
	if !hasPos || !hasNeg {
		return 0, fmt.Errorf("%w: cashFlowsCents needs at least one inflow and one outflow", ErrNoSolution)
	}
	return solveRate(func(rate float64) (float64, float64) {
		// p(v) = Σ cₜ vᵗ and p'(v) by Horner; d NPV/d rate = p'(v)·(−v²).
		v := 1 / (1 + rate)
		p, dp := 0.0, 0.0
		for t := len(cashFlowsCents) - 1; t >= 0; t-- {
			dp = dp*v + p
			p = p*v + float64(cashFlowsCents[t])
		}
		return p, -dp * v * v
	})
}
