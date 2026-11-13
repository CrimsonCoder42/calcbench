package series

import (
	"errors"
	"math"
	"testing"
)

// seqClose: relative 1e-13; each function is a handful of roundings (pow/expm1/log within 1 ulp each).
func seqClose(got, want float64) bool {
	return math.Abs(got-want) <= 1e-13*math.Max(1, math.Abs(want))
}

// Hand computations from the formulas of Abramson, College Algebra 2e, OpenStax, 2021, §9.2–§9.4:
// a₁ = 2, d = 5: a₂₀ = 97, S₂₀ = 20(2 + 97)/2 = 990; a₁ = 2, r = 5: a₁₂ = 2·5¹¹ = 97 656 250,
// S₁₂ = 2(5¹² − 1)/4 = 122 070 312; 1 + 1/2 + … + 1/2⁹ = 2 − 2⁻⁹; r = −1 alternates.
func TestSequences_handExamples(t *testing.T) {
	check := func(name string, got float64, err error, want float64) {
		t.Helper()
		if err != nil || !seqClose(got, want) {
			t.Errorf("%s = %v, %v; want %v", name, got, err, want)
		}
	}
	v, err := ArithmeticSequenceTerm(2, 5, 20)
	check("ArithTerm(2,5,20)", v, err, 97)
	v, err = ArithmeticSequenceTerm(10, -3, 1)
	check("ArithTerm(10,-3,1)", v, err, 10)
	v, err = ArithmeticSeriesSum(2, 5, 20)
	check("ArithSum(2,5,20)", v, err, 990)
	v, err = ArithmeticSeriesSum(1, 1, 100)
	check("ArithSum(1,1,100)", v, err, 5050)
	v, err = ArithmeticSeriesSum(-4, 0, 3)
	check("ArithSum(-4,0,3)", v, err, -12)
	v, err = GeometricSequenceTerm(2, 5, 12)
	check("GeoTerm(2,5,12)", v, err, 97656250)
	v, err = GeometricSequenceTerm(3, -2, 4)
	check("GeoTerm(3,-2,4)", v, err, -24)
	v, err = GeometricSequenceTerm(0, math.MaxFloat64, 50)
	check("GeoTerm(0,…)", v, err, 0)
	v, err = GeometricSeriesSum(2, 5, 12)
	check("GeoSum(2,5,12)", v, err, 122070312)
	v, err = GeometricSeriesSum(1, 0.5, 10)
	check("GeoSum(1,0.5,10)", v, err, 2-math.Ldexp(1, -9))
	v, err = GeometricSeriesSum(7, 1, 4)
	check("GeoSum(7,1,4)", v, err, 28)
	v, err = GeometricSeriesSum(1, -1, 5)
	check("GeoSum(1,-1,5)", v, err, 1)
	v, err = GeometricSeriesSum(1, 0, 5)
	check("GeoSum(1,0,5)", v, err, 1)
	// r = 1 + 1e-12, n = 1000: S = n + n(n−1)/2·ε + … = 1000 + 4.995e-7 (+ O(ε²)); the naive formula
	// loses about 4 digits here.
	eps := (1 + 1e-12) - 1 // the ratio as stored
	v, err = GeometricSeriesSum(1, 1+eps, 1000)
	check("GeoSum near 1", v, err, 1000+499500*eps+166167000*eps*eps)
}

func TestSequences_errors(t *testing.T) {
	nan := math.NaN()
	e := func(_ float64, err error) error { return err }
	tests := []struct {
		name      string
		err, want error
	}{
		{"ArithTerm NaN", e(ArithmeticSequenceTerm(nan, 1, 1)), ErrNotFinite},
		{"ArithTerm n 0", e(ArithmeticSequenceTerm(1, 1, 0)), ErrOutOfRange},
		{"ArithTerm overflow", e(ArithmeticSequenceTerm(1, math.MaxFloat64, 3)), ErrOverflow},
		{"ArithSum Inf", e(ArithmeticSeriesSum(1, math.Inf(1), 2)), ErrNotFinite},
		{"ArithSum n -1", e(ArithmeticSeriesSum(1, 1, -1)), ErrOutOfRange},
		{"ArithSum overflow", e(ArithmeticSeriesSum(math.MaxFloat64, 0, 2)), ErrOverflow},
		{"GeoTerm NaN", e(GeometricSequenceTerm(1, nan, 2)), ErrNotFinite},
		{"GeoTerm n 0", e(GeometricSequenceTerm(1, 2, 0)), ErrOutOfRange},
		{"GeoTerm overflow", e(GeometricSequenceTerm(1, 10, 400)), ErrOverflow},
		{"GeoSum Inf", e(GeometricSeriesSum(math.Inf(-1), 2, 2)), ErrNotFinite},
		{"GeoSum n 0", e(GeometricSeriesSum(1, 2, 0)), ErrOutOfRange},
		{"GeoSum overflow", e(GeometricSeriesSum(1, 10, 400)), ErrOverflow},
	}
	for _, tc := range tests {
		if !errors.Is(tc.err, tc.want) {
			t.Errorf("%s: err = %v, want %v", tc.name, tc.err, tc.want)
		}
	}
}

// The closed forms against a direct loop on small integer cases (exact in float64).
func TestSequences_oracle(t *testing.T) {
	for a := -3.0; a <= 3; a++ {
		for d := -4.0; d <= 4; d++ {
			term, sum := a, 0.0
			gterm, gsum := a, 0.0
			for n := 1; n <= 12; n++ {
				sum += term
				gsum += gterm
				if got, _ := ArithmeticSequenceTerm(a, d, n); got != term {
					t.Fatalf("ArithTerm(%v,%v,%d) = %v, want %v", a, d, n, got, term)
				}
				if got, _ := ArithmeticSeriesSum(a, d, n); got != sum {
					t.Fatalf("ArithSum(%v,%v,%d) = %v, want %v", a, d, n, got, sum)
				}
				if got, _ := GeometricSequenceTerm(a, d, n); got != gterm {
					t.Fatalf("GeoTerm(%v,%v,%d) = %v, want %v", a, d, n, got, gterm)
				}
				if got, _ := GeometricSeriesSum(a, d, n); !seqClose(got, gsum) {
					t.Fatalf("GeoSum(%v,%v,%d) = %v, want %v", a, d, n, got, gsum)
				}
				term += d
				gterm *= d
			}
		}
	}
}

func TestSequences_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() {
		_, _ = ArithmeticSeriesSum(2, 5, 20)
		_, _ = GeometricSeriesSum(2, 5, 12)
	}); a != 0 {
		t.Errorf("allocates %v", a)
	}
}

func BenchmarkArithmeticSequenceTerm(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = ArithmeticSequenceTerm(2, 5, 20)
	}
}

func BenchmarkArithmeticSeriesSum(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = ArithmeticSeriesSum(2, 5, 20)
	}
}

func BenchmarkGeometricSequenceTerm(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = GeometricSequenceTerm(2, 5, 12)
	}
}

func BenchmarkGeometricSeriesSum(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = GeometricSeriesSum(2, 1.0001, 1000)
	}
}
