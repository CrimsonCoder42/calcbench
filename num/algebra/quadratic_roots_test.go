package algebra

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// Worked examples, hand-checked on the quadratic formula of DLMF §1.11(iii): x² + 2x + 5 has D = −16,
// roots (−2 ± 4i)/2 = −1 ± 2i; x² − 3x + 2 has roots 1 and 2 (imaginary parts exactly 0).
func TestQuadraticRoots_workedExample(t *testing.T) {
	r1, r2, err := QuadraticRoots(1, 2, 5)
	if err != nil || r1 != complex(-1, -2) || r2 != complex(-1, 2) {
		t.Errorf("QuadraticRoots(1, 2, 5) = (%v, %v, %v), want (-1-2i, -1+2i, nil)", r1, r2, err)
	}
	r1, r2, err = QuadraticRoots(1, -3, 2)
	if err != nil || r1 != 1 || r2 != 2 || imag(r1) != 0 || imag(r2) != 0 {
		t.Errorf("QuadraticRoots(1, -3, 2) = (%v, %v, %v), want (1+0i, 2+0i, nil)", r1, r2, err)
	}
}

func TestQuadraticRoots_cases(t *testing.T) {
	tests := []struct {
		name    string
		a, b, c float64
		r1, r2  complex128
	}{
		{"pure imaginary", 1, 0, 4, complex(0, -2), complex(0, 2)},
		{"double root returned twice", 1, -2, 1, 1, 1},
		{"negative a, complex", -1, 2, -5, complex(1, -2), complex(1, 2)},
		{"huge coefficients, complex", 0x1p1000, 2 * 0x1p1000, 5 * 0x1p1000, complex(-1, -2), complex(-1, 2)},
		{"tiny coefficients, complex", 0x1p-1000, 2 * 0x1p-1000, 5 * 0x1p-1000, complex(-1, -2), complex(-1, 2)},
		{"huge roots, complex", 1, 0, 0x1p1000, complex(0, -0x1p500), complex(0, 0x1p500)},
		{"zero roots", 3, 0, 0, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r1, r2, err := QuadraticRoots(tt.a, tt.b, tt.c)
			if err != nil || r1 != tt.r1 || r2 != tt.r2 {
				t.Errorf("QuadraticRoots(%v, %v, %v) = (%v, %v, %v), want (%v, %v, nil)", tt.a, tt.b, tt.c, r1, r2, err, tt.r1, tt.r2)
			}
		})
	}
}

func TestQuadraticRoots_errors(t *testing.T) {
	nan, inf := math.NaN(), math.Inf(1)
	tests := []struct {
		name    string
		a, b, c float64
		want    error
	}{
		{"a = 0", 0, 1, 1, ErrDegenerate},
		{"NaN a", nan, 1, 1, ErrNotFinite},
		{"Inf b", 1, inf, 1, ErrNotFinite},
		{"Inf c", 1, 1, inf, ErrNotFinite},
		{"real root overflows", 0x1p-600, 0x1p600, 1, ErrOverflow},
		{"imaginary part overflows", 0x1p-1070, 0, 0x1p1020, ErrOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r1, r2, err := QuadraticRoots(tt.a, tt.b, tt.c)
			if !errors.Is(err, tt.want) || r1 != 0 || r2 != 0 {
				t.Errorf("QuadraticRoots(%v, %v, %v) = (%v, %v, %v), want (0, 0, %v)", tt.a, tt.b, tt.c, r1, r2, err, tt.want)
			}
		})
	}
}

// Oracle. Real case: the same roots as QuadraticRealRoots. Complex case: the real part is −b/(2a)
// correctly rounded (one division of exact power-of-two-scaled values), and the imaginary part is within 4u
// relative of the exact √(−D)/(2|a|) (math/big): D carries at most 2u (NUM-ALGEB-006), the square root
// halves that and adds u/2, the division adds u/2.
func TestQuadraticRoots_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 600 {
		a := (2*rng.Float64() - 1) * math.Ldexp(1, rng.IntN(200)-100)
		b := (2*rng.Float64() - 1) * math.Ldexp(1, rng.IntN(200)-100)
		c := (2*rng.Float64() - 1) * math.Ldexp(1, rng.IntN(200)-100)
		if a == 0 {
			continue
		}
		r1, r2, err := QuadraticRoots(a, b, c)
		if err != nil {
			t.Fatalf("trial %d: %v", trial, err)
		}
		exactD := bigDiscriminant(a, b, c)
		if exactD.Sign() >= 0 {
			lo, hi, n, _ := QuadraticRealRoots(a, b, c)
			if n == 1 {
				hi = lo
			}
			if r1 != complex(lo, 0) || r2 != complex(hi, 0) {
				t.Fatalf("trial %d: (%v, %v), QuadraticRealRoots gives (%v, %v)", trial, r1, r2, lo, hi)
			}
			continue
		}
		wantRe := bigToFloat(new(big.Float).SetPrec(bigPrec).Quo(newBig(-b), newBig(2*a)))
		im := new(big.Float).SetPrec(bigPrec).Neg(exactD)
		im.Sqrt(im)
		im.Quo(im, newBig(2*math.Abs(a)))
		if real(r1) != wantRe || real(r2) != wantRe || imag(r1) != -imag(r2) || !withinRelative(imag(r2), im, 4*unitRoundoff) {
			t.Fatalf("trial %d: QuadraticRoots(%v, %v, %v) = (%v, %v), want re %v, im %v", trial, a, b, c, r1, r2, wantRe, bigToFloat(im))
		}
	}
}

func TestQuadraticRoots_noAllocation(t *testing.T) {
	if allocs := testing.AllocsPerRun(100, func() { benchSinkComplex, _, _ = QuadraticRoots(1, 2, 5) }); allocs != 0 {
		t.Errorf("QuadraticRoots allocates %v times per call, want 0", allocs)
	}
}

func BenchmarkQuadraticRoots(b *testing.B) {
	b.Run("scalar", func(b *testing.B) {
		b.ReportAllocs()
		a := 1.0
		for b.Loop() {
			benchSinkComplex, _, _ = QuadraticRoots(a, 2, 5)
		}
	})
}
