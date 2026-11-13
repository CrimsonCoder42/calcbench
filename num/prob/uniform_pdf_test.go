package prob

import (
	"math"
	"math/big"
	"testing"
)

// NIST §1.3.6.6.2 (f(x) = 1/(B − A) on [A, B]): a = 0, b = 10, x = 3 → 0.1. Hand-checked.
func TestUniformPDF_sourceExample(t *testing.T) {
	if got, err := UniformPDF(3, 0, 10); err != nil || got != 0.1 {
		t.Errorf("UniformPDF(3, 0, 10) = (%v, %v), want (0.1, nil)", got, err)
	}
}

func TestUniformPDF_cases(t *testing.T) {
	for _, c := range []struct{ x, a, b, want float64 }{{-1, 0, 10, 0}, {11, 0, 10, 0}, {0, 0, 10, 0.1}, {10, 0, 10, 0.1}, {0, -1, 1, 0.5}} {
		if got, err := UniformPDF(c.x, c.a, c.b); err != nil || got != c.want {
			t.Errorf("UniformPDF(%v, %v, %v) = (%v, %v), want %v", c.x, c.a, c.b, got, err, c.want)
		}
	}
}

// b − a overflows float64 but the density 1/(b − a) is representable (subnormal for the widest interval):
// it is computed as ½/(b/2 − a/2). Expected values: a = −1e308, b = 1e308 → 1/(2e308) = 5e-309, from
// big.Rat on the float64 inputs, rounded once; [−Max, Max] → 2⁻¹⁰²⁵·(1 − 2⁻⁵³)⁻¹, which rounds to 2⁻¹⁰²⁵.
func TestUniformPDF_wideInterval(t *testing.T) {
	for _, c := range []struct{ a, b float64 }{{-1e308, 1e308}, {-math.MaxFloat64, math.MaxFloat64}, {-1e308, math.MaxFloat64}} {
		w := new(big.Rat).Sub(new(big.Rat).SetFloat64(c.b), new(big.Rat).SetFloat64(c.a))
		want, _ := new(big.Rat).Inv(w).Float64()
		got, err := UniformPDF(0, c.a, c.b)
		if err != nil || got != want || got == 0 {
			t.Errorf("UniformPDF(0, %v, %v) = (%v, %v), want (%v, nil)", c.a, c.b, got, err, want)
		}
	}
}

func TestUniformPDF_errors(t *testing.T) {
	_, err := UniformPDF(0, 1, 1)
	requireErr(t, "a = b", err, ErrOutOfRange)
	_, err = UniformPDF(0, 2, 1)
	requireErr(t, "a > b", err, ErrOutOfRange)
	_, err = UniformPDF(0, 0, 5e-324)
	requireErr(t, "density overflows", err, ErrOverflow)
	_, err = UniformPDF(nan, 0, 1)
	requireErr(t, "x NaN", err, ErrNotFinite)
	_, err = UniformPDF(0, -inf, 1)
	requireErr(t, "a -Inf", err, ErrNotFinite)
}

func TestUniformPDF_noAlloc(t *testing.T) {
	requireNoAlloc(t, "UniformPDF", func() { sinkFloat, _ = UniformPDF(3, 0, 10) })
}

func BenchmarkUniformPDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = UniformPDF(3, 0, 10)
	}
}
