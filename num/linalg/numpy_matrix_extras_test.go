package linalg

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// GvL §5.5.1: the minimum-norm solution of x₁ + x₂ = 2 is (1, 1) (hand-checked).
func TestLeastSquaresMinNorm_workedExample(t *testing.T) {
	res, rank, err := LeastSquaresMinNorm(nil, []float64{1, 1}, 1, 2, []float64{2}, -1)
	if err != nil || rank != 1 || !near(res.X[0], 1, 1e-14) || !near(res.X[1], 1, 1e-14) || res.ResidualNorm > 1e-14 {
		t.Errorf("min-norm = (%v, %d, %v)", res, rank, err)
	}
	// Full rank tall: agrees with Strang §4.3 (5, −3), residual √6.
	res, rank, _ = LeastSquaresMinNorm(nil, []float64{1, 0, 1, 1, 1, 2}, 3, 2, []float64{6, 0, 0}, -1)
	if rank != 2 || !near(res.X[0], 5, 1e-12) || !near(res.X[1], -3, 1e-12) || !near(res.ResidualNorm, math.Sqrt(6), 1e-12) {
		t.Errorf("tall = %v rank %d", res, rank)
	}
	if _, _, err := LeastSquaresMinNorm(nil, []float64{1, 1}, 1, 2, []float64{1, 2}, -1); !errors.Is(err, ErrDimension) {
		t.Errorf("len b: %v", err)
	}
	if _, _, err := LeastSquaresMinNorm(nil, []float64{1, 1}, 1, 2, []float64{math.NaN()}, -1); !errors.Is(err, ErrNotFinite) {
		t.Errorf("NaN: %v", err)
	}
}

// Oracle: for rank-deficient A, the result equals pinv(A)·b.
func TestLeastSquaresMinNorm_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for trial := range 100 {
		m, n, k := 2+trial%6, 2+(trial/6)%6, 1
		k = 1 + trial%min(m, n)
		u, v := make([]float64, m*k), make([]float64, k*n)
		for i := range u {
			u[i] = r.NormFloat64()
		}
		for i := range v {
			v[i] = r.NormFloat64()
		}
		a, _ := MatMul(nil, u, m, k, v, n)
		b := make([]float64, m)
		for i := range b {
			b[i] = r.NormFloat64()
		}
		res, rank, err := LeastSquaresMinNorm(nil, a, m, n, b, 1e-10)
		if err != nil {
			t.Fatal(err)
		}
		if rank != k {
			t.Errorf("trial %d: rank %d, want %d", trial, rank, k)
		}
		pinv, _ := MatPseudoInverse(nil, a, m, n, 1e-10)
		want, _ := MatVecMul(nil, pinv, n, m, b)
		for i := range want {
			if math.Abs(res.X[i]-want[i]) > 1e-8*max(1, math.Abs(want[i])) {
				t.Errorf("trial %d: x = %v, pinv·b = %v", trial, res.X, want)
				break
			}
		}
	}
}

func TestMatLogDet_workedExample(t *testing.T) {
	s, l, err := MatLogDet([]float64{1, 2, 3, 4}, 2)
	if err != nil || s != -1 || !near(l, math.Ln2, 1e-15) {
		t.Errorf("slogdet = (%v, %v, %v)", s, l, err)
	}
	// det = 10^400 overflows MatDet but not MatLogDet.
	n := 400
	a := make([]float64, n*n)
	for i := range n {
		a[i*n+i] = 10
	}
	s, l, _ = MatLogDet(a, n)
	if s != 1 || !near(l, 400*math.Ln10, 1e-13) {
		t.Errorf("large det: %v %v", s, l)
	}
	if _, _, err := MatLogDet([]float64{1, 2, 2, 4}, 2); !errors.Is(err, ErrSingular) {
		t.Errorf("singular: %v", err)
	}
}

func TestMatNormsExtra_workedExample(t *testing.T) {
	a := []float64{3, 0, 4, 5}
	if v, err := MatNormNuclear(a, 2, 2); err != nil || !near(v, 4*math.Sqrt(5), 1e-14) {
		t.Errorf("nuclear = (%v, %v)", v, err)
	}
	if v, err := MatNormNeg2(a, 2, 2); err != nil || !near(v, math.Sqrt(5), 1e-14) {
		t.Errorf("neg2 = (%v, %v)", v, err)
	}
	b := []float64{1, -2, 3, 4}
	if v, err := MatNormNeg1(b, 2, 2); err != nil || v != 4 {
		t.Errorf("neg1 = (%v, %v)", v, err)
	}
	if v, err := MatNormNegInf(b, 2, 2); err != nil || v != 3 {
		t.Errorf("neginf = (%v, %v)", v, err)
	}
	if v, err := VecNorm0([]float64{0, 3, 0, -1}); err != nil || v != 2 {
		t.Errorf("norm0 = (%v, %v)", v, err)
	}
	if _, err := MatNormNeg1(nil, 0, 0); !errors.Is(err, ErrDimension) {
		t.Errorf("empty: %v", err)
	}
	if _, err := MatNormNegInf([]float64{math.Inf(1)}, 1, 1); !errors.Is(err, ErrNotFinite) {
		t.Errorf("Inf: %v", err)
	}
	if _, err := VecNorm0([]float64{math.NaN()}); !errors.Is(err, ErrNotFinite) {
		t.Errorf("NaN: %v", err)
	}
}

// Higham §6.2: κ = ‖A‖‖A⁻¹‖; A = [1 2; 3 4], A⁻¹ = [−2 1; 1.5 −0.5] (hand-checked).
func TestMatConditionExact_workedExample(t *testing.T) {
	a := []float64{1, 2, 3, 4}
	for name, f := range map[string]func([]float64, int) (float64, error){"1": MatCondition1, "inf": MatConditionInf, "fro": MatConditionFrobenius} {
		want := map[string]float64{"1": 21, "inf": 21, "fro": 15}[name]
		if v, err := f(a, 2); err != nil || !near(v, want, 1e-14) {
			t.Errorf("κ_%s = (%v, %v), want %v", name, v, err, want)
		}
		if _, err := f([]float64{1, 2, 2, 4}, 2); !errors.Is(err, ErrSingular) {
			t.Errorf("κ_%s singular: %v", name, err)
		}
	}
}

func TestMatPowerInteger_workedExample(t *testing.T) {
	got, err := MatPowerInteger(nil, []float64{2, 0, 0, 4}, 2, -2)
	if err != nil || !slices.Equal(got, []float64{0.25, 0, 0, 0.0625}) {
		t.Errorf("A⁻² = (%v, %v)", got, err)
	}
	got, _ = MatPowerInteger(nil, []float64{1, 1, 1, 0}, 2, 5)
	if !slices.Equal(got, []float64{8, 5, 5, 3}) {
		t.Errorf("A⁵ = %v", got)
	}
	if _, err := MatPowerInteger(nil, []float64{1, 2, 2, 4}, 2, -1); !errors.Is(err, ErrSingular) {
		t.Errorf("singular: %v", err)
	}
	if _, err := MatPowerInteger(nil, []float64{1}, 1, math.MinInt); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("MinInt: %v", err)
	}
}

// CLRS §15.2: the chain order does not change the product; compare with left-to-right MatMul.
func TestMatMulChain_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for trial := range 50 {
		k := 1 + trial%6
		dims := make([]int, k+1)
		for i := range dims {
			dims[i] = 1 + r.IntN(6)
		}
		mats := make([][]float64, k)
		for i := range mats {
			mats[i] = make([]float64, dims[i]*dims[i+1])
			for j := range mats[i] {
				mats[i][j] = r.NormFloat64()
			}
		}
		got, err := MatMulChain(mats, dims)
		if err != nil {
			t.Fatal(err)
		}
		want := slices.Clone(mats[0])
		for i := 1; i < k; i++ {
			want, _ = MatMul(nil, want, dims[0], dims[i], mats[i], dims[i+1])
		}
		for i := range want {
			if math.Abs(got[i]-want[i]) > 1e-10*max(1, math.Abs(want[i])) {
				t.Fatalf("trial %d: %v vs %v", trial, got, want)
			}
		}
	}
	if _, err := MatMulChain(nil, nil); !errors.Is(err, ErrDimension) {
		t.Errorf("empty: %v", err)
	}
	if _, err := MatMulChain([][]float64{{1, 2}}, []int{1, 3}); !errors.Is(err, ErrDimension) {
		t.Errorf("size: %v", err)
	}
}

func TestVecMatMulDiagonalElementwise_workedExample(t *testing.T) {
	got, err := VecMatMul(nil, []float64{1, 2}, []float64{1, 2, 3, 4, 5, 6}, 2, 3)
	if err != nil || !slices.Equal(got, []float64{9, 12, 15}) {
		t.Errorf("vecmat = (%v, %v)", got, err)
	}
	if _, err := VecMatMul(nil, []float64{1}, []float64{1, 2}, 2, 1); !errors.Is(err, ErrDimension) {
		t.Errorf("vecmat dims: %v", err)
	}
	a := []float64{1, 2, 3, 4, 5, 6}
	for _, tc := range []struct {
		off  int
		want []float64
	}{{0, []float64{1, 5}}, {1, []float64{2, 6}}, {2, []float64{3}}, {-1, []float64{4}}, {3, []float64{}}, {-2, []float64{}}} {
		got, err := MatDiagonal(nil, a, 2, 3, tc.off)
		if err != nil || !slices.Equal(got, tc.want) {
			t.Errorf("diagonal %d = (%v, %v), want %v", tc.off, got, err, tc.want)
		}
	}
	got, err = VecHadamard(nil, []float64{1, 2, 3}, []float64{4, 5, 6})
	if err != nil || !slices.Equal(got, []float64{4, 10, 18}) {
		t.Errorf("hadamard = (%v, %v)", got, err)
	}
	got, err = VecDivElementwise(nil, []float64{4, 10, 18}, []float64{4, 5, 6})
	if err != nil || !slices.Equal(got, []float64{1, 2, 3}) {
		t.Errorf("divide = (%v, %v)", got, err)
	}
	if _, err := VecDivElementwise(nil, []float64{1}, []float64{0}); !errors.Is(err, ErrZeroDivisor) {
		t.Errorf("zero: %v", err)
	}
	if _, err := VecHadamard(nil, []float64{1e200}, []float64{1e200}); !errors.Is(err, ErrOverflow) {
		t.Errorf("overflow: %v", err)
	}
	if _, err := VecHadamard(nil, []float64{1}, []float64{1, 2}); !errors.Is(err, ErrDimension) {
		t.Errorf("length: %v", err)
	}
}

// GvL Theorem 2.4.1: U and V square orthogonal, A = U·Σ·Vᵀ.
func TestSVDFull_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for trial := range 60 {
		m, n := 1+trial%6, 1+(trial/6)%6
		a := make([]float64, m*n)
		for i := range a {
			a[i] = r.NormFloat64()
		}
		if trial%3 == 0 && n > 1 { // make a column repeat: rank deficient
			for i := range m {
				a[i*n+n-1] = a[i*n]
			}
		}
		res, err := SVDFull(a, m, n, -1)
		if err != nil {
			t.Fatal(err)
		}
		for _, q := range []struct {
			mat []float64
			p   int
		}{{res.U, m}, {res.V, n}} {
			qt, _ := MatTranspose(nil, q.mat, q.p, q.p)
			qq, _ := MatMul(nil, qt, q.p, q.p, q.mat, q.p)
			for i := range q.p {
				for j := range q.p {
					want := 0.0
					if i == j {
						want = 1
					}
					if math.Abs(qq[i*q.p+j]-want) > 1e-12 {
						t.Fatalf("trial %d: not orthogonal %v", trial, qq)
					}
				}
			}
		}
		k := min(m, n)
		sig := make([]float64, m*n)
		for i := range k {
			sig[i*n+i] = res.S[i]
		}
		us, _ := MatMul(nil, res.U, m, m, sig, n)
		vt, _ := MatTranspose(nil, res.V, n, n)
		back, _ := MatMul(nil, us, m, n, vt, n)
		for i := range a {
			if math.Abs(back[i]-a[i]) > 1e-12*max(1, math.Abs(a[i]))*10 {
				t.Fatalf("trial %d: U·Σ·Vᵀ ≠ A", trial)
			}
		}
	}
	res, _ := SVDFull([]float64{3, 4, 0}, 3, 1, -1)
	if !near(res.S[0], 5, 1e-15) || !near(math.Abs(res.U[0]), 0.6, 1e-15) || !near(math.Abs(res.U[3]), 0.8, 1e-15) {
		t.Errorf("column SVD = %v", res)
	}
}

func randSquare(n int) []float64 {
	r := rand.New(rand.NewPCG(1, 2))
	a := make([]float64, n*n)
	for i := range a {
		a[i] = r.NormFloat64()
	}
	return a
}

func BenchmarkLeastSquaresMinNorm(b *testing.B) {
	a, rhs := randSquare(50), make([]float64, 50)
	b.ReportAllocs()
	for b.Loop() {
		_, _, _ = LeastSquaresMinNorm(nil, a, 50, 50, rhs, -1)
	}
}

func BenchmarkMatLogDet(b *testing.B) {
	a := randSquare(100)
	b.ReportAllocs()
	for b.Loop() {
		_, _, _ = MatLogDet(a, 100)
	}
}

func BenchmarkMatNormNuclear(b *testing.B) {
	a := randSquare(50)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = MatNormNuclear(a, 50, 50)
	}
}

func BenchmarkMatNormNeg1(b *testing.B) {
	a := randSquare(100)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = MatNormNeg1(a, 100, 100)
	}
}

func BenchmarkMatNormNegInf(b *testing.B) {
	a := randSquare(100)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = MatNormNegInf(a, 100, 100)
	}
}

func BenchmarkMatNormNeg2(b *testing.B) {
	a := randSquare(50)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = MatNormNeg2(a, 50, 50)
	}
}

func BenchmarkVecNorm0(b *testing.B) {
	v := randSquare(100)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = VecNorm0(v)
	}
}

func BenchmarkMatCondition1(b *testing.B) {
	a := randSquare(100)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = MatCondition1(a, 100)
	}
}

func BenchmarkMatConditionInf(b *testing.B) {
	a := randSquare(100)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = MatConditionInf(a, 100)
	}
}

func BenchmarkMatConditionFrobenius(b *testing.B) {
	a := randSquare(100)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = MatConditionFrobenius(a, 100)
	}
}

func BenchmarkMatPowerInteger(b *testing.B) {
	a := randSquare(50)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = MatPowerInteger(nil, a, 50, -5)
	}
}

func BenchmarkMatMulChain(b *testing.B) {
	dims := []int{30, 35, 15, 5, 10, 20, 25}
	mats := make([][]float64, 6)
	for i := range mats {
		mats[i] = make([]float64, dims[i]*dims[i+1])
	}
	b.ReportAllocs()
	for b.Loop() {
		_, _ = MatMulChain(mats, dims)
	}
}

func BenchmarkVecMatMul(b *testing.B) {
	a, x, dst := randSquare(100), make([]float64, 100), make([]float64, 100)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = VecMatMul(dst, x, a, 100, 100)
	}
}

func BenchmarkMatDiagonal(b *testing.B) {
	a, dst := randSquare(100), make([]float64, 100)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = MatDiagonal(dst, a, 100, 100, 1)
	}
}

func BenchmarkVecHadamard(b *testing.B) {
	v, dst := make([]float64, 1000), make([]float64, 1000)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = VecHadamard(dst, v, v)
	}
}

func BenchmarkVecDivElementwise(b *testing.B) {
	v, w, dst := make([]float64, 1000), make([]float64, 1000), make([]float64, 1000)
	for i := range w {
		w[i] = 1
	}
	b.ReportAllocs()
	for b.Loop() {
		_, _ = VecDivElementwise(dst, v, w)
	}
}

func BenchmarkSVDFull(b *testing.B) {
	a := randSquare(30)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = SVDFull(a, 30, 30, -1)
	}
}
