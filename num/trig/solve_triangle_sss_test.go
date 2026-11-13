package trig

import (
	"math"
	"math/rand/v2"
	"testing"
)

// randomTriangle returns a triangle built from two random angles in [0.01, π − 0.02] and a random scale
// k in [1e-10, 1e10]: sides k·sin of the opposite angle (DLMF 4.42.4). Shared by the solver oracles.
func randomTriangle(rng *rand.Rand) Triangle {
	angleA := 0.01 + rng.Float64()*3
	angleB := 0.01 + rng.Float64()*(math.Pi-angleA-0.02)
	angleC := math.Pi - angleA - angleB
	k := math.Pow(10, rng.Float64()*20-10)
	return Triangle{k * math.Sin(angleA), k * math.Sin(angleB), k * math.Sin(angleC), angleA, angleB, angleC}
}

// sameTriangle reports whether every part of got is within tol (relative) of want.
func sameTriangle(got, want Triangle, tol float64) bool {
	return relNear(got.SideA, want.SideA, tol) && relNear(got.SideB, want.SideB, tol) && relNear(got.SideC, want.SideC, tol) &&
		near(got.AngleARad, want.AngleARad, tol) && near(got.AngleBRad, want.AngleBRad, tol) && near(got.AngleCRad, want.AngleCRad, tol)
}

// Worked example: Kahan 2014 Table 1, row 1, (10, 10, 10): every angle 60°. Hand-derived: (3, 4, 5) has
// C = π/2 and A = atan(3/4) = 0.6435011087932844.
func TestSolveTriangleSSS_example(t *testing.T) {
	eq, err := SolveTriangleSSS(10, 10, 10)
	if err != nil || !sameTriangle(eq, Triangle{10, 10, 10, math.Pi / 3, math.Pi / 3, math.Pi / 3}, ulps4) {
		t.Errorf("SolveTriangleSSS(10, 10, 10) = (%+v, %v)", eq, err)
	}
	r, err := SolveTriangleSSS(3, 4, 5)
	if err != nil || !sameTriangle(r, Triangle{3, 4, 5, 0.6435011087932844, math.Pi/2 - 0.6435011087932844, math.Pi / 2}, ulps4) {
		t.Errorf("SolveTriangleSSS(3, 4, 5) = (%+v, %v)", r, err)
	}
}

func TestSolveTriangleSSS_cases(t *testing.T) {
	tests := []struct {
		name    string
		a, b, c float64
		want    Triangle
	}{
		{"flat, angles 0, 0, π (Kahan 2014 Table 1 rows 9–10)", 1, 2, 3, Triangle{1, 2, 3, 0, 0, math.Pi}},
		{"huge equilateral", 1e308, 1e308, 1e308, Triangle{1e308, 1e308, 1e308, math.Pi / 3, math.Pi / 3, math.Pi / 3}},
		{"subnormal equilateral", 1e-320, 1e-320, 1e-320, Triangle{1e-320, 1e-320, 1e-320, math.Pi / 3, math.Pi / 3, math.Pi / 3}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := SolveTriangleSSS(tc.a, tc.b, tc.c); err != nil || !sameTriangle(got, tc.want, ulps4) {
				t.Errorf("got (%+v, %v), want %+v", got, err, tc.want)
			}
		})
	}
}

func TestSolveTriangleSSS_errors(t *testing.T) {
	call := func(a, b, c float64) func() error {
		return func() error { _, err := SolveTriangleSSS(a, b, c); return err }
	}
	runErrCases(t, []errCase{
		{"a NaN", call(math.NaN(), 1, 1), ErrNotFinite},
		{"c +Inf", call(1, 1, math.Inf(1)), ErrNotFinite},
		{"b zero", call(1, 0, 1), ErrNonPositive},
		{"a negative (Kahan 2014 Table 1 row 2)", call(-3, 4, 2), ErrNonPositive},
		{"not a triangle", call(1, 1, 3), ErrNotTriangle},
		{"a too long", call(3, 1, 1), ErrNotTriangle},
	})
}

// Oracle: random triangles, each angle against oracleAngle (1e-13 relative), and the angles sum to π
// within a few ulps of π (1e-14 absolute).
func TestSolveTriangleSSS_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		tr := randomTriangle(rng)
		got, err := SolveTriangleSSS(tr.SideA, tr.SideB, tr.SideC)
		if err != nil {
			continue // rounding the sides can make a near-flat triangle invalid; LawOfCosinesAngle's oracle covers that
		}
		wa, _ := oracleAngle(tr.SideB, tr.SideC, tr.SideA)
		wb, _ := oracleAngle(tr.SideA, tr.SideC, tr.SideB)
		wc, _ := oracleAngle(tr.SideA, tr.SideB, tr.SideC)
		if !relNear(got.AngleARad, wa, 1e-13) || !relNear(got.AngleBRad, wb, 1e-13) || !relNear(got.AngleCRad, wc, 1e-13) ||
			math.Abs(got.AngleARad+got.AngleBRad+got.AngleCRad-math.Pi) > 1e-14 {
			t.Fatalf("SolveTriangleSSS(%v, %v, %v) = %+v; oracle %v %v %v", tr.SideA, tr.SideB, tr.SideC, got, wa, wb, wc)
		}
	}
}

func TestSolveTriangleSSS_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkTriangle, _ = SolveTriangleSSS(3, 4, 5) }); n != 0 {
		t.Errorf("SolveTriangleSSS allocates %v times, want 0", n)
	}
}

func BenchmarkSolveTriangleSSS(b *testing.B) {
	b.ReportAllocs()
	sa, sb, sc := 3.0, 4.0, 5.0
	for b.Loop() {
		sinkTriangle, _ = SolveTriangleSSS(sa, sb, sc)
	}
}
