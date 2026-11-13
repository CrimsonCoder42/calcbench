package trig

import (
	"math"
	"testing"
)

// Worked example, hand-derived from DLMF 4.42.1: hypotenuse 5 and leg 3 give the other leg 4.
func TestSolveRightTriangleHypotenuseLeg_example(t *testing.T) {
	if got, err := SolveRightTriangleHypotenuseLeg(5, 3); err != nil || !sameTriangle(got, rightTriangle345, ulps4) {
		t.Errorf("SolveRightTriangleHypotenuseLeg(5, 3) = (%+v, %v)", got, err)
	}
}

func TestSolveRightTriangleHypotenuseLeg_cases(t *testing.T) {
	tests := []struct {
		name  string
		c, a  float64
		wantB float64
	}{
		// a ≈ c: (c−a)(c+a) keeps the digits that c² − a² would cancel. c = 1, a = 1 − 2^-52:
		// b = √(2^-52·(2 − 2^-52)) exactly representable to rounding.
		{"leg almost the hypotenuse", 1, 1 - 0x1p-52, math.Sqrt(0x1p-52 * (2 - 0x1p-52))},
		{"huge, would overflow c²", 1.5e308, 1e308, math.Sqrt(1.25) * 1e308},
		{"subnormal", 5e-310, 3e-310, 4e-310},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := SolveRightTriangleHypotenuseLeg(tc.c, tc.a)
			if err != nil || !relNear(got.SideB, tc.wantB, ulps4) || got.AngleCRad != math.Pi/2 || got.SideC != tc.c || got.SideA != tc.a {
				t.Errorf("got (%+v, %v), want b = %v", got, err, tc.wantB)
			}
		})
	}
}

func TestSolveRightTriangleHypotenuseLeg_errors(t *testing.T) {
	call := func(c, a float64) func() error {
		return func() error { _, err := SolveRightTriangleHypotenuseLeg(c, a); return err }
	}
	runErrCases(t, []errCase{
		{"c NaN", call(math.NaN(), 1), ErrNotFinite},
		{"a -Inf", call(1, math.Inf(-1)), ErrNotFinite},
		{"c zero", call(0, 1), ErrNonPositive},
		{"a negative", call(1, -1), ErrNonPositive},
		{"leg equals hypotenuse", call(2, 2), ErrNotTriangle},
		{"leg longer than hypotenuse", call(2, 3), ErrNotTriangle},
	})
}

// Oracle: a² + b² = c² in big.Float, to 1e-14 relative (b has a few roundings), and A + B = π/2.
func TestSolveRightTriangleHypotenuseLeg_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		c := math.Pow(10, rng.Float64()*20-10)
		a := c * rng.Float64()
		if a == 0 {
			continue
		}
		got, err := SolveRightTriangleHypotenuseLeg(c, a)
		if err != nil || !pythagorasHolds(got, 1e-14) || !near(got.AngleARad+got.AngleBRad, math.Pi/2, 1e-15) {
			t.Fatalf("SolveRightTriangleHypotenuseLeg(%v, %v) = (%+v, %v)", c, a, got, err)
		}
	}
}

func TestSolveRightTriangleHypotenuseLeg_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkTriangle, _ = SolveRightTriangleHypotenuseLeg(5, 3) }); n != 0 {
		t.Errorf("SolveRightTriangleHypotenuseLeg allocates %v times, want 0", n)
	}
}

func BenchmarkSolveRightTriangleHypotenuseLeg(b *testing.B) {
	b.ReportAllocs()
	c, a := 5.0, 3.0
	for b.Loop() {
		sinkTriangle, _ = SolveRightTriangleHypotenuseLeg(c, a)
	}
}
