package trig

import (
	"math"
	"testing"
)

// Worked example, hand-derived from DLMF 4.42.4: A = atan2(3, 4), B = π/2, a = 3 gives b = 5, c = 4.
func TestSolveTriangleAAS_example(t *testing.T) {
	angleA := math.Atan2(3, 4)
	got, err := SolveTriangleAAS(angleA, math.Pi/2, 3)
	if err != nil || !sameTriangle(got, Triangle{3, 5, 4, angleA, math.Pi / 2, math.Pi/2 - angleA}, ulps4) {
		t.Errorf("SolveTriangleAAS(atan2(3,4), π/2, 3) = (%+v, %v)", got, err)
	}
}

func TestSolveTriangleAAS_cases(t *testing.T) {
	got, err := SolveTriangleAAS(math.Pi/3, math.Pi/3, 7)
	if err != nil || !sameTriangle(got, Triangle{7, 7, 7, math.Pi / 3, math.Pi / 3, math.Pi / 3}, ulps4) {
		t.Errorf("equilateral: (%+v, %v)", got, err)
	}
	// A + B tiny: C rounds to π, but c = a·sin(A + B)/sin A = 2a still.
	got, err = SolveTriangleAAS(1e-200, 1e-200, 1)
	if err != nil || !relNear(got.SideC, 2, ulps4) {
		t.Errorf("C near π: (%+v, %v), want c = 2", got, err)
	}
	got, err = SolveTriangleAAS(math.Pi/6, math.Pi/2, 1e-300)
	if err != nil || !sameTriangle(got, Triangle{1e-300, 2e-300, math.Sqrt(3) * 1e-300, math.Pi / 6, math.Pi / 2, math.Pi / 3}, ulps4) {
		t.Errorf("tiny 30-60-90: (%+v, %v)", got, err)
	}
}

func TestSolveTriangleAAS_errors(t *testing.T) {
	call := func(angleA, angleB, a float64) func() error {
		return func() error { _, err := SolveTriangleAAS(angleA, angleB, a); return err }
	}
	runErrCases(t, []errCase{
		{"B NaN", call(1, math.NaN(), 1), ErrNotFinite},
		{"a -Inf", call(1, 1, math.Inf(-1)), ErrNotFinite},
		{"a negative", call(1, 1, -1), ErrNonPositive},
		{"B zero", call(1, 0, 1), ErrDomain},
		{"angles sum past π", call(2, 2, 1), ErrDomain},
		{"b overflows", call(1e-300, math.Pi/2, 1e10), ErrOverflow},
	})
}

// Oracle: a random triangle solved from (A, B, a) gives back every part, to 1e-12 relative.
func TestSolveTriangleAAS_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		tr := randomTriangle(rng)
		got, err := SolveTriangleAAS(tr.AngleARad, tr.AngleBRad, tr.SideA)
		if err != nil || !sameTriangle(got, tr, 1e-12) {
			t.Fatalf("SolveTriangleAAS(%v, %v, %v) = (%+v, %v), want %+v", tr.AngleARad, tr.AngleBRad, tr.SideA, got, err, tr)
		}
	}
}

func TestSolveTriangleAAS_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkTriangle, _ = SolveTriangleAAS(0.6, 1.2, 3) }); n != 0 {
		t.Errorf("SolveTriangleAAS allocates %v times, want 0", n)
	}
}

func BenchmarkSolveTriangleAAS(b *testing.B) {
	b.ReportAllocs()
	angleA, angleB, a := 0.6, 1.2, 3.0
	for b.Loop() {
		sinkTriangle, _ = SolveTriangleAAS(angleA, angleB, a)
	}
}

// A (or B) near π, δ from 1e-3 to 1e-12, the other angle δ/3, a = 2, against a 256-bit oracle:
// C = π − A − B, b = a sin B / sin A, c = a sin C / sin A (DLMF 4.42.4).
func TestSolveTriangleAAS_nearPi(t *testing.T) {
	for _, d := range nearPiDeltas {
		for _, ang := range [][2]float64{{math.Pi - d, d / 3}, {d / 3, math.Pi - d}} {
			a := 2.0
			A, B := bf(ang[0]), bf(ang[1])
			sinA := bigSin(A)
			wantC := thirdBig(A, B)
			want := Triangle{a, f64(bnew().Quo(bnew().Mul(bf(a), bigSin(B)), sinA)), f64(bnew().Quo(bnew().Mul(bf(a), bigSin(wantC)), sinA)), ang[0], ang[1], f64(wantC)}
			got, err := SolveTriangleAAS(ang[0], ang[1], a)
			if err != nil || !relSameTriangle(got, want, tolNearPi) {
				t.Errorf("SolveTriangleAAS(%v, %v, 2) = (%+v, %v), want %+v", ang[0], ang[1], got, err, want)
			}
		}
	}
}
