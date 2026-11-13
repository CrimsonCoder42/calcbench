package linalg

import (
	"errors"
	"math"
	"testing"
)

// spd returns a random symmetric positive-definite n×n matrix XᵀX + n·I.
func spd(seed int, n int) []float64 {
	r := newRNG()
	for range seed {
		r.Float64()
	}
	x := randMat(r, n, n)
	xt, _ := MatTranspose(nil, x, n, n)
	a, _ := MatMul(nil, xt, n, n, x, n)
	for i := range n {
		a[i*n+i] += float64(n)
	}
	return a
}

// orthoErr returns max |QᵀQ − I| for an m×k matrix Q.
func orthoErr(q []float64, m, k int) float64 {
	e := 0.0
	for i := range k {
		for j := range k {
			s := 0.0
			for r := range m {
				s += q[r*k+i] * q[r*k+j]
			}
			if i == j {
				s--
			}
			e = max(e, math.Abs(s))
		}
	}
	return e
}

func TestTriangularAndTridiagonal_examples(t *testing.T) {
	if x, err := ForwardSubstitution(nil, []float64{2, 0, 3, 4}, 2, []float64{2, 11}); err != nil || !eqSlice(x, []float64{1, 2}, 0) {
		t.Errorf("ForwardSubstitution = %v, %v", x, err)
	}
	if x, err := BackSubstitution(nil, []float64{2, 1, 0, 4}, 2, []float64{4, 8}); err != nil || !eqSlice(x, []float64{1, 2}, 0) {
		t.Errorf("BackSubstitution = %v, %v", x, err)
	}
	if x, err := TridiagonalSolve(nil, []float64{1, 1}, []float64{4, 4, 4}, []float64{1, 1}, []float64{5, 6, 5}); err != nil || !eqSlice(x, []float64{1, 1, 1}, 1e-15) {
		t.Errorf("TridiagonalSolve = %v, %v", x, err)
	}
	if x, err := TridiagonalSolve(nil, nil, []float64{2}, nil, []float64{6}); err != nil || x[0] != 3 {
		t.Errorf("TridiagonalSolve n=1 = %v, %v", x, err)
	}
	// Oracle (fixed seed): Thomas agrees with MatSolve on diagonally dominant tridiagonal systems.
	r := newRNG()
	for trial := range 50 {
		n := 2 + trial%30
		sub, diag, sup, d := randMat(r, 1, n-1), randMat(r, 1, n), randMat(r, 1, n-1), randMat(r, 1, n)
		full := make([]float64, n*n)
		for i := range n {
			diag[i] = 3 + diag[i]
			full[i*n+i] = diag[i]
			if i > 0 {
				full[i*n+i-1] = sub[i-1]
			}
			if i < n-1 {
				full[i*n+i+1] = sup[i]
			}
		}
		x, err := TridiagonalSolve(nil, sub, diag, sup, d)
		y, err2 := MatSolve(nil, full, n, d)
		if err != nil || err2 != nil || !eqSlice(x, y, 1e-12) {
			t.Fatalf("Thomas %v vs MatSolve %v", x, y)
		}
	}
}

func TestCholesky_examplesAndOracle(t *testing.T) {
	l, err := Cholesky(nil, []float64{4, 2, 2, 3}, 2)
	if err != nil || !eqSlice(l, []float64{2, 0, 1, math.Sqrt2}, 1e-15) {
		t.Fatalf("Cholesky = %v, %v", l, err)
	}
	if x, err := CholeskySolve(nil, l, 2, []float64{8, 7}); err != nil || !eqSlice(x, []float64{1.25, 1.5}, 1e-15) {
		t.Errorf("CholeskySolve = %v, %v", x, err)
	}
	if _, err := Cholesky(nil, []float64{1, 2, 2, 1}, 2); !errors.Is(err, ErrNotPositiveDefinite) {
		t.Errorf("indefinite: %v", err)
	}
	// Oracle (fixed seed): A = LLᵀ and the solve's residual is small.
	for trial := range 40 {
		n := 1 + trial%12
		a := spd(trial, n)
		l, err := Cholesky(nil, a, n)
		if err != nil {
			t.Fatal(err)
		}
		lt, _ := MatTranspose(nil, l, n, n)
		llt, _ := MatMul(nil, l, n, n, lt, n)
		if !eqSlice(llt, a, 1e-12) {
			t.Fatalf("LLᵀ ≠ A (n=%d)", n)
		}
		b := randMat(newRNG(), 1, n)
		x, err := CholeskySolve(nil, l, n, b)
		y, err2 := MatSolve(nil, a, n, b)
		if err != nil || err2 != nil || !eqSlice(x, y, 1e-10) {
			t.Fatalf("CholeskySolve %v vs MatSolve %v", x, y)
		}
	}
	// Hilbert matrices are SPD and ill-conditioned; Cholesky succeeds through n = 10 and LLᵀ = H to working accuracy.
	h := hilbert(10)
	l, err = Cholesky(nil, h, 10)
	if err != nil {
		t.Fatalf("Hilbert(10): %v", err)
	}
	lt, _ := MatTranspose(nil, l, 10, 10)
	llt, _ := MatMul(nil, l, 10, 10, lt, 10)
	if !eqSlice(llt, h, 1e-13) {
		t.Errorf("Hilbert LLᵀ ≠ H")
	}
}

func TestQR_examplesAndOracle(t *testing.T) {
	q, rr, err := QRHouseholder(nil, nil, []float64{3, 0, 4, 5}, 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	abs := []float64{math.Abs(rr[0]), math.Abs(rr[1]), rr[2], math.Abs(rr[3])}
	if !eqSlice(abs, []float64{5, 4, 0, 3}, 1e-14) || !eqSlice(rr, []float64{-5, -4, 0, -3}, 1e-14) {
		t.Errorf("R = %v", rr)
	}
	if orthoErr(q, 2, 2) > 1e-15 {
		t.Errorf("QᵀQ ≠ I")
	}
	qg, rg, err := GramSchmidtModified(nil, nil, []float64{1, 1, 1, 0, 0, 1}, 3, 2, -1)
	if err != nil || !eqSlice(rg, []float64{math.Sqrt2, 1 / math.Sqrt2, 0, math.Sqrt(1.5)}, 1e-15) ||
		!eqSlice([]float64{qg[0], qg[2], qg[4]}, []float64{1 / math.Sqrt2, 1 / math.Sqrt2, 0}, 1e-15) ||
		!eqSlice([]float64{qg[1], qg[3], qg[5]}, []float64{1 / math.Sqrt(6), -1 / math.Sqrt(6), 2 / math.Sqrt(6)}, 1e-15) {
		t.Errorf("MGS Q = %v R = %v (%v)", qg, rg, err)
	}
	// Oracle (fixed seed): A = QR, Q orthonormal, R upper triangular, for both methods.
	r := newRNG()
	for trial := range 60 {
		m := 1 + trial%9
		n := 1 + r.IntN(m)
		a := randMat(r, m, n)
		for _, method := range []string{"householder", "mgs"} {
			var q, rr []float64
			var err error
			if method == "householder" {
				q, rr, err = QRHouseholder(nil, nil, a, m, n)
			} else {
				q, rr, err = GramSchmidtModified(nil, nil, a, m, n, -1)
			}
			if err != nil {
				t.Fatal(method, err)
			}
			qr, _ := MatMul(nil, q, m, n, rr, n)
			upper, _ := MatIsUpperTriangular(rr, n, n, 0)
			if !eqSlice(qr, a, 1e-12) || orthoErr(q, m, n) > 1e-12 || !upper {
				t.Fatalf("%s m=%d n=%d: QR=A %v, ortho %v, upper %v", method, m, n, eqSlice(qr, a, 1e-12), orthoErr(q, m, n), upper)
			}
		}
	}
	// Householder keeps orthogonality on the ill-conditioned Hilbert(8); MGS loses it in proportion to κ.
	h := hilbert(8)
	qh, _, _ := QRHouseholder(nil, nil, h, 8, 8)
	if e := orthoErr(qh, 8, 8); e > 1e-14 {
		t.Errorf("Householder Hilbert(8) orthogonality %v", e)
	}
	if _, _, err := GramSchmidtModified(nil, nil, []float64{1, 2, 1, 2}, 2, 2, -1); !errors.Is(err, ErrRankDeficient) {
		t.Errorf("MGS dependent columns: %v", err)
	}
}

func TestLeastSquares_examplesAndOracle(t *testing.T) {
	// Strang §4.3: best line through (0, 6), (1, 0), (2, 0) is 5 − 3t, error e = (1, −2, 1).
	res, err := LeastSquaresQR(nil, []float64{1, 0, 1, 1, 1, 2}, 3, 2, []float64{6, 0, 0}, -1)
	if err != nil || !eqSlice(res.X, []float64{5, -3}, 1e-14) || !near(res.ResidualNorm, math.Sqrt(6), 1e-14) {
		t.Errorf("LeastSquaresQR = %+v, %v", res, err)
	}
	if _, err := LeastSquaresQR(nil, []float64{1, 2, 2, 4, 3, 6}, 3, 2, []float64{1, 2, 3}, -1); !errors.Is(err, ErrRankDeficient) {
		t.Errorf("rank-deficient: %v", err)
	}
	// Oracle (fixed seed): the residual is orthogonal to the columns (Aᵀ(b − Ax) = 0) and matches the pseudo-inverse.
	r := newRNG()
	for trial := range 40 {
		m := 2 + trial%8
		n := 1 + r.IntN(m-1)
		a, b := randMat(r, m, n), randMat(r, 1, m)
		res, err := LeastSquaresQR(nil, a, m, n, b, -1)
		if err != nil {
			t.Fatal(err)
		}
		ax, _ := MatVecMul(nil, a, m, n, res.X)
		e := make([]float64, m)
		for i := range m {
			e[i] = b[i] - ax[i]
		}
		at, _ := MatTranspose(nil, a, m, n)
		g, _ := MatVecMul(nil, at, n, m, e)
		if nrm, _ := VecNorm(g); nrm > 1e-12 {
			t.Fatalf("Aᵀr = %v", g)
		}
		if nrm, _ := VecNorm(e); !near(nrm, res.ResidualNorm, 1e-12) {
			t.Fatalf("residual %v vs reported %v", nrm, res.ResidualNorm)
		}
		pinv, _ := MatPseudoInverse(nil, a, m, n, -1)
		xp, _ := MatVecMul(nil, pinv, n, m, b)
		if !eqSlice(xp, res.X, 1e-9) {
			t.Fatalf("A⁺b %v vs QR %v", xp, res.X)
		}
	}
}

func TestEigenJacobiAndPower(t *testing.T) {
	vals, vecs, err := MatSymmetricEigenJacobi(nil, nil, []float64{2, 1, 1, 2}, 2, -1, true)
	if err != nil || !eqSlice(vals, []float64{1, 3}, 1e-15) {
		t.Fatalf("Jacobi 2×2 = %v, %v", vals, err)
	}
	if math.Abs(math.Abs(vecs[0])-1/math.Sqrt2) > 1e-15 || math.Abs(vecs[0]+vecs[2]) > 1e-15 {
		t.Errorf("eigenvector for 1 = (%v, %v), want ∝ (1, −1)", vecs[0], vecs[2])
	}
	if _, _, err := MatSymmetricEigenJacobi(nil, nil, []float64{1, 2, 3, 4}, 2, -1, false); !errors.Is(err, ErrNotSymmetric) {
		t.Errorf("non-symmetric: %v", err)
	}
	// Oracle (fixed seed): A·V = V·Λ, VᵀV = I, eigenvalues ascending and matching Mat3SymmetricEigenvalues.
	r := newRNG()
	for trial := range 40 {
		n := 1 + trial%10
		x := randMat(r, n, n)
		xt, _ := MatTranspose(nil, x, n, n)
		a, _ := MatAdd(nil, x, xt, n, n)
		vals, v, err := MatSymmetricEigenJacobi(nil, nil, a, n, -1, true)
		if err != nil {
			t.Fatal(err)
		}
		av, _ := MatMul(nil, a, n, n, v, n)
		for k := range n {
			for i := range n {
				if math.Abs(av[i*n+k]-vals[k]*v[i*n+k]) > 1e-12 {
					t.Fatalf("A·v ≠ λv (n=%d, k=%d)", n, k)
				}
			}
			if k > 0 && vals[k] < vals[k-1] {
				t.Fatalf("not ascending: %v", vals)
			}
		}
		if orthoErr(v, n, n) > 1e-13 {
			t.Fatalf("VᵀV ≠ I")
		}
		if n == 3 {
			var m3 Mat3
			copy(m3[:], a)
			e3, _ := Mat3SymmetricEigenvalues(m3)
			if !eqSlice(vals, e3[:], 1e-7) {
				t.Fatalf("Jacobi %v vs Mat3SymmetricEigenvalues %v", vals, e3)
			}
		}
	}
	// Power iteration: dominant eigenpair of [2 1; 1 2] from (1, 0); of an SPD matrix vs Jacobi's largest.
	p, err := MatPowerIteration(nil, []float64{2, 1, 1, 2}, 2, []float64{1, 0}, -1, 1000)
	if err != nil || !near(p.Value, 3, 1e-12) || math.Abs(p.Vector[0]-p.Vector[1]) > 1e-6 {
		t.Errorf("power = %+v, %v", p, err)
	}
	a := spd(7, 6)
	jv, _, _ := MatSymmetricEigenJacobi(nil, nil, a, 6, -1, false)
	pw, err := MatPowerIteration(nil, a, 6, []float64{1, 1, 1, 1, 1, 1}, 1e-12, 100000)
	if err != nil || !near(pw.Value, jv[5], 1e-10) {
		t.Errorf("power %v vs Jacobi max %v (%v)", pw.Value, jv[5], err)
	}
	// Equal-modulus dominant eigenvalues (±1): no convergence.
	if _, err := MatPowerIteration(nil, []float64{0, 1, 1, 0}, 2, []float64{1, 0}, -1, 50); !errors.Is(err, ErrNoConvergence) {
		t.Errorf("±1: %v", err)
	}
}

func TestSVD_examplesAndOracle(t *testing.T) {
	s, err := SVDJacobi([]float64{3, 0, 4, 5}, 2, 2, -1, true)
	if err != nil || !eqSlice(s.S, []float64{3 * math.Sqrt(5), math.Sqrt(5)}, 1e-15) {
		t.Fatalf("SVD = %+v, %v", s, err)
	}
	if c, err := MatCondition2([]float64{3, 0, 4, 5}, 2, 2); err != nil || !near(c, 3, 1e-14) {
		t.Errorf("MatCondition2 = %v, %v", c, err)
	}
	if p, err := MatPseudoInverse(nil, []float64{1, 2, 2, 4}, 2, 2, -1); err != nil || !eqSlice(p, []float64{0.04, 0.08, 0.08, 0.16}, 1e-14) {
		t.Errorf("pinv rank 1 = %v, %v", p, err)
	}
	if _, err := MatCondition2([]float64{1, 2, 2, 4}, 2, 2); !errors.Is(err, ErrSingular) {
		t.Errorf("cond of singular: %v", err)
	}
	// Oracle (fixed seed): A = U·S·Vᵀ, orthonormal U, V, descending S; A·A⁺·A = A; tall and wide shapes.
	r := newRNG()
	for trial := range 60 {
		m, n := 1+r.IntN(7), 1+r.IntN(7)
		a := randMat(r, m, n)
		if trial%5 == 0 && n > 1 { // make a dependent column
			for i := range m {
				a[i*n+n-1] = 2 * a[i*n]
			}
		}
		s, err := SVDJacobi(a, m, n, -1, true)
		if err != nil {
			t.Fatal(err)
		}
		k := min(m, n)
		us := make([]float64, m*k)
		for i := range m {
			for j := range k {
				us[i*k+j] = s.U[i*k+j] * s.S[j]
			}
		}
		vt, _ := MatTranspose(nil, s.V, n, k)
		usv, _ := MatMul(nil, us, m, k, vt, n)
		if !eqSlice(usv, a, 1e-12) {
			t.Fatalf("USVᵀ ≠ A (%d×%d)", m, n)
		}
		if orthoErr(s.V, n, k) > 1e-12 {
			t.Fatalf("VᵀV ≠ I")
		}
		for j := 1; j < k; j++ {
			if s.S[j] > s.S[j-1] {
				t.Fatalf("not descending: %v", s.S)
			}
		}
		p, err := MatPseudoInverse(nil, a, m, n, -1)
		if err != nil {
			t.Fatal(err)
		}
		ap, _ := MatMul(nil, a, m, n, p, m)
		apa, _ := MatMul(nil, ap, m, m, a, n)
		if !eqSlice(apa, a, 1e-10) {
			t.Fatalf("A·A⁺·A ≠ A (%d×%d)", m, n)
		}
	}
	// Hilbert(8): κ₂ ≈ 1.5258e10 (singular values to high relative accuracy).
	if c, err := MatCondition2(hilbert(8), 8, 8); err != nil || !near(c, 1.5257575e10, 1e-6) {
		t.Errorf("κ₂(Hilbert 8) = %v, %v", c, err)
	}
}

func TestEchelon_examples(t *testing.T) {
	a := []float64{1, 2, 2, 4, 3, 8, 6, 16}
	rr, piv, err := MatRREF(nil, a, 2, 4, -1)
	if err != nil || !eqSlice(rr, []float64{1, 0, 2, 0, 0, 1, 0, 2}, 1e-15) || len(piv) != 2 || piv[0] != 0 || piv[1] != 1 {
		t.Errorf("MatRREF = %v %v %v", rr, piv, err)
	}
	ns, k, err := MatNullSpace(nil, a, 2, 4, -1)
	if err != nil || k != 2 || !eqSlice(ns, []float64{-2, 0, 0, -2, 1, 0, 0, 1}, 1e-15) {
		t.Errorf("MatNullSpace = %v %v %v", ns, k, err)
	}
	cs, rk, err := MatColumnSpace(nil, a, 2, 4, -1)
	if err != nil || rk != 2 || !eqSlice(cs, []float64{1, 2, 3, 8}, 0) {
		t.Errorf("MatColumnSpace = %v %v %v", cs, rk, err)
	}
	_, piv, err = MatRowEchelon(nil, []float64{1, 2, 1, 2, 4, 0, 3, 6, 3}, 3, 3, -1)
	if err != nil || len(piv) != 2 || piv[0] != 0 || piv[1] != 2 {
		t.Errorf("MatRowEchelon pivots = %v, %v", piv, err)
	}
	// Oracle (fixed seed): RREF of an invertible matrix is I; A·N = 0; rank from the echelon form equals MatRank.
	r := newRNG()
	for trial := range 60 {
		n := 1 + trial%6
		x := randMat(r, n, n)
		rr, piv, err := MatRREF(nil, x, n, n, -1)
		id, _ := MatIdentity(nil, n)
		if err != nil || len(piv) != n || !eqSlice(rr, id, 1e-12) {
			t.Fatalf("rref(invertible) ≠ I: %v", rr)
		}
		m, c := 1+r.IntN(6), 1+r.IntN(6)
		k := 1 + r.IntN(min(m, c))
		u, v := randMat(r, m, k), randMat(r, k, c)
		a, _ := MatMul(nil, u, m, k, v, c)
		e, piv2, _ := MatRowEchelon(nil, a, m, c, 1e-10)
		_ = e
		rk, _ := MatRank(a, m, c, 1e-10)
		if len(piv2) != rk {
			t.Fatalf("echelon rank %d vs MatRank %d", len(piv2), rk)
		}
		ns, nk, _ := MatNullSpace(nil, a, m, c, 1e-10)
		if nk != c-rk {
			t.Fatalf("nullity %d, want %d", nk, c-rk)
		}
		if nk > 0 {
			z, _ := MatMul(nil, a, m, c, ns, nk)
			if nz, _ := MatNormFrobenius(z, m, nk); nz > 1e-9 {
				t.Fatalf("A·N = %v", nz)
			}
		}
	}
}

func TestMatFactor_errors(t *testing.T) {
	nan := math.NaN()
	cases := []struct {
		name string
		err  error
		want error
	}{
		{"Forward singular", second(ForwardSubstitution(nil, []float64{0, 0, 1, 1}, 2, []float64{1, 1})), ErrSingular},
		{"Forward dims", second(ForwardSubstitution(nil, []float64{1}, 1, []float64{1, 2})), ErrDimension},
		{"Forward NaN", second(ForwardSubstitution(nil, []float64{nan}, 1, []float64{1})), ErrNotFinite},
		{"Forward overflow", second(ForwardSubstitution(nil, []float64{1e-300}, 1, []float64{1e300})), ErrOverflow},
		{"Back singular", second(BackSubstitution(nil, []float64{1, 1, 0, 0}, 2, []float64{1, 1})), ErrSingular},
		{"Back overflow", second(BackSubstitution(nil, []float64{1e-300}, 1, []float64{1e300})), ErrOverflow},
		{"Tridiagonal dims", second(TridiagonalSolve(nil, []float64{1}, []float64{1, 2}, nil, []float64{1, 2})), ErrDimension},
		{"Tridiagonal empty", second(TridiagonalSolve(nil, nil, nil, nil, nil)), ErrDimension},
		{"Tridiagonal NaN", second(TridiagonalSolve(nil, []float64{nan}, []float64{1, 1}, []float64{1}, []float64{1, 1})), ErrNotFinite},
		{"Tridiagonal singular", second(TridiagonalSolve(nil, []float64{1}, []float64{1, 1}, []float64{1}, []float64{1, 1})), ErrSingular},
		{"Tridiagonal overflow", second(TridiagonalSolve(nil, nil, []float64{1e-300}, nil, []float64{1e300})), ErrOverflow},
		{"Cholesky n=0", second(Cholesky(nil, nil, 0)), ErrDimension},
		{"Cholesky NaN", second(Cholesky(nil, []float64{nan}, 1)), ErrNotFinite},
		{"Cholesky not PD", second(Cholesky(nil, []float64{-1}, 1)), ErrNotPositiveDefinite},
		{"CholeskySolve singular", second(CholeskySolve(nil, []float64{0}, 1, []float64{1})), ErrSingular},
		{"QR wide", second3(QRHouseholder(nil, nil, []float64{1, 2}, 1, 2)), ErrDimension},
		{"QR NaN", second3(QRHouseholder(nil, nil, []float64{nan}, 1, 1)), ErrNotFinite},
		{"MGS dims", second3(GramSchmidtModified(nil, nil, []float64{1}, 1, 2, -1)), ErrDimension},
		{"MGS rank", second3(GramSchmidtModified(nil, nil, []float64{0, 0}, 2, 1, -1)), ErrRankDeficient},
		{"LS b len", second(LeastSquaresQR(nil, []float64{1, 1}, 2, 1, []float64{1}, -1)), ErrDimension},
		{"LS NaN", second(LeastSquaresQR(nil, []float64{1, 1}, 2, 1, []float64{nan, 1}, -1)), ErrNotFinite},
		{"LS rank", second(LeastSquaresQR(nil, []float64{0, 0}, 2, 1, []float64{1, 1}, -1)), ErrRankDeficient},
		{"Jacobi dims", second3(MatSymmetricEigenJacobi(nil, nil, []float64{1}, 2, -1, false)), ErrDimension},
		{"Jacobi NaN", second3(MatSymmetricEigenJacobi(nil, nil, []float64{nan}, 1, -1, false)), ErrNotFinite},
		{"Power maxIter", second(MatPowerIteration(nil, []float64{1}, 1, []float64{1}, -1, 0)), ErrOutOfRange},
		{"Power zero x0", second(MatPowerIteration(nil, []float64{1}, 1, []float64{0}, -1, 5)), ErrZeroVector},
		{"Power zero Ax", second(MatPowerIteration(nil, []float64{0, 0, 0, 0}, 2, []float64{1, 0}, -1, 5)), ErrZeroVector},
		{"SVD dims", second(SVDJacobi(nil, 0, 1, -1, false)), ErrDimension},
		{"SVD NaN", second(SVDJacobi([]float64{nan}, 1, 1, -1, false)), ErrNotFinite},
		{"Pinv NaN", second(MatPseudoInverse(nil, []float64{nan}, 1, 1, -1)), ErrNotFinite},
		{"Cond2 NaN", second(MatCondition2([]float64{nan}, 1, 1)), ErrNotFinite},
		{"Norm2 dims", second(MatNorm2([]float64{1}, 1, 2)), ErrDimension},
		{"RowEchelon dims", second3(MatRowEchelon(nil, nil, 0, 1, -1)), ErrDimension},
		{"RREF NaN", second3(MatRREF(nil, []float64{nan}, 1, 1, -1)), ErrNotFinite},
		{"NullSpace dims", second3(MatNullSpace(nil, []float64{1}, 1, 2, -1)), ErrDimension},
		{"ColumnSpace NaN", second3(MatColumnSpace(nil, []float64{nan}, 1, 1, -1)), ErrNotFinite},
	}
	for _, c := range cases {
		if !errors.Is(c.err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, c.err, c.want)
		}
	}
}

// n-dimensional paths reuse dst: the triangular solves, Thomas' result and the echelon outputs allocate only
// their documented work.
func TestMatFactor_dstReuse(t *testing.T) {
	l := []float64{2, 0, 3, 4}
	b := []float64{2, 11}
	dst := make([]float64, 2)
	if n := testing.AllocsPerRun(100, func() { _, _ = ForwardSubstitution(dst, l, 2, b) }); n != 0 {
		t.Errorf("ForwardSubstitution allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { _, _ = BackSubstitution(dst, []float64{2, 1, 0, 4}, 2, b) }); n > 1 {
		t.Errorf("BackSubstitution allocs = %v", n)
	}
	ll, _ := Cholesky(nil, []float64{4, 2, 2, 3}, 2)
	chol := make([]float64, 4)
	if n := testing.AllocsPerRun(100, func() { _, _ = Cholesky(chol, []float64{4, 2, 2, 3}, 2) }); n > 1 {
		t.Errorf("Cholesky allocs = %v", n)
	}
	if n := testing.AllocsPerRun(100, func() { _, _ = CholeskySolve(dst, ll, 2, b) }); n != 0 {
		t.Errorf("CholeskySolve allocs = %v, want 0", n)
	}
}

func benchSPD(n int) []float64 { return spd(3, n) }

func BenchmarkTriangular(b *testing.B) {
	for _, n := range []int{32, 100, 316} {
		l, _ := Cholesky(nil, benchSPD(n), n)
		rhs, dst := randMat(newRNG(), 1, n), make([]float64, n)
		lt, _ := MatTranspose(nil, l, n, n)
		b.Run("Forward/"+sizeName(n*n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = ForwardSubstitution(dst, l, n, rhs)
			}
		})
		b.Run("Back/"+sizeName(n*n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = BackSubstitution(dst, lt, n, rhs)
			}
		})
	}
}

func BenchmarkTridiagonalSolve(b *testing.B) {
	for _, n := range []int{1000, 10000, 100000} {
		r := newRNG()
		sub, diag, sup, d := randMat(r, 1, n-1), randMat(r, 1, n), randMat(r, 1, n-1), randMat(r, 1, n)
		for i := range diag {
			diag[i] += 4
		}
		dst := make([]float64, n)
		b.Run(sizeName(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = TridiagonalSolve(dst, sub, diag, sup, d)
			}
		})
	}
}

func BenchmarkCholesky(b *testing.B) {
	for _, n := range []int{10, 30, 100} {
		a, dst := benchSPD(n), make([]float64, n*n)
		rhs, x := randMat(newRNG(), 1, n), make([]float64, n)
		b.Run(sizeName(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = Cholesky(dst, a, n)
			}
		})
		l, _ := Cholesky(nil, a, n)
		b.Run("Solve/"+sizeName(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = CholeskySolve(x, l, n, rhs)
			}
		})
	}
}

func BenchmarkQR(b *testing.B) {
	for _, n := range []int{10, 30, 100} {
		a := randMat(newRNG(), n, n)
		q, r := make([]float64, n*n), make([]float64, n*n)
		rhs := randMat(newRNG(), 1, n)
		b.Run("Householder/"+sizeName(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = QRHouseholder(q, r, a, n, n)
			}
		})
		b.Run("MGS/"+sizeName(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = GramSchmidtModified(q, r, a, n, n, -1)
			}
		})
		b.Run("LeastSquares/"+sizeName(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = LeastSquaresQR(nil, a, n, n, rhs, -1)
			}
		})
	}
}

func BenchmarkEigen(b *testing.B) {
	for _, n := range []int{10, 30} {
		a := benchSPD(n)
		vals, vecs := make([]float64, n), make([]float64, n*n)
		b.Run("Jacobi/"+sizeName(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = MatSymmetricEigenJacobi(vals, vecs, a, n, -1, true)
			}
		})
		x0 := make([]float64, n)
		for i := range x0 {
			x0[i] = 1
		}
		v := make([]float64, n)
		b.Run("Power/"+sizeName(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = MatPowerIteration(v, a, n, x0, 1e-10, 100000)
			}
		})
	}
}

func BenchmarkSVD(b *testing.B) {
	for _, n := range []int{10, 30} {
		a := randMat(newRNG(), n, n)
		b.Run("Jacobi/"+sizeName(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = SVDJacobi(a, n, n, -1, true)
			}
		})
		b.Run("PseudoInverse/"+sizeName(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = MatPseudoInverse(nil, a, n, n, -1)
			}
		})
		b.Run("Condition2/"+sizeName(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = MatCondition2(a, n, n)
			}
		})
		b.Run("Norm2/"+sizeName(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = MatNorm2(a, n, n)
			}
		})
	}
}

func BenchmarkEchelon(b *testing.B) {
	for _, n := range []int{10, 30, 100} {
		a := randMat(newRNG(), n, n)
		dst := make([]float64, n*n)
		b.Run("RowEchelon/"+sizeName(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = MatRowEchelon(dst, a, n, n, -1)
			}
		})
		b.Run("RREF/"+sizeName(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = MatRREF(dst, a, n, n, -1)
			}
		})
		b.Run("NullSpace/"+sizeName(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = MatNullSpace(nil, a, n, n, -1)
			}
		})
		b.Run("ColumnSpace/"+sizeName(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = MatColumnSpace(nil, a, n, n, -1)
			}
		})
	}
}
