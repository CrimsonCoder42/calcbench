package trig

import (
	"math"
	"math/big"
	"testing"
)

// Independent 256-bit oracles for the near-π tests of the triangle solvers. sin, cos and atan are summed
// from their Maclaurin series (DLMF 4.19.1, 4.19.2, 4.24.3) in big.Float, with π from bigPi; nothing here
// calls the code under test or math.Sin/Cos/Atan.

const bigPrec = 256

func bf(x float64) *big.Float { return new(big.Float).SetPrec(bigPrec).SetFloat64(x) }
func bnew() *big.Float        { return new(big.Float).SetPrec(bigPrec) }

// bigSin is sin x for |x| ≤ π, reduced to [0, π/2] by sin(−x) = −sin x and sin x = sin(π − x).
func bigSin(x *big.Float) *big.Float {
	if x.Sign() < 0 {
		return bnew().Neg(bigSin(bnew().Neg(x)))
	}
	halfPi := bnew().Quo(bigPi, bf(2))
	if x.Cmp(halfPi) > 0 {
		x = bnew().Sub(bigPi, x)
	}
	return bigSeries(x, bnew().Set(x), 1)
}

// bigCos is cos x = sin(π/2 − x) for 0 ≤ x ≤ π.
func bigCos(x *big.Float) *big.Float {
	return bigSin(bnew().Sub(bnew().Quo(bigPi, bf(2)), x))
}

// bigSeries sums term_k = (−1)^k x^(2k+start)/(2k+start)!, starting from term0.
func bigSeries(x, term0 *big.Float, start int) *big.Float {
	x2 := bnew().Mul(x, x)
	sum, term := bnew().Set(term0), bnew().Set(term0)
	eps := bnew().SetMantExp(bf(1), -bigPrec)
	for n := start; ; n += 2 {
		term.Mul(term, x2)
		term.Quo(term, bf(float64((n+1)*(n+2))))
		term.Neg(term)
		sum.Add(sum, term)
		if bnew().Abs(term).Cmp(bnew().Mul(eps, bnew().Abs(sum))) <= 0 {
			return sum
		}
	}
}

// bigAtan is atan x for x ≥ 0: halved by atan x = 2 atan(x / (1 + √(1 + x²))) until x < 1/8, then summed.
func bigAtan(x *big.Float) *big.Float {
	k := 0
	for x.Cmp(bf(0.125)) > 0 {
		r := bnew().Sqrt(bnew().Add(bf(1), bnew().Mul(x, x)))
		x = bnew().Quo(x, bnew().Add(bf(1), r))
		k++
	}
	x2 := bnew().Mul(x, x)
	sum, pow := bnew().Set(x), bnew().Set(x)
	eps := bnew().SetMantExp(bf(1), -bigPrec)
	for n := 3; ; n += 2 {
		pow.Mul(pow, x2)
		pow.Neg(pow)
		term := bnew().Quo(pow, bf(float64(n)))
		sum.Add(sum, term)
		if bnew().Abs(term).Cmp(bnew().Mul(eps, bnew().Abs(sum))) <= 0 {
			return sum.SetMantExp(sum, k)
		}
	}
}

func f64(x *big.Float) float64 { v, _ := x.Float64(); return v }

// thirdBig is π − x − y.
func thirdBig(x, y *big.Float) *big.Float { return bnew().Sub(bnew().Sub(bigPi, x), y) }

// relSameTriangle compares all six parts purely relatively (tiny angles included).
func relSameTriangle(got, want Triangle, tol float64) bool {
	return relNear(got.SideA, want.SideA, tol) && relNear(got.SideB, want.SideB, tol) && relNear(got.SideC, want.SideC, tol) &&
		relNear(got.AngleARad, want.AngleARad, tol) && relNear(got.AngleBRad, want.AngleBRad, tol) && relNear(got.AngleCRad, want.AngleCRad, tol)
}

// nearPiDeltas are the gaps π − angle the near-π tests sweep.
var nearPiDeltas = []float64{1e-3, 1e-4, 1e-5, 1e-6, 1e-7, 1e-8, 1e-9, 1e-10, 1e-11, 1e-12}

// tolNearPi is the relative tolerance of the near-π tests: 16 ulps of a value near 1.
const tolNearPi = 16 * 0x1p-52

// The oracles check themselves on values with known answers.
func TestBigOracles(t *testing.T) {
	if got := f64(bigSin(bnew().Quo(bigPi, bf(6)))); !relNear(got, 0.5, 0x1p-52) {
		t.Errorf("bigSin(π/6) = %v", got)
	}
	if got := f64(bigCos(bnew().Quo(bigPi, bf(3)))); !relNear(got, 0.5, 0x1p-52) {
		t.Errorf("bigCos(π/3) = %v", got)
	}
	if got := f64(bigAtan(bf(1))); !relNear(got, math.Pi/4, 0x1p-52) {
		t.Errorf("bigAtan(1) = %v", got)
	}
	if got := f64(bigSin(bf(math.Pi))); !relNear(got, piLo, 1e-15) {
		t.Errorf("bigSin(math.Pi) = %v, want π − math.Pi", got)
	}
}

// thirdAngle against π − x − y in 256 bits, near π, near 0 and in between.
func TestThirdAngle(t *testing.T) {
	cases := [][2]float64{{1, 1}, {0.3, 2.5}, {1e-20, 1e-20}, {math.Pi / 2, math.Pi / 2 * (1 - 1e-12)}}
	for _, d := range nearPiDeltas {
		cases = append(cases, [2]float64{math.Pi - d, d / 3}, [2]float64{d / 7, math.Pi - d})
	}
	for _, c := range cases {
		want := thirdBig(bf(c[0]), bf(c[1]))
		z, s := thirdAngle(c[0], c[1])
		if !relNear(z, f64(want), 4*0x1p-52) || !relNear(s, f64(bigSin(want)), 4*0x1p-52) {
			t.Errorf("thirdAngle(%v, %v) = (%v, %v), want (%v, %v)", c[0], c[1], z, s, f64(want), f64(bigSin(want)))
		}
	}
}
