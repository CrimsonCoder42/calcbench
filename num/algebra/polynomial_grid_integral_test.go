package algebra

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Worked example hand-checked: 1 + 2y + 3x + 4xy on xs = (0, 2), ys = (1, 3) (Knuth §4.6.4 nested Horner).
func TestPolynomialEvaluateGrid2D_workedExample(t *testing.T) {
	got, err := PolynomialEvaluateGrid2D(nil, []float64{1, 2, 3, 4}, 2, 2, []float64{0, 2}, []float64{1, 3})
	if err != nil || !slices.Equal(got, []float64{3, 7, 17, 37}) {
		t.Errorf("grid = (%v, %v), want [3 7 17 37]", got, err)
	}
}

// Every grid point agrees with the pointwise NUM-ALGEB-021 on fixed-seed random grids.
func TestPolynomialEvaluateGrid2D_matchesPointwise(t *testing.T) {
	r := rand.New(rand.NewPCG(7, 11))
	for trial := 0; trial < 50; trial++ {
		rows, cols := 1+r.IntN(5), 1+r.IntN(5)
		c := make([]float64, rows*cols)
		for i := range c {
			c[i] = r.Float64()*4 - 2
		}
		xs := make([]float64, r.IntN(6))
		ys := make([]float64, r.IntN(6))
		for i := range xs {
			xs[i] = r.Float64()*2 - 1
		}
		for i := range ys {
			ys[i] = r.Float64()*2 - 1
		}
		got, err := PolynomialEvaluateGrid2D(nil, c, rows, cols, xs, ys)
		if err != nil || len(got) != len(xs)*len(ys) {
			t.Fatalf("trial %d: (%v, %v)", trial, got, err)
		}
		for i, x := range xs {
			for j, y := range ys {
				want, _ := PolynomialEvaluate2D(c, rows, cols, x, y)
				if math.Abs(got[i*len(ys)+j]-want) > 1e-13*(1+math.Abs(want)) {
					t.Errorf("trial %d (%d,%d): %v, want %v", trial, i, j, got[i*len(ys)+j], want)
				}
			}
		}
	}
}

func TestPolynomialEvaluateGrid2D_errors(t *testing.T) {
	c := []float64{1, 2, 3, 4}
	if _, err := PolynomialEvaluateGrid2D(nil, c, 0, 2, nil, nil); !errors.Is(err, ErrEmpty) {
		t.Errorf("rows 0: %v", err)
	}
	if _, err := PolynomialEvaluateGrid2D(nil, c, 2, 3, nil, nil); !errors.Is(err, ErrDegenerate) {
		t.Errorf("shape: %v", err)
	}
	if _, err := PolynomialEvaluateGrid2D(nil, []float64{1, math.NaN(), 3, 4}, 2, 2, nil, nil); !errors.Is(err, ErrNotFinite) {
		t.Errorf("NaN coeff: %v", err)
	}
	if _, err := PolynomialEvaluateGrid2D(nil, c, 2, 2, []float64{math.Inf(1)}, []float64{1}); !errors.Is(err, ErrNotFinite) {
		t.Errorf("Inf x: %v", err)
	}
	if _, err := PolynomialEvaluateGrid2D(nil, c, 2, 2, []float64{1}, []float64{math.NaN()}); !errors.Is(err, ErrNotFinite) {
		t.Errorf("NaN y: %v", err)
	}
	if _, err := PolynomialEvaluateGrid2D(nil, []float64{0, 0, 0, math.MaxFloat64}, 2, 2, []float64{1e200}, []float64{1e200}); !errors.Is(err, ErrOverflow) {
		t.Errorf("overflow: %v", err)
	}
	got, err := PolynomialEvaluateGrid2D(nil, c, 2, 2, nil, []float64{1})
	if err != nil || len(got) != 0 {
		t.Errorf("empty xs: (%v, %v)", got, err)
	}
}

// Hand-checked: ∫∫6 from lower bound 1 with zero constants is 3(x − 1)²; numpy polyint semantics.
func TestPolynomialIntegralOrder_workedExample(t *testing.T) {
	got, err := PolynomialIntegralOrder([]float64{6}, 2, 1, nil)
	if err != nil || !slices.Equal(got, []float64{3, -6, 3}) {
		t.Errorf("order 2 = (%v, %v), want [3 -6 3]", got, err)
	}
	// Order 1, bound 0 must equal PolynomialIntegral (NUM-ALGEB-016).
	a := []float64{2, 6, -3}
	want, _ := PolynomialIntegral(a, 1.5)
	got, err = PolynomialIntegralOrder(a, 1, 0, []float64{1.5})
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("order 1 = (%v, %v), want %v", got, err, want)
	}
	got, err = PolynomialIntegralOrder(a, 0, 0, nil)
	if err != nil || !slices.Equal(got, a) {
		t.Errorf("order 0 = (%v, %v)", got, err)
	}
}

// Each step's constant is met at the lower bound, and differentiating m times recovers a.
func TestPolynomialIntegralOrder_properties(t *testing.T) {
	a := []float64{1, -2, 0.5, 4}
	k := []float64{2, -1, 0.25}
	lb := -0.75
	got, err := PolynomialIntegralOrder(a, 3, lb, k)
	if err != nil || len(got) != 7 {
		t.Fatalf("(%v, %v)", got, err)
	}
	cur := got
	for step := 2; step >= 0; step-- {
		v, _ := PolynomialEvaluate(cur, lb)
		if math.Abs(v-k[step]) > 1e-13 {
			t.Errorf("step %d: P(lb) = %v, want %v", step, v, k[step])
		}
		cur, _ = PolynomialDerivative(cur, 1)
	}
	for i := range a {
		if math.Abs(cur[i]-a[i]) > 1e-13 {
			t.Errorf("derivative³ = %v, want %v", cur, a)
			break
		}
	}
}

func TestPolynomialIntegralOrder_errors(t *testing.T) {
	if _, err := PolynomialIntegralOrder(nil, 1, 0, nil); !errors.Is(err, ErrEmpty) {
		t.Errorf("empty: %v", err)
	}
	if _, err := PolynomialIntegralOrder([]float64{1}, -1, 0, nil); !errors.Is(err, ErrDegenerate) {
		t.Errorf("negative order: %v", err)
	}
	if _, err := PolynomialIntegralOrder([]float64{1}, 1, 0, []float64{1, 2}); !errors.Is(err, ErrDegenerate) {
		t.Errorf("too many constants: %v", err)
	}
	if _, err := PolynomialIntegralOrder([]float64{1}, 1, math.NaN(), nil); !errors.Is(err, ErrNotFinite) {
		t.Errorf("NaN bound: %v", err)
	}
	if _, err := PolynomialIntegralOrder([]float64{1}, 1, 0, []float64{math.Inf(-1)}); !errors.Is(err, ErrNotFinite) {
		t.Errorf("Inf constant: %v", err)
	}
	if _, err := PolynomialIntegralOrder([]float64{math.MaxFloat64}, 2, 1e300, nil); !errors.Is(err, ErrOverflow) {
		t.Errorf("overflow: %v", err)
	}
}

func BenchmarkPolynomialEvaluateGrid2D(b *testing.B) {
	c := make([]float64, 6*6)
	for i := range c {
		c[i] = float64(i%7) - 3
	}
	xs := make([]float64, 64)
	ys := make([]float64, 64)
	for i := range xs {
		xs[i] = float64(i)/64 - 0.5
		ys[i] = float64(i)/32 - 1
	}
	dst := make([]float64, 64*64)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = PolynomialEvaluateGrid2D(dst, c, 6, 6, xs, ys)
	}
}

func BenchmarkPolynomialIntegralOrder(b *testing.B) {
	a := []float64{1, -2, 0.5, 4, 3, -1, 2, 0.25}
	k := []float64{1, 2, 3}
	b.ReportAllocs()
	for b.Loop() {
		_, _ = PolynomialIntegralOrder(a, 3, 0.5, k)
	}
}
