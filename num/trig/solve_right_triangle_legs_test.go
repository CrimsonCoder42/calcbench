package trig

import (
	"math"
	"math/big"
	"testing"
)

// rightTriangle345 is the 3-4-5 right triangle: A = atan(3/4) = 0.6435011087932844 (36.8698976458°).
var rightTriangle345 = Triangle{3, 4, 5, 0.6435011087932844, math.Pi/2 - 0.6435011087932844, math.Pi / 2}

// pythagorasHolds reports whether a² + b² = c² to tol relative, in 256-bit big.Float.
func pythagorasHolds(t Triangle, tol float64) bool {
	sq := func(x float64) *big.Float { f := new(big.Float).SetPrec(256).SetFloat64(x); return f.Mul(f, f) }
	lhs, _ := new(big.Float).SetPrec(256).Add(sq(t.SideA), sq(t.SideB)).Float64()
	rhs, _ := sq(t.SideC).Float64()
	return relNear(lhs, rhs, tol)
}

// Worked example, hand-derived from DLMF 4.42.1–4.42.3: legs (3, 4) give c = 5, A = 0.6435011087932844.
func TestSolveRightTriangleLegs_example(t *testing.T) {
	if got, err := SolveRightTriangleLegs(3, 4); err != nil || !sameTriangle(got, rightTriangle345, ulps4) {
		t.Errorf("SolveRightTriangleLegs(3, 4) = (%+v, %v)", got, err)
	}
}

func TestSolveRightTriangleLegs_cases(t *testing.T) {
	tests := []struct {
		name string
		a, b float64
		want Triangle
	}{
		{"isosceles", 1, 1, Triangle{1, 1, math.Sqrt2, math.Pi / 4, math.Pi / 4, math.Pi / 2}},
		{"huge, squares would overflow", 3e200, 4e200, Triangle{3e200, 4e200, 5e200, rightTriangle345.AngleARad, rightTriangle345.AngleBRad, math.Pi / 2}},
		{"tiny, squares would underflow", 3e-200, 4e-200, Triangle{3e-200, 4e-200, 5e-200, rightTriangle345.AngleARad, rightTriangle345.AngleBRad, math.Pi / 2}},
		{"very thin", 1, 1e-300, Triangle{1, 1e-300, 1, math.Pi / 2, 1e-300, math.Pi / 2}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := SolveRightTriangleLegs(tc.a, tc.b); err != nil || !sameTriangle(got, tc.want, ulps4) {
				t.Errorf("got (%+v, %v), want %+v", got, err, tc.want)
			}
		})
	}
}

func TestSolveRightTriangleLegs_errors(t *testing.T) {
	call := func(a, b float64) func() error {
		return func() error { _, err := SolveRightTriangleLegs(a, b); return err }
	}
	runErrCases(t, []errCase{
		{"a NaN", call(math.NaN(), 1), ErrNotFinite},
		{"b +Inf", call(1, math.Inf(1)), ErrNotFinite},
		{"a zero", call(0, 1), ErrNonPositive},
		{"b negative", call(1, -1), ErrNonPositive},
		{"hypotenuse overflows", call(1.5e308, 1.5e308), ErrOverflow},
	})
}

// Oracle: a² + b² = c² in big.Float (one rounding in hypot: 1e-15 relative) and A + B = π/2 (1e-15).
func TestSolveRightTriangleLegs_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		a, b := math.Pow(10, rng.Float64()*20-10), math.Pow(10, rng.Float64()*20-10)
		got, err := SolveRightTriangleLegs(a, b)
		if err != nil || !pythagorasHolds(got, 1e-15) || !near(got.AngleARad+got.AngleBRad, math.Pi/2, 1e-15) {
			t.Fatalf("SolveRightTriangleLegs(%v, %v) = (%+v, %v)", a, b, got, err)
		}
	}
}

func TestSolveRightTriangleLegs_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkTriangle, _ = SolveRightTriangleLegs(3, 4) }); n != 0 {
		t.Errorf("SolveRightTriangleLegs allocates %v times, want 0", n)
	}
}

func BenchmarkSolveRightTriangleLegs(b *testing.B) {
	b.ReportAllocs()
	la, lb := 3.0, 4.0
	for b.Loop() {
		sinkTriangle, _ = SolveRightTriangleLegs(la, lb)
	}
}
