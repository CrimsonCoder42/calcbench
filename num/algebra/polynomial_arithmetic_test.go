package algebra

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Worked examples are hand-checked against the definitions in Knuth, TAOCP vol. 2, 3rd ed., §4.6.
func TestPolynomialAddSubtract_workedExample(t *testing.T) {
	got, err := PolynomialAdd([]float64{1, 2}, []float64{3, 4, 5})
	if err != nil || !slices.Equal(got, []float64{4, 6, 5}) {
		t.Errorf("add = (%v, %v)", got, err)
	}
	got, err = PolynomialSubtract([]float64{1, 2}, []float64{3, 4, 5})
	if err != nil || !slices.Equal(got, []float64{-2, -2, -5}) {
		t.Errorf("sub = (%v, %v)", got, err)
	}
	if _, err := PolynomialAdd(nil, []float64{1}); !errors.Is(err, ErrEmpty) {
		t.Errorf("empty: %v", err)
	}
	if _, err := PolynomialSubtract([]float64{math.NaN()}, []float64{1}); !errors.Is(err, ErrNotFinite) {
		t.Errorf("NaN: %v", err)
	}
	if _, err := PolynomialAdd([]float64{math.MaxFloat64}, []float64{math.MaxFloat64}); !errors.Is(err, ErrOverflow) {
		t.Errorf("overflow: %v", err)
	}
}

// Knuth §4.6.1 Algorithm D; example hand-checked: x³ − 2x² − 4 = (x − 3)(x² + x + 3) + 5.
func TestPolynomialDivide_workedExample(t *testing.T) {
	q, r, err := PolynomialDivide([]float64{-4, 0, -2, 1}, []float64{-3, 1})
	if err != nil || !slices.Equal(q, []float64{3, 1, 1}) || !slices.Equal(r, []float64{5}) {
		t.Errorf("divide = (%v, %v, %v)", q, r, err)
	}
	q, r, _ = PolynomialDivide([]float64{1, 2}, []float64{1, 1, 1})
	if !slices.Equal(q, []float64{0}) || !slices.Equal(r, []float64{1, 2}) {
		t.Errorf("low degree: %v %v", q, r)
	}
	q, r, _ = PolynomialDivide([]float64{2, 4}, []float64{2})
	if !slices.Equal(q, []float64{1, 2}) || !slices.Equal(r, []float64{0}) {
		t.Errorf("constant divisor: %v %v", q, r)
	}
	if _, _, err := PolynomialDivide([]float64{1}, []float64{1, 0}); !errors.Is(err, ErrDegenerate) {
		t.Errorf("zero lead: %v", err)
	}
}

// Oracle: a = q·b + r on random inputs.
func TestPolynomialDivide_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		a := make([]float64, 1+rng.IntN(8))
		b := make([]float64, 1+rng.IntN(5))
		for i := range a {
			a[i] = rng.NormFloat64()
		}
		for i := range b {
			b[i] = rng.NormFloat64()
		}
		b[len(b)-1] = 1 + rng.Float64()
		q, r, err := PolynomialDivide(a, b)
		if err != nil {
			t.Fatal(err)
		}
		qb, _ := PolynomialMultiply(q, b)
		back, _ := PolynomialAdd(qb, r)
		for i := range back {
			want := 0.0
			if i < len(a) {
				want = a[i]
			}
			if math.Abs(back[i]-want) > 1e-9 {
				t.Fatalf("a = %v, q·b + r = %v", a, back)
			}
		}
	}
}

func TestPolynomialPower_workedExample(t *testing.T) {
	got, err := PolynomialPower([]float64{1, 1}, 3)
	if err != nil || !slices.Equal(got, []float64{1, 3, 3, 1}) {
		t.Errorf("(1+x)³ = (%v, %v)", got, err)
	}
	got, _ = PolynomialPower([]float64{2, 5}, 0)
	if !slices.Equal(got, []float64{1}) {
		t.Errorf("p⁰ = %v", got)
	}
	if _, err := PolynomialPower([]float64{1}, -1); !errors.Is(err, ErrDegenerate) {
		t.Errorf("negative k: %v", err)
	}
}

func TestPolynomialDerivativeIntegral_workedExample(t *testing.T) {
	got, err := PolynomialDerivative([]float64{1, 2, 3}, 1)
	if err != nil || !slices.Equal(got, []float64{2, 6}) {
		t.Errorf("d/dx = (%v, %v)", got, err)
	}
	got, _ = PolynomialDerivative([]float64{1, 2, 3}, 2)
	if !slices.Equal(got, []float64{6}) {
		t.Errorf("d²/dx² = %v", got)
	}
	got, _ = PolynomialDerivative([]float64{1, 2, 3}, 5)
	if !slices.Equal(got, []float64{0}) {
		t.Errorf("high order = %v", got)
	}
	if _, err := PolynomialDerivative([]float64{1}, -1); !errors.Is(err, ErrDegenerate) {
		t.Errorf("negative order: %v", err)
	}
	got, err = PolynomialIntegral([]float64{2, 6}, 1)
	if err != nil || !slices.Equal(got, []float64{1, 2, 3}) {
		t.Errorf("∫ = (%v, %v)", got, err)
	}
	if _, err := PolynomialIntegral([]float64{1}, math.Inf(1)); !errors.Is(err, ErrNotFinite) {
		t.Errorf("Inf constant: %v", err)
	}
	// Derivative undoes the integral.
	back, _ := PolynomialDerivative(got, 1)
	if !slices.Equal(back, []float64{2, 6}) {
		t.Errorf("d/dx ∫ = %v", back)
	}
}

func TestPolynomialFromRoots_workedExample(t *testing.T) {
	got, err := PolynomialFromRoots([]complex128{1, 2}, -1)
	if err != nil || !slices.Equal(got, []float64{2, -3, 1}) {
		t.Errorf("roots 1, 2 = (%v, %v)", got, err)
	}
	got, _ = PolynomialFromRoots([]complex128{1i, -1i}, -1)
	if !slices.Equal(got, []float64{1, 0, 1}) {
		t.Errorf("roots ±i = %v", got)
	}
	got, _ = PolynomialFromRoots(nil, -1)
	if !slices.Equal(got, []float64{1}) {
		t.Errorf("no roots = %v", got)
	}
	if _, err := PolynomialFromRoots([]complex128{1i}, -1); !errors.Is(err, ErrDegenerate) {
		t.Errorf("unpaired: %v", err)
	}
	v, err := PolynomialValueFromRoots([]float64{1, 2}, 4)
	if err != nil || v != 6 {
		t.Errorf("value from roots = (%v, %v)", v, err)
	}
	// Oracle: PolynomialEvaluate of the expanded form agrees with the product form.
	rng := rand.New(rand.NewPCG(1, 2))
	for range 200 {
		roots := make([]float64, 1+rng.IntN(6))
		croots := make([]complex128, len(roots))
		for i := range roots {
			roots[i] = rng.Float64()*4 - 2
			croots[i] = complex(roots[i], 0)
		}
		c, _ := PolynomialFromRoots(croots, -1)
		x := rng.Float64()*4 - 2
		p, _ := PolynomialEvaluate(c, x)
		q, _ := PolynomialValueFromRoots(roots, x)
		if math.Abs(p-q) > 1e-10 {
			t.Fatalf("roots %v at %v: %v vs %v", roots, x, p, q)
		}
	}
}

func TestPolynomialCompanionTrim2D_workedExample(t *testing.T) {
	got, err := PolynomialCompanion(nil, []float64{2, -3, 1})
	if err != nil || !slices.Equal(got, []float64{0, -2, 1, 3}) {
		t.Errorf("companion = (%v, %v)", got, err)
	}
	if _, err := PolynomialCompanion(nil, []float64{1}); !errors.Is(err, ErrEmpty) {
		t.Errorf("constant: %v", err)
	}
	if _, err := PolynomialCompanion(nil, []float64{1, 0}); !errors.Is(err, ErrDegenerate) {
		t.Errorf("zero lead: %v", err)
	}
	tr, err := PolynomialTrim([]float64{1, 2, 0, 1e-20}, 1e-15)
	if err != nil || !slices.Equal(tr, []float64{1, 2}) {
		t.Errorf("trim = (%v, %v)", tr, err)
	}
	tr, _ = PolynomialTrim([]float64{0, 0}, 0)
	if !slices.Equal(tr, []float64{0}) {
		t.Errorf("trim zero = %v", tr)
	}
	if _, err := PolynomialTrim([]float64{1}, -1); !errors.Is(err, ErrDegenerate) {
		t.Errorf("negative tol: %v", err)
	}
	v, err := PolynomialEvaluate2D([]float64{1, 2, 3, 4}, 2, 2, 2, 3)
	if err != nil || v != 37 {
		t.Errorf("2D = (%v, %v)", v, err)
	}
	if _, err := PolynomialEvaluate2D([]float64{1, 2, 3}, 2, 2, 0, 0); !errors.Is(err, ErrDegenerate) {
		t.Errorf("grid size: %v", err)
	}
}

// Knuth TAOCP vol. 1 §1.2.4: x mod y has the sign of y; IEEE 754 §5.3.1 ties to even.
func TestScalarAlgebra_workedExample(t *testing.T) {
	for _, tc := range []struct{ x, y, q, r float64 }{{7, 3, 2, 1}, {-7, 3, -3, 2}, {7, -3, -3, -2}, {-7, -3, 2, -1}, {5.5, 2, 2, 1.5}} {
		q, r, err := FloorDivMod(tc.x, tc.y)
		if err != nil || q != tc.q || r != tc.r {
			t.Errorf("FloorDivMod(%v, %v) = (%v, %v, %v)", tc.x, tc.y, q, r, err)
		}
		m, _ := FloorMod(tc.x, tc.y)
		if m != tc.r {
			t.Errorf("FloorMod(%v, %v) = %v", tc.x, tc.y, m)
		}
	}
	if _, _, err := FloorDivMod(1, 0); !errors.Is(err, ErrDegenerate) {
		t.Errorf("y = 0: %v", err)
	}
	for _, tc := range []struct {
		x    float64
		d    int
		want float64
	}{{2.5, 0, 2}, {3.5, 0, 4}, {-2.5, 0, -2}, {1234.5678, 2, 1234.57}, {1234.5678, -2, 1200}, {1e300, 20, 1e300}} {
		got, err := RoundDecimals(tc.x, tc.d)
		if err != nil || got != tc.want {
			t.Errorf("RoundDecimals(%v, %d) = (%v, %v), want %v", tc.x, tc.d, got, err, tc.want)
		}
	}
	if _, err := RoundDecimals(math.NaN(), 0); !errors.Is(err, ErrNotFinite) {
		t.Errorf("NaN: %v", err)
	}
	for _, tc := range []struct{ x, want float64 }{{-2.5, -1}, {0, 0}, {3, 1}, {math.Inf(-1), -1}} {
		if got, err := Sign(tc.x); err != nil || got != tc.want {
			t.Errorf("Sign(%v) = (%v, %v)", tc.x, got, err)
		}
	}
	if _, err := Sign(math.NaN()); !errors.Is(err, ErrNotFinite) {
		t.Errorf("Sign NaN: %v", err)
	}
	for _, tc := range []struct{ x, want float64 }{{-1, 0}, {2, 1}, {0, 0.5}} {
		if got, err := Heaviside(tc.x, 0.5); err != nil || got != tc.want {
			t.Errorf("Heaviside(%v) = (%v, %v)", tc.x, got, err)
		}
	}
	if _, err := Heaviside(math.NaN(), 0.5); !errors.Is(err, ErrNotFinite) {
		t.Errorf("Heaviside NaN: %v", err)
	}
}

// Oracle: x = y·q + r and |r| < |y| with r carrying y's sign.
func TestFloorDivMod_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 500 {
		x, y := rng.NormFloat64()*100, rng.NormFloat64()*10
		q, r, err := FloorDivMod(x, y)
		if err != nil || math.Abs(y*q+r-x) > 1e-9*math.Abs(x)+1e-12 || math.Abs(r) >= math.Abs(y) || (r != 0 && (r < 0) != (y < 0)) {
			t.Fatalf("FloorDivMod(%v, %v) = (%v, %v, %v)", x, y, q, r, err)
		}
	}
}

func BenchmarkPolynomialAdd(b *testing.B) {
	p := make([]float64, 1000)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = PolynomialAdd(p, p)
	}
}

func BenchmarkPolynomialSubtract(b *testing.B) {
	p := make([]float64, 1000)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = PolynomialSubtract(p, p)
	}
}

func BenchmarkPolynomialDivide(b *testing.B) {
	a, d := make([]float64, 1000), []float64{1, 2, 3, 1}
	b.ReportAllocs()
	for b.Loop() {
		_, _, _ = PolynomialDivide(a, d)
	}
}

func BenchmarkPolynomialPower(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = PolynomialPower([]float64{1, 1}, 64)
	}
}

func BenchmarkPolynomialDerivative(b *testing.B) {
	p := make([]float64, 1000)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = PolynomialDerivative(p, 2)
	}
}

func BenchmarkPolynomialIntegral(b *testing.B) {
	p := make([]float64, 1000)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = PolynomialIntegral(p, 0)
	}
}

func BenchmarkPolynomialFromRoots(b *testing.B) {
	r := make([]complex128, 100)
	for i := range r {
		r[i] = complex(float64(i)/100, 0)
	}
	b.ReportAllocs()
	for b.Loop() {
		_, _ = PolynomialFromRoots(r, -1)
	}
}

func BenchmarkPolynomialValueFromRoots(b *testing.B) {
	r := make([]float64, 1000)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = PolynomialValueFromRoots(r, 0.5)
	}
}

func BenchmarkPolynomialCompanion(b *testing.B) {
	p := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1}
	dst := make([]float64, 100)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = PolynomialCompanion(dst, p)
	}
}

func BenchmarkPolynomialTrim(b *testing.B) {
	p := make([]float64, 1000)
	p[0] = 1
	b.ReportAllocs()
	for b.Loop() {
		_, _ = PolynomialTrim(p, 0)
	}
}

func BenchmarkPolynomialEvaluate2D(b *testing.B) {
	c := make([]float64, 100)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = PolynomialEvaluate2D(c, 10, 10, 0.5, 0.25)
	}
}

func BenchmarkRoundDecimals(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = RoundDecimals(1234.5678, 2)
	}
}

func BenchmarkFloorDivMod(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _, _ = FloorDivMod(-7.5, 2)
	}
}

func BenchmarkFloorMod(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = FloorMod(-7.5, 2)
	}
}

func BenchmarkSign(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = Sign(-2.5)
	}
}

func BenchmarkHeaviside(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = Heaviside(-2.5, 0.5)
	}
}
