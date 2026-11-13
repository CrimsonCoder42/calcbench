package stats

import (
	"fmt"
	"math"
	"math/bits"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// compensatedSum is a running sum with Neumaier's correction: the true total of the added terms is
// sum + comp up to the error of comp itself, O(n·u²) for n terms and unit roundoff u. The zero value is
// an empty sum. Each step is algo/reduce.NeumaierAdd (ALG-RED-002).
// Source: Neumaier, ZAMM, 1974, 54(1):39–51, doi:10.1002/zamm.19740540106 (doc/sources/neumaier_1974.md).
type compensatedSum struct {
	sum, comp float64
}

// add adds x and keeps the exact rounding error of that addition in comp.
func (s *compensatedSum) add(x float64) {
	s.sum, s.comp = reduce.NeumaierAdd(s.sum, s.comp, x)
}

// total returns the corrected sum. It is NaN or ±Inf when a partial sum overflowed.
func (s *compensatedSum) total() float64 {
	return s.sum + s.comp
}

// checkFinite returns ErrNotFinite, naming the parameter and the index, for the first NaN or ±Inf.
func checkFinite(name string, values []float64) error {
	for i, v := range values {
		if !reduce.AllFinite(v) {
			return fmt.Errorf("%w: %s[%d] = %v", ErrNotFinite, name, i, v)
		}
	}
	return nil
}

// checkScalar returns ErrNotFinite, naming the parameter, when x is NaN or ±Inf.
func checkScalar(name string, x float64) error {
	if !reduce.AllFinite(x) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, x)
	}
	return nil
}

// meanOf returns the mean of finite, non-empty values with a compensated sum. A sum beyond the float64
// range is redone with every value scaled by 2^-k, 2^k > n (exact for normal numbers), so that the sum
// cannot overflow; the mean of finite values lies between their minimum and maximum, so it is finite
// (a rounding of the result past ±MaxFloat64 is clamped to ±MaxFloat64).
func meanOf(values []float64) float64 {
	n := float64(len(values))
	var s compensatedSum
	for _, v := range values {
		s.add(v)
	}
	if t := s.total(); reduce.AllFinite(t) {
		return t / n
	}
	k := bitLen(len(values))
	scale := math.Ldexp(1, -k)
	s = compensatedSum{}
	for _, v := range values {
		s.add(v * scale)
	}
	m := math.Ldexp(s.total()/n, k)
	if math.IsInf(m, 0) {
		return math.Copysign(math.MaxFloat64, m)
	}
	return m
}

// bitLen returns the number of bits needed to write n ≥ 0, so that 2^bitLen(n) > n.
func bitLen(n int) int {
	return bits.Len(uint(n))
}

// overflowScaleExp is the power of two by which every deviation is scaled down when a deviation itself
// overflows: deviations are below 2^1025, so scaled ones are below 2^425, their squares below 2^850, and
// a sum of up to 2^63 of them below 2^913, inside the float64 range. Values below 2^-422 lose bits on
// that path; they are negligible there, since some deviation exceeds 2^1023.
const overflowScaleExp = 600

// smallDeviation bounds the unscaled pass from below: when the largest |deviation| is at least
// smallDeviation = 2^-480, its square is at least 2^-960, so a square or product that underflows (below
// 2^-1022, absolute error at most 2^-1075 each) changes the sum by under 2^-114 relative per term.
// A smaller largest deviation triggers the rescaled pass.
const smallDeviation = 0x1p-480

// minScaleExp clamps the rescaling exponent from below so that 2^-k stays representable: with the
// largest |deviation| at least 2^-1074 and k ≥ -1000, scaled deviations reach at least 2^-74.
const minScaleExp = -1000

// deviationScaleExp returns k such that the deviations, multiplied by 2^-k, have a largest magnitude in
// [1/2, 1), or of at least 2^-74 when k is clamped to minScaleExp. dMax is the largest |deviation| from
// the unscaled pass; +Inf (a deviation overflowed) gives overflowScaleExp.
func deviationScaleExp(dMax float64) int {
	if !reduce.AllFinite(dMax) {
		return overflowScaleExp
	}
	_, k := math.Frexp(dMax) // dMax = f·2^k, f in [1/2, 1)
	return max(k, minScaleExp)
}

// needsRescale reports whether an unscaled pass with largest |deviation| dMax and sum of squares ss must
// be redone on scaled deviations: the sum overflowed, or the deviations are so small that their squares
// underflow. dMax = 0 (every deviation exactly 0) never needs it.
func needsRescale(ss, dMax float64) bool {
	if dMax == 0 {
		return false
	}
	return !reduce.AllFinite(ss) || dMax < smallDeviation
}

// centredSquares returns ss and e such that Σ(v − mean)² = ss·4^e, by the corrected two-pass formula
// Σd² − (Σd)²/n with d = v − mean, both sums compensated. e is 0 unless the unscaled sum overflows or
// the largest |d| is below smallDeviation (so squares would underflow); the sum is then redone with
// every value and the mean scaled by 2^-e (exact: a power of two), e from deviationScaleExp, so the
// largest scaled |d| is about 1. Hence the result is accurate from subnormal to overflowing spreads.
// ss is 0 only when every deviation is exactly 0. values are finite, non-empty.
func centredSquares(values []float64, mean float64) (ss float64, e int) {
	ss, dMax := correctedSquares(values, mean, 1)
	if !needsRescale(ss, dMax) {
		return ss, 0
	}
	e = deviationScaleExp(dMax)
	ss, _ = correctedSquares(values, mean, math.Ldexp(1, -e))
	return ss, e
}

// correctedSquares is centredSquares' pass at one scale (a power of two). It also returns the largest
// |scaled deviation|. When rounding drives the corrected sum to 0 or below although the deviations differ,
// it returns the lower bound 2·((hi − lo)/2)² of Σd² (two deviations hi − lo apart contribute at least
// that), so a zero sum always means every deviation is equal, that is every value is equal.
func correctedSquares(values []float64, mean, scale float64) (ss, dMax float64) {
	m := mean * scale
	var sq, dev compensatedSum
	lo, hi := math.Inf(1), math.Inf(-1)
	for _, v := range values {
		d := v*scale - m
		sq.add(d * d)
		dev.add(d)
		lo, hi = min(lo, d), max(hi, d)
	}
	c := dev.total()
	ss = sq.total() - c*c/float64(len(values))
	return floorSquares(ss, lo, hi), max(-lo, hi)
}

// floorSquares returns ss, or 2·((hi − lo)/2)² when ss ≤ 0 and hi > lo (the exact sum of squares is
// then positive and at least that), or 0 when ss ≤ 0 and every deviation is equal.
func floorSquares(ss, lo, hi float64) float64 {
	if ss > 0 || math.IsNaN(ss) {
		return ss
	}
	if hi > lo {
		h := hi/2 - lo/2
		return 2 * h * h
	}
	return 0
}

// centredProducts returns sxx, syy, sxy, ex and ey such that Σ(x − xMean)² = sxx·4^ex,
// Σ(y − yMean)² = syy·4^ey and Σ(x − xMean)(y − yMean) = sxy·2^(ex+ey), each by the corrected two-pass
// formula, compensated. ex (ey) is 0 unless the x (y) sum of squares overflows or its largest deviation
// is below smallDeviation; then all three sums are recomputed with that variable and its mean
// scaled by 2^-ex (2^-ey), chosen as in centredSquares. Scaling x and y separately keeps a small spread
// in one from underflowing when the other is huge. Since |sxy| ≤ √(sxx·syy), sxy is finite whenever
// both squares are. sxx (syy) is 0 only when every x (y) is equal. x and y are finite, non-empty, of
// equal length.
func centredProducts(x, y []float64, xMean, yMean float64) (sxx, syy, sxy float64, ex, ey int) {
	sxx, syy, sxy, dxMax, dyMax := correctedProducts(x, y, xMean, yMean, 1, 1)
	rx, ry := needsRescale(sxx, dxMax), needsRescale(syy, dyMax)
	if !rx && !ry && reduce.AllFinite(sxy) {
		return sxx, syy, sxy, 0, 0
	}
	if rx {
		ex = deviationScaleExp(dxMax)
	}
	if ry {
		ey = deviationScaleExp(dyMax)
	}
	sxx, syy, sxy, _, _ = correctedProducts(x, y, xMean, yMean, math.Ldexp(1, -ex), math.Ldexp(1, -ey))
	return sxx, syy, sxy, ex, ey
}

// correctedProducts is centredProducts' pass with x scaled by xScale and y by yScale (powers of two). It
// also returns the largest |scaled deviation| of x and of y.
func correctedProducts(x, y []float64, xMean, yMean, xScale, yScale float64) (sxx, syy, sxy, dxMax, dyMax float64) {
	mx, my := xMean*xScale, yMean*yScale
	var sx2, sy2, sxyc, dx, dy compensatedSum
	xLo, xHi := math.Inf(1), math.Inf(-1)
	yLo, yHi := math.Inf(1), math.Inf(-1)
	for i, xv := range x {
		a := xv*xScale - mx
		b := y[i]*yScale - my
		sx2.add(a * a)
		sy2.add(b * b)
		sxyc.add(a * b)
		dx.add(a)
		dy.add(b)
		xLo, xHi = min(xLo, a), max(xHi, a)
		yLo, yHi = min(yLo, b), max(yHi, b)
	}
	n := float64(len(x))
	cx, cy := dx.total(), dy.total()
	sxx = floorSquares(sx2.total()-cx*cx/n, xLo, xHi)
	syy = floorSquares(sy2.total()-cy*cy/n, yLo, yHi)
	sxy = sxyc.total() - cx*cy/n
	return sxx, syy, sxy, max(-xLo, xHi), max(-yLo, yHi)
}

// midpoint returns (a + b)/2 without overflow for any finite a and b.
func midpoint(a, b float64) float64 {
	if (a >= 0) == (b >= 0) {
		return a + (b-a)/2 // same sign: b − a cannot overflow
	}
	return (a + b) / 2 // opposite signs: a + b cannot overflow
}

// lerp returns a + t·(b − a), the point a fraction t of the way from a to b, for finite a, b and t in
// [0, 1]: exactly a at t = 0 and exactly b at t = 1, and without overflow when b − a would overflow.
func lerp(a, b, t float64) float64 {
	switch t {
	case 0:
		return a
	case 1:
		return b
	}
	if d := b - a; reduce.AllFinite(d) {
		return a + t*d
	}
	return (1-t)*a + t*b // a and b have opposite signs here, so the sum cannot overflow
}
