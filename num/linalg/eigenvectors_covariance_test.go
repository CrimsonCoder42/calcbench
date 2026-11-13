package linalg

import (
	"errors"
	"math"
	"math/cmplx"
	"math/rand/v2"
	"strconv"
	"testing"
)

// eigResidual returns max over k of ‖A·vₖ − λₖ·vₖ‖₂ / (‖A‖₁·‖vₖ‖₂) and the largest |‖vₖ‖₂ − 1|.
func eigResidual(a []float64, n int, vals, vecs []complex128) (res, normErr float64) {
	anorm := 0.0
	for j := range n {
		s := 0.0
		for i := range n {
			s += math.Abs(a[i*n+j])
		}
		anorm = max(anorm, s)
	}
	if anorm == 0 {
		anorm = 1
	}
	for k := range n {
		r, vn := 0.0, 0.0
		for i := range n {
			var s complex128
			for j := range n {
				s += complex(a[i*n+j], 0) * vecs[j*n+k]
			}
			s -= vals[k] * vecs[i*n+k]
			r = math.Hypot(r, cmplx.Abs(s))
			vn = math.Hypot(vn, cmplx.Abs(vecs[i*n+k]))
		}
		res = max(res, r/(anorm*vn))
		normErr = max(normErr, math.Abs(vn-1))
	}
	return res, normErr
}

// Worked examples hand-checked (GvL §7.6.1 setting): [2 1; 0 3] has λ = 2, 3 with v = (1, 0) and
// (1, 1)/√2; the rotation [0 −1; 1 0] has λ = ±i with v = (1, ∓i)/√2.
func TestMatEigenvectors_workedExample(t *testing.T) {
	vals, vecs, err := MatEigenvectors(nil, nil, []float64{2, 1, 0, 3}, 2)
	if err != nil {
		t.Fatal(err)
	}
	for k := range 2 {
		v := []complex128{vecs[k], vecs[2+k]}
		var want []complex128
		switch {
		case cmplx.Abs(vals[k]-2) < 1e-12:
			want = []complex128{1, 0}
		case cmplx.Abs(vals[k]-3) < 1e-12:
			want = []complex128{complex(math.Sqrt2/2, 0), complex(math.Sqrt2/2, 0)}
		default:
			t.Fatalf("eigenvalue %v", vals[k])
		}
		for i := range 2 {
			if cmplx.Abs(v[i]-want[i]) > 1e-12 {
				t.Errorf("λ = %v: v = %v, want %v", vals[k], v, want)
			}
		}
	}
	vals, vecs, err = MatEigenvectors(nil, nil, []float64{0, -1, 1, 0}, 2)
	if err != nil {
		t.Fatal(err)
	}
	for k := range 2 {
		// For λ = i: (1, −i)/√2 up to a unit phase; check v₂/v₁ = −i·sign(Im λ).
		ratio := vecs[2+k] / vecs[k]
		want := complex(0, -imag(vals[k]))
		if cmplx.Abs(ratio-want) > 1e-12 {
			t.Errorf("λ = %v: v₂/v₁ = %v, want %v", vals[k], ratio, want)
		}
	}
	if k1, k0 := vecs[1], vecs[0]; cmplx.Abs(k1-cmplx.Conj(k0)) > 0 {
		t.Errorf("conjugate pair vectors not conjugate: %v %v", k0, k1)
	}
}

// A known 3×3 nonsymmetric example: A = [1 2 0; 0 3 0; 2 −4 2] has λ = 1, 2, 3 with eigenvectors
// (−1, 0, 2), (0, 0, 1), (−1, −1, 2) (hand-checked: A·v = λ·v for each).
func TestMatEigenvectors_known3x3(t *testing.T) {
	a := []float64{1, 2, 0, 0, 3, 0, 2, -4, 2}
	vals, vecs, err := MatEigenvectors(nil, nil, a, 3)
	if err != nil {
		t.Fatal(err)
	}
	want := map[int][]float64{1: {-1, 0, 2}, 2: {0, 0, 1}, 3: {-1, -1, 2}}
	for k := range 3 {
		lam := int(math.Round(real(vals[k])))
		w, ok := want[lam]
		if !ok || math.Abs(imag(vals[k])) > 1e-12 {
			t.Fatalf("eigenvalue %v", vals[k])
		}
		wn := math.Sqrt(w[0]*w[0] + w[1]*w[1] + w[2]*w[2])
		// Compare up to sign: |vᴴw|/‖w‖ = 1.
		var dot complex128
		for i := range 3 {
			if imag(vecs[i*3+k]) != 0 {
				t.Errorf("λ = %d: vector not real: %v", lam, vecs[i*3+k])
			}
			dot += cmplx.Conj(vecs[i*3+k]) * complex(w[i]/wn, 0)
		}
		if math.Abs(cmplx.Abs(dot)-1) > 1e-12 {
			t.Errorf("λ = %d: |⟨v, w⟩| = %v, want 1", lam, cmplx.Abs(dot))
		}
	}
	if res, ne := eigResidual(a, 3, vals, vecs); res > 1e-13 || ne > 1e-13 {
		t.Errorf("residual %g, norm error %g", res, ne)
	}
}

// A·v = λ·v (complex) on fixed-seed random nonsymmetric matrices of order ≤ 24.
func TestMatEigenvectors_randomResidual(t *testing.T) {
	r := rand.New(rand.NewPCG(2026, 1113))
	for trial := 0; trial < 60; trial++ {
		n := 1 + r.IntN(24)
		a := make([]float64, n*n)
		for i := range a {
			a[i] = r.NormFloat64()
		}
		vals, vecs, err := MatEigenvectors(nil, nil, a, n)
		if err != nil {
			t.Fatalf("trial %d n=%d: %v", trial, n, err)
		}
		if res, ne := eigResidual(a, n, vals, vecs); res > 1e-11 || ne > 1e-12 {
			t.Errorf("trial %d n=%d: residual %g, norm error %g", trial, n, res, ne)
		}
	}
}

// Repeated eigenvalues: the identity and a diagonalisable matrix with a double eigenvalue get
// independent vectors; a Jordan block (defective) still satisfies A·v = λ·v; badly scaled input works.
func TestMatEigenvectors_cases(t *testing.T) {
	id := []float64{1, 0, 0, 0, 1, 0, 0, 0, 1}
	vals, vecs, err := MatEigenvectors(nil, nil, id, 3)
	if err != nil {
		t.Fatal(err)
	}
	if res, _ := eigResidual(id, 3, vals, vecs); res > 1e-14 {
		t.Errorf("identity residual %g", res)
	}
	for k := range 3 {
		for l := k + 1; l < 3; l++ {
			var dot complex128
			for i := range 3 {
				dot += cmplx.Conj(vecs[i*3+k]) * vecs[i*3+l]
			}
			if cmplx.Abs(dot) > 1e-10 {
				t.Errorf("identity vectors %d, %d not independent: ⟨⟩ = %v", k, l, dot)
			}
		}
	}
	// S·diag(2, 2, 5)·S⁻¹ with S = [1 1 0; 0 1 1; 1 0 1] (det 2): a double eigenvalue, full eigenspace.
	d := []float64{2, 0, 0, -1.5, 3.5, 1.5, -1.5, 1.5, 3.5}
	vals, vecs, err = MatEigenvectors(nil, nil, d, 3)
	if err != nil {
		t.Fatal(err)
	}
	if res, _ := eigResidual(d, 3, vals, vecs); res > 1e-10 {
		t.Errorf("double eigenvalue residual %g (values %v)", res, vals)
	}
	var twos []int
	for k := range 3 {
		if cmplx.Abs(vals[k]-2) < 1e-8 {
			twos = append(twos, k)
		}
	}
	if len(twos) != 2 {
		t.Fatalf("double eigenvalue: values %v", vals)
	}
	var dot complex128
	for i := range 3 {
		dot += cmplx.Conj(vecs[i*3+twos[0]]) * vecs[i*3+twos[1]]
	}
	if cmplx.Abs(dot) > 1-1e-6 {
		t.Errorf("double eigenvalue vectors parallel: |⟨⟩| = %v", cmplx.Abs(dot))
	}
	jordan := []float64{2, 1, 0, 2}
	vals, vecs, err = MatEigenvectors(nil, nil, jordan, 2)
	if err != nil {
		t.Fatal(err)
	}
	if res, _ := eigResidual(jordan, 2, vals, vecs); res > 1e-7 {
		t.Errorf("Jordan residual %g", res)
	}
	zero := make([]float64, 4)
	vals, vecs, err = MatEigenvectors(nil, nil, zero, 2)
	if err != nil || vals[0] != 0 || vals[1] != 0 {
		t.Fatalf("zero: (%v, %v)", vals, err)
	}
	if res, ne := eigResidual(zero, 2, vals, vecs); res > 0 || ne > 1e-15 {
		t.Errorf("zero matrix residual %g, norm %g", res, ne)
	}
	scaled := []float64{1e300, 2e300, 3e300, 4e300}
	vals, vecs, err = MatEigenvectors(nil, nil, scaled, 2)
	if err != nil {
		t.Fatal(err)
	}
	if res, _ := eigResidual([]float64{1, 2, 3, 4}, 2, []complex128{vals[0] / 1e300, vals[1] / 1e300}, vecs); res > 1e-14 {
		t.Errorf("1e300-scaled residual %g", res)
	}
}

func TestMatEigenvectors_errors(t *testing.T) {
	if _, _, err := MatEigenvectors(nil, nil, nil, 0); !errors.Is(err, ErrDimension) {
		t.Errorf("n = 0: %v", err)
	}
	if _, _, err := MatEigenvectors(nil, nil, []float64{1, 2, 3}, 2); !errors.Is(err, ErrDimension) {
		t.Errorf("len: %v", err)
	}
	if _, _, err := MatEigenvectors(nil, nil, []float64{1, math.NaN(), 0, 1}, 2); !errors.Is(err, ErrNotFinite) {
		t.Errorf("NaN: %v", err)
	}
}

// Hand-checked: y = 5 − 3x through (0, 6), (1, 0), (2, 0) has σ² = 6 and (AᵀA)⁻¹ = [5/6 −1/2; −1/2 1/2],
// so cov = [5 −3; −3 3] (Björck 1996 §1.2: cov = σ²(AᵀA)⁻¹).
func TestPolynomialFitCovariance_workedExample(t *testing.T) {
	c, cov, err := PolynomialFitCovariance(nil, nil, []float64{0, 1, 2}, []float64{6, 0, 0}, nil, 1, -1)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(c[0]-5) > 1e-13 || math.Abs(c[1]+3) > 1e-13 {
		t.Errorf("c = %v, want (5, −3)", c)
	}
	want := []float64{5, -3, -3, 3}
	for i := range want {
		if math.Abs(cov[i]-want[i]) > 1e-12 {
			t.Errorf("cov = %v, want %v", cov, want)
			break
		}
	}
}

// Against the explicit normal-equations formula on a well-conditioned weighted fit: cov =
// σ²·(VᵀW²V)⁻¹, with (VᵀW²V)⁻¹ checked by multiplying back to the identity.
func TestPolynomialFitCovariance_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(5, 3))
	m, deg := 30, 3
	d := deg + 1
	x := make([]float64, m)
	y := make([]float64, m)
	w := make([]float64, m)
	for i := range m {
		x[i] = float64(i)/float64(m-1)*2 - 1
		y[i] = 1 - 2*x[i] + 0.5*x[i]*x[i]*x[i] + 0.1*r.NormFloat64()
		w[i] = 0.5 + r.Float64()
	}
	c, cov, err := PolynomialFitCovariance(nil, nil, x, y, w, deg, -1)
	if err != nil {
		t.Fatal(err)
	}
	ssr := 0.0
	ata := make([]float64, d*d)
	for i := range m {
		p, pw := 0.0, 1.0
		for j := range d {
			p += c[j] * pw
			pw *= x[i]
		}
		ssr += w[i] * w[i] * (p - y[i]) * (p - y[i])
		for j := range d {
			for k := range d {
				ata[j*d+k] += w[i] * w[i] * math.Pow(x[i], float64(j+k))
			}
		}
	}
	sigma2 := ssr / float64(m-d)
	for i := range d {
		for k := range d {
			s := 0.0
			for j := range d {
				s += ata[i*d+j] * cov[j*d+k] / sigma2
			}
			want := 0.0
			if i == k {
				want = 1
			}
			if math.Abs(s-want) > 1e-10 {
				t.Errorf("(VᵀW²V)·cov/σ² [%d][%d] = %v, want %v", i, k, s, want)
			}
		}
		if cov[i*d+i] <= 0 {
			t.Errorf("cov[%d][%d] = %v ≤ 0", i, i, cov[i*d+i])
		}
	}
}

func TestPolynomialFitCovariance_errors(t *testing.T) {
	if _, _, err := PolynomialFitCovariance(nil, nil, []float64{0, 1}, []float64{1, 2}, nil, 1, -1); !errors.Is(err, ErrDimension) {
		t.Errorf("m = d: %v", err)
	}
	if _, _, err := PolynomialFitCovariance(nil, nil, []float64{0, 1, 2}, []float64{1, 2}, nil, 1, -1); !errors.Is(err, ErrDimension) {
		t.Errorf("lengths: %v", err)
	}
	if _, _, err := PolynomialFitCovariance(nil, nil, []float64{0, 1, math.NaN()}, []float64{1, 2, 3}, nil, 1, -1); !errors.Is(err, ErrNotFinite) {
		t.Errorf("NaN: %v", err)
	}
	if _, _, err := PolynomialFitCovariance(nil, nil, []float64{1, 1, 1}, []float64{1, 2, 3}, nil, 1, -1); !errors.Is(err, ErrRankDeficient) {
		t.Errorf("rank: %v", err)
	}
}

func BenchmarkPolynomialFitCovariance(b *testing.B) {
	m := 200
	x := make([]float64, m)
	y := make([]float64, m)
	for i := range m {
		x[i] = float64(i) / float64(m)
		y[i] = math.Sin(3 * x[i])
	}
	c := make([]float64, 6)
	cov := make([]float64, 36)
	b.ReportAllocs()
	for b.Loop() {
		_, _, _ = PolynomialFitCovariance(c, cov, x, y, nil, 5, -1)
	}
}

func BenchmarkMatEigenvectors(b *testing.B) {
	for _, n := range []int{8, 32} {
		r := rand.New(rand.NewPCG(1, uint64(n)))
		a := make([]float64, n*n)
		for i := range a {
			a[i] = r.NormFloat64()
		}
		vals := make([]complex128, n)
		vecs := make([]complex128, n*n)
		b.Run("n="+strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = MatEigenvectors(vals, vecs, a, n)
			}
		})
	}
}
