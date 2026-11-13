package trig

import (
	"math"
	"testing"
)

// Worked example, hand-derived: (3, 4, C = 90°) is the 3-4-5 triangle, A = 0.6435011087932844.
func TestSolveTriangleSAS_example(t *testing.T) {
	got, err := SolveTriangleSAS(3, 4, math.Pi/2)
	if err != nil || !sameTriangle(got, Triangle{3, 4, 5, 0.6435011087932844, math.Pi/2 - 0.6435011087932844, math.Pi / 2}, ulps4) {
		t.Errorf("SolveTriangleSAS(3, 4, π/2) = (%+v, %v)", got, err)
	}
}

func TestSolveTriangleSAS_cases(t *testing.T) {
	tests := []struct {
		name         string
		a, b, angleC float64
		want         Triangle
	}{
		{"equilateral", 10, 10, math.Pi / 3, Triangle{10, 10, 10, math.Pi / 3, math.Pi / 3, math.Pi / 3}},
		// Needle: a = b, C tiny: c = 2 sin(C/2), A = B = (π − C)/2.
		{"needle", 1, 1, 1e-10, Triangle{1, 1, 2 * math.Sin(0.5e-10), (math.Pi - 1e-10) / 2, (math.Pi - 1e-10) / 2, 1e-10}},
		{"huge", 3e307, 4e307, math.Pi / 2, Triangle{3e307, 4e307, 5e307, 0.6435011087932844, math.Pi/2 - 0.6435011087932844, math.Pi / 2}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := SolveTriangleSAS(tc.a, tc.b, tc.angleC); err != nil || !sameTriangle(got, tc.want, ulps4) {
				t.Errorf("got (%+v, %v), want %+v", got, err, tc.want)
			}
		})
	}
}

func TestSolveTriangleSAS_errors(t *testing.T) {
	call := func(a, b, angleC float64) func() error {
		return func() error { _, err := SolveTriangleSAS(a, b, angleC); return err }
	}
	runErrCases(t, []errCase{
		{"a NaN", call(math.NaN(), 1, 1), ErrNotFinite},
		{"C -Inf", call(1, 1, math.Inf(-1)), ErrNotFinite},
		{"b zero", call(1, 0, 1), ErrNonPositive},
		{"C zero", call(1, 1, 0), ErrDomain},
		{"C π", call(1, 1, math.Pi), ErrDomain},
		{"c overflows", call(1e308, 1e308, 3), ErrOverflow},
	})
}

// Oracle: a random triangle solved from (a, b, C) gives back c, A and B. The sides were rounded when the
// triangle was built, so the comparison allows 1e-9 (a needle's angles are sensitive to that rounding);
// and the returned parts satisfy the law of sines (DLMF 4.42.4) to 1e-12.
func TestSolveTriangleSAS_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		tr := randomTriangle(rng)
		got, err := SolveTriangleSAS(tr.SideA, tr.SideB, tr.AngleCRad)
		if err != nil || !sameTriangle(got, tr, 1e-9) {
			t.Fatalf("SolveTriangleSAS(%v, %v, %v) = (%+v, %v), want %+v", tr.SideA, tr.SideB, tr.AngleCRad, got, err, tr)
		}
		if !relNear(got.SideA/math.Sin(got.AngleARad), got.SideC/math.Sin(got.AngleCRad), 1e-12) {
			t.Fatalf("SolveTriangleSAS: law of sines fails for %+v", got)
		}
	}
}

func TestSolveTriangleSAS_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkTriangle, _ = SolveTriangleSAS(3, 4, 1.2) }); n != 0 {
		t.Errorf("SolveTriangleSAS allocates %v times, want 0", n)
	}
}

func BenchmarkSolveTriangleSAS(b *testing.B) {
	b.ReportAllocs()
	sa, sb, angleC := 3.0, 4.0, 1.2
	for b.Loop() {
		sinkTriangle, _ = SolveTriangleSAS(sa, sb, angleC)
	}
}

// C near π (reviewer case a = 2, b = 1, both orders, and a = b), δ = π − C from 1e-3 to 1e-12, against
// a 256-bit oracle: B = atan(b sin C / (a − b cos C)) (DLMF 4.42.4 with 4.42.6), A = π − C − B,
// c = √(a² + b² − 2ab cos C) (4.42.5). Before the fix A and B had relative error up to 100% at δ = 1e-10.
func TestSolveTriangleSAS_nearPi(t *testing.T) {
	for _, sides := range [][2]float64{{2, 1}, {1, 2}, {1, 1}, {1e300, 3e299}} {
		for _, d := range nearPiDeltas {
			a, b, angleC := sides[0], sides[1], math.Pi-d
			A, B, C := bf(a), bf(b), bf(angleC)
			cosC := bigCos(C)
			// angle opposite b: atan2(b sin C, a − b cos C), positive denominator since cos C < 0
			wantB := bigAtan(bnew().Quo(bnew().Mul(B, bigSin(C)), bnew().Sub(A, bnew().Mul(B, cosC))))
			wantA := thirdBig(C, wantB)
			c2 := bnew().Sub(bnew().Add(bnew().Mul(A, A), bnew().Mul(B, B)), bnew().Mul(bnew().Mul(bf(2), bnew().Mul(A, B)), cosC))
			want := Triangle{a, b, f64(bnew().Sqrt(c2)), f64(wantA), f64(wantB), angleC}
			got, err := SolveTriangleSAS(a, b, angleC)
			if err != nil || !relSameTriangle(got, want, tolNearPi) {
				t.Errorf("SolveTriangleSAS(%v, %v, π − %v) = (%+v, %v), want %+v", a, b, d, got, err, want)
			}
		}
	}
}
