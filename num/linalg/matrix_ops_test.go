package linalg

import (
	"errors"
	"fmt"
	"math"
	"math/cmplx"
	"math/rand/v2"
	"testing"
)

// randMat returns an m×n matrix with entries uniform in [−1, 1).
func randMat(r *rand.Rand, m, n int) []float64 {
	a := make([]float64, m*n)
	for i := range a {
		a[i] = sym(r)
	}
	return a
}

// Worked examples (hand-checked, as in the doc comments).
func TestMatElementwise_examples(t *testing.T) {
	a := []float64{1, 2, 3, 4}
	b := []float64{5, 6, 7, 8}
	cases := []struct {
		name string
		got  func() ([]float64, error)
		want []float64
	}{
		{"MatAdd", func() ([]float64, error) { return MatAdd(nil, a, b, 2, 2) }, []float64{6, 8, 10, 12}},
		{"MatSub", func() ([]float64, error) { return MatSub(nil, b, a, 2, 2) }, []float64{4, 4, 4, 4}},
		{"MatScale", func() ([]float64, error) { return MatScale(nil, a, 2, 2, 3) }, []float64{3, 6, 9, 12}},
		{"MatHadamard", func() ([]float64, error) { return MatHadamard(nil, a, b, 2, 2) }, []float64{5, 12, 21, 32}},
		{"VecOuter", func() ([]float64, error) { return VecOuter(nil, []float64{1, 2}, []float64{3, 4, 5}) }, []float64{3, 4, 5, 6, 8, 10}},
		{"MatKronecker", func() ([]float64, error) {
			return MatKronecker(nil, a, 2, 2, []float64{0, 5, 6, 7}, 2, 2)
		}, []float64{0, 5, 0, 10, 6, 7, 12, 14, 0, 15, 0, 20, 18, 21, 24, 28}},
		{"MatPower Fibonacci", func() ([]float64, error) { return MatPower(nil, []float64{1, 1, 1, 0}, 2, 5) }, []float64{8, 5, 5, 3}},
		{"MatPower k=0", func() ([]float64, error) { return MatPower(nil, a, 2, 0) }, []float64{1, 0, 0, 1}},
		{"MatPower k=1", func() ([]float64, error) { return MatPower(nil, a, 2, 1) }, a},
	}
	for _, c := range cases {
		got, err := c.got()
		if err != nil || !eqSlice(got, c.want, 0) {
			t.Errorf("%s = (%v, %v), want %v", c.name, got, err, c.want)
		}
	}
}

// In place: dst may be an input for the elementwise operations.
func TestMatElementwise_inPlace(t *testing.T) {
	a := []float64{1, 2, 3, 4}
	if _, err := MatAdd(a, a, []float64{1, 1, 1, 1}, 2, 2); err != nil || !eqSlice(a, []float64{2, 3, 4, 5}, 0) {
		t.Errorf("MatAdd in place → %v, %v", a, err)
	}
}

// Oracle (fixed seed): Aᵏ by squaring equals k repeated products; Kronecker mixed product
// (A⊗B)(C⊗D) = (AC)⊗(BD) (Van Loan 2000, §1).
func TestMatProducts_oracle(t *testing.T) {
	r := newRNG()
	for trial := range 50 {
		n := 1 + trial%5
		a := randMat(r, n, n)
		k := trial % 9
		want, _ := MatIdentity(nil, n)
		for range k {
			want, _ = MatMul(nil, want, n, n, a, n)
		}
		got, err := MatPower(nil, a, n, k)
		if err != nil || !eqSlice(got, want, 1e-12) {
			t.Fatalf("MatPower n=%d k=%d = %v, want %v (%v)", n, k, got, want, err)
		}
	}
	for range 20 {
		a, b, c, d := randMat(r, 2, 3), randMat(r, 2, 2), randMat(r, 3, 2), randMat(r, 2, 3)
		ab, _ := MatKronecker(nil, a, 2, 3, b, 2, 2)    // 4×6
		cd, _ := MatKronecker(nil, c, 3, 2, d, 2, 3)    // 6×6
		lhs, _ := MatMul(nil, ab, 4, 6, cd, 6)          // 4×6
		ac, _ := MatMul(nil, a, 2, 3, c, 2)             // 2×2
		bd, _ := MatMul(nil, b, 2, 2, d, 3)             // 2×3
		rhs, _ := MatKronecker(nil, ac, 2, 2, bd, 2, 3) // 4×6
		if !eqSlice(lhs, rhs, 1e-12) {
			t.Fatalf("mixed product: %v vs %v", lhs, rhs)
		}
	}
}

func TestMatNorms_examples(t *testing.T) {
	a := []float64{1, -2, 3, 4}
	if got, err := MatNormFrobenius([]float64{1, 2, 3, 4}, 2, 2); err != nil || !near(got, math.Sqrt(30), 1e-15) {
		t.Errorf("Frobenius = %v, %v", got, err)
	}
	if got, err := MatNorm1(a, 2, 2); err != nil || got != 6 {
		t.Errorf("MatNorm1 = %v, %v, want 6", got, err)
	}
	if got, err := MatNormInf(a, 2, 2); err != nil || got != 7 {
		t.Errorf("MatNormInf = %v, %v, want 7", got, err)
	}
	if got, err := MatNorm2([]float64{3, 0, 4, 5}, 2, 2); err != nil || !near(got, 3*math.Sqrt(5), 1e-14) {
		t.Errorf("MatNorm2 = %v, %v, want 3√5", got, err)
	}
	// Frobenius does not overflow for huge entries (power-of-two scaling).
	if got, err := MatNormFrobenius([]float64{1e300, 1e300}, 1, 2); err != nil || !near(got, math.Sqrt2*1e300, 1e-15) {
		t.Errorf("Frobenius huge = %v, %v", got, err)
	}
	// ‖A‖₂ ≤ ‖A‖_F ≤ √r‖A‖₂ and ‖A‖₂² ≤ ‖A‖₁‖A‖∞ (Golub–Van Loan §2.3) on random matrices.
	r := newRNG()
	for range 40 {
		m, n := 1+r.IntN(5), 1+r.IntN(5)
		x := randMat(r, m, n)
		f, _ := MatNormFrobenius(x, m, n)
		s, _ := MatNorm2(x, m, n)
		o, _ := MatNorm1(x, m, n)
		i, _ := MatNormInf(x, m, n)
		if s > f*(1+1e-12) || f > math.Sqrt(float64(min(m, n)))*s*(1+1e-12) || s*s > o*i*(1+1e-12) {
			t.Fatalf("norm inequalities fail: 2=%v F=%v 1=%v inf=%v", s, f, o, i)
		}
	}
}

func TestMatFixed_examples(t *testing.T) {
	m3 := Mat3{1, 2, 3, 4, 5, 6, 7, 8, 10}
	if got, err := Mat3Minor(m3, 0, 0); err != nil || got != 2 {
		t.Errorf("Mat3Minor(0,0) = %v, %v, want 2", got, err)
	}
	if got, err := MatMinor(m3[:], 3, 1, 2); err != nil || !near(got, -6, 1e-14) {
		t.Errorf("MatMinor(1,2) = %v, %v, want −6", got, err)
	}
	if got, err := Mat2Adjugate(Mat2{1, 2, 3, 4}); err != nil || got != (Mat2{4, -2, -3, 1}) {
		t.Errorf("Mat2Adjugate = %v, %v", got, err)
	}
	if got, err := Mat3Adjugate(Mat3{1, 2, 3, 0, 1, 4, 5, 6, 0}); err != nil || got != (Mat3{-24, 18, 5, 20, -15, -4, -5, 4, 1}) {
		t.Errorf("Mat3Adjugate = %v, %v", got, err)
	}
	if got, err := Mat4Adjugate(Mat4{1, 0, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, 0, 0, 0, 4}); err != nil || got != (Mat4{24, 0, 0, 0, 0, 12, 0, 0, 0, 0, 8, 0, 0, 0, 0, 6}) {
		t.Errorf("Mat4Adjugate = %v, %v", got, err)
	}
	if got, err := Mat2SolveCramer(Mat2{2, 1, 1, 3}, Vec2{5, 10}); err != nil || !near(got[0], 1, 1e-15) || !near(got[1], 3, 1e-15) {
		t.Errorf("Mat2SolveCramer = %v, %v", got, err)
	}
	if got, err := Mat3SolveCramer(Mat3{1, 1, 1, 0, 2, 5, 2, 5, -1}, Vec3{6, -4, 27}); err != nil || !nearV3(got, Vec3{5, 3, -2}, 1e-14) {
		t.Errorf("Mat3SolveCramer = %v, %v", got, err)
	}
	if got, err := Mat2CharPoly(Mat2{4, 1, 2, 3}); err != nil || got != [3]float64{10, -7, 1} {
		t.Errorf("Mat2CharPoly = %v, %v", got, err)
	}
	if got, err := Mat3CharPoly(Mat3{1, 0, 0, 0, 2, 0, 0, 0, 3}); err != nil || got != [4]float64{-6, 11, -6, 1} {
		t.Errorf("Mat3CharPoly = %v, %v", got, err)
	}
	if got, err := MatCharPolyFaddeevLeVerrier(nil, []float64{2, 1, 1, 2}, 2); err != nil || !eqSlice(got, []float64{3, -4, 1}, 1e-15) {
		t.Errorf("Faddeev = %v, %v", got, err)
	}
	if got, err := MatCofactorMatrix(nil, []float64{1, 2, 3, 4}, 2); err != nil || !eqSlice(got, []float64{4, -3, -2, 1}, 1e-14) {
		t.Errorf("MatCofactorMatrix = %v, %v", got, err)
	}
	if got, err := MatAdjugate(nil, []float64{1, 2, 2, 4}, 2); err != nil || !eqSlice(got, []float64{4, -2, -2, 1}, 1e-14) {
		t.Errorf("MatAdjugate singular = %v, %v", got, err)
	}
	if got, err := MatAdjugate(nil, []float64{7}, 1); err != nil || got[0] != 1 {
		t.Errorf("MatAdjugate 1×1 = %v, %v", got, err)
	}
}

// Oracle (fixed seed): Cramer agrees with MatSolve; det(adj A) = det(A)^(n−1); A·adj A = det(A)·I; the
// Faddeev–LeVerrier polynomial equals the closed forms for n = 2, 3; adjugates of fixed sizes equal MatAdjugate.
func TestMatFixed_oracle(t *testing.T) {
	r := newRNG()
	for range 200 {
		m3 := randM3(r)
		b := Vec3{sym(r), sym(r), sym(r)}
		x, err := Mat3SolveCramer(m3, b)
		y, err2 := MatSolve(nil, m3[:], 3, b[:])
		if err != nil || err2 != nil || !eqSlice(x[:], y, 1e-8) {
			t.Fatalf("Cramer %v vs MatSolve %v (%v, %v)", x, y, err, err2)
		}
		adj, _ := Mat3Adjugate(m3)
		gen, _ := MatAdjugate(nil, m3[:], 3)
		if !eqSlice(adj[:], gen, 1e-9) {
			t.Fatalf("Mat3Adjugate %v vs MatAdjugate %v", adj, gen)
		}
		d, _ := Mat3Det(m3)
		da, _ := Mat3Det(adj)
		if !near(da, d*d, 1e-10) {
			t.Fatalf("det adj = %v, det² = %v", da, d*d)
		}
		prod, _ := Mat3Mul(m3, adj)
		want := Mat3{d, 0, 0, 0, d, 0, 0, 0, d}
		if !eqSlice(prod[:], want[:], 1e-12) {
			t.Fatalf("A·adj A = %v, want det·I", prod)
		}
		p, _ := Mat3CharPoly(m3)
		f, _ := MatCharPolyFaddeevLeVerrier(nil, m3[:], 3)
		if !eqSlice(p[:], f, 1e-12) {
			t.Fatalf("Mat3CharPoly %v vs Faddeev %v", p, f)
		}
	}
	for range 100 {
		var m4 Mat4
		for i := range m4 {
			m4[i] = sym(r)
		}
		adj, err := Mat4Adjugate(m4)
		gen, err2 := MatAdjugate(nil, m4[:], 4)
		if err != nil || err2 != nil || !eqSlice(adj[:], gen, 1e-9) {
			t.Fatalf("Mat4Adjugate %v vs MatAdjugate %v", adj, gen)
		}
		d, _ := Mat4Det(m4)
		da, _ := Mat4Det(adj)
		if !near(da, d*d*d, 1e-9) {
			t.Fatalf("det adj(4×4) = %v, det³ = %v", da, d*d*d)
		}
	}
	// Singular n×n: the cofactor path; rank-(n−1) matrices have adj of rank 1 with A·adj A = 0.
	for range 30 {
		n := 3 + r.IntN(2)
		u, v := randMat(r, n, n-1), randMat(r, n-1, n)
		a, _ := MatMul(nil, u, n, n-1, v, n)
		adj, err := MatAdjugate(nil, a, n)
		if err != nil {
			t.Fatal(err)
		}
		z, _ := MatMul(nil, a, n, n, adj, n)
		if nz, _ := MatNormFrobenius(z, n, n); nz > 1e-9 {
			t.Fatalf("A·adj A for singular A has norm %v", nz)
		}
	}
}

// The characteristic polynomial's roots are the eigenvalues.
func TestCharPoly_rootsAreEigenvalues(t *testing.T) {
	s := Mat3{2, -1, 0, -1, 2, -1, 0, -1, 2}
	info, err := Mat3Characteristics(s, -1)
	if err != nil {
		t.Fatal(err)
	}
	want := []float64{2 - math.Sqrt2, 2, 2 + math.Sqrt2}
	for k, w := range want {
		if !near(real(info.Eigenvalues[k]), w, 1e-12) || imag(info.Eigenvalues[k]) != 0 {
			t.Errorf("eigenvalue %d = %v, want %v", k, info.Eigenvalues[k], w)
		}
	}
	jac, _, _ := MatSymmetricEigenJacobi(nil, nil, s[:], 3, -1, false)
	if !eqSlice(jac, want, 1e-14) {
		t.Errorf("Jacobi = %v, want %v", jac, want)
	}
}

func TestMatCharacteristics_examples(t *testing.T) {
	i2, err := Mat2Characteristics(Mat2{4, 1, 2, 3}, -1)
	if err != nil || i2.Det != 10 || i2.Trace != 7 || i2.Rank != 2 || !i2.Invertible ||
		!eqSlice(i2.Inverse[:], []float64{0.3, -0.1, -0.2, 0.4}, 1e-15) || i2.Eigenvalues != [2]complex128{2, 5} {
		t.Errorf("Mat2Characteristics = %+v, %v", i2, err)
	}
	s2, err := Mat2Characteristics(Mat2{1, 2, 2, 4}, -1)
	if err != nil || s2.Invertible || s2.Rank != 1 || s2.Det != 0 {
		t.Errorf("singular Mat2Characteristics = %+v, %v", s2, err)
	}
	i3, err := Mat3Characteristics(Mat3{1, 0, 0, 0, 2, 0, 0, 0, 3}, -1)
	if err != nil || i3.Det != 6 || i3.Trace != 6 || i3.Rank != 3 || !i3.Invertible ||
		!near(real(i3.Eigenvalues[0]), 1, 1e-14) || !near(real(i3.Eigenvalues[2]), 3, 1e-14) {
		t.Errorf("Mat3Characteristics = %+v, %v", i3, err)
	}
}

func TestMat2Eigen_cases(t *testing.T) {
	res, err := Mat2Eigen(Mat2{4, 1, 2, 3})
	if err != nil || res.Values != [2]complex128{2, 5} || res.Defective {
		t.Fatalf("Mat2Eigen = %+v, %v", res, err)
	}
	rot, err := Mat2Eigen(Mat2{0, -1, 1, 0})
	if err != nil || rot.Values[0] != complex(0, -1) || rot.Values[1] != complex(0, 1) {
		t.Errorf("rotation eigenvalues = %v, %v", rot.Values, err)
	}
	def, err := Mat2Eigen(Mat2{1, 1, 0, 1})
	if err != nil || !def.Defective {
		t.Errorf("Jordan block: Defective = %v, %v", def.Defective, err)
	}
	id, err := Mat2Eigen(Mat2{3, 0, 0, 3})
	if err != nil || id.Defective || id.Vectors[0] != [2]complex128{1, 0} || id.Vectors[1] != [2]complex128{0, 1} {
		t.Errorf("3I = %+v, %v", id, err)
	}
	// Oracle (fixed seed): A·v = λ·v and ‖v‖ = 1, real or complex pairs.
	r := newRNG()
	for range 300 {
		m := Mat2{sym(r), sym(r), sym(r), sym(r)}
		e, err := Mat2Eigen(m)
		if err != nil {
			t.Fatal(err)
		}
		for k := range 2 {
			lam, v := e.Values[k], e.Vectors[k]
			r0 := complex(m[0], 0)*v[0] + complex(m[1], 0)*v[1] - lam*v[0]
			r1 := complex(m[2], 0)*v[0] + complex(m[3], 0)*v[1] - lam*v[1]
			if cmplx.Abs(r0) > 1e-12 || cmplx.Abs(r1) > 1e-12 || math.Abs(cnorm2(v)-1) > 1e-14 {
				t.Fatalf("m=%v λ=%v v=%v residual %v %v", m, lam, v, r0, r1)
			}
		}
	}
}

func TestMatChecks_examples(t *testing.T) {
	type check struct {
		name string
		got  func() (bool, error)
		want bool
	}
	cases := []check{
		{"symmetric", func() (bool, error) { return MatIsSymmetric([]float64{1, 2, 2, 3}, 2, -1) }, true},
		{"not symmetric", func() (bool, error) { return MatIsSymmetric([]float64{1, 2, 3, 4}, 2, -1) }, false},
		{"diagonal", func() (bool, error) { return MatIsDiagonal([]float64{2, 0, 0, 3}, 2, 2, -1) }, true},
		{"not diagonal tol 0", func() (bool, error) { return MatIsDiagonal([]float64{2, 1e-20, 0, 3}, 2, 2, 0) }, false},
		{"upper", func() (bool, error) { return MatIsUpperTriangular([]float64{1, 2, 0, 3}, 2, 2, -1) }, true},
		{"not upper", func() (bool, error) { return MatIsUpperTriangular([]float64{1, 2, 4, 3}, 2, 2, -1) }, false},
		{"lower", func() (bool, error) { return MatIsLowerTriangular([]float64{1, 0, 2, 3}, 2, 2, -1) }, true},
		{"not lower", func() (bool, error) { return MatIsLowerTriangular([]float64{1, 5, 2, 3}, 2, 2, -1) }, false},
		{"orthogonal", func() (bool, error) { return MatIsOrthogonal([]float64{0, -1, 1, 0}, 2, -1) }, true},
		{"not orthogonal", func() (bool, error) { return MatIsOrthogonal([]float64{1, 1, 0, 1}, 2, -1) }, false},
		{"positive definite", func() (bool, error) { return MatIsPositiveDefinite([]float64{2, -1, -1, 2}, 2, -1) }, true},
		{"indefinite", func() (bool, error) { return MatIsPositiveDefinite([]float64{1, 2, 2, 1}, 2, -1) }, false},
		{"semidefinite", func() (bool, error) { return MatIsPositiveSemidefinite([]float64{1, 1, 1, 1}, 2, -1) }, true},
		{"semidefinite but not definite", func() (bool, error) { return MatIsPositiveDefinite([]float64{1, 1, 1, 1}, 2, -1) }, false},
		{"not semidefinite", func() (bool, error) { return MatIsPositiveSemidefinite([]float64{1, 2, 2, 1}, 2, -1) }, false},
		{"singular", func() (bool, error) { return MatIsSingular([]float64{1, 2, 2, 4}, 2, -1) }, true},
		{"nonsingular", func() (bool, error) { return MatIsSingular([]float64{1, 2, 3, 4}, 2, -1) }, false},
	}
	for _, c := range cases {
		got, err := c.got()
		if err != nil || got != c.want {
			t.Errorf("%s = (%v, %v), want %v", c.name, got, err, c.want)
		}
	}
	// Random rotations are orthogonal; their QR's Q is too.
	r := newRNG()
	for range 50 {
		q := randRotation(r)
		if ok, err := MatIsOrthogonal(q[:], 3, -1); err != nil || !ok {
			t.Fatalf("rotation %v not orthogonal (%v)", q, err)
		}
	}
}

func TestCovariance_examples(t *testing.T) {
	corr, std, err := CovarianceToCorrelation(nil, nil, []float64{4, 2, 2, 9}, 2)
	if err != nil || !eqSlice(std, []float64{2, 3}, 0) || !eqSlice(corr, []float64{1, 1.0 / 3, 1.0 / 3, 1}, 1e-15) {
		t.Errorf("CovarianceToCorrelation = %v, %v, %v", corr, std, err)
	}
	cov, err := CorrelationToCovariance(nil, corr, 2, std)
	if err != nil || !eqSlice(cov, []float64{4, 2, 2, 9}, 1e-15) {
		t.Errorf("CorrelationToCovariance = %v, %v", cov, err)
	}
	// NIST e-Handbook §6.5.4.1 example covariance matrix: the round trip reproduces it.
	s := []float64{0.025, 0.0075, 0.00175, 0.0075, 0.0070, 0.00135, 0.00175, 0.00135, 0.00043}
	c, sd, err := CovarianceToCorrelation(nil, nil, s, 3)
	if err != nil {
		t.Fatal(err)
	}
	back, err := CorrelationToCovariance(nil, c, 3, sd)
	if err != nil || !eqSlice(back, s, 1e-14) {
		t.Errorf("round trip = %v, %v", back, err)
	}
	if _, _, err := CovarianceToCorrelation(nil, nil, []float64{0, 0, 0, 1}, 2); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("zero variance: %v", err)
	}
	if _, err := CorrelationToCovariance(nil, corr, 2, []float64{1, -1}); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("negative std: %v", err)
	}
}

// Every documented error of the operations in this file.
func TestMatOps_errors(t *testing.T) {
	nan, inf := math.NaN(), math.Inf(1)
	cases := []struct {
		name string
		err  error
		want error
	}{
		{"MatAdd dims", second(MatAdd(nil, []float64{1}, []float64{1, 2}, 1, 1)), ErrDimension},
		{"MatAdd NaN", second(MatAdd(nil, []float64{nan}, []float64{1}, 1, 1)), ErrNotFinite},
		{"MatAdd overflow", second(MatAdd(nil, []float64{math.MaxFloat64}, []float64{math.MaxFloat64}, 1, 1)), ErrOverflow},
		{"MatSub overflow", second(MatSub(nil, []float64{math.MaxFloat64}, []float64{-math.MaxFloat64}, 1, 1)), ErrOverflow},
		{"MatScale k Inf", second(MatScale(nil, []float64{1}, 1, 1, inf)), ErrNotFinite},
		{"MatScale overflow", second(MatScale(nil, []float64{1e300}, 1, 1, 1e10)), ErrOverflow},
		{"MatHadamard overflow", second(MatHadamard(nil, []float64{1e200}, []float64{1e200}, 1, 1)), ErrOverflow},
		{"VecOuter NaN", second(VecOuter(nil, []float64{nan}, []float64{1})), ErrNotFinite},
		{"VecOuter overflow", second(VecOuter(nil, []float64{1e200}, []float64{1e200})), ErrOverflow},
		{"MatKronecker dims", second(MatKronecker(nil, []float64{1}, 1, 2, []float64{1}, 1, 1)), ErrDimension},
		{"MatKronecker NaN", second(MatKronecker(nil, []float64{1}, 1, 1, []float64{nan}, 1, 1)), ErrNotFinite},
		{"MatKronecker overflow", second(MatKronecker(nil, []float64{1e200}, 1, 1, []float64{1e200}, 1, 1)), ErrOverflow},
		{"MatPower k<0", second(MatPower(nil, []float64{1}, 1, -1)), ErrOutOfRange},
		{"MatPower n=0", second(MatPower(nil, nil, 0, 1)), ErrDimension},
		{"MatPower NaN", second(MatPower(nil, []float64{nan}, 1, 2)), ErrNotFinite},
		{"MatPower overflow", second(MatPower(nil, []float64{1e200}, 1, 2)), ErrOverflow},
		{"MatNorm1 NaN", second(MatNorm1([]float64{nan}, 1, 1)), ErrNotFinite},
		{"MatNorm1 overflow", second(MatNorm1([]float64{math.MaxFloat64, math.MaxFloat64}, 2, 1)), ErrOverflow},
		{"MatNormInf overflow", second(MatNormInf([]float64{math.MaxFloat64, math.MaxFloat64}, 1, 2)), ErrOverflow},
		{"MatNormFrobenius dims", second(MatNormFrobenius([]float64{1}, 2, 2)), ErrDimension},
		{"Mat3Minor index", second(Mat3Minor(Mat3{}, 3, 0)), ErrOutOfRange},
		{"Mat3Minor NaN", second(Mat3Minor(Mat3{nan}, 0, 1)), ErrNotFinite},
		{"MatMinor n=1", second(MatMinor([]float64{1}, 1, 0, 0)), ErrDimension},
		{"MatMinor index", second(MatMinor([]float64{1, 2, 3, 4}, 2, 2, 0)), ErrOutOfRange},
		{"Mat2Adjugate NaN", second(Mat2Adjugate(Mat2{nan})), ErrNotFinite},
		{"Mat3Adjugate overflow", second(Mat3Adjugate(Mat3{1e200, 0, 0, 0, 1e200, 0, 0, 0, 1e200})), ErrOverflow},
		{"Mat4Adjugate NaN", second(Mat4Adjugate(Mat4{inf})), ErrNotFinite},
		{"Mat2SolveCramer singular", second(Mat2SolveCramer(Mat2{1, 2, 2, 4}, Vec2{1, 1})), ErrSingular},
		{"Mat2SolveCramer NaN", second(Mat2SolveCramer(Mat2{1, 0, 0, 1}, Vec2{nan, 1})), ErrNotFinite},
		{"Mat3SolveCramer singular", second(Mat3SolveCramer(Mat3{1, 2, 3, 2, 4, 6, 1, 1, 1}, Vec3{1, 1, 1})), ErrSingular},
		{"Mat3SolveCramer overflow", second(Mat3SolveCramer(Mat3{1e-300, 0, 0, 0, 1, 0, 0, 0, 1}, Vec3{1e300, 0, 0})), ErrOverflow},
		{"Mat2CharPoly NaN", second(Mat2CharPoly(Mat2{nan})), ErrNotFinite},
		{"Mat2CharPoly overflow", second(Mat2CharPoly(Mat2{math.MaxFloat64, 0, 0, math.MaxFloat64})), ErrOverflow},
		{"Mat3CharPoly overflow", second(Mat3CharPoly(Mat3{1e200, 0, 0, 0, 1e200, 0, 0, 0, 1e200})), ErrOverflow},
		{"Faddeev NaN", second(MatCharPolyFaddeevLeVerrier(nil, []float64{nan}, 1)), ErrNotFinite},
		{"Faddeev overflow", second(MatCharPolyFaddeevLeVerrier(nil, []float64{1e200, 0, 0, 1e200}, 2)), ErrOverflow},
		{"MatAdjugate n=0", second(MatAdjugate(nil, nil, 0)), ErrDimension},
		{"MatCofactorMatrix NaN", second(MatCofactorMatrix(nil, []float64{nan, 0, 0, 1}, 2)), ErrNotFinite},
		{"Mat2Eigen NaN", second(Mat2Eigen(Mat2{nan})), ErrNotFinite},
		{"Mat2Eigen overflow", second(Mat2Eigen(Mat2{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64, math.MaxFloat64})), ErrOverflow},
		{"Mat2Characteristics NaN", second(Mat2Characteristics(Mat2{nan}, -1)), ErrNotFinite},
		{"Mat3Characteristics NaN", second(Mat3Characteristics(Mat3{nan}, -1)), ErrNotFinite},
		{"MatIsSymmetric dims", second(MatIsSymmetric([]float64{1}, 2, -1)), ErrDimension},
		{"MatIsDiagonal NaN", second(MatIsDiagonal([]float64{nan}, 1, 1, -1)), ErrNotFinite},
		{"MatIsUpperTriangular dims", second(MatIsUpperTriangular(nil, 0, 1, -1)), ErrDimension},
		{"MatIsLowerTriangular NaN", second(MatIsLowerTriangular([]float64{inf}, 1, 1, -1)), ErrNotFinite},
		{"MatIsOrthogonal NaN", second(MatIsOrthogonal([]float64{nan}, 1, -1)), ErrNotFinite},
		{"MatIsPositiveDefinite NaN", second(MatIsPositiveDefinite([]float64{nan}, 1, -1)), ErrNotFinite},
		{"MatIsPositiveSemidefinite dims", second(MatIsPositiveSemidefinite([]float64{1}, 0, -1)), ErrDimension},
		{"MatIsSingular NaN", second(MatIsSingular([]float64{nan}, 1, -1)), ErrNotFinite},
		{"CovarianceToCorrelation NaN", second3(CovarianceToCorrelation(nil, nil, []float64{nan}, 1)), ErrNotFinite},
		{"CorrelationToCovariance std len", second(CorrelationToCovariance(nil, []float64{1}, 1, nil)), ErrDimension},
		{"CorrelationToCovariance overflow", second(CorrelationToCovariance(nil, []float64{1}, 1, []float64{1e200})), ErrOverflow},
	}
	for _, c := range cases {
		if !errors.Is(c.err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, c.err, c.want)
		}
	}
}

// Fixed-size paths allocate nothing; n-dimensional paths reuse dst.
func TestMatOps_noAllocation(t *testing.T) {
	m2, m3 := Mat2{4, 1, 2, 3}, Mat3{1, 2, 3, 0, 1, 4, 5, 6, 0}
	var m4 Mat4
	for i := range m4 {
		m4[i] = float64(i%5) + 1
	}
	fixed := map[string]func(){
		"Mat3Minor":       func() { _, _ = Mat3Minor(m3, 1, 1) },
		"Mat2Adjugate":    func() { _, _ = Mat2Adjugate(m2) },
		"Mat3Adjugate":    func() { _, _ = Mat3Adjugate(m3) },
		"Mat4Adjugate":    func() { _, _ = Mat4Adjugate(m4) },
		"Mat2SolveCramer": func() { _, _ = Mat2SolveCramer(m2, Vec2{1, 2}) },
		"Mat3SolveCramer": func() { _, _ = Mat3SolveCramer(m3, Vec3{1, 2, 3}) },
		"Mat2CharPoly":    func() { _, _ = Mat2CharPoly(m2) },
		"Mat3CharPoly":    func() { _, _ = Mat3CharPoly(m3) },
		"Mat2Eigen":       func() { _, _ = Mat2Eigen(m2) },
	}
	for name, f := range fixed {
		if n := testing.AllocsPerRun(100, f); n != 0 {
			t.Errorf("%s allocs = %v, want 0", name, n)
		}
	}
	a, b := []float64{1, 2, 3, 4}, []float64{5, 6, 7, 8}
	dst := make([]float64, 16)
	reuse := map[string]func(){
		"MatAdd":       func() { _, _ = MatAdd(dst, a, b, 2, 2) },
		"MatSub":       func() { _, _ = MatSub(dst, a, b, 2, 2) },
		"MatScale":     func() { _, _ = MatScale(dst, a, 2, 2, 2) },
		"MatHadamard":  func() { _, _ = MatHadamard(dst, a, b, 2, 2) },
		"VecOuter":     func() { _, _ = VecOuter(dst, a, b) },
		"MatKronecker": func() { _, _ = MatKronecker(dst, a, 2, 2, b, 2, 2) },
		"MatNorm1":     func() { _, _ = MatNorm1(a, 2, 2) },
		"MatNormInf":   func() { _, _ = MatNormInf(a, 2, 2) },
		"MatNormFrob":  func() { _, _ = MatNormFrobenius(a, 2, 2) },
		"MatIsSym":     func() { _, _ = MatIsSymmetric(a, 2, -1) },
		"MatIsOrth":    func() { _, _ = MatIsOrthogonal(a, 2, -1) },
	}
	for name, f := range reuse {
		if n := testing.AllocsPerRun(100, f); n != 0 {
			t.Errorf("%s with dst allocs = %v, want 0", name, n)
		}
	}
}

func BenchmarkMatAdd(b *testing.B) {
	for _, n := range []int{32, 100, 316} { // 1e3, 1e4, 1e5 elements
		x, y, dst := randMat(newRNG(), n, n), randMat(newRNG(), n, n), make([]float64, n*n)
		b.Run(sizeName(n*n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = MatAdd(dst, x, y, n, n)
			}
		})
	}
}

func BenchmarkMatHadamard(b *testing.B) {
	for _, n := range []int{32, 100, 316} {
		x, y, dst := randMat(newRNG(), n, n), randMat(newRNG(), n, n), make([]float64, n*n)
		b.Run(sizeName(n*n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = MatHadamard(dst, x, y, n, n)
			}
		})
	}
}

func BenchmarkMatScale(b *testing.B) {
	x, dst := randMat(newRNG(), 100, 100), make([]float64, 10000)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = MatScale(dst, x, 100, 100, 1.5)
	}
}

func BenchmarkMatSub(b *testing.B) {
	x, y, dst := randMat(newRNG(), 100, 100), randMat(newRNG(), 100, 100), make([]float64, 10000)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = MatSub(dst, x, y, 100, 100)
	}
}

func BenchmarkVecOuter(b *testing.B) {
	u, v, dst := randMat(newRNG(), 1, 100), randMat(newRNG(), 1, 100), make([]float64, 10000)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = VecOuter(dst, u, v)
	}
}

func BenchmarkMatKronecker(b *testing.B) {
	x, y, dst := randMat(newRNG(), 10, 10), randMat(newRNG(), 10, 10), make([]float64, 10000)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = MatKronecker(dst, x, 10, 10, y, 10, 10)
	}
}

func BenchmarkMatPower(b *testing.B) {
	x, dst := randMat(newRNG(), 20, 20), make([]float64, 400)
	for i := range x {
		x[i] /= 20
	}
	for _, k := range []int{8, 64, 512} {
		b.Run(sizeName(k), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = MatPower(dst, x, 20, k)
			}
		})
	}
}

func BenchmarkMatNorms(b *testing.B) {
	x := randMat(newRNG(), 100, 100)
	b.Run("Frobenius", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = MatNormFrobenius(x, 100, 100)
		}
	})
	b.Run("One", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = MatNorm1(x, 100, 100)
		}
	})
	b.Run("Inf", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = MatNormInf(x, 100, 100)
		}
	})
}

func BenchmarkMatFixed(b *testing.B) {
	m2, m3 := Mat2{4, 1, 2, 3}, Mat3{1, 2, 3, 0, 1, 4, 5, 6, 0}
	var m4 Mat4
	for i := range m4 {
		m4[i] = float64(i%5) + 1
	}
	b.ReportAllocs()
	b.Run("Mat3Minor", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = Mat3Minor(m3, 1, 1)
		}
	})
	b.Run("Mat2Adjugate", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = Mat2Adjugate(m2)
		}
	})
	b.Run("Mat3Adjugate", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = Mat3Adjugate(m3)
		}
	})
	b.Run("Mat4Adjugate", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = Mat4Adjugate(m4)
		}
	})
	b.Run("Mat2SolveCramer", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = Mat2SolveCramer(m2, Vec2{1, 2})
		}
	})
	b.Run("Mat3SolveCramer", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = Mat3SolveCramer(m3, Vec3{1, 2, 3})
		}
	})
	b.Run("Mat2CharPoly", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = Mat2CharPoly(m2)
		}
	})
	b.Run("Mat3CharPoly", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = Mat3CharPoly(m3)
		}
	})
	b.Run("Mat2Eigen", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = Mat2Eigen(m2)
		}
	})
	b.Run("Mat2Characteristics", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = Mat2Characteristics(m2, -1)
		}
	})
	b.Run("Mat3Characteristics", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = Mat3Characteristics(m3, -1)
		}
	})
}

func BenchmarkMatCofactorN(b *testing.B) {
	for _, n := range []int{10, 30} {
		x := randMat(newRNG(), n, n)
		dst := make([]float64, n*n)
		b.Run(sizeName(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = MatAdjugate(dst, x, n)
			}
		})
	}
	x := randMat(newRNG(), 10, 10)
	b.Run("MatMinor/10", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = MatMinor(x, 10, 3, 4)
		}
	})
	b.Run("MatCofactorMatrix/10", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = MatCofactorMatrix(nil, x, 10)
		}
	})
	b.Run("Faddeev/10", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = MatCharPolyFaddeevLeVerrier(nil, x, 10)
		}
	})
}

func BenchmarkMatChecks(b *testing.B) {
	x := randMat(newRNG(), 100, 100)
	s, _ := MatMul(nil, x, 100, 100, first(MatTranspose(nil, x, 100, 100)), 100) // symmetric PSD
	for i := 0; i < 100; i++ {
		s[i*100+i] += 1
	}
	b.ReportAllocs()
	b.Run("Symmetric", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = MatIsSymmetric(s, 100, -1)
		}
	})
	b.Run("Diagonal", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = MatIsDiagonal(s, 100, 100, -1)
		}
	})
	b.Run("Upper", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = MatIsUpperTriangular(s, 100, 100, -1)
		}
	})
	b.Run("Lower", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = MatIsLowerTriangular(s, 100, 100, -1)
		}
	})
	b.Run("Orthogonal", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = MatIsOrthogonal(s, 100, -1)
		}
	})
	b.Run("PositiveDefinite", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = MatIsPositiveDefinite(s, 100, -1)
		}
	})
	b.Run("Singular", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = MatIsSingular(s, 100, -1)
		}
	})
	s20 := s[:400]
	for i := range 20 {
		for j := range 20 {
			s20[i*20+j] = s[i*100+j]
		}
	}
	b.Run("PositiveSemidefinite/20", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = MatIsPositiveSemidefinite(s20, 20, -1)
		}
	})
}

func BenchmarkCovariance(b *testing.B) {
	r := newRNG()
	x := randMat(r, 100, 100)
	cov, _ := MatMul(nil, x, 100, 100, first(MatTranspose(nil, x, 100, 100)), 100)
	corr, std := make([]float64, 10000), make([]float64, 100)
	dst := make([]float64, 10000)
	b.ReportAllocs()
	b.Run("ToCorrelation", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _, _ = CovarianceToCorrelation(corr, std, cov, 100)
		}
	})
	b.Run("ToCovariance", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = CorrelationToCovariance(dst, corr, 100, std)
		}
	})
}

func sizeName(n int) string { return fmt.Sprint(n) }
