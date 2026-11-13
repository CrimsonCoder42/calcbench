package trig

import (
	"math"
	"testing"
)

// Worked example, hand-derived from DLMF 4.42.4: A = B = 60°, c = 10 is the equilateral triangle.
func TestSolveTriangleASA_example(t *testing.T) {
	got, err := SolveTriangleASA(radOf(60), radOf(60), 10)
	if err != nil || !sameTriangle(got, Triangle{10, 10, 10, math.Pi / 3, math.Pi / 3, math.Pi / 3}, ulps4) {
		t.Errorf("SolveTriangleASA(60°, 60°, 10) = (%+v, %v)", got, err)
	}
}

func TestSolveTriangleASA_cases(t *testing.T) {
	tests := []struct {
		name              string
		angleA, angleB, c float64
		want              Triangle
	}{
		{"3-4-5", 0.6435011087932844, math.Pi/2 - 0.6435011087932844, 5, Triangle{3, 4, 5, 0.6435011087932844, math.Pi/2 - 0.6435011087932844, math.Pi / 2}},
		{"30-60-90", math.Pi / 6, math.Pi / 3, 2, Triangle{1, math.Sqrt(3), 2, math.Pi / 6, math.Pi / 3, math.Pi / 2}},
	}
	// A + B tiny: C = π − A − B rounds to π, but a = c·sin A / sin(A + B) = c/2 still.
	if got, err := SolveTriangleASA(1e-200, 1e-200, 2); err != nil || !relNear(got.SideA, 1, ulps4) || !relNear(got.SideB, 1, ulps4) {
		t.Errorf("C near π: (%+v, %v), want a = b = 1", got, err)
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := SolveTriangleASA(tc.angleA, tc.angleB, tc.c); err != nil || !sameTriangle(got, tc.want, ulps4) {
				t.Errorf("got (%+v, %v), want %+v", got, err, tc.want)
			}
		})
	}
}

func TestSolveTriangleASA_errors(t *testing.T) {
	call := func(angleA, angleB, c float64) func() error {
		return func() error { _, err := SolveTriangleASA(angleA, angleB, c); return err }
	}
	runErrCases(t, []errCase{
		{"A NaN", call(math.NaN(), 1, 1), ErrNotFinite},
		{"c +Inf", call(1, 1, math.Inf(1)), ErrNotFinite},
		{"c zero", call(1, 1, 0), ErrNonPositive},
		{"A zero", call(0, 1, 1), ErrDomain},
		{"angles sum to π", call(2, math.Pi-2, 1), ErrDomain},
		{"B past π", call(0.1, 4, 1), ErrDomain},
		{"a overflows", call(1, math.Pi-1-1e-10, 1e300), ErrOverflow},
	})
}

// Oracle: a random triangle solved from (A, B, c) gives back every part, to 1e-12 relative (the angle
// C = π − A − B and the sine ratios round a few times).
func TestSolveTriangleASA_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		tr := randomTriangle(rng)
		got, err := SolveTriangleASA(tr.AngleARad, tr.AngleBRad, tr.SideC)
		if err != nil || !sameTriangle(got, tr, 1e-12) {
			t.Fatalf("SolveTriangleASA(%v, %v, %v) = (%+v, %v), want %+v", tr.AngleARad, tr.AngleBRad, tr.SideC, got, err, tr)
		}
	}
}

func TestSolveTriangleASA_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkTriangle, _ = SolveTriangleASA(0.6, 1.2, 4) }); n != 0 {
		t.Errorf("SolveTriangleASA allocates %v times, want 0", n)
	}
}

func BenchmarkSolveTriangleASA(b *testing.B) {
	b.ReportAllocs()
	angleA, angleB, c := 0.6, 1.2, 4.0
	for b.Loop() {
		sinkTriangle, _ = SolveTriangleASA(angleA, angleB, c)
	}
}

// A (or B) near π, δ from 1e-3 to 1e-12, the other angle δ/3, c = 1, against a 256-bit oracle:
// C = π − A − B, a = c sin A / sin C, b = c sin B / sin C (DLMF 4.42.4).
func TestSolveTriangleASA_nearPi(t *testing.T) {
	for _, d := range nearPiDeltas {
		for _, ang := range [][2]float64{{math.Pi - d, d / 3}, {d / 3, math.Pi - d}} {
			A, B := bf(ang[0]), bf(ang[1])
			wantC := thirdBig(A, B)
			sinC := bigSin(wantC)
			want := Triangle{f64(bnew().Quo(bigSin(A), sinC)), f64(bnew().Quo(bigSin(B), sinC)), 1, ang[0], ang[1], f64(wantC)}
			got, err := SolveTriangleASA(ang[0], ang[1], 1)
			if err != nil || !relSameTriangle(got, want, tolNearPi) {
				t.Errorf("SolveTriangleASA(%v, %v, 1) = (%+v, %v), want %+v", ang[0], ang[1], got, err, want)
			}
		}
	}
}
