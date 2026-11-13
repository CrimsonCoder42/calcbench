package algebra

import (
	"fmt"
	"math"
	"math/cmplx"
	"slices"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

const (
	// aberthBaseSweeps and aberthSweepsPerDegree set the sweep limit: 100 + n for degree n.
	aberthBaseSweeps      = 100
	aberthSweepsPerDegree = 1
	// aberthAngleOffset is σ, the angle (radians) added to every start point. Our own choice, not taken
	// from a source: any σ that keeps the start points off the real axis and not symmetric under conjugation
	// works (a real start set that is closed under conjugation can never leave it); the roots returned do
	// not depend on σ beyond the convergence tolerance.
	aberthAngleOffset = 0.7
	// aberthStepUlps is the convergence test on the correction: |Δ| ≤ 4u·|z|, a few ulps of z.
	aberthStepUlps = 4 * unitRoundoff
)

// PolynomialRoots returns all n complex roots of the real polynomial
// p(x) = coeffs[0] + coeffs[1]·x + … + coeffs[n]·xⁿ (ascending order), sorted by real part, then by
// imaginary part, each repeated by its multiplicity. coeffs is read, never modified.
//
// Method: Aberth–Ehrlich simultaneous iteration. Leading zero coefficients coeffs[0] = … = coeffs[m−1] = 0
// give m exact zero roots and are factored out first. Start points lie on circles whose radii come from
// the upper convex hull of the points (i, log|coeffs[i]|) (Bini 1996), at angles offset by σ = 0.7 (our
// own choice), so none is random. Each sweep updates every unconverged approximation zₖ in place
// (Gauss–Seidel order) by zₖ ← zₖ − Δₖ, Δₖ = 1/(p′(zₖ)/p(zₖ) − Σ_{j≠k} 1/(zₖ − zⱼ)), with p and p′ on
// the original coefficients by the compensated Horner scheme (Graillat and Ménissier-Morain 2012), about
// as accurate as Horner's rule in twice the working precision (on the reversed polynomial when |zₖ| > 1,
// so nothing overflows). Convergence test (our own, not from a source): zₖ is converged when p(zₖ) = 0,
// when |Δₖ| ≤ 4u·|zₖ| (u = 2⁻⁵³; the step is applied), or when the correction has stopped decreasing
// (|Δₖ| ≥ its previous value) while the residual is at the second-order rounding level of the compensated
// evaluation, |p(zₖ)| ≤ (6n·u)²·Σ|coeffs[i]|·|zₖ|ⁱ (Graillat and Ménissier-Morain 2012, Theorem 5.5); the
// last case is the accuracy limit of an ill-conditioned or multiple root, and the step is not applied.
// So a simple root is accurate to a few ulps plus (6n·u)²·Σ|coeffs[i]||z|ⁱ/|p′(z)|: Wilkinson's
// ∏(x − k), k = 1…20, gives the roots of its float64 coefficients to about 1e-13. Afterwards a root is
// made real (imaginary part exactly 0) when p at its real part is no larger in modulus than p at the
// root, and each non-real root is paired with the root nearest its conjugate (when that is nearer than
// the real axis) and the two are made exact conjugates. Convergence is cubic at simple roots and linear
// at multiple roots; a root of multiplicity m is accurate to about the m-th root of the evaluation error
// (relative), which is the conditioning of the problem. A group of roots whose start radius
// (|coeffs[i₀]|/|coeffs[i₁]|)^(1/(i₁−i₀)) underflows to 0 is returned as 0, like any float64 underflow:
// those roots are below the float64 range, which is not an error.
// Errors: ErrEmpty when len(coeffs) < 2 (a constant has no roots); ErrDegenerate when coeffs[n] = 0
// (the degree is lower: trim it); ErrNotFinite when a coefficient is NaN or ±Inf; ErrOverflow when
// n·Σ|coeffs[i]| or a root's modulus is above the float64 range; ErrNoConvergence when some root has
// not met the convergence test after 100 + n sweeps, in which case the best approximations are still
// returned, sorted, with the error. The result is nil with any other error.
//
// ID:      NUM-ALGEB-010
// Time:    O(n²·k)   n = degree = len(coeffs)−1, k = sweeps ≤ 100 + n
// Space:   O(n) for the result (one allocation), O(n) extra (hull indices, last corrections and converged flags, at most three allocations)
// Source:  Aberth, "Iteration methods for finding all zeros of a polynomial simultaneously", Mathematics of Computation, 1973, 27(122):339–344, doi:10.1090/S0025-5718-1973-0329236-7
// Source:  Ehrlich, "A modified Newton method for polynomials", Communications of the ACM, 1967, 10(2):107–108, doi:10.1145/363067.363115
// Source:  Bini, "Numerical computation of polynomial zeros by means of Aberth's method", Numerical Algorithms, 1996, 13(2):179–200, doi:10.1007/BF02207694 (start points equispaced on circles of radii from the Newton polygon)
// Source:  Graillat, Ménissier-Morain, "Accurate summation, dot product and polynomial evaluation in complex floating point arithmetic", Information and Computation, 2012, 216:57–71, doi:10.1016/j.ic.2011.09.003, Algorithms 2.1, 3.1, 3.4, 5.2, 5.4 and Theorem 5.5, eq. (5.17) (compensated Horner evaluation and its error bound)
// Tags:    polynomial, roots, aberth, ehrlich, quartic, quintic, algebra
// Units:   coefficients and roots are dimensionless
// Example: z⁴ − 4z³ + 5z + 2 → ½(1 − √5), ½(3 − √17), ½(1 + √5), ½(3 + √17) (DLMF §1.11(iii), Example)
func PolynomialRoots(coeffs []float64) ([]complex128, error) {
	return polynomialRoots(coeffs, aberthBaseSweeps+aberthSweepsPerDegree*(len(coeffs)-1))
}

// polynomialRoots is PolynomialRoots with the sweep limit as a parameter (the tests lower it).
func polynomialRoots(coeffs []float64, maxSweeps int) ([]complex128, error) {
	n := len(coeffs) - 1
	if n < 1 {
		return nil, fmt.Errorf("%w: coeffs has %d elements, want at least 2", ErrEmpty, len(coeffs))
	}
	for i, v := range coeffs {
		if !reduce.AllFinite(v) {
			return nil, fmt.Errorf("%w: coeffs[%d] = %v", ErrNotFinite, i, v)
		}
	}
	if coeffs[n] == 0 {
		return nil, fmt.Errorf("%w: coeffs[%d] = 0 (leading coefficient)", ErrDegenerate, n)
	}
	weighted := 0.0
	for _, v := range coeffs {
		weighted += math.Abs(v)
	}
	weighted *= float64(n)
	if !reduce.AllFinite(weighted) {
		return nil, fmt.Errorf("%w: coeffs are too large to evaluate p′ without overflow", ErrOverflow)
	}

	roots := make([]complex128, n)
	zeros := 0
	for coeffs[zeros] == 0 {
		zeros++ // stops at n at the latest: coeffs[n] != 0
	}
	poly := coeffs[zeros:] // p(x) = x^zeros · poly(x), poly(0) != 0
	work := roots[zeros:]  // roots[:zeros] stay exactly 0
	deg := len(poly) - 1
	if deg == 0 {
		return roots, nil
	}
	converged := make([]bool, deg)
	underflowed, err := aberthStartPoints(poly, work)
	if err != nil {
		return nil, err
	}
	// The first underflowed start points are exactly 0 (radius below the float64 range): they stay 0.
	for k := range underflowed {
		converged[k] = true
	}
	remaining := deg - underflowed
	lastStep := make([]float64, deg) // |Δₖ| of the previous sweep; 0 before the first
	for sweep := 0; sweep < maxSweeps && remaining > 0; sweep++ {
		for k := range work {
			if converged[k] {
				continue
			}
			z := work[k]
			ratio, exact, atNoise := newtonRatio(poly, z)
			if exact {
				converged[k] = true
				remaining--
				continue
			}
			var sum complex128
			for j, zj := range work {
				if j != k {
					sum += 1 / (z - zj)
				}
			}
			delta := 1 / (ratio - sum)
			next := z - delta
			if !reduce.AllFinite(real(next)) || !reduce.AllFinite(imag(next)) {
				continue
			}
			step := cmplx.Abs(delta)
			switch {
			case step <= aberthStepUlps*cmplx.Abs(next):
				work[k] = next
				converged[k] = true
				remaining--
			case atNoise && lastStep[k] > 0 && step >= lastStep[k]:
				converged[k] = true // the correction is rounding noise: keep z
				remaining--
			default:
				work[k] = next
				lastStep[k] = step
			}
		}
	}

	snapAndPairConjugates(poly, work, converged)
	slices.SortFunc(roots, compareRoots)
	if remaining > 0 {
		return roots, fmt.Errorf("%w: %d of %d roots after %d sweeps", ErrNoConvergence, remaining, n, maxSweeps)
	}
	return roots, nil
}

// aberthStartPoints fills z with the start points for poly (poly[0] != 0, poly[deg] != 0): for each edge
// (i₀, i₁) of the upper convex hull of (i, log|poly[i]|), i₁ − i₀ points equally spaced on the circle of
// radius (|poly[i₀]|/|poly[i₁]|)^(1/(i₁−i₀)), at angles 2πj/(i₁−i₀) + 2π·edge/deg + σ. The radii increase
// along the hull, so the edges whose radius underflows to 0 come first: their points are exactly 0 and
// their count is returned as underflowed. A radius above the float64 range is ErrOverflow.
func aberthStartPoints(poly []float64, z []complex128) (underflowed int, err error) {
	deg := len(poly) - 1
	logAbs := func(i int) float64 { return math.Log(math.Abs(poly[i])) }
	hull := make([]int, 0, deg+1)
	for i, v := range poly {
		if v == 0 {
			continue
		}
		// Keep only right turns: pop the last vertex while it lies on or below the chord.
		for len(hull) >= 2 {
			o, m := hull[len(hull)-2], hull[len(hull)-1]
			cross := float64(m-o)*(logAbs(i)-logAbs(o)) - (logAbs(m)-logAbs(o))*float64(i-o)
			if cross < 0 {
				break
			}
			hull = hull[:len(hull)-1]
		}
		hull = append(hull, i)
	}
	next := 0
	for edge := 1; edge < len(hull); edge++ {
		i0, i1 := hull[edge-1], hull[edge]
		count := i1 - i0
		radius := math.Exp((logAbs(i0) - logAbs(i1)) / float64(count))
		if !reduce.AllFinite(radius) {
			return 0, fmt.Errorf("%w: a root's modulus is above the float64 range", ErrOverflow)
		}
		if radius == 0 {
			underflowed += count // z[next:next+count] stay 0
			next += count
			continue
		}
		for j := range count {
			angle := 2*math.Pi*float64(j)/float64(count) + 2*math.Pi*float64(edge)/float64(deg) + aberthAngleOffset
			z[next] = cmplx.Rect(radius, angle)
			next++
		}
	}
	return underflowed, nil
}

// newtonRatio returns p′(z)/p(z) for the polynomial with ascending coefficients poly; exact = true when
// the computed p(z) is 0 (the ratio is then not needed), and atNoise = true when
// |p(z)| ≤ (6n·u)²·Σ|poly[i]|·|z|ⁱ, the second-order rounding level of the compensated evaluation. For
// |z| > 1 it evaluates the reversed polynomial q(w) = wⁿ·p(1/w) at w = 1/z, so no power of |z| above 1 is
// formed: p′(z)/p(z) = w·(n − w·q′(w)/q(w)), and the residual test on q is the same test scaled by |w|ⁿ.
func newtonRatio(poly []float64, z complex128) (ratio complex128, exact, atNoise bool) {
	deg := len(poly) - 1
	if cmplx.Abs(z) <= 1 {
		p, dp, bound := compensatedHorner(poly, z, false)
		if p == 0 {
			return 0, true, true
		}
		return dp / p, false, cmplx.Abs(p) <= noiseLevel(deg, bound)
	}
	w := 1 / z
	q, dq, bound := compensatedHorner(poly, w, true)
	if q == 0 {
		return 0, true, true
	}
	return w * (complex(float64(deg), 0) - w*dq/q), false, cmplx.Abs(q) <= noiseLevel(deg, bound)
}

// noiseLevel returns (6n·u)²·bound: the second-order term γ̃₂ₙ²·Σ|aᵢ||z|ⁱ of the error bound of the
// compensated Horner scheme (Graillat and Ménissier-Morain 2012, Theorem 5.5, eq. (5.17)), with
// γ̃₂ₙ = 2n·√2·γ₂/(1 − 2n·√2·γ₂) ≈ 4√2·n·u rounded up to 6n·u.
func noiseLevel(deg int, bound float64) float64 {
	g := 6 * float64(deg) * unitRoundoff
	return g * g * bound
}

// compensatedHorner returns p(z) and p′(z) by the compensated Horner scheme (Graillat and
// Ménissier-Morain 2012, Algorithms 3.1, 3.4, 5.2 and 5.4): each step's rounding errors are recovered
// exactly by TwoSum and the FMA form of TwoProduct and evaluated as a second polynomial, which is added at
// the end, so the value is about as accurate as Horner's rule in twice the working precision. Two
// simplifications and one extension are our own: the error terms of a step are summed in plain floating
// point instead of by AccSum; and the derivative recurrence d ← d·z + s is compensated the same way, with
// the running error of s added to the error of d, because near a multiple root p′ is as ill conditioned
// as p. It also returns Σ|aᵢ|·|z|ⁱ. The coefficients are poly in ascending order, or in descending order
// when reversed is true (the reversed polynomial).
func compensatedHorner(poly []float64, z complex128, reversed bool) (p, dp complex128, bound float64) {
	deg := len(poly) - 1
	coef := func(i int) float64 { // coefficient of zⁱ
		if reversed {
			return poly[deg-i]
		}
		return poly[i]
	}
	r := cmplx.Abs(z)
	s := complex(coef(deg), 0) // Horner value of p
	var c complex128           // Horner value of the error polynomial of p
	var d, dc complex128       // the same for p′
	bound = math.Abs(real(s))
	for i := deg - 1; i >= 0; i-- {
		// d ← d·z + s, with its rounding errors and the error c of s carried into dc.
		dProd, dErr := twoProductComplex(d, z)
		dRe, e1 := reduce.TwoSum(real(dProd), real(s))
		dIm, e2 := reduce.TwoSum(imag(dProd), imag(s))
		dc = dc*z + dErr + complex(e1, e2) + c
		d = complex(dRe, dIm)
		// s ← s·z + aᵢ (TwoSumCplx with a real aᵢ: the imaginary part is exact).
		sProd, sErr := twoProductComplex(s, z)
		a := coef(i)
		sRe, e3 := reduce.TwoSum(real(sProd), a)
		c = c*z + sErr + complex(e3, 0)
		s = complex(sRe, imag(sProd))
		bound = bound*r + math.Abs(a)
	}
	return s + c, d + dc, bound
}

// twoProductComplex returns p = fl(x·y) and the sum of its exact error terms e + f + g, rounded
// (TwoProductFMACplx, Graillat and Ménissier-Morain 2012, Algorithm 3.4: x·y = p + e + f + g exactly).
// float64(…) stops the compiler fusing a product into a later addition.
func twoProductComplex(x, y complex128) (p, err complex128) {
	a, b, c, d := real(x), imag(x), real(y), imag(y)
	z1 := float64(a * c)
	h1 := math.FMA(a, c, -z1)
	z2 := float64(b * d)
	h2 := math.FMA(b, d, -z2)
	z3 := float64(a * d)
	h3 := math.FMA(a, d, -z3)
	z4 := float64(b * c)
	h4 := math.FMA(b, c, -z4)
	z5, h5 := reduce.TwoSum(z1, -z2)
	z6, h6 := reduce.TwoSum(z3, z4)
	return complex(z5, z6), complex(h1-h2+h5, h3+h4+h6)
}

// logAbsValue returns log|p(z)| for the polynomial with ascending coefficients poly, −Inf when p(z) = 0,
// evaluated as newtonRatio does (reversed for |z| > 1: log|p(z)| = log|q(w)| − n·log|w|), so it never
// overflows.
func logAbsValue(poly []float64, z complex128) float64 {
	deg := len(poly) - 1
	if cmplx.Abs(z) <= 1 {
		p, _, _ := compensatedHorner(poly, z, false)
		return math.Log(cmplx.Abs(p))
	}
	w := 1 / z
	q, _, _ := compensatedHorner(poly, w, true)
	return math.Log(cmplx.Abs(q)) - float64(deg)*math.Log(cmplx.Abs(w))
}

// snapAndPairConjugates makes a root real when p at its real part is no larger in modulus than p at the
// root, then pairs each root of positive imaginary part with the unpaired root of negative imaginary part
// nearest its conjugate, if that one is closer to the conjugate than the root is to the real axis, and
// replaces both by the mean conjugate pair. paired is scratch space of len(z).
func snapAndPairConjugates(poly []float64, z []complex128, paired []bool) {
	for k, zk := range z {
		paired[k] = false
		if imag(zk) == 0 {
			continue
		}
		if re := complex(real(zk), 0); logAbsValue(poly, re) <= logAbsValue(poly, zk) {
			z[k] = re
		}
	}
	for k, zk := range z {
		if imag(zk) <= 0 || paired[k] {
			continue
		}
		best, bestDist := -1, math.Inf(1)
		for j, zj := range z {
			if paired[j] || imag(zj) >= 0 {
				continue
			}
			if d := cmplx.Abs(zj - cmplx.Conj(zk)); d < bestDist {
				best, bestDist = j, d
			}
		}
		// A true conjugate partner is much nearer conj(zₖ) than zₖ is to the real axis; anything farther is
		// a different root and is left alone.
		if best < 0 || bestDist >= imag(zk) {
			continue
		}
		zj := z[best]
		re := (real(zk) + real(zj)) / 2
		im := (imag(zk) - imag(zj)) / 2
		z[k] = complex(re, im)
		z[best] = complex(re, -im)
		paired[k], paired[best] = true, true
	}
}
