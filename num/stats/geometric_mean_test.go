package stats

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// Hand computations of (∏xᵢ)^(1/n) (NIST Dataplot "GEOMETRIC MEAN"): √(2·8) = 4; ∛27 = 3; one value is
// itself; a zero makes the product 0; 1e300 and 1e-300 → 1;
// 2⁶⁰⁰ and 2⁶⁰⁰ → 2⁶⁰⁰ (product overflows); 2⁻¹⁰⁰⁰·2⁻¹⁰⁶⁰·2⁻¹⁰⁷⁰·2¹⁰⁰⁰ = 2⁻²¹³⁰, fourth root 2⁻⁵³²·⁵ = √2·2⁻⁵³³
// (two subnormal factors); 2⁻³·2⁻⁴ = 2⁻⁷, square root 2⁻³·⁵ = √2/16 (floor division of a negative E).
func TestGeometricMean_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   float64
	}{
		{"two", []float64{2, 8}, 4},
		{"three", []float64{1, 3, 9}, 3},
		{"one", []float64{7.5}, 7.5},
		{"zero", []float64{3, 0, 5}, 0},
		{"extremes", []float64{1e300, 1e-300}, 1},
		{"product overflows", []float64{0x1p600, 0x1p600}, 0x1p600},
		{"equal", []float64{5, 5, 5, 5}, 5},
		{"subnormal factors", []float64{0x1p-1000, 0x1p-1060, 0x1p-1070, 0x1p1000}, math.Ldexp(math.Sqrt2, -533)},
		{"negative odd exponent sum", []float64{0x1p-3, 0x1p-4}, math.Sqrt2 / 16},
	}
	// Tolerance: relative 2e-15 (about 9 ulps). The scaled product adds one rounding per factor, divided by n
	// by the root, plus one rounding each from Pow, Exp2 and the final multiply; Ldexp is exact. Unlike
	// exp(mean(ln x)), the bound does not grow with the magnitude of the values (the 2⁶⁰⁰ case).
	for _, tc := range tests {
		got, err := GeometricMean(tc.values)
		if err != nil || math.Abs(got-tc.want) > 2e-15*math.Abs(tc.want) {
			t.Errorf("%s: GeometricMean(%v) = %v, %v; want %v", tc.name, tc.values, got, err, tc.want)
		}
	}
}

func TestGeometricMean_errors(t *testing.T) {
	if _, err := GeometricMean(nil); !errors.Is(err, ErrEmpty) {
		t.Errorf("empty: %v", err)
	}
	if _, err := GeometricMean([]float64{1, math.NaN()}); !errors.Is(err, ErrNotFinite) {
		t.Errorf("NaN: %v", err)
	}
	if _, err := GeometricMean([]float64{1, math.Inf(1)}); !errors.Is(err, ErrNotFinite) {
		t.Errorf("Inf: %v", err)
	}
	if _, err := GeometricMean([]float64{1, -2}); !errors.Is(err, ErrNegative) {
		t.Errorf("negative: %v", err)
	}
}

// Oracle: the direct product root on small positive values (no overflow at n ≤ 20, values in [0.5, 2)).
func TestGeometricMean_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		n := 1 + r.IntN(20)
		v := make([]float64, n)
		p := 1.0
		for i := range v {
			v[i] = 0.5 + 1.5*r.Float64()
			p *= v[i]
		}
		want := math.Pow(p, 1/float64(n))
		if got, err := GeometricMean(v); err != nil || math.Abs(got-want) > 1e-13*want {
			t.Fatalf("GeometricMean(%v) = %v, %v; want %v", v, got, err, want)
		}
	}
}

func TestGeometricMean_noAlloc(t *testing.T) {
	v := []float64{10, 2, 38, 23, 38, 23, 21}
	if a := testing.AllocsPerRun(100, func() { _, _ = GeometricMean(v) }); a != 0 {
		t.Errorf("allocates %v", a)
	}
}

func BenchmarkGeometricMean(b *testing.B) {
	for _, n := range []int{1000, 10000, 100000} {
		v := make([]float64, n)
		for i := range v {
			v[i] = 1 + float64(i%97)
		}
		b.Run(map[int]string{1000: "1e3", 10000: "1e4", 100000: "1e5"}[n], func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = GeometricMean(v)
			}
		})
	}
}
