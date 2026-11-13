package linalg

import (
	"errors"
	"math"
	"math/cmplx"
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"

	"github.com/CrimsonCoder42/calcbench/num/algebra"
)

func TestVandermonde_workedExample(t *testing.T) {
	got, err := Vandermonde(nil, []float64{1, 2, 3}, 3)
	if err != nil || !slices.Equal(got, []float64{1, 1, 1, 1, 2, 4, 1, 3, 9}) {
		t.Errorf("Vandermonde = (%v, %v)", got, err)
	}
	if _, err := Vandermonde(nil, []float64{1}, 0); !errors.Is(err, ErrDimension) {
		t.Errorf("cols 0: %v", err)
	}
	if _, err := Vandermonde(nil, []float64{math.Inf(1)}, 2); !errors.Is(err, ErrNotFinite) {
		t.Errorf("Inf: %v", err)
	}
	if _, err := Vandermonde(nil, []float64{1e200}, 3); !errors.Is(err, ErrOverflow) {
		t.Errorf("overflow: %v", err)
	}
	got, _ = Vandermonde(nil, nil, 2)
	if len(got) != 0 {
		t.Errorf("empty x: %v", got)
	}
}

// Strang, Introduction to Linear Algebra, 5th ed., §4.3 worked example: the best line through
// (0, 6), (1, 0), (2, 0) is 5 − 3t.
func TestPolynomialFit_workedExample(t *testing.T) {
	c, err := PolynomialFit(nil, []float64{0, 1, 2}, []float64{6, 0, 0}, nil, 1, -1)
	if err != nil || math.Abs(c[0]-5) > 1e-13 || math.Abs(c[1]+3) > 1e-13 {
		t.Errorf("PolynomialFit = (%v, %v), want (5, −3)", c, err)
	}
}

func TestPolynomialFit_cases(t *testing.T) {
	// Exact cubic recovered; weights on exact data change nothing.
	x := []float64{-2, -1, 0, 1, 2, 3}
	y := make([]float64, len(x))
	for i, v := range x {
		y[i] = 1 - 2*v + 0.5*v*v*v
	}
	w := []float64{1, 2, 3, 4, 5, 6}
	for _, weights := range [][]float64{nil, w} {
		c, err := PolynomialFit(nil, x, y, weights, 3, -1)
		if err != nil {
			t.Fatal(err)
		}
		for i, want := range []float64{1, -2, 0, 0.5} {
			if math.Abs(c[i]-want) > 1e-12 {
				t.Errorf("weights %v: c = %v", weights, c)
			}
		}
	}
	// Zero weight drops a point: fitting a line through 3 points with an outlier weighted 0.
	c, _ := PolynomialFit(nil, []float64{0, 1, 2}, []float64{1, 3, 100}, []float64{1, 1, 0}, 1, -1)
	if math.Abs(c[0]-1) > 1e-12 || math.Abs(c[1]-2) > 1e-12 {
		t.Errorf("zero weight: %v", c)
	}
	// Degree 0 is the mean.
	c, _ = PolynomialFit(nil, []float64{0, 1, 2}, []float64{1, 2, 6}, nil, 0, -1)
	if math.Abs(c[0]-3) > 1e-14 {
		t.Errorf("degree 0: %v", c)
	}
}

func TestPolynomialFit_errors(t *testing.T) {
	x := []float64{0, 1, 2}
	cases := []struct {
		name    string
		x, y, w []float64
		deg     int
		want    error
	}{
		{"negative degree", x, x, nil, -1, ErrDimension},
		{"length", x, x[:2], nil, 1, ErrDimension},
		{"weights length", x, x, []float64{1}, 1, ErrDimension},
		{"too few points", x, x, nil, 3, ErrDimension},
		{"NaN y", x, []float64{0, math.NaN(), 1}, nil, 1, ErrNotFinite},
		{"negative weight", x, x, []float64{1, -1, 1}, 1, ErrOutOfRange},
		{"repeated x", []float64{1, 1, 1}, x, nil, 1, ErrRankDeficient},
	}
	for _, tc := range cases {
		if _, err := PolynomialFit(nil, tc.x, tc.y, tc.w, tc.deg, -1); !errors.Is(err, tc.want) {
			t.Errorf("%s: %v, want %v", tc.name, err, tc.want)
		}
	}
}

// Oracle: the fitted residual is orthogonal to the weighted columns (normal equations hold).
func TestPolynomialFit_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for trial := range 100 {
		m, deg := 10+trial%20, trial%5
		x, y := make([]float64, m), make([]float64, m)
		for i := range x {
			x[i], y[i] = r.Float64()*4-2, r.NormFloat64()
		}
		c, err := PolynomialFit(nil, x, y, nil, deg, -1)
		if err != nil {
			t.Fatal(err)
		}
		for j := 0; j <= deg; j++ {
			s := 0.0
			for i := range x {
				p, _ := algebra.PolynomialEvaluate(c, x[i])
				s += math.Pow(x[i], float64(j)) * (p - y[i])
			}
			if math.Abs(s) > 1e-9 {
				t.Errorf("trial %d: column %d · residual = %g", trial, j, s)
			}
		}
	}
}

// Companion eigenvalues of a polynomial built from known roots are those roots.
func TestPolynomialRootsCompanion_oracle(t *testing.T) {
	got, err := PolynomialRootsCompanion(nil, []float64{2, -3, 1})
	sortComplex(got)
	if err != nil || cmplx.Abs(got[0]-1) > 1e-14 || cmplx.Abs(got[1]-2) > 1e-14 {
		t.Errorf("x² − 3x + 2 roots = (%v, %v)", got, err)
	}
	got, _ = PolynomialRootsCompanion(nil, []float64{1, 0, 1})
	sortComplex(got)
	if cmplx.Abs(got[0]+1i) > 1e-14 || cmplx.Abs(got[1]-1i) > 1e-14 {
		t.Errorf("x² + 1 roots = %v", got)
	}
	r := rand.New(rand.NewPCG(1, 2))
	for trial := range 100 {
		n := 1 + trial%8
		roots := make([]complex128, 0, n)
		for len(roots) < n {
			if n-len(roots) >= 2 && r.IntN(2) == 0 {
				z := complex(r.Float64()*4-2, r.Float64()*2+0.5)
				roots = append(roots, z, cmplx.Conj(z))
			} else {
				roots = append(roots, complex(r.Float64()*4-2, 0))
			}
		}
		coeffs, err := algebra.PolynomialFromRoots(roots, -1)
		if err != nil {
			t.Fatal(err)
		}
		got, err := PolynomialRootsCompanion(nil, coeffs)
		if err != nil {
			t.Fatal(err)
		}
		for _, want := range roots {
			best := math.Inf(1)
			for _, z := range got {
				best = min(best, cmplx.Abs(z-want))
			}
			if best > 1e-6 {
				t.Errorf("trial %d: root %v not found in %v", trial, want, got)
			}
		}
	}
	if _, err := PolynomialRootsCompanion(nil, []float64{1}); !errors.Is(err, algebra.ErrEmpty) {
		t.Errorf("constant: %v", err)
	}
	if _, err := PolynomialRootsCompanion(nil, []float64{1, 0}); !errors.Is(err, algebra.ErrDegenerate) {
		t.Errorf("zero lead: %v", err)
	}
}

func BenchmarkVandermonde(b *testing.B) {
	for _, m := range []int{1000, 10000, 100000} {
		x := make([]float64, m)
		for i := range x {
			x[i] = float64(i) / float64(m)
		}
		dst := make([]float64, m*4)
		b.Run(strconv.Itoa(m), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = Vandermonde(dst, x, 4)
			}
		})
	}
}

func BenchmarkPolynomialFit(b *testing.B) {
	for _, m := range []int{1000, 10000, 100000} {
		r := rand.New(rand.NewPCG(1, 2))
		x, y := make([]float64, m), make([]float64, m)
		for i := range x {
			x[i], y[i] = r.Float64(), r.NormFloat64()
		}
		b.Run(strconv.Itoa(m), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = PolynomialFit(nil, x, y, nil, 3, -1)
			}
		})
	}
}

func BenchmarkPolynomialRootsCompanion(b *testing.B) {
	for _, n := range []int{10, 50} {
		c := make([]float64, n+1)
		for i := range c {
			c[i] = float64(i%7) + 1
		}
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = PolynomialRootsCompanion(nil, c)
			}
		})
	}
}
