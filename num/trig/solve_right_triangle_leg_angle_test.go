package trig

import (
	"math"
	"testing"
)

// Worked example, hand-derived from DLMF 4.42.1 and 4.42.3: leg 3 opposite A = atan2(3, 4) gives b = 4, c = 5.
func TestSolveRightTriangleLegAngle_example(t *testing.T) {
	if got, err := SolveRightTriangleLegAngle(3, math.Atan2(3, 4)); err != nil || !sameTriangle(got, rightTriangle345, ulps4) {
		t.Errorf("SolveRightTriangleLegAngle(3, atan2(3, 4)) = (%+v, %v)", got, err)
	}
}

func TestSolveRightTriangleLegAngle_cases(t *testing.T) {
	want := Triangle{1, math.Sqrt(3), 2, math.Pi / 6, math.Pi / 3, math.Pi / 2}
	if got, err := SolveRightTriangleLegAngle(1, math.Pi/6); err != nil || !sameTriangle(got, want, ulps4) {
		t.Errorf("30-60-90: (%+v, %v)", got, err)
	}
	// The adjacent-angle form: the caller passes π/2 − (angle at the leg).
	if got, err := SolveRightTriangleLegAngle(4, math.Pi/2-math.Atan2(3, 4)); err != nil || !relNear(got.SideB, 3, ulps4) || !relNear(got.SideC, 5, ulps4) {
		t.Errorf("adjacent angle: (%+v, %v)", got, err)
	}
	if got, err := SolveRightTriangleLegAngle(1e-300, math.Pi/4); err != nil || !relNear(got.SideC, math.Sqrt2*1e-300, ulps4) {
		t.Errorf("tiny leg: (%+v, %v)", got, err)
	}
}

func TestSolveRightTriangleLegAngle_errors(t *testing.T) {
	call := func(a, angleA float64) func() error {
		return func() error { _, err := SolveRightTriangleLegAngle(a, angleA); return err }
	}
	runErrCases(t, []errCase{
		{"a NaN", call(math.NaN(), 1), ErrNotFinite},
		{"A -Inf", call(1, math.Inf(-1)), ErrNotFinite},
		{"a zero", call(0, 1), ErrNonPositive},
		{"A zero", call(1, 0), ErrDomain},
		{"A right", call(1, math.Pi/2), ErrDomain},
		{"sides overflow", call(1e300, 1e-10), ErrOverflow},
	})
}

// Oracle: a² + b² = c² in big.Float (1e-14) and a/c = sin A (1e-15).
func TestSolveRightTriangleLegAngle_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		a := math.Pow(10, rng.Float64()*20-10)
		angleA := (0.001 + rng.Float64()*0.998) * math.Pi / 2
		got, err := SolveRightTriangleLegAngle(a, angleA)
		if err != nil || !pythagorasHolds(got, 1e-14) || !relNear(a/got.SideC, math.Sin(angleA), 1e-15) {
			t.Fatalf("SolveRightTriangleLegAngle(%v, %v) = (%+v, %v)", a, angleA, got, err)
		}
	}
}

func TestSolveRightTriangleLegAngle_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkTriangle, _ = SolveRightTriangleLegAngle(3, 0.6) }); n != 0 {
		t.Errorf("SolveRightTriangleLegAngle allocates %v times, want 0", n)
	}
}

func BenchmarkSolveRightTriangleLegAngle(b *testing.B) {
	b.ReportAllocs()
	a, angleA := 3.0, 0.6
	for b.Loop() {
		sinkTriangle, _ = SolveRightTriangleLegAngle(a, angleA)
	}
}
