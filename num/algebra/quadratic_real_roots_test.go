package algebra

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// ulpsApart reports whether got is within k units in the last place of want.
func ulpsApart(got, want float64, k float64) bool {
	ulp := math.Nextafter(math.Abs(want), math.Inf(1)) - math.Abs(want)
	return math.Abs(got-want) <= k*ulp
}

// Worked examples, hand-checked on the q form of Numerical Recipes in C, 2nd ed., §5.6, eqs. (5.6.4)–(5.6.5):
// x² − 3x + 2: D = 1, q = −½(−3 − 1) = 2, x₁ = q/a = 2, x₂ = c/q = 1.
// x² − 10⁸x + 1 (the cancellation case of Higham 2002, §1.8): hi = 10⁸ and lo = 10⁻⁸ to within 2 ulps,
// where the textbook formula (−b − √D)/(2a) loses about half the digits of lo.
func TestQuadraticRealRoots_workedExample(t *testing.T) {
	lo, hi, n, err := QuadraticRealRoots(1, -3, 2)
	if err != nil || lo != 1 || hi != 2 || n != 2 {
		t.Errorf("QuadraticRealRoots(1, -3, 2) = (%v, %v, %d, %v), want (1, 2, 2, nil)", lo, hi, n, err)
	}
	lo, hi, n, err = QuadraticRealRoots(1, -1e8, 1)
	if err != nil || n != 2 || !ulpsApart(lo, 1e-8, 2) || !ulpsApart(hi, 1e8, 2) {
		t.Errorf("QuadraticRealRoots(1, -1e8, 1) = (%v, %v, %d, %v), want (1e-8, 1e8, 2, nil)", lo, hi, n, err)
	}
}

func TestQuadraticRealRoots_cases(t *testing.T) {
	tests := []struct {
		name    string
		a, b, c float64
		lo, hi  float64
		n       int
	}{
		{"double root", 1, -2, 1, 1, 1, 1},
		{"no real roots is a result", 1, 0, 1, 0, 0, 0},
		{"b = 0", 1, 0, -4, -2, 2, 2},
		{"c = 0 gives a zero root", 2, -6, 0, 0, 3, 2},
		{"b = c = 0", 5, 0, 0, 0, 0, 1},
		{"negative a", -1, 3, -2, 1, 2, 2},
		{"negative roots", 1, 5, 6, -3, -2, 2},
		{"huge coefficients (b² would overflow)", 0x1p1000, -3 * 0x1p1000, 0x1p1001, 1, 2, 2},
		{"tiny coefficients (b² would underflow)", 0x1p-1000, -3 * 0x1p-1000, 0x1p-999, 1, 2, 2},
		{"huge roots", 1, -3 * 0x1p500, 0x1p1001, 0x1p500, 0x1p501, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lo, hi, n, err := QuadraticRealRoots(tt.a, tt.b, tt.c)
			if err != nil || lo != tt.lo || hi != tt.hi || n != tt.n {
				t.Errorf("QuadraticRealRoots(%v, %v, %v) = (%v, %v, %d, %v), want (%v, %v, %d, nil)",
					tt.a, tt.b, tt.c, lo, hi, n, err, tt.lo, tt.hi, tt.n)
			}
		})
	}
}

func TestQuadraticRealRoots_errors(t *testing.T) {
	nan, inf := math.NaN(), math.Inf(1)
	tests := []struct {
		name    string
		a, b, c float64
		want    error
	}{
		{"a = 0", 0, 1, 1, ErrDegenerate},
		{"NaN a", nan, 1, 1, ErrNotFinite},
		{"Inf b", 1, -inf, 1, ErrNotFinite},
		{"NaN c", 1, 1, nan, ErrNotFinite},
		{"root overflows", 0x1p-600, 0x1p600, 1, ErrOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lo, hi, n, err := QuadraticRealRoots(tt.a, tt.b, tt.c)
			if !errors.Is(err, tt.want) || lo != 0 || hi != 0 || n != 0 {
				t.Errorf("QuadraticRealRoots(%v, %v, %v) = (%v, %v, %d, %v), want (0, 0, 0, %v)", tt.a, tt.b, tt.c, lo, hi, n, err, tt.want)
			}
		})
	}
}

// Oracle: n must match the sign of the exact discriminant (math/big) exactly, since NUM-ALGEB-006 has
// relative error below 1 and so never changes a sign. For n = 2 the roots are checked against Vieta's
// formulas lo + hi = −b/a and lo·hi = c/a. Tolerance 16u relative to |lo| + |hi| and to |lo·hi|: the q form
// computes each root with a few roundings (√D carries at most 2u from D) and never cancels.
func TestQuadraticRealRoots_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 600 {
		a := (2*rng.Float64() - 1) * math.Ldexp(1, rng.IntN(200)-100)
		r1 := (2*rng.Float64() - 1) * math.Ldexp(1, rng.IntN(60)-30)
		r2 := (2*rng.Float64() - 1) * math.Ldexp(1, rng.IntN(60)-30)
		b := -a * (r1 + r2)
		c := a * r1 * r2
		if trial%3 == 0 {
			c = (2*rng.Float64() - 1) * math.Ldexp(1, rng.IntN(200)-100) // often D < 0
		}
		if a == 0 {
			continue
		}
		lo, hi, n, err := QuadraticRealRoots(a, b, c)
		if err != nil {
			t.Fatalf("trial %d: %v", trial, err)
		}
		wantN := bigDiscriminant(a, b, c).Sign() + 1
		if n != wantN {
			t.Fatalf("trial %d: QuadraticRealRoots(%v, %v, %v) has n = %d, want %d", trial, a, b, c, n, wantN)
		}
		if n != 2 {
			continue
		}
		if lo > hi {
			t.Fatalf("trial %d: lo %v > hi %v", trial, lo, hi)
		}
		tol := 16 * unitRoundoff
		if sum := lo + hi; math.Abs(sum-(-b/a)) > tol*(math.Abs(lo)+math.Abs(hi)) {
			t.Fatalf("trial %d: lo + hi = %v, −b/a = %v", trial, sum, -b/a)
		}
		if prod := lo * hi; math.Abs(prod-c/a) > tol*math.Abs(prod) {
			t.Fatalf("trial %d: lo·hi = %v, c/a = %v", trial, prod, c/a)
		}
	}
}

func TestQuadraticRealRoots_noAllocation(t *testing.T) {
	if allocs := testing.AllocsPerRun(100, func() { benchSinkFloat, _, _, _ = QuadraticRealRoots(1, -3, 2) }); allocs != 0 {
		t.Errorf("QuadraticRealRoots allocates %v times per call, want 0", allocs)
	}
}

func BenchmarkQuadraticRealRoots(b *testing.B) {
	b.Run("scalar", func(b *testing.B) {
		b.ReportAllocs()
		a := 1.0
		for b.Loop() {
			benchSinkFloat, _, _, _ = QuadraticRealRoots(a, -1e8, 1)
		}
	})
}
