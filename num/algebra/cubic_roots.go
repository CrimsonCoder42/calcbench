package algebra

import (
	"cmp"
	"fmt"
	"math"
	"math/cmplx"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// cubicNewtonSteps is the most Newton steps used to polish a closed-form root of a cubic.
const cubicNewtonSteps = 3

// CubicRoots returns the three roots of a·x³ + b·x² + c·x + d = 0, sorted by real part, then by
// imaginary part. Real roots have imaginary part exactly 0; a multiple root is repeated.
//
// Method: x is scaled by an exact power of two, x = 2ᵏ·y, chosen from the coefficients' exponents so
// that the monic cubic y³ + A·y² + B·y + C has |A|, |B|, |C| < 2 (Q and R below cannot overflow).
// With Q = (A² − 3B)/9 and R = (2A³ − 9AB + 27C)/54: if R² < Q³ the three real roots are
// −2√Q·cos((θ + 2πj)/3) − A/3, θ = arccos(R/√Q³), j = 0, 1, 2, and only the one of largest modulus is
// kept (it is free of cancellation); otherwise the one real root is S + Q/S − A/3,
// S = −sgn(R)·∛(|R| + √(R² − Q³)) (Q/S taken as 0 when S = 0). That root x₁ is polished by at most 3
// Newton steps on the scaled cubic and again on the original coefficients (divided by a power of two so
// that the largest is below 1; PolynomialValueAndDerivative), a step kept only if it lowers |p|. The
// cubic is then deflated to a·x² + f1·x + f0 in the direction that does not cancel: from the constant
// term (f0 = −d/x₁, f1 = (f0 − c)/x₁) when |a|·|x₁|³ ≥ |d|, i.e. x₁ is at least the geometric mean of the
// other two roots, else from the top (f1 = b + a·x₁, f0 = c + x₁·f1); forward division by the largest
// root would form f1 = −a(x₂ + x₃) as the difference of two numbers of size |a·x₁|. The quadratic is
// solved by QuadraticRoots (the smaller root as Vieta's c/(a·x), no cancellation), and each of its roots
// is polished by at most 3 Newton steps on the original coefficients (complex Newton for a conjugate pair,
// which stays exactly conjugate). So every root is accurate to its own conditioning, also when the roots
// are many orders of magnitude apart. d = 0 gives the root 0 exactly and the rest from
// QuadraticRoots(a, b, c). Near a root of multiplicity m the error is of order u^(1/m) (u = 2⁻⁵³): that is
// the conditioning of the problem, not a defect.
// Errors: ErrDegenerate when a = 0 (use QuadraticRoots); ErrNotFinite when a coefficient is NaN or
// ±Inf; ErrOverflow when a root is outside the float64 range. The roots are 0 with any error.
//
// ID:      NUM-ALGEB-009
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Press, Teukolsky, Vetterling, Flannery, Numerical Recipes in C: The Art of Scientific Computing, 2nd ed., Cambridge University Press, 1992, §5.6, eqs. (5.6.10)–(5.6.12) and (5.6.15)–(5.6.17)
// Source:  NIST Digital Library of Mathematical Functions, §1.11(iii) (cubic equations, Example), https://dlmf.nist.gov/1.11 (read 2026-11-12)
// Source:  Kahan, "To Solve a Real Cubic Equation", lecture notes, University of California, Berkeley, 1986, https://people.eecs.berkeley.edu/~wkahan/Math128/Cubic.pdf (read 2026-11-12) (Newton polishing, deflation to a quadratic)
// Tags:    cubic, roots, cardano, viete, polynomial, algebra
// Units:   a, b, c, d and the roots are dimensionless
// Example: z³ − 6z² + 6z − 2 → 2 + ∛4 + ∛2 ≈ 4.847322101863 and 2 + ∛4·ρ + ∛2·ρ², its conjugate,
// ≈ 0.576338949068 ± 0.283606001027i, ρ = e^(2πi/3) (DLMF §1.11(iii), Example)
func CubicRoots(a, b, c, d float64) ([3]complex128, error) {
	var roots [3]complex128
	if err := checkQuadratic(a, b, c); err != nil {
		return roots, err
	}
	if !reduce.AllFinite(d) {
		return roots, fmt.Errorf("%w: d = %v", ErrNotFinite, d)
	}
	if a == 0 {
		return roots, fmt.Errorf("%w: a = 0", ErrDegenerate)
	}
	if d == 0 {
		r1, r2, err := QuadraticRoots(a, b, c)
		if err != nil {
			return [3]complex128{}, err
		}
		roots = [3]complex128{0, r1, r2}
		sortRoots3(&roots)
		return roots, nil
	}

	// Scale: x = 2^k·y. With a = fa·2^ka and so on, A = b/(a·2^k), B = c/(a·2^2k), C = d/(a·2^3k).
	fa, ka := math.Frexp(a)
	fb, kb := math.Frexp(b)
	fc, kc := math.Frexp(c)
	fd, kd := math.Frexp(d)
	k := ceilDiv(kd-ka, 3)
	if b != 0 {
		k = max(k, kb-ka)
	}
	if c != 0 {
		k = max(k, ceilDiv(kc-ka, 2))
	}
	coefA := math.Ldexp(fb/fa, kb-ka-k)
	coefB := math.Ldexp(fc/fa, kc-ka-2*k)
	coefC := math.Ldexp(fd/fa, kd-ka-3*k)
	monic := [4]float64{coefC, coefB, coefA, 1}

	q := (coefA*coefA - 3*coefB) / 9
	r := (2*coefA*coefA*coefA - 9*coefA*coefB + 27*coefC) / 54
	third := coefA / 3
	// Explicit conversions round each product, so the compiler cannot fuse them into an FMA: the sign test
	// and the square root below must see the same value of R² − Q³.
	q3 := float64(q * q * q)
	disc := float64(r*r) - q3
	var y float64
	if disc < 0 {
		// Three real roots (Viète's trigonometric form). Only the one of largest modulus is kept: it is free
		// of cancellation, and the other two come from the deflated quadratic, where Vieta's c/(a·x) gives
		// the small one to full relative precision.
		cosArg := max(-1, min(1, r/math.Sqrt(q3)))
		theta := math.Acos(cosArg)
		scale := -2 * math.Sqrt(q)
		for j := range 3 {
			yj := scale*math.Cos((theta+2*math.Pi*float64(j))/3) - third
			if math.Abs(yj) > math.Abs(y) {
				y = yj
			}
		}
	} else {
		s := math.Cbrt(math.Abs(r) + math.Sqrt(disc))
		if r > 0 {
			s = -s
		}
		t := 0.0
		if s != 0 {
			t = q / s
		}
		y = s + t - third
	}
	y = polishCubicRoot(monic[:], y)
	x := math.Ldexp(y, k)
	if !reduce.AllFinite(x) {
		return [3]complex128{}, fmt.Errorf("%w: a root of the cubic", ErrOverflow)
	}

	// Deflate and polish on the original coefficients, divided by an exact power of two so that the
	// largest has magnitude in [½, 1): the scaled cubic above may have lost a tiny coefficient to underflow.
	top := kd
	for _, v := range [3]struct {
		f float64
		k int
	}{{fa, ka}, {fb, kb}, {fc, kc}} {
		if v.f != 0 {
			top = max(top, v.k)
		}
	}
	p := [4]float64{math.Ldexp(d, -top), math.Ldexp(c, -top), math.Ldexp(b, -top), math.Ldexp(a, -top)}
	x = polishCubicRoot(p[:], x)
	// a·x³ + b·x² + c·x + d = (x − x₁)(a·x² + f1·x + f0). Backward (from d) when x₁ is the larger:
	// |a|·|x₁|³ ≥ |d| ⇔ |x₁|² ≥ |x₂·x₃|; forward (from a) otherwise. Forward division by a large x₁ would
	// form f1 = b + a·x₁ = −a(x₂ + x₃) by cancellation; backward division by a small x₁ would do the same
	// in f1 = (f0 − c)/x₁.
	var f1, f0 float64
	if x != 0 && math.Log2(math.Abs(p[3]))+3*math.Log2(math.Abs(x)) >= math.Log2(math.Abs(p[0])) {
		f0 = -p[0] / x
		f1 = (f0 - p[1]) / x
	} else {
		f1 = p[2] + p[3]*x
		f0 = p[1] + x*f1
	}
	if !reduce.AllFinite(f1) || !reduce.AllFinite(f0) {
		return [3]complex128{}, fmt.Errorf("%w: a root of the cubic", ErrOverflow)
	}
	r1, r2, err := QuadraticRoots(p[3], f1, f0)
	if err != nil {
		return [3]complex128{}, err
	}
	if imag(r1) == 0 {
		r1 = complex(polishCubicRoot(p[:], real(r1)), 0)
		r2 = complex(polishCubicRoot(p[:], real(r2)), 0)
	} else {
		r2 = polishComplexCubicRoot(p, r2) // imag(r2) > 0
		r1 = cmplx.Conj(r2)
	}
	roots = [3]complex128{complex(x, 0), r1, r2}
	sortRoots3(&roots)
	return roots, nil
}

// polishCubicRoot improves y as a root of the polynomial with ascending coefficients monic by at most
// cubicNewtonSteps Newton steps, keeping a step only when it lowers |p(y)|.
func polishCubicRoot(monic []float64, y float64) float64 {
	p, dp, err := PolynomialValueAndDerivative(monic, y)
	if err != nil {
		return y
	}
	for range cubicNewtonSteps {
		if p == 0 || dp == 0 {
			break
		}
		next := y - p/dp
		pNext, dpNext, err := PolynomialValueAndDerivative(monic, next)
		if err != nil || math.Abs(pNext) >= math.Abs(p) {
			break
		}
		y, p, dp = next, pNext, dpNext
	}
	return y
}

// polishComplexCubicRoot improves z as a root of p[3]·z³ + p[2]·z² + p[1]·z + p[0] by at most
// cubicNewtonSteps complex Newton steps (Horner for p and p′), keeping a step only when it lowers |p(z)|.
func polishComplexCubicRoot(p [4]float64, z complex128) complex128 {
	eval := func(z complex128) (v, dv complex128, ok bool) {
		v = complex(p[3], 0)
		for i := 2; i >= 0; i-- {
			dv = dv*z + v
			v = v*z + complex(p[i], 0)
		}
		ok = reduce.AllFinite(real(v)) && reduce.AllFinite(imag(v)) && reduce.AllFinite(real(dv)) && reduce.AllFinite(imag(dv))
		return v, dv, ok
	}
	v, dv, ok := eval(z)
	if !ok {
		return z
	}
	for range cubicNewtonSteps {
		if v == 0 || dv == 0 {
			break
		}
		next := z - v/dv
		vNext, dvNext, ok := eval(next)
		if !ok || cmplx.Abs(vNext) >= cmplx.Abs(v) {
			break
		}
		z, v, dv = next, vNext, dvNext
	}
	return z
}

// ceilDiv returns ⌈x/n⌉ for n > 0.
func ceilDiv(x, n int) int {
	q := x / n // truncates toward zero, which is the ceiling for x < 0
	if x%n > 0 {
		q++
	}
	return q
}

// compareRoots orders complex numbers by real part, then by imaginary part.
func compareRoots(x, y complex128) int {
	if c := cmp.Compare(real(x), real(y)); c != 0 {
		return c
	}
	return cmp.Compare(imag(x), imag(y))
}

// sortRoots3 sorts three roots with compareRoots (a three-element sorting network, no allocation).
func sortRoots3(r *[3]complex128) {
	if compareRoots(r[0], r[1]) > 0 {
		r[0], r[1] = r[1], r[0]
	}
	if compareRoots(r[1], r[2]) > 0 {
		r[1], r[2] = r[2], r[1]
	}
	if compareRoots(r[0], r[1]) > 0 {
		r[0], r[1] = r[1], r[0]
	}
}
