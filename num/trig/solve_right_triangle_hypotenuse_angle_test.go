package trig

import (
	"math"
	"testing"
)

// Worked example, hand-derived from DLMF 4.42.1–4.42.2: hypotenuse 2, A = π/6 gives a = 1, b = √3.
func TestSolveRightTriangleHypotenuseAngle_example(t *testing.T) {
	want := Triangle{1, math.Sqrt(3), 2, math.Pi / 6, math.Pi / 3, math.Pi / 2}
	if got, err := SolveRightTriangleHypotenuseAngle(2, math.Pi/6); err != nil || !sameTriangle(got, want, ulps4) {
		t.Errorf("SolveRightTriangleHypotenuseAngle(2, π/6) = (%+v, %v)", got, err)
	}
}

func TestSolveRightTriangleHypotenuseAngle_cases(t *testing.T) {
	if got, err := SolveRightTriangleHypotenuseAngle(5, rightTriangle345.AngleARad); err != nil || !sameTriangle(got, rightTriangle345, ulps4) {
		t.Errorf("3-4-5: (%+v, %v)", got, err)
	}
	if got, err := SolveRightTriangleHypotenuseAngle(1, 1e-300); err != nil || !relNear(got.SideA, 1e-300, ulps4) || got.SideB != 1 {
		t.Errorf("tiny angle: (%+v, %v)", got, err)
	}
	if got, err := SolveRightTriangleHypotenuseAngle(math.MaxFloat64, math.Pi/4); err != nil || !relNear(got.SideA, math.MaxFloat64*math.Sqrt2/2, ulps4) {
		t.Errorf("largest hypotenuse: (%+v, %v)", got, err)
	}
}

func TestSolveRightTriangleHypotenuseAngle_errors(t *testing.T) {
	call := func(c, angleA float64) func() error {
		return func() error { _, err := SolveRightTriangleHypotenuseAngle(c, angleA); return err }
	}
	runErrCases(t, []errCase{
		{"c NaN", call(math.NaN(), 1), ErrNotFinite},
		{"A +Inf", call(1, math.Inf(1)), ErrNotFinite},
		{"c negative", call(-1, 1), ErrNonPositive},
		{"A zero", call(1, 0), ErrDomain},
		{"A right", call(1, math.Pi/2), ErrDomain},
		{"A obtuse", call(1, 2), ErrDomain},
	})
}

// Oracle: a² + b² = c² in big.Float (1e-15) and a/c = sin A (DLMF 4.42.1, 1e-15).
func TestSolveRightTriangleHypotenuseAngle_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		c := math.Pow(10, rng.Float64()*20-10)
		angleA := (0.001 + rng.Float64()*0.998) * math.Pi / 2
		got, err := SolveRightTriangleHypotenuseAngle(c, angleA)
		if err != nil || !pythagorasHolds(got, 1e-15) || !relNear(got.SideA/c, math.Sin(angleA), 1e-15) {
			t.Fatalf("SolveRightTriangleHypotenuseAngle(%v, %v) = (%+v, %v)", c, angleA, got, err)
		}
	}
}

func TestSolveRightTriangleHypotenuseAngle_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkTriangle, _ = SolveRightTriangleHypotenuseAngle(5, 0.6) }); n != 0 {
		t.Errorf("SolveRightTriangleHypotenuseAngle allocates %v times, want 0", n)
	}
}

func BenchmarkSolveRightTriangleHypotenuseAngle(b *testing.B) {
	b.ReportAllocs()
	c, angleA := 5.0, 0.6
	for b.Loop() {
		sinkTriangle, _ = SolveRightTriangleHypotenuseAngle(c, angleA)
	}
}
