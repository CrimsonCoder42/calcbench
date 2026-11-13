package trig

import (
	"math"
	"testing"
)

// Worked example: Kahan 2014 Table 2, row 1: a = 3, b = 17, A = 10° has two triangles, B = 79.73894317°
// or 100.26105683° (tolerance 1e-8°, as in TestLawOfSinesAngle_kahanTable2). Each solution is checked
// to close: its three sides give back its angles (SolveTriangleSSS) to 1e-12.
func TestSolveTriangleSSA_kahanExample(t *testing.T) {
	sol, n, err := SolveTriangleSSA(3, 17, radOf(10))
	if err != nil || n != 2 {
		t.Fatalf("SolveTriangleSSA(3, 17, 10°) = (%+v, %d, %v), want 2 triangles", sol, n, err)
	}
	for i, want := range []float64{79.73894317, 100.26105683} {
		if math.Abs(degOf(sol[i].AngleBRad)-want) > 1e-8 {
			t.Errorf("solution %d: B = %.10f°, want %.8f°", i, degOf(sol[i].AngleBRad), want)
		}
		check, err := SolveTriangleSSS(sol[i].SideA, sol[i].SideB, sol[i].SideC)
		if err != nil || !sameTriangle(check, sol[i], 1e-12) {
			t.Errorf("solution %d does not close: %+v vs %+v", i, sol[i], check)
		}
	}
}

func TestSolveTriangleSSA_cases(t *testing.T) {
	// No triangle: Kahan 2014 Table 2, (0.9999999998, 2, 30°). Slots stay zero.
	sol, n, err := SolveTriangleSSA(0.9999999998, 2, radOf(30))
	if err != nil || n != 0 || sol != [2]Triangle{} {
		t.Errorf("no triangle: (%+v, %d, %v)", sol, n, err)
	}
	// One triangle (a > b): 3-4-5 with a = 5 opposite the right angle, b = 4.
	sol, n, err = SolveTriangleSSA(5, 4, math.Pi/2)
	if err != nil || n != 1 || sol[1] != (Triangle{}) ||
		!sameTriangle(sol[0], Triangle{5, 4, 3, math.Pi / 2, math.Atan2(4, 3), math.Atan2(3, 4)}, ulps4) {
		t.Errorf("one triangle: (%+v, %d, %v)", sol, n, err)
	}
	// A tiny, a = b: B = A and C = π − 2A rounds to π, but c = a·sin(2A)/sin A = 2a still.
	sol, n, err = SolveTriangleSSA(1, 1, 1e-200)
	if err != nil || n != 1 || !relNear(sol[0].SideC, 2, ulps4) {
		t.Errorf("C near π: (%+v, %d, %v), want c = 2", sol, n, err)
	}
	// a = b with A obtuse: no triangle.
	if sol, n, err = SolveTriangleSSA(5, 5, radOf(91)); err != nil || n != 0 {
		t.Errorf("isosceles obtuse: (%+v, %d, %v)", sol, n, err)
	}
}

func TestSolveTriangleSSA_errors(t *testing.T) {
	call := func(a, b, angleA float64) func() error {
		return func() error { _, _, err := SolveTriangleSSA(a, b, angleA); return err }
	}
	runErrCases(t, []errCase{
		{"a NaN", call(math.NaN(), 1, 1), ErrNotFinite},
		{"A +Inf", call(1, 1, math.Inf(1)), ErrNotFinite},
		{"b zero", call(1, 0, 1), ErrNonPositive},
		{"A π", call(1, 1, math.Pi), ErrDomain},
		{"c overflows", call(1e308, 1e308, 1e-300), ErrOverflow},
	})
}

// Oracle: a random triangle solved from (a, b, A): the original is among the solutions (1e-6, for the
// nearly double roots near B = π/2, Kahan 2014 §9), the count is 2 when a < b and 1 when a > b, and
// every solution closes (SolveTriangleSSS gives back its angles to 1e-9).
func TestSolveTriangleSSA_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		tr := randomTriangle(rng)
		sol, n, err := SolveTriangleSSA(tr.SideA, tr.SideB, tr.AngleARad)
		wantN := 1
		if tr.SideA < tr.SideB {
			wantN = 2
		}
		if err != nil || n != wantN {
			t.Fatalf("SolveTriangleSSA(%v, %v, %v) = (%d, %v), want %d", tr.SideA, tr.SideB, tr.AngleARad, n, err, wantN)
		}
		found := false
		for _, s := range sol[:n] {
			found = found || sameTriangle(s, tr, 1e-6)
			check, err := SolveTriangleSSS(s.SideA, s.SideB, s.SideC)
			if err != nil || !near(check.AngleBRad, s.AngleBRad, 1e-9) || !near(check.AngleCRad, s.AngleCRad, 1e-9) {
				t.Fatalf("solution %+v does not close: %+v", s, check)
			}
		}
		if !found {
			t.Fatalf("SolveTriangleSSA(%v, %v, %v) = %+v, missing %+v", tr.SideA, tr.SideB, tr.AngleARad, sol[:n], tr)
		}
	}
}

func TestSolveTriangleSSA_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sol, _, _ := SolveTriangleSSA(3, 17, 0.17); sinkTriangle = sol[0] }); n != 0 {
		t.Errorf("SolveTriangleSSA allocates %v times, want 0", n)
	}
}

func BenchmarkSolveTriangleSSA(b *testing.B) {
	b.ReportAllocs()
	sa, sb, angleA := 3.0, 17.0, 0.17
	for b.Loop() {
		sol, _, _ := SolveTriangleSSA(sa, sb, angleA)
		sinkTriangle = sol[0]
	}
}

// A near π (reviewer case a = 2, b = 1), δ = π − A from 1e-3 to 1e-12, against a 256-bit oracle:
// sin B = b sin A / a (DLMF 4.42.4), B = atan(s/√(1 − s²)), C = π − A − B, c = a sin C / sin A.
// Before the fix c was 1.0000012246 (true 1.0000000000) and C had relative error 1.1e-6 at δ = 1e-10.
func TestSolveTriangleSSA_nearPi(t *testing.T) {
	for _, d := range nearPiDeltas {
		a, b, angleA := 2.0, 1.0, math.Pi-d
		A := bf(angleA)
		sinA := bigSin(A)
		s := bnew().Quo(bnew().Mul(bf(b), sinA), bf(a))
		wantB := bigAtan(bnew().Quo(s, bnew().Sqrt(bnew().Sub(bf(1), bnew().Mul(s, s)))))
		wantC := thirdBig(A, wantB)
		wantc := bnew().Quo(bnew().Mul(bf(a), bigSin(wantC)), sinA)
		want := Triangle{a, b, f64(wantc), angleA, f64(wantB), f64(wantC)}
		sol, n, err := SolveTriangleSSA(a, b, angleA)
		if err != nil || n != 1 || !relSameTriangle(sol[0], want, tolNearPi) {
			t.Errorf("SolveTriangleSSA(2, 1, π − %v) = (%+v, %d, %v), want %+v", d, sol[0], n, err, want)
		}
	}
}
