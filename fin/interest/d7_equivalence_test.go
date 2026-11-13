package interest

import (
	"math"
	"testing"
)

// D7 (doc/numpy_build_2026-11-13/BRIEF.md): npvFloat now evaluates the NPV polynomial in v = 1/(1+r) with
// NUM-ALGEB-001 algebra.PolynomialEvaluate. oldNpvFloat is the implementation it replaced, copied verbatim
// (the original is in _delete/fin/interest/net_present_value.go). It is a test-only oracle.
//
// The new path is not bit-identical: PolynomialEvaluate does each Horner step as one fused multiply-add
// (one rounding) where the old loop rounded the product and the sum separately. Both are backward stable
// Horner (Higham §5.1); they differ by at most about 2n·u·Σ|cₜ|vᵗ (u = 2⁻⁵³). So the 1e-12 relative
// requirement is asserted where the polynomial does not cancel (|old| ≥ 0.1·Σ|cₜ|vᵗ, n ≤ 200, which
// keeps 4n·u·Σ|cₜ|vᵗ ≤ 9e-13·|old|); where the flows cancel (an NPV near 0, as at an IRR) the assertion is
// the Horner bound itself, since there no evaluation order is closer to the exact value than that.
func oldNpvFloat(rate float64, cashFlowsCents []int64) float64 {
	v := 1 / (1 + rate)
	acc := 0.0
	for t := len(cashFlowsCents) - 1; t >= 0; t-- {
		acc = acc*v + float64(cashFlowsCents[t])
	}
	return acc
}

// fmaHorner is one uninterrupted FMA Horner pass: the chunked npvFloat must equal it bit for bit.
func fmaHorner(rate float64, cashFlowsCents []int64) float64 {
	v := 1 / (1 + rate)
	acc := 0.0
	for t := len(cashFlowsCents) - 1; t >= 0; t-- {
		acc = math.FMA(acc, v, float64(cashFlowsCents[t]))
	}
	return acc
}

func d7FlowGrid() [][]int64 {
	grid := [][]int64{
		{-3000000, 600000, 800000, 1000000, 1200000, 1400000, 1600000, 1800000},
		{-1600000, 200000, 400000, 500000, 500000, 500000, 500000},
		{12345}, {-100, 40, 40, 40}, {-1000, 1100}, {0, 100}, {0, 0, 0}, {maxExactCents},
		{-1000, 1000}, {-1000, 500}, {1000, -1100}, {-100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1000000000},
		{-1000000, 1}, {-35000000000, 95000000000, -62000000000}, {-100, 300, -300},
		{maxExactCents, -maxExactCents, maxExactCents},
	}
	// Annuities and mixed series across the chunk boundaries (64, 65, 128, 129) of npvFloat.
	for _, n := range []int{1, 2, 63, 64, 65, 66, 127, 128, 129, 130, 200} {
		level := make([]int64, n)
		mixed := make([]int64, n)
		for t := range n {
			level[t] = 250000
			mixed[t] = int64((t*7919)%20001) - 7000
		}
		mixed[0] = -5000000
		grid = append(grid, level, mixed)
	}
	return grid
}

func TestD7_npvFloatMatchesOld(t *testing.T) {
	const u = 0x1p-53
	rates := []float64{-0.5, -0.1, 0, 1e-9, 0.01, 0.05, 0.09, 0.1, 0.5, 1, 10}
	strict := 0
	for _, flows := range d7FlowGrid() {
		for _, r := range rates {
			oldV, newV := oldNpvFloat(r, flows), npvFloat(r, flows)
			if fh := fmaHorner(r, flows); newV != fh && !(math.IsNaN(newV) && math.IsInf(fh, 0)) {
				t.Errorf("n=%d r=%v: chunked %v != one FMA Horner pass %v", len(flows), r, newV, fh)
			}
			if math.IsInf(oldV, 0) || math.IsNaN(oldV) {
				if !math.IsNaN(newV) {
					t.Errorf("n=%d r=%v: old overflowed (%v), new %v; want NaN (NetPresentValue: ErrOverflow both)", len(flows), r, oldV, newV)
				}
				continue
			}
			v := 1 / (1 + r)
			s, p := 0.0, 1.0
			for _, c := range flows {
				s += math.Abs(float64(c)) * p
				p *= v
			}
			diff := math.Abs(newV - oldV)
			if bound := 4 * float64(len(flows)+1) * u * s; diff > bound && diff > 1e-12*math.Max(1, math.Abs(oldV)) {
				t.Errorf("n=%d r=%v: |new−old| = %v exceeds the Horner bound %v (old %v, new %v)", len(flows), r, diff, bound, oldV, newV)
			}
			if math.Abs(oldV) >= 0.1*s && len(flows) <= 200 {
				strict++
				if diff > 1e-12*math.Max(1, math.Abs(oldV)) {
					t.Errorf("n=%d r=%v: |new−old| = %v > 1e-12·max(1,|old|) (old %v, new %v)", len(flows), r, diff, oldV, newV)
				}
			}
			// The public result: the same cents unless old sits within 1e-6 cent of a half cent.
			if frac := math.Abs(oldV - math.Trunc(oldV)); math.Abs(frac-0.5) > 1e-6 {
				oc, oerr := roundCents(oldV)
				nc, nerr := NetPresentValue(r, flows)
				if (oerr == nil) != (nerr == nil) || oc != nc {
					t.Errorf("n=%d r=%v: NetPresentValue %d (%v), old %d (%v)", len(flows), r, nc, nerr, oc, oerr)
				}
			}
		}
	}
	if strict < 100 {
		t.Fatalf("only %d grid points under the strict 1e-12 check; the grid is too weak", strict)
	}
}
