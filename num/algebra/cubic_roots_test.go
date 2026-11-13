package algebra

import (
	"errors"
	"math"
	"math/big"
	"math/cmplx"
	"math/rand/v2"
	"testing"
)

// closeTo reports whether |got − want| ≤ tol·max(1, |want|).
func closeTo(got, want complex128, tol float64) bool {
	return cmplx.Abs(got-want) <= tol*max(1, cmplx.Abs(want))
}

// Worked example: DLMF §1.11(iii), Example: z³ − 6z² + 6z − 2 = 0 has the roots 2 + ∛4 + ∛2 and
// 2 + ∛4·ρ + ∛2·ρ² with its conjugate, ρ = e^(2πi/3) = −½ + i·√3/2. The expected values are built with
// math.Cbrt exactly as DLMF writes them. Tolerance 1e-14 relative: the polished real root has a residual
// of a few u·Σ|aᵢ||x|ⁱ ≈ 300u, i.e. an error up to about 300u/|p′(x)| ≈ 2e-15, and deflation passes a
// similar error to the complex pair; 1e-14 leaves a margin.
func TestCubicRoots_workedExample(t *testing.T) {
	c4, c2 := math.Cbrt(4), math.Cbrt(2)
	realRoot := complex(2+c4+c2, 0)
	pair := complex(2-(c4+c2)/2, (c4-c2)*math.Sqrt(3)/2)
	want := [3]complex128{cmplx.Conj(pair), pair, realRoot}
	got, err := CubicRoots(1, -6, 6, -2)
	if err != nil {
		t.Fatal(err)
	}
	for i := range got {
		if !closeTo(got[i], want[i], 1e-14) {
			t.Errorf("CubicRoots(1, -6, 6, -2)[%d] = %v, want %v", i, got[i], want[i])
		}
	}
	if imag(got[2]) != 0 {
		t.Errorf("real root has imaginary part %v, want exactly 0", imag(got[2]))
	}
}

func TestCubicRoots_cases(t *testing.T) {
	s3 := math.Sqrt(3) / 2
	tests := []struct {
		name       string
		a, b, c, d float64
		want       [3]complex128
		tol        float64
	}{
		// (x − 1)(x − 2)(x − 3), hand-checked.
		// Simple roots: a polished root is accurate to about γ·Σ|aᵢ||x|ⁱ/|p′(x)|, a few ulps of the scale;
		// 1e-14 relative leaves a margin.
		{"three real roots", 1, -6, 11, -6, [3]complex128{1, 2, 3}, 1e-14},
		{"d = 0 gives an exact zero root", 1, 0, -1, 0, [3]complex128{-1, 0, 1}, 0},
		{"x³ = 0", 1, 0, 0, 0, [3]complex128{0, 0, 0}, 0},
		{"cube roots of −1", 1, 0, 0, 1, [3]complex128{-1, complex(0.5, -s3), complex(0.5, s3)}, 1e-14},
		{"negative leading coefficient", -2, 12, -22, 12, [3]complex128{1, 2, 3}, 1e-14},
		{"double root", 1, -4, 5, -2, [3]complex128{1, 1, 2}, 1e-7}, // (x−1)²(x−2): error ~ √u
		{"triple root", 1, -3, 3, -1, [3]complex128{1, 1, 1}, 1e-4}, // (x−1)³: error ~ u^(1/3)
		{"huge coefficients", 0x1p1000, -6 * 0x1p1000, 11 * 0x1p1000, -6 * 0x1p1000, [3]complex128{1, 2, 3}, 1e-14},
		{"huge roots", 1, -6 * 0x1p300, 11 * 0x1p600, -6 * 0x1p900, [3]complex128{0x1p300, 2 * 0x1p300, 3 * 0x1p300}, 1e-14},
		{"tiny roots", 1, -6 * 0x1p-300, 11 * 0x1p-600, -6 * 0x1p-900, [3]complex128{0x1p-300, 2 * 0x1p-300, 3 * 0x1p-300}, 1e-14},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CubicRoots(tt.a, tt.b, tt.c, tt.d)
			if err != nil {
				t.Fatal(err)
			}
			// Tolerance relative to the largest root: the "huge" and "tiny" cases are scaled copies.
			// tol = 0 demands exact roots.
			scale := 0.0
			for _, w := range tt.want {
				scale = max(scale, cmplx.Abs(w))
			}
			for i := range got {
				if cmplx.Abs(got[i]-tt.want[i]) > tt.tol*scale {
					t.Errorf("CubicRoots(%v, %v, %v, %v)[%d] = %v, want %v", tt.a, tt.b, tt.c, tt.d, i, got[i], tt.want[i])
				}
			}
		})
	}
}

// bigCubicRoots returns the roots of a·x³ + b·x² + c·x + d (d != 0) computed in bigPrec-bit arithmetic,
// independently of CubicRoots: one real root by bisection on the Cauchy bound [−B, B],
// B = 1 + max(|b|, |c|, |d|)/|a| (2500 halvings: the bound is below 2¹⁰²⁶ and the roots the tests use
// need an absolute width below 2⁻¹⁴⁰⁰), then exact forward division to a·x² + f1·x + f0 and the quadratic
// formula, with the cancellation-free q form for real roots. With 8192 bits every rounding is far below
// float64 resolution, so the results, rounded to float64, are the correctly rounded roots.
func bigCubicRoots(a, b, c, d float64) [3]complex128 {
	nb := func() *big.Float { return new(big.Float).SetPrec(bigPrec) }
	eval := func(x *big.Float) int {
		v := newBig(a)
		for _, coef := range []float64{b, c, d} {
			v.Mul(v, x)
			v.Add(v, newBig(coef))
		}
		return v.Sign()
	}
	bound := newBig(math.Max(math.Abs(b), math.Max(math.Abs(c), math.Abs(d))))
	bound.Quo(bound, newBig(math.Abs(a)))
	bound.Add(bound, newBig(1))
	lo, hi := nb().Neg(bound), nb().Set(bound)
	signLo := eval(lo)
	for range 2500 {
		mid := nb().Add(lo, hi)
		mid.Quo(mid, newBig(2))
		if eval(mid) == signLo {
			lo = mid
		} else {
			hi = mid
		}
	}
	x := nb().Add(lo, hi)
	x.Quo(x, newBig(2))
	f1 := nb().Mul(newBig(a), x)
	f1.Add(f1, newBig(b))
	f0 := nb().Mul(f1, x)
	f0.Add(f0, newBig(c))
	disc := nb().Mul(f1, f1)
	fourAC := nb().Mul(newBig(a), f0)
	fourAC.Mul(fourAC, newBig(4))
	disc.Sub(disc, fourAC)
	xr, _ := x.Float64()
	var r1, r2 complex128
	if disc.Sign() >= 0 {
		q := nb().Sqrt(disc)
		if f1.Sign() < 0 {
			q.Neg(q)
		}
		q.Add(q, f1)
		q.Quo(q, newBig(-2)) // q = −(f1 + sgn(f1)·√disc)/2
		x1 := nb().Quo(q, newBig(a))
		x2 := nb().Quo(f0, q)
		v1, _ := x1.Float64()
		v2, _ := x2.Float64()
		r1, r2 = complex(v1, 0), complex(v2, 0)
	} else {
		twoA := newBig(2 * a)
		re := nb().Quo(f1, twoA)
		re.Neg(re)
		im := nb().Sqrt(disc.Neg(disc))
		im.Quo(im, twoA)
		im.Abs(im)
		vr, _ := re.Float64()
		vi, _ := im.Float64()
		r1, r2 = complex(vr, -vi), complex(vr, vi)
	}
	roots := [3]complex128{complex(xr, 0), r1, r2}
	sortRoots3(&roots)
	return roots
}

// Widely spread roots (reviewer's cases, 2026-11-12; reference roots from bigCubicRoots, not from the
// code). Each root must be within 8u of its own modulus: these roots are simple and well separated
// relative to their size, so their relative condition number Σ|aᵢ||x|ⁱ/(|x|·|p′(x)|) is at most about 3,
// a polished root is within a few u·cond of the true one, and the reference is correctly rounded (½u).
// The third case is also checked against its exact factorisation (x − 1)(x² − 10⁸x + 1): the roots are 1
// and 5·10⁷ ± √(25·10¹⁴ − 1).
func TestCubicRoots_widelySpreadRoots(t *testing.T) {
	tests := []struct {
		name       string
		a, b, c, d float64
	}{
		{"x³ + 1e10x² + x + 1e-10", 1, 1e10, 1, 1e-10},
		{"tiny leading coefficient", 1e-300, 1, 1, 1},
		{"(x − 1)(x² − 1e8x + 1)", 1, -(1e8 + 1), 1e8 + 1, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := bigCubicRoots(tt.a, tt.b, tt.c, tt.d)
			got, err := CubicRoots(tt.a, tt.b, tt.c, tt.d)
			if err != nil {
				t.Fatal(err)
			}
			for i := range got {
				if cmplx.Abs(got[i]-want[i]) > 8*unitRoundoff*cmplx.Abs(want[i]) {
					t.Errorf("root %d = %v, want %v (all: %v)", i, got[i], want[i], got)
				}
			}
		})
	}
	// Exact factorisation of the third case, independent of bigCubicRoots.
	sq := new(big.Float).SetPrec(bigPrec).SetFloat64(25e14 - 1) // exact in float64
	sq.Sqrt(sq)
	hiRoot := new(big.Float).SetPrec(bigPrec).Add(sq, newBig(5e7))
	loRoot := new(big.Float).SetPrec(bigPrec).Quo(newBig(1), hiRoot) // product of the pair is 1
	lo, _ := loRoot.Float64()
	hi, _ := hiRoot.Float64()
	got, err := CubicRoots(1, -(1e8 + 1), 1e8+1, -1)
	if err != nil {
		t.Fatal(err)
	}
	for i, w := range [3]float64{lo, 1, hi} {
		if imag(got[i]) != 0 || math.Abs(real(got[i])-w) > 8*unitRoundoff*w {
			t.Errorf("(x − 1)(x² − 1e8x + 1): root %d = %v, want %v", i, got[i], w)
		}
	}
}

func TestCubicRoots_errors(t *testing.T) {
	nan, inf := math.NaN(), math.Inf(1)
	tests := []struct {
		name       string
		a, b, c, d float64
		want       error
	}{
		{"a = 0", 0, 1, 1, 1, ErrDegenerate},
		{"NaN a", nan, 1, 1, 1, ErrNotFinite},
		{"Inf b", 1, inf, 1, 1, ErrNotFinite},
		{"NaN c", 1, 1, nan, 1, ErrNotFinite},
		{"Inf d", 1, 1, 1, -inf, ErrNotFinite},
		{"root overflows", math.SmallestNonzeroFloat64, 0x1p1000, 0, 1, ErrOverflow},
		{"root overflows with d = 0", math.SmallestNonzeroFloat64, 0x1p1000, 0, 0, ErrOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CubicRoots(tt.a, tt.b, tt.c, tt.d)
			if !errors.Is(err, tt.want) || got != [3]complex128{} {
				t.Errorf("CubicRoots(%v, %v, %v, %v) = (%v, %v), want (zeros, %v)", tt.a, tt.b, tt.c, tt.d, got, err, tt.want)
			}
		})
	}
}

// rebuildFromRoots returns the monic real polynomial with the given roots (ascending coefficients): a
// linear factor per real root, a quadratic x² − 2·Re z·x + |z|² per root with positive imaginary part
// (its conjugate is skipped). ok is false when a non-real root has no exact conjugate in roots.
func rebuildFromRoots(roots []complex128) (coeffs []float64, ok bool) {
	coeffs = []float64{1}
	for _, z := range roots {
		var factor []float64
		switch {
		case imag(z) == 0:
			factor = []float64{-real(z), 1}
		case imag(z) > 0:
			found := false
			for _, w := range roots {
				if w == cmplx.Conj(z) {
					found = true
					break
				}
			}
			if !found {
				return nil, false
			}
			factor = []float64{real(z)*real(z) + imag(z)*imag(z), -2 * real(z), 1}
		default:
			continue
		}
		var err error
		if coeffs, err = PolynomialMultiply(coeffs, factor); err != nil {
			return nil, false
		}
	}
	return coeffs, true
}

// randomRoots returns degree roots, real or in conjugate pairs, with real and imaginary parts in
// [−spread, spread] and every two roots at least minGap apart (the tests measure accuracy, not the
// conditioning of clustered roots, which has its own cases).
func randomRoots(rng *rand.Rand, degree int, spread, minGap float64) []complex128 {
	roots := make([]complex128, 0, degree)
	for len(roots) < degree {
		var cand []complex128
		re := (2*rng.Float64() - 1) * spread
		if degree-len(roots) >= 2 && rng.IntN(2) == 0 {
			im := (0.1 + 0.9*rng.Float64()) * spread
			cand = []complex128{complex(re, im), complex(re, -im)}
		} else {
			cand = []complex128{complex(re, 0)}
		}
		farEnough := true
		for _, c := range cand {
			for _, r := range roots {
				if cmplx.Abs(c-r) < minGap {
					farEnough = false
				}
			}
		}
		if farEnough {
			roots = append(roots, cand...)
		}
	}
	return roots
}

// rootTolerance returns the first-order bound on how far a computed root may sit from roots[k] when the
// polynomial's value carries an error of c·u·Σ|monic[i]|·|r|ⁱ (r = roots[k]): that error divided by
// |p′(r)| = Π_{j≠k}|r − roots[j]| (monic has leading coefficient 1), plus 4u·|r| for the root's own
// rounding.
func rootTolerance(monic []float64, roots []complex128, k int, c float64) float64 {
	r := roots[k]
	s := 0.0
	for i := len(monic) - 1; i >= 0; i-- {
		s = s*cmplx.Abs(r) + math.Abs(monic[i])
	}
	deriv := 1.0
	for j, rj := range roots {
		if j != k {
			deriv *= cmplx.Abs(r - rj)
		}
	}
	return c*unitRoundoff*s/deriv + 4*unitRoundoff*cmplx.Abs(r)
}

// checkAgainstKnownRoots fails unless every known root has a computed root within rootTolerance of it
// (the known roots are at least 0.1 apart, far more than any tolerance, so the nearest one is its match).
func checkAgainstKnownRoots(t *testing.T, trial int, monic []float64, known, got []complex128, c float64) {
	t.Helper()
	for k, r := range known {
		nearest := math.Inf(1)
		for _, z := range got {
			nearest = min(nearest, cmplx.Abs(z-r))
		}
		if tol := rootTolerance(monic, known, k, c); nearest > tol {
			t.Fatalf("trial %d: known root %v: nearest computed root is %g away, tolerance %g (computed %v)", trial, r, nearest, tol, got)
		}
	}
}

// Oracle: random cubics built from known roots (PolynomialMultiply of the factors), then multiplied by a
// random leading coefficient. Each known root must have a computed root within rootTolerance with
// c = 1000: the polished closed-form root has a residual of a few u·Σ|aᵢ||x|ⁱ, the deflated quadratic adds
// a few more, and building and scaling the coefficients adds about 4u; 1000 is a wide margin over all of
// these. The roots must be sorted and be real or exact conjugate pairs.
func TestCubicRoots_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		known := randomRoots(rng, 3, 10, 0.1)
		monic, _ := rebuildFromRoots(known)
		lead := (2*rng.Float64() - 1) * math.Ldexp(1, rng.IntN(100)-50)
		if lead == 0 {
			continue
		}
		got, err := CubicRoots(lead*monic[3], lead*monic[2], lead*monic[1], lead*monic[0])
		if err != nil {
			t.Fatalf("trial %d: %v", trial, err)
		}
		for i := 1; i < 3; i++ {
			if compareRoots(got[i-1], got[i]) > 0 {
				t.Fatalf("trial %d: roots not sorted: %v", trial, got)
			}
		}
		if _, ok := rebuildFromRoots(got[:]); !ok {
			t.Fatalf("trial %d: roots %v are not real or exact conjugate pairs", trial, got)
		}
		checkAgainstKnownRoots(t, trial, monic, known, got[:], 1000)
	}
}

func TestCubicRoots_noAllocation(t *testing.T) {
	var sink [3]complex128
	if allocs := testing.AllocsPerRun(100, func() { sink, _ = CubicRoots(1, -6, 6, -2) }); allocs != 0 {
		t.Errorf("CubicRoots allocates %v times per call, want 0", allocs)
	}
	benchSinkComplex = sink[0]
}

func BenchmarkCubicRoots(b *testing.B) {
	b.Run("scalar", func(b *testing.B) {
		b.ReportAllocs()
		a := 1.0
		for b.Loop() {
			r, _ := CubicRoots(a, -6, 6, -2)
			benchSinkComplex = r[0]
		}
	})
}
