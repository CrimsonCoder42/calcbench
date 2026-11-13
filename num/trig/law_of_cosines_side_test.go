package trig

import (
	"math"
	"math/big"
	"testing"
)

// Worked example, hand-derived from DLMF 4.42.5 (c² = a² + b² − 2ab cos C): a = b = 10, C = 60° gives
// c = 10; a = 3, b = 4, C = 90° gives c = 5.
func TestLawOfCosinesSide_dlmfExample(t *testing.T) {
	for _, tc := range []struct{ a, b, angleC, want float64 }{{10, 10, radOf(60), 10}, {3, 4, math.Pi / 2, 5}} {
		if got, err := LawOfCosinesSide(tc.a, tc.b, tc.angleC); err != nil || !relNear(got, tc.want, ulps4) {
			t.Errorf("LawOfCosinesSide(%v, %v, %v) = (%v, %v), want %v", tc.a, tc.b, tc.angleC, got, err, tc.want)
		}
	}
}

func TestLawOfCosinesSide_cases(t *testing.T) {
	tests := []struct {
		name         string
		a, b, angleC float64
		want         float64
	}{
		// Needle: a = b, tiny C: c = 2a·sin(C/2) exactly in the formula (|a−b| = 0). The textbook form
		// gives 0 here (1 − cos 1e-10 rounds to 0).
		{"needle, a = b", 1, 1, 1e-10, 2 * math.Sin(0.5e-10)},
		{"near-flat C, c ≈ a + b", 1, 1, math.Pi - 1e-10, 2 * math.Sin((math.Pi-1e-10)/2)},
		{"huge sides that fit", 1e308, 1e308, math.Pi / 2, math.Sqrt2 * 1e308},
		{"tiny sides", 3e-310, 4e-310, math.Pi / 2, 5e-310},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := LawOfCosinesSide(tc.a, tc.b, tc.angleC)
			tol := ulps4
			if tc.want < 1e-300 {
				tol = 1e-9 // subnormal result: about 34 bits of precision left
			}
			if err != nil || !relNear(got, tc.want, tol) {
				t.Errorf("got (%v, %v), want %v", got, err, tc.want)
			}
		})
	}
}

func TestLawOfCosinesSide_errors(t *testing.T) {
	call := func(a, b, angleC float64) func() error {
		return func() error { _, err := LawOfCosinesSide(a, b, angleC); return err }
	}
	runErrCases(t, []errCase{
		{"a NaN", call(math.NaN(), 1, 1), ErrNotFinite},
		{"b -Inf", call(1, math.Inf(-1), 1), ErrNotFinite},
		{"C +Inf", call(1, 1, math.Inf(1)), ErrNotFinite},
		{"a zero", call(0, 1, 1), ErrNonPositive},
		{"b negative", call(1, -2, 1), ErrNonPositive},
		{"C zero", call(1, 1, 0), ErrDomain},
		{"C π", call(1, 1, math.Pi), ErrDomain},
		{"C past π", call(1, 1, 4), ErrDomain},
		{"result overflows", call(1e308, 1e308, 3), ErrOverflow},
	})
}

// Oracle 1: c² = a² + b² − 2ab·cos C in 256-bit big.Float (cos C from math.Cos, one rounding) for
// C in [0.5, π − 0.5], where the textbook form is well conditioned: 1e-13 relative.
// Oracle 2: the round trip LawOfCosinesAngle(a, b, c) gives C back, for C down to 3e-8, to 1e-14 times
// (1 + the condition number of C in c).
func TestLawOfCosinesSide_oracle(t *testing.T) {
	rng := newRNG()
	for range 500 {
		a := math.Pow(10, rng.Float64()*6-3)
		b := a * math.Pow(10, rng.Float64()*2-1)
		angleC := 0.5 + rng.Float64()*(math.Pi-1)
		c, err := LawOfCosinesSide(a, b, angleC)
		if err != nil {
			t.Fatalf("LawOfCosinesSide(%v, %v, %v): %v", a, b, angleC, err)
		}
		A, B := new(big.Float).SetPrec(256).SetFloat64(a), new(big.Float).SetPrec(256).SetFloat64(b)
		sq := new(big.Float).SetPrec(256).Add(new(big.Float).Mul(A, A), new(big.Float).Mul(B, B))
		ab2 := new(big.Float).SetPrec(256).Mul(new(big.Float).Mul(A, B), big.NewFloat(2*math.Cos(angleC)))
		want, _ := sq.Sub(sq, ab2).Float64()
		if !relNear(c*c, want, 1e-13) {
			t.Fatalf("LawOfCosinesSide(%v, %v, %v) = %v, c² = %v, want %v", a, b, angleC, c, c*c, want)
		}
	}
	for range 500 {
		a := math.Pow(10, rng.Float64()*6-3)
		b := a * math.Pow(10, rng.Float64()*2-1)
		angleC := math.Pow(10, -rng.Float64()*8) * 3
		c, err := LawOfCosinesSide(a, b, angleC)
		back, err2 := LawOfCosinesAngle(a, b, c)
		// Rounding c by one ulp moves the angle by (c²/(ab·C·sin C)) ulps relative: the condition number.
		cond := c * c / (a * b * angleC * math.Sin(angleC))
		if err != nil || err2 != nil || !relNear(back, angleC, 1e-14*(1+cond)) {
			t.Fatalf("LawOfCosinesSide(%v, %v, %v) = %v; angle back = (%v, %v), cond %v", a, b, angleC, c, back, err2, cond)
		}
	}
}

func TestLawOfCosinesSide_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = LawOfCosinesSide(3, 4, 1.2) }); n != 0 {
		t.Errorf("LawOfCosinesSide allocates %v times, want 0", n)
	}
}

func BenchmarkLawOfCosinesSide(b *testing.B) {
	b.ReportAllocs()
	sa, sb, angleC := 3.0, 4.0, 1.2
	for b.Loop() {
		sink, _ = LawOfCosinesSide(sa, sb, angleC)
	}
}
