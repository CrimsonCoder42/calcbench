package prob

import (
	"fmt"
	"math"
	"testing"
)

// The private incomplete gamma that ChiSquareCDF used before it moved to num/special (NUM-SPEC-001/002),
// kept verbatim from _delete/num/prob/chi_square_cdf.go as a test-only oracle.

// gammaMaxTerms caps the terms of the series and of the continued fraction (a defect guard: for
// a ≤ chiSquareMaxDF/2 both converge in O(√a) terms, under 10⁴).
const gammaMaxTerms = 100000

// gammaPrefactor returns z^a·e^(−z)/Γ(a + 1) for a > 0, z > 0. Above a = 15 it uses Loader's saddle-point
// form exp(−δ(a) − D₀(a, z))/√(2πa) (δ the Stirling remainder, DLMF 5.11.1, valid for real a; D₀ the
// deviance), which avoids the cancellation of a·ln z − z − ln Γ(a + 1) when a and z are large and close;
// at or below 15 the direct form is exact to a few ulps.
func gammaPrefactor(a, z float64) float64 {
	if a > 15 {
		return math.Exp(-stirlingRemainder(a)-deviance(a, z)) / math.Sqrt(2*math.Pi*a)
	}
	lg, _ := math.Lgamma(a + 1)
	return math.Exp(a*math.Log(z) - z - lg)
}

// regularizedGamma returns P(a, z) = γ(a, z)/Γ(a) and Q(a, z) = Γ(a, z)/Γ(a) = 1 − P (DLMF 8.2.4–8.2.5) for
// a > 0, z ≥ 0. For z < a + 1 it sums the series of DLMF 8.7.1, P = z^a·e^(−z)·Σₖ zᵏ/Γ(a + k + 1), and
// Q = 1 − P; otherwise it evaluates the continued fraction of DLMF 8.9.2 for z^(−a)·e^z·Γ(a, z) by the
// modified Lentz method, and P = 1 − Q. The lower tail (z small) is always the series and the upper tail
// (z large) always the fraction, so each tail keeps its relative accuracy. On the series side Q ≥ 1 − P(1, 2)
// = e⁻² ≈ 0.135 when a ≥ 1; for a < 1, Q = 1 − P carries an absolute error of a few 2⁻⁵³.
func regularizedGamma(a, z float64) (p, q float64, err error) {
	if z == 0 {
		return 0, 1, nil
	}
	pre := gammaPrefactor(a, z)
	if z < a+1 {
		// Σₖ zᵏ/((a + 1)(a + 2)…(a + k)): z^a e^(−z)/Γ(a + k + 1) = pre·zᵏ/((a + 1)…(a + k)).
		sum, term := 1.0, 1.0
		for k := 1; ; k++ {
			if k > gammaMaxTerms {
				return 0, 0, fmt.Errorf("%w: incomplete gamma series did not converge for a = %v, z = %v", ErrOutOfRange, a, z)
			}
			term *= z / (a + float64(k))
			sum += term
			if term <= sum*0x1p-53 {
				break
			}
		}
		p = min(1, pre*sum)
		return p, 1 - p, nil
	}
	// DLMF 8.9.2: z^(−a)·e^z·Γ(a, z) = 1/z / (1 + (1 − a)/z / (1 + 1/z / (1 + (2 − a)/z / (1 + 2/z / (1 + …))))),
	// partial numerators a₁ = 1/z, a₂ₘ = (m − a)/z, a₂ₘ₊₁ = m/z, partial denominators all 1.
	// Modified Lentz: f = b₀ + a₁/(b₁ + a₂/(b₂ + …)) with b₀ = 0 replaced by tiny.
	const tiny = 0x1p-1000
	f, c, d := tiny, tiny, 0.0
	for j := 1; ; j++ {
		if j > 2*gammaMaxTerms {
			return 0, 0, fmt.Errorf("%w: incomplete gamma continued fraction did not converge for a = %v, z = %v", ErrOutOfRange, a, z)
		}
		var num float64
		switch {
		case j == 1:
			num = 1 / z
		case j%2 == 0:
			num = (float64(j/2) - a) / z
		default:
			num = float64(j/2) / z
		}
		d = 1 + num*d
		if d == 0 {
			d = tiny
		}
		c = 1 + num/c
		if c == 0 {
			c = tiny
		}
		d = 1 / d
		delta := c * d
		f *= delta
		if math.Abs(delta-1) <= 0x1p-53 {
			break
		}
	}
	// Γ(a, z)/Γ(a) = z^a e^(−z)/Γ(a)·f = pre·a·f.
	q = min(1, pre*a*f)
	return 1 - q, q, nil
}

// chiSquarePrivateGammaTol is the relative tolerance between ChiSquareCDF (num/special) and the retired
// private code at a = df/2, z = x/2. Reading the two codes: δ(a) and the deviance D₀ are the same
// expressions in both (num/prob stirlingRemainder/deviance, num/special stirlingRemainder/gammaDeviance), so
// for a > 15 the prefactors differ only by a few ulps (·a and /√(2πa) against ·√(a/2π)). Three places
// cannot agree to 1e-14 and get the allowance the analysis supports:
//   - a ≤ 15, far tails: both use the direct form exp(a·ln z − z − ln Γ(·)) but with ln Γ(a + 1) against
//     ln Γ(a); each rounds an exponent of size E = |a·ln z| + z + |ln Γ(a)| independently, so the results
//     differ by up to about 4·E·2⁻⁵³ relative (E ≈ 700 near underflow → about 3e-13);
//   - large a near the bulk (z ≈ a): the series stops at 2⁻⁵³ (old) against 2⁻⁵² (new) of the sum while
//     its terms fall off like exp(−k²/2a), so the truncated tails differ by about √a·2⁻⁵²; the two
//     continued fractions (odd/even Lentz against the even form) take O(√a) steps each, the same order;
//   - 1e-14 everywhere else (bulk, small and moderate a).
//
// Below 1e-290 (near and inside the subnormal range) relative accuracy is not defined; there both must be
// under 1e-280.
// tol = 1e-14 + 4·E·2⁻⁵³ (a ≤ 15 only) + 4·√a·2⁻⁵².
func chiSquarePrivateGammaTol(a, z float64) float64 {
	tol := 1e-14 + 4*math.Sqrt(a)*0x1p-52
	if a <= 15 {
		lg, _ := math.Lgamma(a)
		tol += 4 * (math.Abs(a*math.Log(z)) + z + math.Abs(lg)) * 0x1p-53
	}
	return tol
}

func TestChiSquareCDF_matchesPrivateGamma(t *testing.T) {
	dfs := []float64{1, 1.5, 2, 3, 5, 10, 29, 30, 31, 50, 100, 333, 1000, 1e4, 1e5, 1e6}
	// x as multiples of df across the bulk, plus the far tails as offsets in standard deviations √(2df).
	ratios := []float64{1e-6, 1e-3, 0.01, 0.1, 0.3, 0.5, 0.8, 0.9, 0.99, 1, 1.01, 1.1, 1.5, 2, 3, 5, 10, 30}
	sds := []float64{-30, -10, -5, -2, -1, -0.1, 0.1, 1, 2, 5, 10, 30}
	for _, df := range dfs {
		var xs []float64
		for _, r := range ratios {
			xs = append(xs, r*df)
		}
		for _, s := range sds {
			if x := df + s*math.Sqrt(2*df); x > 0 {
				xs = append(xs, x)
			}
		}
		xs = append(xs, 1e-300, 1e-10, df+2, df+1.999, 1400)
		for _, x := range xs {
			t.Run(fmt.Sprintf("df=%g/x=%g", df, x), func(t *testing.T) {
				want, _, err := regularizedGamma(df/2, x/2)
				if err != nil {
					t.Fatalf("oracle: %v", err)
				}
				got, err := ChiSquareCDF(x, df)
				if err != nil {
					t.Fatalf("ChiSquareCDF(%v, %v): %v", x, df, err)
				}
				if want < 1e-290 || got < 1e-290 {
					if math.Max(want, got) > 1e-280 {
						t.Fatalf("ChiSquareCDF(%v, %v) = %v, private %v (one underflowed)", x, df, got, want)
					}
					return
				}
				tol := chiSquarePrivateGammaTol(df/2, x/2)
				if rel := math.Abs(got-want) / want; rel > tol {
					t.Fatalf("ChiSquareCDF(%v, %v) = %.17g, private %.17g, rel %.3g > %.3g", x, df, got, want, rel, tol)
				}
			})
		}
	}
}
