package series

import (
	"errors"
	"math"
	"testing"
)

// ar1Series returns n values of x_t = c + φ·x_{t−1} + ε_t, ε_t standard normal from the fixed-seed source
// newRand(), started at the process mean c/(1 − φ) (|φ| < 1). Test-only generator.
func ar1Series(n int, c, phi float64) []float64 {
	r := newRand()
	x := make([]float64, n)
	prev := c / (1 - phi)
	for i := range x {
		prev = c + phi*prev + r.NormFloat64()
		x[i] = prev
	}
	return x
}

// Hand-worked from BJRL §2.1.5 (divisor N): 1, 2, 3, 4 → z̄ = 2.5, deviations −1.5, −0.5, 0.5, 1.5,
// N·c_0 = 5, N·c_1 = 1.25, N·c_2 = −1.5, N·c_3 = −2.25 → 1, 0.25, −0.3, −0.45. Tolerance 1e-15 (rounding).
func TestAutocorrelation_handWorked(t *testing.T) {
	want := []float64{1, 0.25, -0.3, -0.45}
	for lag, w := range want {
		got, err := Autocorrelation([]float64{1, 2, 3, 4}, lag)
		if err != nil || !closeRel(got, w, 1e-15) {
			t.Errorf("lag %d: = (%v, %v), want %v", lag, got, err, w)
		}
	}
}

func TestAutocorrelation_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		lag    int
		want   float64
	}{
		{"two values", []float64{0, 1}, 1, -0.5},
		{"alternating", []float64{1, -1, 1, -1}, 1, -0.75},
		{"negative values", []float64{-4, -3, -2, -1}, 1, 0.25},
		{"huge scale (rescaled exactly)", []float64{1e300, 2e300, 3e300, 4e300}, 1, 0.25},
		{"tiny scale (rescaled exactly)", []float64{1e-310, 2e-310, 3e-310, 4e-310}, 1, 0.25},
		{"subnormal exact multiples", []float64{0x1p-1060, 0x1p-1059, 3 * 0x1p-1060, 0x1p-1058}, 1, 0.25},
		{"smallest subnormal", []float64{0, 0x1p-1074}, 1, -0.5},
		{"opposite extremes", []float64{-math.MaxFloat64, math.MaxFloat64}, 1, -0.5},
	}
	for _, tc := range tests {
		got, err := Autocorrelation(tc.values, tc.lag)
		if err != nil || !closeRel(got, tc.want, 1e-12) {
			t.Errorf("%s: = (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

func TestAutocorrelation_errors(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		lag    int
		want   error
	}{
		{"empty", nil, 0, ErrTooFew},
		{"one", []float64{1}, 0, ErrTooFew},
		{"NaN", []float64{1, math.NaN()}, 0, ErrNotFinite},
		{"Inf", []float64{math.Inf(-1), 1}, 0, ErrNotFinite},
		{"negative lag", []float64{1, 2}, -1, ErrOutOfRange},
		{"lag = n", []float64{1, 2}, 2, ErrOutOfRange},
		{"constant", []float64{3, 3, 3}, 1, ErrZeroVariance},
	}
	for _, tc := range tests {
		if _, err := Autocorrelation(tc.values, tc.lag); !errors.Is(err, tc.want) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.want)
		}
	}
}

func TestACF_handWorked(t *testing.T) {
	got, err := ACF(nil, []float64{1, 2, 3, 4}, 3)
	want := []float64{1, 0.25, -0.3, -0.45}
	if err != nil || len(got) != len(want) {
		t.Fatalf("= (%v, %v), want %v", got, err, want)
	}
	for i := range want {
		if !closeRel(got[i], want[i], 1e-15) {
			t.Errorf("[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}

// The ACF of a stationary AR(1) is φ^k. With n = 20000 the standard error of r_k is about 1/√n ≈ 0.007
// (Bartlett), so 0.04 is more than five of them.
func TestACF_ar1(t *testing.T) {
	const phi = 0.6
	got, err := ACF(nil, ar1Series(20000, 1, phi), 3)
	if err != nil {
		t.Fatal(err)
	}
	for k := 0; k <= 3; k++ {
		if w := math.Pow(phi, float64(k)); math.Abs(got[k]-w) > 0.04 {
			t.Errorf("r_%d = %v, want ≈ %v", k, got[k], w)
		}
	}
}

// ACF must equal Autocorrelation at every lag, and reuse a long enough dst without allocating.
func TestACF_matchesAutocorrelationAndReusesDst(t *testing.T) {
	values := ar1Series(300, 0, -0.4)
	dst := make([]float64, 0, 21)
	got, err := ACF(dst, values, 20)
	if err != nil {
		t.Fatal(err)
	}
	if &got[0] != &dst[:1][0] {
		t.Errorf("dst storage not reused")
	}
	for k := range got {
		w, _ := Autocorrelation(values, k)
		if got[k] != w {
			t.Errorf("r_%d = %v, Autocorrelation = %v", k, got[k], w)
		}
	}
	if a := testing.AllocsPerRun(10, func() { _, _ = ACF(dst, values, 20) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func TestACF_errors(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		maxLag int
		want   error
	}{
		{"empty", nil, 0, ErrTooFew},
		{"NaN", []float64{math.NaN(), 1}, 1, ErrNotFinite},
		{"negative maxLag", []float64{1, 2}, -1, ErrOutOfRange},
		{"maxLag = n", []float64{1, 2}, 2, ErrOutOfRange},
		{"constant", []float64{0, 0}, 1, ErrZeroVariance},
	}
	for _, tc := range tests {
		if _, err := ACF(nil, tc.values, tc.maxLag); !errors.Is(err, tc.want) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.want)
		}
	}
}

func BenchmarkAutocorrelation(b *testing.B) {
	for _, n := range []int{1000, 10000, 100000} {
		values := ar1Series(n, 0, 0.5)
		b.Run(seriesSizeName(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = Autocorrelation(values, 1)
			}
		})
	}
}

func BenchmarkACF(b *testing.B) {
	for _, n := range []int{1000, 10000, 100000} {
		values := ar1Series(n, 0, 0.5)
		dst := make([]float64, 21)
		b.Run(seriesSizeName(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = ACF(dst, values, 20)
			}
		})
	}
}
