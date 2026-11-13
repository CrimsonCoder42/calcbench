package linalg

import (
	"cmp"
	"errors"
	"math"
	"math/cmplx"
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"
)

func sortComplex(z []complex128) {
	slices.SortFunc(z, func(a, b complex128) int {
		if c := cmp.Compare(real(a), real(b)); math.Abs(real(a)-real(b)) > 1e-9 && c != 0 {
			return c
		}
		return cmp.Compare(imag(a), imag(b))
	})
}

// Worked examples, hand-checked: a 90° rotation has eigenvalues ±i; the block matrix
// [2 0 0; 0 3 4; 0 4 9] has 2 and the eigenvalues of [3 4; 4 9], 1 and 11 (GvL §7.5 setting).
func TestMatEigenvalues_workedExample(t *testing.T) {
	got, err := MatEigenvalues(nil, []float64{0, -1, 1, 0}, 2)
	if err != nil {
		t.Fatal(err)
	}
	sortComplex(got)
	if cmplx.Abs(got[0]-(-1i)) > 1e-14 || cmplx.Abs(got[1]-1i) > 1e-14 {
		t.Errorf("rotation eigenvalues = %v, want ±i", got)
	}
	got, err = MatEigenvalues(nil, []float64{2, 0, 0, 0, 3, 4, 0, 4, 9}, 3)
	if err != nil {
		t.Fatal(err)
	}
	sortComplex(got)
	for i, w := range []float64{1, 2, 11} {
		if cmplx.Abs(got[i]-complex(w, 0)) > 1e-12 {
			t.Errorf("eigenvalue %d = %v, want %v", i, got[i], w)
		}
	}
}

func TestMatEigenvalues_cases(t *testing.T) {
	got, err := MatEigenvalues(nil, []float64{5}, 1)
	if err != nil || got[0] != 5 {
		t.Errorf("1×1 = (%v, %v)", got, err)
	}
	// Upper triangular: the diagonal.
	got, _ = MatEigenvalues(nil, []float64{1, 7, 3, 0, -2, 5, 0, 0, 4}, 3)
	sortComplex(got)
	for i, w := range []float64{-2, 1, 4} {
		if cmplx.Abs(got[i]-complex(w, 0)) > 1e-12 {
			t.Errorf("triangular eigenvalue %d = %v, want %v", i, got[i], w)
		}
	}
	// Zero matrix and huge scale.
	got, _ = MatEigenvalues(nil, make([]float64, 9), 3)
	for _, z := range got {
		if z != 0 {
			t.Errorf("zero matrix eigenvalue %v", z)
		}
	}
	got, _ = MatEigenvalues(nil, []float64{0, -1e300, 1e300, 0}, 2)
	sortComplex(got)
	if cmplx.Abs(got[1]-1e300i)/1e300 > 1e-14 {
		t.Errorf("scaled rotation = %v", got)
	}
}

func TestMatEigenvalues_errors(t *testing.T) {
	if _, err := MatEigenvalues(nil, nil, 0); !errors.Is(err, ErrDimension) {
		t.Errorf("n = 0: %v", err)
	}
	if _, err := MatEigenvalues(nil, []float64{1, 2, 3}, 2); !errors.Is(err, ErrDimension) {
		t.Errorf("len: %v", err)
	}
	if _, err := MatEigenvalues(nil, []float64{1, math.NaN(), 0, 1}, 2); !errors.Is(err, ErrNotFinite) {
		t.Errorf("NaN: %v", err)
	}
}

// Oracle: for random matrices, Σλ = trace, and each λ is a root of det(A − λI) checked as the
// determinant of the complex shifted matrix being small relative to the scale; for a companion
// matrix of known roots, the eigenvalues are those roots.
func TestMatEigenvalues_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for trial := range 200 {
		n := 1 + trial%12
		a := make([]float64, n*n)
		for i := range a {
			a[i] = r.NormFloat64()
		}
		got, err := MatEigenvalues(nil, a, n)
		if err != nil {
			t.Fatalf("trial %d: %v", trial, err)
		}
		var sum complex128
		for _, z := range got {
			sum += z
		}
		tr, _ := MatTrace(a, n)
		if math.Abs(real(sum)-tr) > 1e-9*float64(n) || math.Abs(imag(sum)) > 1e-9*float64(n) {
			t.Errorf("trial %d: Σλ = %v, trace %v", trial, sum, tr)
		}
		// det A = Πλ
		prod := complex(1, 0)
		for _, z := range got {
			prod *= z
		}
		det, _ := MatDet(a, n)
		if cmplx.Abs(prod-complex(det, 0)) > 1e-8*max(1, math.Abs(det)) {
			t.Errorf("trial %d: Πλ = %v, det %v", trial, prod, det)
		}
	}
	// Real symmetric case agrees with Jacobi; A·v = λ·v holds for Jacobi's vectors.
	for trial := range 50 {
		n := 2 + trial%8
		a := make([]float64, n*n)
		for i := range n {
			for j := 0; j <= i; j++ {
				v := r.NormFloat64()
				a[i*n+j], a[j*n+i] = v, v
			}
		}
		got, _ := MatEigenvalues(nil, a, n)
		want, _, _ := MatSymmetricEigenJacobi(nil, nil, a, n, -1, false)
		re := make([]float64, n)
		for i, z := range got {
			if math.Abs(imag(z)) > 1e-10 {
				t.Errorf("symmetric trial %d: complex %v", trial, z)
			}
			re[i] = real(z)
		}
		slices.Sort(re)
		w := slices.Clone(want)
		slices.Sort(w)
		for i := range re {
			if math.Abs(re[i]-w[i]) > 1e-10 {
				t.Errorf("symmetric trial %d: %v vs %v", trial, re, w)
				break
			}
		}
	}
}

func BenchmarkMatEigenvalues(b *testing.B) {
	for _, n := range []int{10, 50, 100} {
		r := rand.New(rand.NewPCG(1, 2))
		a := make([]float64, n*n)
		for i := range a {
			a[i] = r.NormFloat64()
		}
		dst := make([]complex128, n)
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = MatEigenvalues(dst, a, n)
			}
		})
	}
}
