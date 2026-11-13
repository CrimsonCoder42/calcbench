package geometry

import (
	"errors"
	"math"
	"math/big"
	"testing"
)

func TestProduct_noSpuriousOverflow(t *testing.T) {
	tests := []struct {
		name       string
		x, y, z, w float64
		want       float64
	}{
		{"large then small", 1e200, 1e200, 1e-300, 1, 1e100},
		{"small then large", 1e-300, 1e-300, 1e300, 1, 1e-300},
		{"zero", 0, 1e300, 1e300, 1e300, 0},
		{"true overflow", 1e300, 1e300, 1, 1, math.Inf(1)},
	}
	for _, tt := range tests {
		got := product(tt.x, tt.y, tt.z, tt.w)
		if !(got == tt.want || relNear(got, tt.want, tolFew)) {
			t.Errorf("%s: product = %v, want %v", tt.name, got, tt.want)
		}
	}
}

// Worked example: CRC 31st ed. §4.5.1 (area = ½h_c·c), hand-derived: 6, 4 → 12.
func TestTriangleAreaBaseHeight_example(t *testing.T) {
	runValCases(t, []valCase{
		{"TriangleAreaBaseHeight(6, 4)", func() (float64, error) { return TriangleAreaBaseHeight(6, 4) }, 12, 0},
	})
}

func TestTriangleAreaBaseHeight_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"no spurious overflow", func() (float64, error) { return TriangleAreaBaseHeight(1e300, 1e-300) }, 0.5, tolFew},
		{"tiny", func() (float64, error) { return TriangleAreaBaseHeight(1e-160, 1e-160) }, 5e-321, 1e-3},
	})
}

func TestTriangleAreaBaseHeight_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { return e1(TriangleAreaBaseHeight(nan, 1)) }, ErrNotFinite},
		{"Inf", func() error { return e1(TriangleAreaBaseHeight(1, inf)) }, ErrNotFinite},
		{"zero", func() error { return e1(TriangleAreaBaseHeight(0, 1)) }, ErrNonPositive},
		{"negative", func() error { return e1(TriangleAreaBaseHeight(1, -2)) }, ErrNonPositive},
		{"overflow", func() error { return e1(TriangleAreaBaseHeight(1e200, 1e200)) }, ErrOverflow},
	})
}

// Worked example: Kahan 2014, §5 Table 1 ("Accurate ∆" column, p. 4). The table's data are decimal; in
// binary they are rounded, and that rounding (not the formula) sets each tolerance: it moves the small
// factor c − (a − b) by up to one ulp of the largest side, relative to the factor's size. The HP-15C values
// are correct to at least nine digits (Kahan, same page). Row 7 (5278.64055, …) is omitted: its decimal
// data make c − (a − b) exactly 0, which binary rounding of the inputs turns into ±1e-11.
func TestTriangleAreaHeron_kahanTable1(t *testing.T) {
	runValCases(t, []valCase{
		{"(10, 10, 10)", func() (float64, error) { return TriangleAreaHeron(10, 10, 10) }, 43.30127020, 5e-9},
		{"(100000, 99999.99979, 0.00029)", func() (float64, error) { return TriangleAreaHeron(100000, 99999.99979, 0.00029) }, 9.999999990, 2e-7},
		{"(100000, 100000, 1.00005)", func() (float64, error) { return TriangleAreaHeron(100000, 100000, 1.00005) }, 50002.50003, 5e-9},
		{"(99999.99996, 99999.99994, 0.00003)", func() (float64, error) { return TriangleAreaHeron(99999.99996, 99999.99994, 0.00003) }, 1.118033988, 2e-6},
		{"(99999.99996, 0.00003, 99999.99994)", func() (float64, error) { return TriangleAreaHeron(99999.99996, 0.00003, 99999.99994) }, 1.118033988, 2e-6},
		{"(10000, 5000.000001, 15000)", func() (float64, error) { return TriangleAreaHeron(10000, 5000.000001, 15000) }, 612.3724358, 1e-6},
		{"(100002, 100002, 200004)", func() (float64, error) { return TriangleAreaHeron(100002, 100002, 200004) }, 0, 0},
		{"(31622.77662, 0.000023, 31622.77661)", func() (float64, error) { return TriangleAreaHeron(31622.77662, 0.000023, 31622.77661) }, 0.327490458, 2e-6},
		{"(31622.77662, 0.0155555, 31622.77661)", func() (float64, error) { return TriangleAreaHeron(31622.77662, 0.0155555, 31622.77661) }, 245.9540000, 1e-8},
	})
	if _, err := TriangleAreaHeron(99999.99999, 99999.99999, 200000); !errors.Is(err, ErrNotTriangle) {
		t.Errorf("(99999.99999, 99999.99999, 200000): err = %v, want ErrNotTriangle (Table 1: Error)", err)
	}
}

func TestTriangleAreaHeron_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"3-4-5", func() (float64, error) { return TriangleAreaHeron(3, 4, 5) }, 6, tolFew},
		{"13-14-15", func() (float64, error) { return TriangleAreaHeron(13, 14, 15) }, 84, tolFew},
		{"equilateral, side 2", func() (float64, error) { return TriangleAreaHeron(2, 2, 2) }, math.Sqrt(3), tolFew},
		{"huge sides, no overflow of the factors", func() (float64, error) { return TriangleAreaHeron(3e150, 4e150, 5e150) }, 6e300, tolSome},
		{"needle with tiny factors, no underflow", func() (float64, error) { return TriangleAreaHeron(1, 1, 1e-200) }, 5e-201, tolFew},
	})
}

func TestTriangleAreaHeron_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN", func() error { return e1(TriangleAreaHeron(1, nan, 1)) }, ErrNotFinite},
		{"Inf", func() error { return e1(TriangleAreaHeron(inf, 1, 1)) }, ErrNotFinite},
		{"zero", func() error { return e1(TriangleAreaHeron(1, 1, 0)) }, ErrNonPositive},
		{"negative (Table 1 row 2)", func() error { return e1(TriangleAreaHeron(-3, 4, 2)) }, ErrNonPositive},
		{"not a triangle", func() error { return e1(TriangleAreaHeron(1, 1, 3)) }, ErrNotTriangle},
		{"overflow", func() error { return e1(TriangleAreaHeron(1e308, 1e308, 1e308)) }, ErrOverflow},
	})
}

// heronExact returns ¼√((a+b+c)(−a+b+c)(a−b+c)(a+b−c)) from exact big.Float arithmetic (2048 bits hold
// every sum and product of these float64 sides exactly), rounded once; NaN when the product is negative.
func heronExact(a, b, c float64) float64 {
	const prec = 2048
	nf := func(x float64) *big.Float { return new(big.Float).SetPrec(prec).SetFloat64(x) }
	add := func(x, y *big.Float) *big.Float { return new(big.Float).SetPrec(prec).Add(x, y) }
	sub := func(x, y *big.Float) *big.Float { return new(big.Float).SetPrec(prec).Sub(x, y) }
	mul := func(x, y *big.Float) *big.Float { return new(big.Float).SetPrec(prec).Mul(x, y) }
	A, B, C := nf(a), nf(b), nf(c)
	p := mul(mul(add(add(A, B), C), sub(add(B, C), A)), mul(sub(add(A, C), B), sub(add(A, B), C)))
	if p.Sign() < 0 {
		return math.NaN()
	}
	r := new(big.Float).SetPrec(prec).Sqrt(p)
	v, _ := r.Float64()
	return v / 4
}

// Oracle: exact Heron in big.Float on random triangles, needles included. Kahan bounds each of the four
// factors to an ulp or two; with the product, the two square roots and the scaling, 18 ulps (4e-15) is a
// safe ceiling.
func TestTriangleAreaHeron_oracle(t *testing.T) {
	rng := newRNG()
	for range 600 {
		a := 1 + rng.Float64()*999
		b := 1 + rng.Float64()*999
		u := math.Pow(rng.Float64(), 1+rng.Float64()*20) // clustered near 0: needles
		if rng.IntN(2) == 0 {
			u = 1 - u // clustered near 1: flat triangles
		}
		c := math.Abs(a-b) + u*(a+b-math.Abs(a-b))
		if c <= 0 {
			continue
		}
		want := heronExact(a, b, c)
		got, err := TriangleAreaHeron(a, b, c)
		if math.IsNaN(want) {
			if !errors.Is(err, ErrNotTriangle) {
				t.Fatalf("Heron(%v, %v, %v) = (%v, %v), want ErrNotTriangle", a, b, c, got, err)
			}
			continue
		}
		if err != nil || !relNear(got, want, 4e-15) {
			t.Fatalf("Heron(%v, %v, %v) = (%v, %v), want %v", a, b, c, got, err, want)
		}
	}
}

// Oracle: base × height for right triangles, k·(3, 4, 5) with area 6k².
func TestTriangleAreaHeron_vsBaseHeight(t *testing.T) {
	rng := newRNG()
	for range 200 {
		k := math.Ldexp(1, rng.IntN(200)-100) // powers of two keep k·(3, 4, 5) exact
		want, _ := TriangleAreaBaseHeight(3*k, 4*k)
		if got, err := TriangleAreaHeron(3*k, 4*k, 5*k); err != nil || !relNear(got, want, 4e-15) {
			t.Fatalf("Heron(3k, 4k, 5k), k = %v: (%v, %v), want %v", k, got, err, want)
		}
	}
}

// Worked example: CRC 31st ed. §4.5.1 (area = ½ab sin C), hand-derived: 3, 4, π/6 → 3.
func TestTriangleAreaSAS_example(t *testing.T) {
	runValCases(t, []valCase{
		{"TriangleAreaSAS(3, 4, π/6)", func() (float64, error) { return TriangleAreaSAS(3, 4, math.Pi/6) }, 3, tolFew},
	})
}

func TestTriangleAreaSAS_cases(t *testing.T) {
	runValCases(t, []valCase{
		{"right angle", func() (float64, error) { return TriangleAreaSAS(3, 4, math.Pi/2) }, 6, 0},
		{"no spurious overflow", func() (float64, error) { return TriangleAreaSAS(1e300, 1e-300, math.Pi/2) }, 0.5, tolFew},
	})
}

func TestTriangleAreaSAS_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"NaN side", func() error { return e1(TriangleAreaSAS(nan, 1, 1)) }, ErrNotFinite},
		{"Inf angle", func() error { return e1(TriangleAreaSAS(1, 1, inf)) }, ErrNotFinite},
		{"zero side", func() error { return e1(TriangleAreaSAS(1, 0, 1)) }, ErrNonPositive},
		{"zero angle", func() error { return e1(TriangleAreaSAS(1, 1, 0)) }, ErrOutOfRange},
		{"straight angle", func() error { return e1(TriangleAreaSAS(1, 1, math.Pi)) }, ErrOutOfRange},
		{"overflow", func() error { return e1(TriangleAreaSAS(1e200, 1e200, 1)) }, ErrOverflow},
	})
}

// Oracle: Heron on the third side from the law of cosines (angles kept away from 0 and π, so the
// textbook law of cosines does not cancel).
func TestTriangleAreaSAS_vsHeron(t *testing.T) {
	rng := newRNG()
	for range 400 {
		a, b := 1+rng.Float64()*99, 1+rng.Float64()*99
		angle := 0.3 + rng.Float64()*2.5
		c := math.Sqrt(a*a + b*b - 2*a*b*math.Cos(angle))
		want, err := TriangleAreaHeron(a, b, c)
		if err != nil {
			t.Fatal(err)
		}
		if got, err := TriangleAreaSAS(a, b, angle); err != nil || !relNear(got, want, 1e-12) {
			t.Fatalf("SAS(%v, %v, %v) = (%v, %v), Heron %v", a, b, angle, got, err, want)
		}
	}
}

func TestTriangleArea_allocs(t *testing.T) {
	assertNoAllocs(t, "TriangleAreaBaseHeight", func() { sink, _ = TriangleAreaBaseHeight(6, 4) })
	assertNoAllocs(t, "TriangleAreaHeron", func() { sink, _ = TriangleAreaHeron(13, 14, 15) })
	assertNoAllocs(t, "TriangleAreaSAS", func() { sink, _ = TriangleAreaSAS(3, 4, 0.5) })
}

func BenchmarkTriangleAreaBaseHeight(b *testing.B) {
	b.ReportAllocs()
	base, height := 6.0, 4.0
	for b.Loop() {
		sink, _ = TriangleAreaBaseHeight(base, height)
	}
}

func BenchmarkTriangleAreaHeron(b *testing.B) {
	b.ReportAllocs()
	x, y, z := 13.0, 14.0, 15.0
	for b.Loop() {
		sink, _ = TriangleAreaHeron(x, y, z)
	}
}

func BenchmarkTriangleAreaSAS(b *testing.B) {
	b.ReportAllocs()
	x, y, angle := 3.0, 4.0, 0.5
	for b.Loop() {
		sink, _ = TriangleAreaSAS(x, y, angle)
	}
}
