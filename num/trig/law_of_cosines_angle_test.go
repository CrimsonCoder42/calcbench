package trig

import (
	"math"
	"math/big"
	"testing"
)

// oracleAngle is the angle opposite c computed independently of Kahan's arrangement: the exact
// tan²(C/2) = (c² − (a−b)²)/((a+b)² − c²) (from DLMF 4.42.5 and the half-angle formula 4.21.23) in
// 512-bit big.Float, rounded once, then 2·atan(√·). atan and √ are well conditioned, so the result is
// within a few ulps of the true C for the float64 sides. ok is false when the sides are not a triangle.
func oracleAngle(a, b, c float64) (angle float64, ok bool) {
	const p = 512
	A := new(big.Float).SetPrec(p).SetFloat64(a)
	B := new(big.Float).SetPrec(p).SetFloat64(b)
	C := new(big.Float).SetPrec(p).SetFloat64(c)
	d := new(big.Float).SetPrec(p).Sub(A, B)
	s := new(big.Float).SetPrec(p).Add(A, B)
	c2 := new(big.Float).SetPrec(p).Mul(C, C)
	num := new(big.Float).SetPrec(p).Sub(c2, new(big.Float).SetPrec(p).Mul(d, d))
	den := new(big.Float).SetPrec(p).Sub(new(big.Float).SetPrec(p).Mul(s, s), c2)
	if num.Sign() < 0 || den.Sign() < 0 {
		return 0, false
	}
	if den.Sign() == 0 {
		return math.Pi, true
	}
	t, _ := new(big.Float).SetPrec(p).Quo(num, den).Float64()
	return 2 * math.Atan(math.Sqrt(t)), true
}

// kahanTable1 is Kahan 2014, §5, Table 1 (p. 4): sides a, b, c and the accurate angle C in degrees,
// which Kahan checked to at least nine significant digits. tolDeg is the allowed difference: 2e-8
// relative where rounding the decimal data to float64 moves C by less than that, and wider on the rows
// where it does not (checked against oracleAngle, which agrees with the code to 1e-13 on every row).
var kahanTable1 = []struct{ a, b, c, wantDeg, relTol float64 }{
	{10, 10, 10, 60, 2e-8},
	{100000, 99999.99979, 0.00029, 1.145915591e-7, 2e-8},
	{100000, 100000, 1.00005, 5.72986443e-4, 2e-8},
	// a − b = 2e-5 rounds by up to an ulp of 1e5 (1.5e-11), so C moves by up to about 1e-6 relative.
	{99999.99996, 99999.99994, 0.00003, 1.281172578e-8, 1e-6},
	{99999.99996, 0.00003, 99999.99994, 48.18968510, 1e-6},
	{10000, 5000.000001, 15000, 179.9985965, 2e-8},
	// Flat in decimals (5278.64055 + 94721.35941 = 99999.99996); the float64 data are not quite flat
	// and give 179.9999978°: tolerance 2e-8 relative of 180, i.e. 3.6e-6°.
	{5278.64055, 94721.35941, 99999.99996, 180, 2e-8},
	{100002, 100002, 200004, 180, 0},
	{31622.77662, 0.000023, 31622.77661, 64.22853824, 2e-8},
	{31622.77662, 0.0155555, 31622.77661, 89.96315276, 2e-8},
}

// Worked example: Kahan 2014 Table 1, the "Accurate C°" column.
func TestLawOfCosinesAngle_kahanTable1(t *testing.T) {
	for _, row := range kahanTable1 {
		got, err := LawOfCosinesAngle(row.a, row.b, row.c)
		if err != nil || !relNear(degOf(got), row.wantDeg, row.relTol) {
			t.Errorf("LawOfCosinesAngle(%v, %v, %v) = (%.10g°, %v), want %.10g°", row.a, row.b, row.c, degOf(got), err, row.wantDeg)
		}
		if want, _ := oracleAngle(row.a, row.b, row.c); !relNear(got, want, 1e-13) {
			t.Errorf("LawOfCosinesAngle(%v, %v, %v) = %v, oracle %v", row.a, row.b, row.c, got, want)
		}
	}
}

// Kahan 2014 Table 1, rows marked "Error": (99999.99999, 99999.99999, 200000) is not a triangle, and
// (−3, 4, 2) has a negative side.
func TestLawOfCosinesAngle_kahanErrorRows(t *testing.T) {
	runErrCases(t, []errCase{
		{"c > a + b", func() error { _, err := LawOfCosinesAngle(99999.99999, 99999.99999, 200000); return err }, ErrNotTriangle},
		{"negative side", func() error { _, err := LawOfCosinesAngle(-3, 4, 2); return err }, ErrNonPositive},
	})
}

func TestLawOfCosinesAngle_cases(t *testing.T) {
	tests := []struct {
		name          string
		a, b, c, want float64
	}{
		{"3-4-5 right angle", 3, 4, 5, math.Pi / 2},
		{"swapped a and b", 4, 3, 5, math.Pi / 2},
		{"flat, c = a + b", 1, 2, 3, math.Pi},
		{"flat, c = a − b", 3, 2, 1, 0},
		{"flat, c = b − a", 2, 3, 1, 0},
		{"equilateral, huge", 1e308, 1e308, 1e308, math.Pi / 3},
		{"equilateral, subnormal", 1e-320, 1e-320, 1e-320, math.Pi / 3},
		{"3-4-5 scaled huge", 3e307, 4e307, 5e307, math.Pi / 2},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := LawOfCosinesAngle(tc.a, tc.b, tc.c); err != nil || !near(got, tc.want, ulps4) {
				t.Errorf("LawOfCosinesAngle(%v, %v, %v) = (%v, %v), want %v", tc.a, tc.b, tc.c, got, err, tc.want)
			}
		})
	}
}

func TestLawOfCosinesAngle_errors(t *testing.T) {
	call := func(a, b, c float64) func() error {
		return func() error { _, err := LawOfCosinesAngle(a, b, c); return err }
	}
	runErrCases(t, []errCase{
		{"a NaN", call(math.NaN(), 1, 1), ErrNotFinite},
		{"b +Inf", call(1, math.Inf(1), 1), ErrNotFinite},
		{"c -Inf", call(1, 1, math.Inf(-1)), ErrNotFinite},
		{"a zero", call(0, 1, 1), ErrNonPositive},
		{"c negative", call(1, 1, -1), ErrNonPositive},
		{"a too long", call(5, 1, 1), ErrNotTriangle},
		{"b too long", call(1, 5, 1), ErrNotTriangle},
		{"c too long", call(1, 1, 5), ErrNotTriangle},
	})
}

// Oracle: random triangles, many of them needle-like (c near |a−b| or a+b), against oracleAngle.
// Tolerance 1e-13 relative: Kahan proves a few ulps (§4).
func TestLawOfCosinesAngle_oracle(t *testing.T) {
	rng := newRNG()
	for trial := range 1000 {
		a := math.Pow(10, rng.Float64()*8-4)
		b := a * math.Pow(10, rng.Float64()*4-2)
		lo, hi := math.Abs(a-b), a+b
		var c float64
		switch trial % 3 {
		case 0:
			c = lo + rng.Float64()*(hi-lo)
		case 1:
			c = lo + (hi-lo)*math.Pow(10, -rng.Float64()*12) // near |a − b|: C tiny
		default:
			c = hi - (hi-lo)*math.Pow(10, -rng.Float64()*12) // near a + b: C near π
		}
		want, ok := oracleAngle(a, b, c)
		got, err := LawOfCosinesAngle(a, b, c)
		if !ok {
			if err == nil {
				t.Fatalf("LawOfCosinesAngle(%v, %v, %v) = %v, want ErrNotTriangle", a, b, c, got)
			}
			continue
		}
		if err != nil || !relNear(got, want, 1e-13) {
			t.Fatalf("LawOfCosinesAngle(%v, %v, %v) = (%v, %v), want %v", a, b, c, got, err, want)
		}
	}
}

func TestLawOfCosinesAngle_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = LawOfCosinesAngle(3, 4, 5) }); n != 0 {
		t.Errorf("LawOfCosinesAngle allocates %v times, want 0", n)
	}
}

func BenchmarkLawOfCosinesAngle(b *testing.B) {
	b.ReportAllocs()
	sa, sb, sc := 3.0, 4.0, 5.0
	for b.Loop() {
		sink, _ = LawOfCosinesAngle(sa, sb, sc)
	}
}
