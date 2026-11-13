package trig

import (
	"math"
	"testing"
)

// kahanTable2 is Kahan 2014, §9, Table 2 (p. 8): data a, b, A° and the true B° (computed by Kahan to at
// least 30 digits and printed to 10). A NaN B means no triangle. The row (1, 2, 30°) is tested apart:
// it is an exact double root, see TestLawOfSinesAngle_doubleRoot.
var kahanTable2 = []struct{ a, b, angleADeg, trueBDeg float64 }{
	{3, 17, 10, 79.73894317},
	{3, 17, 10.16424862, 89.99895044},
	{0.9999999998, 2, 28, 69.87481897},
	{0.9999999998, 2, 29.9916, 88.70873768},
	{0.9999999998, 2, 30, math.NaN()},
	{18817, 21728, 59.9, 87.42569088},
	{38620, 38673, 86.9, 89.21861758},
	{4.999999999, 5.000000001, 88.3, 88.30000077},
	{4.999999999, 5.000000001, 89.999007, math.NaN()},
	{5, 5, 89.999007, 89.999007},
	{8.000000001, 7.999999999, 89.999007, 89.99837906},
	{8.000000001, 7.999999999, 90.000993, 89.99837906},
	{8.000000001, 7.999999999, 91, 88.99999918},
	{8.000000001, 7.999999999, 5.739170509, 5.739170508},
}

// Worked example: Kahan 2014 Table 2. Tolerance 1e-8 degrees: Kahan prints B to 10 significant digits
// (half a unit is at most 5e-9°), and rounding the decimal data to float64 moves B by less than that on
// these rows. a < b gives two triangles, B and 180° − B; a ≥ b gives one.
func TestLawOfSinesAngle_kahanTable2(t *testing.T) {
	for _, row := range kahanTable2 {
		b1, b2, n, err := LawOfSinesAngle(row.a, row.b, radOf(row.angleADeg))
		if err != nil {
			t.Errorf("%v: err = %v", row, err)
			continue
		}
		if math.IsNaN(row.trueBDeg) {
			if n != 0 || b1 != 0 || b2 != 0 {
				t.Errorf("%v: (%v, %v, %d), want no triangle", row, b1, b2, n)
			}
			continue
		}
		wantN := 1
		if row.a < row.b {
			wantN = 2
		}
		if n != wantN || math.Abs(degOf(b1)-row.trueBDeg) > 1e-8 {
			t.Errorf("%v: B1 = %.10f°, count %d; want %.10f°, count %d", row, degOf(b1), n, row.trueBDeg, wantN)
		}
		if n == 2 && !near(b1+b2, math.Pi, ulps4) {
			t.Errorf("%v: B1 + B2 = %v, want π", row, b1+b2)
		}
	}
}

// Kahan 2014 Table 2, row (1, 2, 30°): in exact decimals B = 90° is a double root. math.Pi/6 is below π/6,
// so for the float64 data there are two triangles with B within about √ε of 90°: tolerance 1e-6°
// (Kahan §9: a nearly double root can lose half the digits carried).
func TestLawOfSinesAngle_doubleRoot(t *testing.T) {
	b1, b2, n, err := LawOfSinesAngle(1, 2, radOf(30))
	if err != nil || n == 0 || math.Abs(degOf(b1)-90) > 1e-6 || (n == 2 && math.Abs(degOf(b2)-90) > 1e-6) {
		t.Errorf("LawOfSinesAngle(1, 2, 30°) = (%v°, %v°, %d, %v), want B ≈ 90°", degOf(b1), degOf(b2), n, err)
	}
}

func TestLawOfSinesAngle_cases(t *testing.T) {
	tests := []struct {
		name         string
		a, b, angleA float64
		wantN        int
		want1, want2 float64
	}{
		{"a > b: one acute B", 5, 3, math.Pi / 2, 1, math.Asin(0.6), 0},
		{"a = b, A acute: B = A", 2, 2, 1, 1, 1, 0},
		{"a = b, A right: no triangle", 2, 2, math.Pi / 2, 0, 0, 0},
		{"a = b, A obtuse: no triangle", 2, 2, 2, 0, 0, 0},
		{"a < b, A obtuse: no triangle", 1, 2, 2, 0, 0, 0},
		{"a < b, A right: no triangle", 1, 2, math.Pi / 2, 0, 0, 0},
		{"a < b·sin A: no triangle", 0.4, 1, math.Pi / 6, 0, 0, 0},
		{"a < b: two triangles", 3, 4, math.Pi / 6, 2, math.Asin(4.0 / 6), math.Pi - math.Asin(4.0/6)},
		{"huge sides", 3e300, 4e300, math.Pi / 6, 2, math.Asin(4.0 / 6), math.Pi - math.Asin(4.0/6)},
		{"tiny sides", 3e-300, 4e-300, math.Pi / 6, 2, math.Asin(4.0 / 6), math.Pi - math.Asin(4.0/6)},
		// A tiny: a/tan A is beyond 2^500 and its square would overflow; sin B = (b/a) sin A, so B = (b/a)·A.
		{"tiny A, a > b", 1, 0.9, 1e-200, 1, 0.9e-200, 0},
		{"tiny A, a < b", 1, 1.1, 1e-200, 2, 1.1e-200, math.Pi},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b1, b2, n, err := LawOfSinesAngle(tc.a, tc.b, tc.angleA)
			if err != nil || n != tc.wantN || !near(b1, tc.want1, 1e-14) || !near(b2, tc.want2, 1e-14) {
				t.Errorf("got (%v, %v, %d, %v), want (%v, %v, %d)", b1, b2, n, err, tc.want1, tc.want2, tc.wantN)
			}
		})
	}
}

func TestLawOfSinesAngle_errors(t *testing.T) {
	call := func(a, b, angleA float64) func() error {
		return func() error { _, _, _, err := LawOfSinesAngle(a, b, angleA); return err }
	}
	runErrCases(t, []errCase{
		{"a NaN", call(math.NaN(), 1, 1), ErrNotFinite},
		{"b +Inf", call(1, math.Inf(1), 1), ErrNotFinite},
		{"A NaN", call(1, 1, math.NaN()), ErrNotFinite},
		{"a zero", call(0, 1, 1), ErrNonPositive},
		{"b negative", call(1, -1, 1), ErrNonPositive},
		{"A zero", call(1, 1, 0), ErrDomain},
		{"A π", call(1, 1, math.Pi), ErrDomain},
	})
}

// Oracle: a random triangle from its angles (sides k·sin A, k·sin B). The true B must be among the
// solutions, the count must be 2 when a < b and 1 when a > b, and every solution must satisfy
// sin B / b = sin A / a to 1e-13. Tolerance on B 1e-6 rad: near B = π/2 the root is nearly double and
// up to half the digits can go (Kahan 2014 §9); elsewhere the error is far smaller.
func TestLawOfSinesAngle_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		angleA := 0.01 + rng.Float64()*3
		angleB := 0.01 + rng.Float64()*(math.Pi-angleA-0.02)
		k := math.Pow(10, rng.Float64()*20-10)
		a, b := k*math.Sin(angleA), k*math.Sin(angleB)
		b1, b2, n, err := LawOfSinesAngle(a, b, angleA)
		wantN := 1
		if a < b {
			wantN = 2
		}
		if err != nil || n != wantN {
			t.Fatalf("LawOfSinesAngle(%v, %v, %v) = (%v, %v, %d, %v), want count %d", a, b, angleA, b1, b2, n, err, wantN)
		}
		found := math.Abs(b1-angleB) < 1e-6 || (n == 2 && math.Abs(b2-angleB) < 1e-6)
		if !found || !relNear(math.Sin(b1)/b, math.Sin(angleA)/a, 1e-13) {
			t.Fatalf("LawOfSinesAngle(%v, %v, %v) = (%v, %v, %d), want B = %v", a, b, angleA, b1, b2, n, angleB)
		}
	}
}

func TestLawOfSinesAngle_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _, _, _ = LawOfSinesAngle(3, 17, 0.17) }); n != 0 {
		t.Errorf("LawOfSinesAngle allocates %v times, want 0", n)
	}
}

func BenchmarkLawOfSinesAngle(b *testing.B) {
	b.ReportAllocs()
	sa, sb, angleA := 3.0, 17.0, 0.17
	for b.Loop() {
		sink, _, _, _ = LawOfSinesAngle(sa, sb, angleA)
	}
}
