package algebra

import (
	"errors"
	"math"
	"math/cmplx"
	"math/rand/v2"
	"slices"
	"testing"
)

// checkRoots fails unless got has len(want) roots, each within tol·max(1, |want[i]|) of want[i] (both
// sorted by real part, then imaginary part).
func checkRoots(t *testing.T, label string, got, want []complex128, tol float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s: %d roots %v, want %d", label, len(got), got, len(want))
	}
	for i := range want {
		if !closeTo(got[i], want[i], tol) {
			t.Errorf("%s: root %d = %v, want %v (all: %v)", label, i, got[i], want[i], got)
		}
	}
}

// Worked example: DLMF §1.11(iii), Example (quartic): z⁴ − 4z³ + 5z + 2 = 0 has the roots ½(3 ± √17) and
// ½(1 ± √5). Tolerance 1e-13 relative: a simple root is accepted once the correction is below 4u·|z|
// (≈ 4e-15 for the largest root here); the margin covers the last step's own rounding.
func TestPolynomialRoots_workedExample(t *testing.T) {
	s5, s17 := math.Sqrt(5), math.Sqrt(17)
	want := []complex128{complex((1-s5)/2, 0), complex((3-s17)/2, 0), complex((1+s5)/2, 0), complex((3+s17)/2, 0)}
	got, err := PolynomialRoots([]float64{2, 5, 0, -4, 1})
	if err != nil {
		t.Fatal(err)
	}
	checkRoots(t, "DLMF quartic", got, want, 1e-13)
	for i, z := range got {
		if imag(z) != 0 {
			t.Errorf("root %d = %v: a real root must have imaginary part exactly 0", i, z)
		}
	}
}

func TestPolynomialRoots_cases(t *testing.T) {
	s3 := math.Sqrt(3) / 2
	// Tolerances relative to the largest root (at least 1). Simple roots: the iteration stops when the
	// correction is below 4u·|z|, a few ulps, so 1e-14. Multiple roots: the conditioning, noted per case.
	tests := []struct {
		name   string
		coeffs []float64
		want   []complex128
		tol    float64
	}{
		{"degree 1", []float64{6, 3}, []complex128{-2}, 1e-14},
		{"x² + 1", []float64{1, 0, 1}, []complex128{complex(0, -1), complex(0, 1)}, 1e-14},
		{"x⁴ − 1", []float64{-1, 0, 0, 0, 1}, []complex128{-1, complex(0, -1), complex(0, 1), 1}, 1e-14},
		{"x³ + 1", []float64{1, 0, 0, 1}, []complex128{-1, complex(0.5, -s3), complex(0.5, s3)}, 1e-14},
		{"leading zeros give exact zero roots", []float64{0, 0, -1, 1}, []complex128{0, 0, 1}, 1e-14},
		{"xⁿ: every root exactly 0", []float64{0, 0, 0, 0, 7}, []complex128{0, 0, 0, 0}, 0},
		{"double root", []float64{2, -3, 0, 1}, []complex128{-2, 1, 1}, 1e-7},         // (x−1)²(x+2): error ~ √u
		{"fourfold root", []float64{1, -4, 6, -4, 1}, []complex128{1, 1, 1, 1}, 1e-3}, // (x−1)⁴: error ~ u^(1/4)
		{"huge coefficients", []float64{-0x1p900, 0, 0x1p900}, []complex128{-1, 1}, 1e-14},
		{"huge roots", []float64{-0x1p800, 0, 1}, []complex128{-0x1p400, 0x1p400}, 1e-14},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PolynomialRoots(tt.coeffs)
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("PolynomialRoots(%v) = %v, want %v", tt.coeffs, got, tt.want)
			}
			scale := 0.0
			for _, w := range tt.want {
				scale = max(scale, cmplx.Abs(w))
			}
			for i := range got {
				if cmplx.Abs(got[i]-tt.want[i]) > tt.tol*max(scale, 1) {
					t.Errorf("PolynomialRoots(%v) = %v, want %v", tt.coeffs, got, tt.want)
					break
				}
			}
		})
	}
}

// (x − 1e-6)(x − 1e6) = x² − (1e6 + 1e-6)·x + 1: roots twelve orders of magnitude apart. The start
// points (one per circle of the Newton polygon) find both to full relative precision.
func TestPolynomialRoots_widelySpreadRoots(t *testing.T) {
	got, err := PolynomialRoots([]float64{1, -(1e6 + 1e-6), 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || math.Abs(real(got[0])-1e-6) > 1e-20 || math.Abs(real(got[1])-1e6) > 1e-8 ||
		imag(got[0]) != 0 || imag(got[1]) != 0 {
		t.Errorf("PolynomialRoots = %v, want [1e-06 1e+06]", got)
	}
}

// Wilkinson 1959 (Numerische Mathematik 1:150–166, doi:10.1007/BF01386381): the roots of (x−1)(x−2)…(x−10)
// are ill conditioned; perturbations of order u in the coefficients move them by up to about 1e-10 here.
// Tolerance 1e-7 absolute.
func TestPolynomialRoots_wilkinson10(t *testing.T) {
	coeffs := []float64{1}
	want := make([]complex128, 10)
	for k := 1; k <= 10; k++ {
		var err error
		if coeffs, err = PolynomialMultiply(coeffs, []float64{-float64(k), 1}); err != nil {
			t.Fatal(err)
		}
		want[k-1] = complex(float64(k), 0)
	}
	got, err := PolynomialRoots(coeffs)
	if err != nil {
		t.Fatal(err)
	}
	for i := range want {
		if cmplx.Abs(got[i]-want[i]) > 1e-7 {
			t.Errorf("Wilkinson root %d = %v, want %v", i, got[i], want[i])
		}
	}
}

// Wilkinson 1959: the coefficients of (x−1)(x−2)…(x−20), rounded to float64 by PolynomialMultiply. The
// roots of that float polynomial are all real but up to about 1.7e-4 from the integers (the rounding of
// the coefficients is a relative perturbation of order u, magnified by condition numbers up to ~1e12).
// Reference: each root of the float polynomial by Newton's method in bigPrec-bit arithmetic from the
// integer k (60 steps; the roots are within 2e-4 of k and 1 apart, so Newton converges quadratically).
// Tolerance per root: the compensated evaluation's error bound γ̃₂ₙ²·Σ|aᵢ||r|ⁱ (Graillat and
// Ménissier-Morain 2012, eq. (5.17), γ̃₂ₙ ≤ 6n·u) divided by |p′(r)|, plus 4u·|r|, times 10 for the
// stagnation stop (the iteration stops inside the noise region, not at its centre). That is below 1e-11
// for every root, far tighter than the 2e-3 a plain-Horner method reaches here.
func TestPolynomialRoots_wilkinson20(t *testing.T) {
	coeffs := []float64{1}
	for k := 1; k <= 20; k++ {
		var err error
		if coeffs, err = PolynomialMultiply(coeffs, []float64{-float64(k), 1}); err != nil {
			t.Fatal(err)
		}
	}
	ref := make([]float64, 20)
	for k := 1; k <= 20; k++ {
		x := newBig(float64(k))
		for range 60 {
			p, dp := newBig(0), newBig(0)
			for i := len(coeffs) - 1; i >= 0; i-- {
				dp.Mul(dp, x)
				dp.Add(dp, p)
				p.Mul(p, x)
				p.Add(p, newBig(coeffs[i]))
			}
			x.Sub(x, p.Quo(p, dp))
		}
		ref[k-1] = bigToFloat(x)
		if math.Abs(ref[k-1]-float64(k)) > 2e-4 {
			t.Fatalf("reference root %d = %v, more than 2e-4 from %d", k, ref[k-1], k)
		}
	}
	got, err := PolynomialRoots(coeffs)
	if err != nil {
		t.Fatal(err)
	}
	g := 6 * 20 * unitRoundoff
	for k, r := range ref {
		s, deriv := 0.0, 1.0
		for i := len(coeffs) - 1; i >= 0; i-- {
			s = s*math.Abs(r) + math.Abs(coeffs[i])
		}
		for j, rj := range ref {
			if j != k {
				deriv *= math.Abs(r - rj)
			}
		}
		tol := 10 * (g*g*s/deriv + 4*unitRoundoff*math.Abs(r))
		if imag(got[k]) != 0 || math.Abs(real(got[k])-r) > tol {
			t.Errorf("root %d = %v, want %v ± %g", k+1, got[k], r, tol)
		}
	}
}

// A root below the float64 range is not an error: 1e300·x + 1e-300 has the root −1e-600, which underflows
// to 0; the start radius is exactly 0 and the root is returned as 0.
func TestPolynomialRoots_underflowingRoot(t *testing.T) {
	got, err := PolynomialRoots([]float64{1e-300, 1e300})
	if err != nil || len(got) != 1 || got[0] != 0 {
		t.Errorf("PolynomialRoots([1e-300 1e300]) = (%v, %v), want ([0], nil)", got, err)
	}
	// Mixed: (x − 1)·(1e300·x + 1e-300) = 1e300·x² + (1e-300 − 1e300)·x − 1e-300 keeps the root 1.
	got, err = PolynomialRoots([]float64{-1e-300, 1e-300 - 1e300, 1e300})
	if err != nil || len(got) != 2 || got[0] != 0 || math.Abs(real(got[1])-1) > 4*unitRoundoff || imag(got[1]) != 0 {
		t.Errorf("PolynomialRoots((x − 1)(1e300·x + 1e-300)) = (%v, %v), want ([0 1], nil)", got, err)
	}
}

// snapAndPairConjugates on fixed approximations: a conjugate-looking pair at a double root whose real part
// passes the stopping rule becomes real; a near-conjugate pair becomes an exact conjugate pair (the mean).
func TestPolynomialRoots_snapAndPair(t *testing.T) {
	doubleRoot := []float64{1, -2, 1} // (x − 1)²: p(1) = 0 exactly
	z := []complex128{complex(1, 1e-9), complex(1, -1e-9)}
	snapAndPairConjugates(doubleRoot, z, make([]bool, 2))
	if z[0] != 1 || z[1] != 1 {
		t.Errorf("double root: got %v, want [1 1] with imaginary parts exactly 0", z)
	}
	pair := []float64{5, -4, 1} // x² − 4x + 5, roots 2 ± i
	z = []complex128{complex(2, 1), complex(2+0x1p-40, -1-0x1p-40)}
	snapAndPairConjugates(pair, z, make([]bool, 2))
	if z[0] != cmplx.Conj(z[1]) || real(z[0]) != 2+0x1p-41 || imag(z[0]) != 1+0x1p-41 {
		t.Errorf("pair: got %v, want exact conjugates 2+2⁻⁴¹ ± (1+2⁻⁴¹)i", z)
	}
}

func TestPolynomialRoots_errors(t *testing.T) {
	nan, inf := math.NaN(), math.Inf(1)
	tests := []struct {
		name   string
		coeffs []float64
		want   error
	}{
		{"nil", nil, ErrEmpty},
		{"constant has no roots", []float64{5}, ErrEmpty},
		{"zero leading coefficient", []float64{1, 2, 0}, ErrDegenerate},
		{"all zero", []float64{0, 0}, ErrDegenerate},
		{"NaN coefficient", []float64{1, nan, 1}, ErrNotFinite},
		{"Inf leading coefficient", []float64{1, inf}, ErrNotFinite},
		{"coefficients too large to evaluate", []float64{math.MaxFloat64, math.MaxFloat64}, ErrOverflow},
		{"root modulus overflows", []float64{1e300, 1e-300}, ErrOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PolynomialRoots(tt.coeffs)
			if !errors.Is(err, tt.want) || got != nil {
				t.Errorf("PolynomialRoots(%v) = (%v, %v), want (nil, %v)", tt.coeffs, got, err, tt.want)
			}
		})
	}
}

// ErrNoConvergence: with the sweep limit lowered to 1 the iteration cannot finish; the best approximations
// are still returned, sorted, with the error.
func TestPolynomialRoots_noConvergence(t *testing.T) {
	coeffs := []float64{-1, 0, 0, 0, 0, 0, 0, 0, 1} // x⁸ − 1
	got, err := polynomialRoots(coeffs, 1)
	if !errors.Is(err, ErrNoConvergence) || len(got) != 8 {
		t.Fatalf("polynomialRoots(x⁸ − 1, 1 sweep) = (%v, %v), want 8 roots and ErrNoConvergence", got, err)
	}
	for i := 1; i < len(got); i++ {
		if compareRoots(got[i-1], got[i]) > 0 {
			t.Errorf("roots not sorted: %v", got)
		}
	}
}

func TestPolynomialRoots_doesNotModifyInput(t *testing.T) {
	coeffs := []float64{0, 2, 5, 0, -4, 1}
	if _, err := PolynomialRoots(coeffs); err != nil {
		t.Fatal(err)
	}
	if want := []float64{0, 2, 5, 0, -4, 1}; !slices.Equal(coeffs, want) {
		t.Errorf("input modified: %v", coeffs)
	}
}

// Oracle: random polynomials of degree 2–14 built from known roots (real, or conjugate pairs) in
// [−3, 3] at least 0.1 apart. Each known root must have a computed root within rootTolerance with
// c = 1000: the iteration stops at a correction below 4u·|z|, or in the compensated evaluation's noise
// region, whose residual (6n·u)²·Σ|aᵢ||z|ⁱ is far below u·Σ|aᵢ||z|ⁱ; building the coefficients adds
// about n·u·Σ|aᵢ||z|ⁱ (n ≤ 14), so 1000 leaves a wide margin. The roots must be sorted and be real or exact conjugate pairs.
func TestPolynomialRoots_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		degree := 2 + rng.IntN(13)
		known := randomRoots(rng, degree, 3, 0.1)
		monic, _ := rebuildFromRoots(known)
		got, err := PolynomialRoots(monic)
		if err != nil {
			t.Fatalf("trial %d (degree %d): %v", trial, degree, err)
		}
		for i := 1; i < len(got); i++ {
			if compareRoots(got[i-1], got[i]) > 0 {
				t.Fatalf("trial %d: roots not sorted: %v", trial, got)
			}
		}
		if rebuilt, ok := rebuildFromRoots(got); !ok || len(rebuilt) != len(monic) {
			t.Fatalf("trial %d: roots %v are not real or exact conjugate pairs", trial, got)
		}
		checkAgainstKnownRoots(t, trial, monic, known, got, 1000)
	}
}

// At most four allocations: the result, the hull indices, the converged flags and the last corrections
// (the compiler may place a small scratch slice on the stack, so fewer is allowed).
func TestPolynomialRoots_atMostFourAllocations(t *testing.T) {
	coeffs := []float64{2, 5, 0, -4, 1}
	if allocs := testing.AllocsPerRun(50, func() { _, _ = PolynomialRoots(coeffs) }); allocs > 4 {
		t.Errorf("PolynomialRoots allocates %v times per call, want at most 4 (result, hull, flags, corrections)", allocs)
	}
}

func BenchmarkPolynomialRoots(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	for _, size := range []struct {
		name string
		n    int
	}{{"1e1", 10}, {"1e2", 100}, {"1e3", 1_000}} {
		coeffs := randomCoeffs(rng, size.n+1)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				r, _ := PolynomialRoots(coeffs)
				benchSinkComplex = r[0]
			}
		})
	}
}
