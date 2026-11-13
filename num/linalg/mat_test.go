package linalg

import (
	"errors"
	"math"
	"strconv"
	"testing"
)

func TestMatN_workedExamples(t *testing.T) {
	// Hand-checked examples from the doc comments (Strang 2016 §2.4–§2.7, §3.2, §5.1; Golub, Van Loan §3.2–§3.4).
	if m, err := MatIdentity(nil, 2); err != nil || !nearM(m, []float64{1, 0, 0, 1}, 0) {
		t.Errorf("MatIdentity = %v, %v", m, err)
	}
	if m, err := MatMul(nil, []float64{1, 2, 3, 4, 5, 6}, 2, 3, []float64{7, 8, 9, 10, 11, 12}, 2); err != nil || !nearM(m, []float64{58, 64, 139, 154}, 0) {
		t.Errorf("MatMul = %v, %v", m, err)
	}
	if v, err := MatVecMul(nil, []float64{1, 2, 3, 4, 5, 6}, 2, 3, []float64{1, 0, -1}); err != nil || !nearM(v, []float64{-2, -2}, 0) {
		t.Errorf("MatVecMul = %v, %v", v, err)
	}
	if m, err := MatTranspose(nil, []float64{1, 2, 3, 4, 5, 6}, 2, 3); err != nil || !nearM(m, []float64{1, 4, 2, 5, 3, 6}, 0) {
		t.Errorf("MatTranspose = %v, %v", m, err)
	}
	if tr, err := MatTrace([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, 3); err != nil || tr != 15 {
		t.Errorf("MatTrace = %v, %v", tr, err)
	}
	f, err := LUDecompose([]float64{0, 1, 2, 3}, 2)
	if err != nil || f.sign != -1 || !nearM(f.lu, []float64{2, 3, 0, 1}, 0) || f.piv[0] != 1 {
		t.Errorf("LUDecompose = %+v, %v", f, err)
	}
	g, _ := LUDecompose([]float64{2, 1, 1, 3}, 2)
	if x, err := LUSolve(nil, g, []float64{3, 5}); err != nil || !nearM(x, []float64{0.8, 1.4}, ulp1) {
		t.Errorf("LUSolve = %v, %v", x, err)
	}
	if x, err := MatSolve(nil, []float64{2, 1, 1, 3}, 2, []float64{3, 5}); err != nil || !nearM(x, []float64{0.8, 1.4}, ulp1) {
		t.Errorf("MatSolve = %v, %v", x, err)
	}
	a := []float64{1, 2, 3, 0, 1, 4, 5, 6, 0}
	if d, err := MatDet(a, 3); err != nil || !near(d, 1, 8*ulp1) {
		t.Errorf("MatDet = %v, %v", d, err)
	}
	// κ₁ of this matrix is about 2e3, so LU loses ~3 of 16 digits: 1e−13 relative.
	if m, err := MatInverse(nil, a, 3); err != nil || !nearM(m, []float64{-24, 18, 5, 20, -15, -4, -5, 4, 1}, 1e-13) {
		t.Errorf("MatInverse = %v, %v", m, err)
	}
	if r, err := MatRank([]float64{1, 2, 3, 2, 4, 6, 1, 0, 1}, 3, 3, -1); err != nil || r != 2 {
		t.Errorf("MatRank = %v, %v", r, err)
	}
	if k, err := MatConditionEstimate1([]float64{1, 0, 0, 1e-8}, 2); err != nil || !near(k, 1e8, ulp1) {
		t.Errorf("MatConditionEstimate1 = %v, %v", k, err)
	}
}

func TestMatN_cases(t *testing.T) {
	// Empty and one-element shapes.
	if m, err := MatMul(nil, nil, 0, 3, make([]float64, 6), 2); err != nil || len(m) != 0 {
		t.Errorf("MatMul 0×3 · 3×2 = %v, %v", m, err)
	}
	if d, err := MatDet([]float64{-7}, 1); err != nil || d != -7 {
		t.Errorf("MatDet 1×1 = %v, %v", d, err)
	}
	if r, _ := MatRank(make([]float64, 6), 2, 3, -1); r != 0 {
		t.Errorf("rank of zero = %v", r)
	}
	if r, _ := MatRank(nil, 0, 0, -1); r != 0 {
		t.Errorf("rank of empty = %v", r)
	}
	if r, _ := MatRank([]float64{1, 0, 0, 0, 1, 0}, 2, 3, -1); r != 2 {
		t.Errorf("rank of a wide matrix = %v", r)
	}
	// A pivot of 1e-20 counts with tol = 0 and not with tol = 1e-10.
	nearlySingular := []float64{1, 0, 0, 1e-20}
	if r, _ := MatRank(nearlySingular, 2, 2, 0); r != 2 {
		t.Errorf("rank with tol 0 = %v", r)
	}
	if r, _ := MatRank(nearlySingular, 2, 2, 1e-10); r != 1 {
		t.Errorf("rank with tol 1e−10 = %v", r)
	}
	// A singular matrix has determinant 0, and the product of pivots does not overflow part-way.
	if d, err := MatDet([]float64{1, 2, 2, 4}, 2); err != nil || d != 0 {
		t.Errorf("singular det = %v, %v", d, err)
	}
	if d, err := MatDet([]float64{1e200, 0, 0, 0, 1e200, 0, 0, 0, 1e-300}, 3); err != nil || !near(d/1e100, 1, 4*ulp1) {
		t.Errorf("det without intermediate overflow = %v, %v", d, err)
	}
	// Partial pivoting handles a zero leading entry.
	if x, err := MatSolve(nil, []float64{0, 1, 1, 0}, 2, []float64{2, 3}); err != nil || !nearM(x, []float64{3, 2}, 0) {
		t.Errorf("MatSolve with a pivot swap = %v, %v", x, err)
	}
}

func TestMatN_errors(t *testing.T) {
	nan, big := math.NaN(), math.MaxFloat64
	sing := []float64{1, 2, 2, 4}
	fs, _ := LUDecompose(sing, 2)
	f2, _ := LUDecompose([]float64{2, 1, 1, 3}, 2)
	ftiny, _ := LUDecompose([]float64{1e-300, 0, 0, 1}, 2)
	cases := []struct {
		name string
		err  error
		want error
	}{
		{"Identity n<0", second(MatIdentity(nil, -1)), ErrDimension},
		{"Identity huge", second(MatIdentity(nil, math.MaxInt)), ErrDimension},
		{"Mul a length", second(MatMul(nil, []float64{1}, 2, 2, []float64{1, 2, 3, 4}, 2)), ErrDimension},
		{"Mul b length", second(MatMul(nil, []float64{1, 2, 3, 4}, 2, 2, []float64{1}, 2)), ErrDimension},
		{"Mul negative", second(MatMul(nil, nil, -1, 0, nil, 0)), ErrDimension},
		{"Mul NaN", second(MatMul(nil, []float64{nan}, 1, 1, []float64{1}, 1)), ErrNotFinite},
		{"Mul NaN b", second(MatMul(nil, []float64{1}, 1, 1, []float64{nan}, 1)), ErrNotFinite},
		{"Mul overflow", second(MatMul(nil, []float64{big, big}, 1, 2, []float64{2, 2}, 1)), ErrOverflow},
		{"VecMul x length", second(MatVecMul(nil, []float64{1, 2}, 1, 2, []float64{1})), ErrDimension},
		{"VecMul NaN", second(MatVecMul(nil, []float64{1, nan}, 1, 2, []float64{1, 1})), ErrNotFinite},
		{"VecMul NaN x", second(MatVecMul(nil, []float64{1, 1}, 1, 2, []float64{1, nan})), ErrNotFinite},
		{"VecMul overflow", second(MatVecMul(nil, []float64{big, big}, 1, 2, []float64{1, 1})), ErrOverflow},
		{"Transpose length", second(MatTranspose(nil, []float64{1, 2, 3}, 2, 2)), ErrDimension},
		{"Transpose NaN", second(MatTranspose(nil, []float64{nan}, 1, 1)), ErrNotFinite},
		{"Trace length", second(MatTrace([]float64{1, 2}, 2)), ErrDimension},
		{"Trace NaN", second(MatTrace([]float64{nan}, 1)), ErrNotFinite},
		{"Trace overflow", second(MatTrace([]float64{big, 0, 0, big}, 2)), ErrOverflow},
		{"LU n = 0", second(LUDecompose(nil, 0)), ErrDimension},
		{"LU length", second(LUDecompose([]float64{1, 2, 3}, 2)), ErrDimension},
		{"LU NaN", second(LUDecompose([]float64{nan}, 1)), ErrNotFinite},
		{"LU growth overflow", second(LUDecompose([]float64{1, big, 1, -big}, 2)), ErrOverflow},
		{"LUSolve singular", second(LUSolve(nil, fs, []float64{1, 1})), ErrSingular},
		{"LUSolve length", second(LUSolve(nil, f2, []float64{1})), ErrDimension},
		{"LUSolve empty", second(LUSolve(nil, LU{}, nil)), ErrDimension},
		{"LUSolve NaN", second(LUSolve(nil, f2, []float64{nan, 1})), ErrNotFinite},
		{"LUSolve overflow", second(LUSolve(nil, ftiny, []float64{1e10, 1})), ErrOverflow},
		{"Solve singular", second(MatSolve(nil, sing, 2, []float64{1, 1})), ErrSingular},
		{"Det length", second(MatDet([]float64{1, 2}, 2)), ErrDimension},
		{"Det overflow", second(MatDet([]float64{1e200, 0, 0, 1e200}, 2)), ErrOverflow},
		{"Inverse singular", second(MatInverse(nil, sing, 2)), ErrSingular},
		{"Inverse overflow", second(MatInverse(nil, []float64{1e-310, 0, 0, 1}, 2)), ErrOverflow},
		{"Inverse NaN", second(MatInverse(nil, []float64{nan}, 1)), ErrNotFinite},
		{"Rank length", second(MatRank([]float64{1}, 2, 2, 0)), ErrDimension},
		{"Rank NaN tol", second(MatRank([]float64{1}, 1, 1, nan)), ErrNotFinite},
		{"Rank NaN", second(MatRank([]float64{nan}, 1, 1, 0)), ErrNotFinite},
		{"Cond singular", second(MatConditionEstimate1(sing, 2)), ErrSingular},
		{"Cond length", second(MatConditionEstimate1([]float64{1}, 2)), ErrDimension},
		{"Cond overflow", second(MatConditionEstimate1([]float64{big, 0, 0, 1e-300}, 2)), ErrOverflow},
	}
	for _, c := range cases {
		if !errors.Is(c.err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, c.err, c.want)
		}
	}
}

// hilbert returns the n×n Hilbert matrix 1/(i+j+1), the classic ill-conditioned test matrix.
func hilbert(n int) []float64 {
	h := make([]float64, n*n)
	for i := range n {
		for j := range n {
			h[i*n+j] = 1 / float64(i+j+1)
		}
	}
	return h
}

// Oracle (fixed seed): MatSolve's residual is small relative to the condition number; MatInverse gives
// A·A⁻¹ ≈ I; MatMul matches Mat3Mul on 3×3 inputs; the condition estimate is a lower bound on
// ‖A‖₁‖A⁻¹‖₁ computed from the explicit inverse, and within a factor of 10 of it (Higham reports "usually
// within 3"; 10 keeps the test from depending on that being true for every random matrix).
func TestMatN_oracle(t *testing.T) {
	r := newRNG()
	for trial := range 200 {
		n := 1 + r.IntN(12)
		a := make([]float64, n*n)
		for i := range a {
			a[i] = sym(r)
		}
		b := make([]float64, n)
		for i := range b {
			b[i] = sym(r)
		}
		cond, err := MatConditionEstimate1(a, n)
		if err != nil {
			continue
		}
		x, _ := MatSolve(nil, a, n, b)
		ax, _ := MatVecMul(nil, a, n, n, x)
		for i := range b {
			if math.Abs(ax[i]-b[i]) > 1e-13*float64(n)*cond {
				t.Fatalf("trial %d: residual %v (κ ≈ %v)", trial, ax[i]-b[i], cond)
			}
		}
		inv, _ := MatInverse(nil, a, n)
		ai, _ := MatMul(nil, a, n, n, inv, n)
		id, _ := MatIdentity(nil, n)
		if !nearM(ai, id, 1e-13*float64(n)*cond) {
			t.Fatalf("trial %d: A·A⁻¹ ≠ I (κ ≈ %v)", trial, cond)
		}
		exact := norm1Mat(a, n) * norm1Mat(inv, n)
		if cond > exact*(1+1e-10) || cond < exact/10 {
			t.Fatalf("trial %d (n = %d): estimate %v vs ‖A‖‖A⁻¹‖ = %v", trial, n, cond, exact)
		}
		if n == 3 {
			m3 := Mat3(a)
			p3, _ := Mat3Mul(m3, m3)
			pn, _ := MatMul(nil, a, 3, 3, a, 3)
			if !nearM(p3[:], pn, 4*ulp1) {
				t.Fatalf("trial %d: MatMul ≠ Mat3Mul", trial)
			}
		}
		rk, _ := MatRank(a, n, n, -1)
		if rk != n {
			t.Fatalf("trial %d: rank %d of a random %d×%d matrix", trial, rk, n, n)
		}
	}
	// Hilbert 5: κ₁ = 943656 (computed from its exact integer inverse); the estimate must be within 3×.
	k, _ := MatConditionEstimate1(hilbert(5), 5)
	if k > 943656*(1+1e-6) || k < 943656/3.0 {
		t.Errorf("κ₁(H₅) estimate = %v, true 943656", k)
	}
}

func norm1Mat(a []float64, n int) float64 {
	m := 0.0
	for j := range n {
		s := 0.0
		for i := range n {
			s += math.Abs(a[i*n+j])
		}
		m = max(m, s)
	}
	return m
}

func TestMatN_dstReuseNoAllocation(t *testing.T) {
	a := []float64{4, 1, 0, 1, 4, 1, 0, 1, 4}
	x := []float64{1, 2, 3}
	dst9, dst3 := make([]float64, 9), make([]float64, 3)
	f, _ := LUDecompose(a, 3)
	if n := testing.AllocsPerRun(50, func() {
		sinkS, _ = MatIdentity(dst9, 3)
		sinkS, _ = MatMul(dst9, a, 3, 3, a, 3)
		sinkS, _ = MatVecMul(dst3, a, 3, 3, x)
		sinkS, _ = MatTranspose(dst9, a, 3, 3)
		sinkF, _ = MatTrace(a, 3)
		sinkS, _ = LUSolve(dst3, f, x)
	}); n != 0 {
		t.Errorf("allocations with a reused dst = %v, want 0", n)
	}
}

// benchMatrix returns a diagonally dominant (well-conditioned) n×n matrix, a vector and a reusable dst.
func benchMatrix(n int) (a, x, dst []float64) {
	r := newRNG()
	a, x, dst = make([]float64, n*n), make([]float64, n), make([]float64, n*n)
	for i := range n {
		for j := range n {
			a[i*n+j] = sym(r)
		}
		a[i*n+i] += float64(n)
		x[i] = sym(r)
	}
	return a, x, dst
}

func BenchmarkMatIdentity(b *testing.B) {
	for _, n := range []int{32, 64, 128} {
		a, x, dst := benchMatrix(n)
		_, _, _ = a, x, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkS, _ = MatIdentity(dst, n)
			}
		})
	}
}

func BenchmarkMatMul(b *testing.B) {
	for _, n := range []int{32, 64, 128} {
		a, x, dst := benchMatrix(n)
		_, _, _ = a, x, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkS, _ = MatMul(dst, a, n, n, a, n)
			}
		})
	}
}

func BenchmarkMatVecMul(b *testing.B) {
	for _, n := range []int{32, 64, 128} {
		a, x, dst := benchMatrix(n)
		_, _, _ = a, x, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkS, _ = MatVecMul(dst, a, n, n, x)
			}
		})
	}
}

func BenchmarkMatTranspose(b *testing.B) {
	for _, n := range []int{32, 64, 128} {
		a, x, dst := benchMatrix(n)
		_, _, _ = a, x, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkS, _ = MatTranspose(dst, a, n, n)
			}
		})
	}
}

func BenchmarkMatTrace(b *testing.B) {
	for _, n := range []int{32, 64, 128} {
		a, x, dst := benchMatrix(n)
		_, _, _ = a, x, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkF, _ = MatTrace(a, n)
			}
		})
	}
}

func BenchmarkLUDecompose(b *testing.B) {
	for _, n := range []int{32, 64, 128} {
		a, x, dst := benchMatrix(n)
		_, _, _ = a, x, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkLU, _ = LUDecompose(a, n)
			}
		})
	}
}

func BenchmarkLUSolve(b *testing.B) {
	for _, n := range []int{32, 64, 128} {
		a, x, dst := benchMatrix(n)
		_, _, _ = a, x, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			f, _ := LUDecompose(a, n)
			b.ReportAllocs()
			for b.Loop() {
				sinkS, _ = LUSolve(dst, f, x)
			}
		})
	}
}

func BenchmarkMatSolve(b *testing.B) {
	for _, n := range []int{32, 64, 128} {
		a, x, dst := benchMatrix(n)
		_, _, _ = a, x, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkS, _ = MatSolve(dst, a, n, x)
			}
		})
	}
}

func BenchmarkMatDet(b *testing.B) {
	for _, n := range []int{32, 64, 128} {
		a, x, dst := benchMatrix(n)
		_, _, _ = a, x, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkF, _ = MatDet(a, n)
			}
		})
	}
}

func BenchmarkMatInverse(b *testing.B) {
	for _, n := range []int{32, 64, 128} {
		a, x, dst := benchMatrix(n)
		_, _, _ = a, x, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkS, _ = MatInverse(dst, a, n)
			}
		})
	}
}

func BenchmarkMatRank(b *testing.B) {
	for _, n := range []int{32, 64, 128} {
		a, x, dst := benchMatrix(n)
		_, _, _ = a, x, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkI, _ = MatRank(a, n, n, -1)
			}
		})
	}
}

func BenchmarkMatConditionEstimate1(b *testing.B) {
	for _, n := range []int{32, 64, 128} {
		a, x, dst := benchMatrix(n)
		_, _, _ = a, x, dst
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkF, _ = MatConditionEstimate1(a, n)
			}
		})
	}
}

// Review regression (finding 10): rank near the float64 maximum.
func TestMatRank_nearMaxFloat(t *testing.T) {
	if r, err := MatRank([]float64{1e308, 1e308, -1e308, 1e308}, 2, 2, -1); err != nil || r != 2 {
		t.Errorf("rank = %v, %v; want 2", r, err)
	}
	if _, err := MatSolve(nil, []float64{1}, 2, []float64{1, 1}); !errors.Is(err, ErrDimension) {
		t.Errorf("MatSolve dimension error = %v", err)
	}
}
